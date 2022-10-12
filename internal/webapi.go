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
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/webapi"
	"net/http"
	"time"
)

const (
	DefaultWebApiMode       = gin.ReleaseMode
	DefaultWebApiHost       = ":8080"
	GracefulShutdownTimeout = 10 * time.Second
)

// RouteApiV1 is the path to the route group
// used for the first version of the applications
// web api.
const RouteApiV1 = "/v1"

// engine is the gin.Engine that runs the API
// webserver.
var engine *gin.Engine

// v1ApiRouter is the gin.RouterGroup that holds
// the entire first version of the applications API.
var v1ApiRouter *gin.RouterGroup

// httpServer is the http.Server started by the initialization routine
// of the application.
var httpServer *http.Server

func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, api.Components)
}

// enrichMiddlewares enriches the passed gin.Engine with
// the default middleware for the application.
func enrichMiddlewares(e *gin.Engine) {
	e.Use(WebApiLogger(), gin.Recovery())
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

	v1ApiRouter = engine.Group(RouteApiV1)

	httpServer = &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	webApiLogger.Info("Starting api webserver on \"%s\" in mode %s...", Config.webApiHost, Config.webApiMode)
	go func() {
		err := httpServer.ListenAndServe()
		if nil == err || errors.Is(err, http.ErrServerClosed) {
			webApiLogger.Info("The api webserver has been shutdown!")

			return
		}

		ExitFatal(fmt.Sprintf("The api webserver quit unexpectedly: %v", err))
	}()

	err := webapi.Init(v1ApiRouter)
	if nil != err {
		ExitFatal(fmt.Sprintf("Failed to initialize the api framework for the web api: %v", err))
	}
}

// shutdownApiWebserver tries to gracefully shut down
// the api webserver run by the bot.
func shutdownApiWebserver() {
	if nil == httpServer {
		webApiLogger.Warn("Tried to shut down the apüi webserver before it has been started!")

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	webApiLogger.Info("Gracefully shutting down api webserver...")

	if err := httpServer.Shutdown(ctx); err != nil {
		webApiLogger.Err(err, "Failed to gracefully shutdown the api webserv!")

		return
	}

	webApiLogger.Info("Successfully shutdown the api webserver!")
}
