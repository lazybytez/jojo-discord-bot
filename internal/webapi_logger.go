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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lazybytez/jojo-discord-bot/services/logger"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

// webApiLoggerPrefix is the prefix appended to all log messages
// that are logged by Gin.
const webApiLoggerPrefix = "web_api"

// webApiLogger is the logger used for logging
// in Gin.
var webApiLogger *logger.Logger

func init() {
	webApiLogger = logger.New(webApiLoggerPrefix, nil)
}

// WebApiLogger creates a new gin.HandlerFunc that
// acts as a logger middleware and replaces the original logger.
// This is done to match our standard logging format and support our standard
// api.
func WebApiLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)
		clientIP := c.ClientIP()
		requestMethod := c.Request.Method
		statusCode := c.Writer.Status()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// Create log entry
		baseEvent := findLoggingBaseEvent(c)
		baseEvent.Str(logger.ComponentLogPrefix, webApiLoggerPrefix)

		// Enrich fields with data
		baseEvent.Str("method", requestMethod)
		baseEvent.Str("client-ip", clientIP)
		baseEvent.Int("status", statusCode)
		baseEvent.Str("route", path)

		strErrorPart := "without any error"
		if errMsg != "" {
			strErrorPart = fmt.Sprintf("with error message: \"%s\"", errMsg)
		}

		baseEvent.Msgf("served %d bytes in %s %s", bodySize, latency, strErrorPart)
	}
}

func findLoggingBaseEvent(c *gin.Context) *zerolog.Event {
	responseStatusCode := c.Writer.Status()
	switch {
	case responseStatusCode >= http.StatusOK && responseStatusCode < http.StatusMultipleChoices:
		return webApiLogger.Logger().Info()
	case responseStatusCode >= http.StatusMultipleChoices && responseStatusCode < http.StatusBadRequest:
		return webApiLogger.Logger().Info()
	case responseStatusCode >= http.StatusBadRequest && responseStatusCode < http.StatusInternalServerError:
		return webApiLogger.Logger().Warn()
	default:
		return webApiLogger.Logger().Error()
	}
}
