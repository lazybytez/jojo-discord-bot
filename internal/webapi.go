/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lazybytez/jojo-discord-bot/docs"
	"github.com/lazybytez/jojo-discord-bot/webapi"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultWebApiMode       = gin.ReleaseMode
	DefaultWebApiBind       = ":8080"
	DefaultWebApiHost       = "localhost:8080"
	DefaultWebApiBasePath   = "/"
	DefaultWebApiSchemes    = "https,http"
	GracefulShutdownTimeout = 10 * time.Second
)

// The root routes that are available on the running bot.
const (
	RouteApiV1        = "/v1"
	RouteSwagger      = "/swagger"
	RouteSwaggerIndex = "/swagger/index.html"
)

// engine is the gin.Engine that runs the API
// webserver.
var engine *gin.Engine

// v1ApiRouter is the gin.RouterGroup that holds
// the entire first version of the applications API.
var v1ApiRouter *gin.RouterGroup

// httpServer is the http.Server started by the initialization routine
// of the application.
var httpServer *http.Server

// handlePanic handles the response when a panic occurs.
// Unlike the default recovery function of Gin, this function returns a JSON response.
func handlePanic(g *gin.Context, recovered interface{}) {
	if Config.webApiMode == gin.DebugMode {
		if err, ok := recovered.(error); ok {
			webapi.RespondWithError(g, webapi.ErrorResponse{
				Status:    http.StatusInternalServerError,
				Error:     "An unexpected error occurred",
				Message:   err.Error(),
				Timestamp: time.Now(),
			})

			return
		}

		if err, ok := recovered.(string); ok {
			webapi.RespondWithError(g, webapi.ErrorResponse{
				Status:    http.StatusInternalServerError,
				Error:     "An unexpected error occurred",
				Message:   err,
				Timestamp: time.Now(),
			})

			return
		}
	}

	webapi.RespondWithError(g, webapi.ErrorResponse{
		Status: http.StatusInternalServerError,
		Error:  "An unexpected error occurred",
		Message: "An unexpected error occurred, " +
			"please contact the administrator of the bot to obtain further information.",
		Timestamp: time.Now(),
	})
}

// handleNoRoute handles requests which do not match any route.
func handleNoRoute(g *gin.Context) {
	path := url.PathEscape(g.Request.URL.Path)
	webapi.RespondWithError(g, webapi.ErrorResponse{
		Status:    http.StatusNotFound,
		Error:     "Resource not found",
		Message:   fmt.Sprintf("There is no resource matching the path \"%s\"", path),
		Timestamp: time.Now(),
	})
}

// enrichMiddlewares enriches the passed gin.Engine with
// the default middleware for the application.
func enrichMiddlewares(e *gin.Engine) {
	e.Use(WebApiLogger(), gin.CustomRecovery(handlePanic))
}

// addSwaggerRedirect adds a middleware that redirects calls to "/swagger"
// to "/swagger/index.html".
// This way it is easier to access the swagger.
func addSwaggerRedirect(originalHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(strings.TrimSuffix(c.Request.RequestURI, "/"), buildRoutePath(RouteSwagger)) {
			c.Redirect(http.StatusMovedPermanently, buildRoutePath(RouteSwaggerIndex))

			return
		}

		originalHandler(c)
	}
}

// configureGeneralSwaggerMeta configures some general options
// like the API base path or the current version.
func configureGeneralSwaggerMeta(basePath string) {
	docs.SwaggerInfo.Host = Config.webApiHost
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = strings.Split(Config.webApiSchemes, ",")
}

// initSwagger cares about initializing the Swagger
// that can be used to find information about the web api.
func initSwagger() {
	configureGeneralSwaggerMeta(buildRoutePath(RouteApiV1))
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	routeGroup := engine.Group(buildRoutePath(RouteSwagger))
	routeGroup.GET("/*any", addSwaggerRedirect(swaggerHandler))
}

// initWebApi initializes the web api.
// This means:
//   - the webserver (gin) is started
//   - default routes are registered
//   - the API is prepared to be used by components
func initWebApi() {
	if nil != httpServer {
		ExitFatal("Tried to initialize the api webserver more than once!")
	}

	gin.SetMode(Config.webApiMode)

	engine = gin.New()
	enrichMiddlewares(engine)
	engine.NoRoute(handleNoRoute)

	v1ApiRouter = engine.Group(buildRoutePath(RouteApiV1))

	httpServer = &http.Server{
		Addr:    Config.webApiBind,
		Handler: engine,
	}

	webApiLogger.Info("Starting api webserver on \"%s\" in mode %s...", Config.webApiBind, Config.webApiMode)
	go func() {
		err := httpServer.ListenAndServe()
		if nil == err || errors.Is(err, http.ErrServerClosed) {
			webApiLogger.Info("The api webserver has been shutdown!")

			return
		}

		ExitFatal(fmt.Sprintf("The api webserver quit unexpectedly: %v", err))
	}()

	initSwagger()
	err := webapi.Init(v1ApiRouter)
	if nil != err {
		ExitFatal(fmt.Sprintf("Failed to initialize the api framework for the web api: %v", err))
	}
}

// buildRoutePath prepares the path of a route to be registered
func buildRoutePath(route string) string {
	return fmt.Sprintf("%s/%s",
		strings.TrimSuffix(Config.webApiBasePath, "/"),
		strings.TrimPrefix(route, "/"))
}

// shutdownApiWebserver tries to gracefully shut down
// the api webserver run by the bot.
func shutdownApiWebserver() {
	if nil == httpServer {
		webApiLogger.Warn("Tried to shut down the api webserver before it has been started!")

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	webApiLogger.Info("Gracefully shutting down api webserver...")

	if err := httpServer.Shutdown(ctx); err != nil {
		webApiLogger.Err(err, "Failed to gracefully shutdown the api webserver!")

		return
	}

	webApiLogger.Info("Successfully shutdown the api webserver!")
}
