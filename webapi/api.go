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

package webapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// routerGroup is the base gin.RouterGroup
// used to create API endpoints.
var routerGroup *gin.RouterGroup

// Init initializes the webapi and makes
// it ready to be used.
func Init(apiRouterGroup *gin.RouterGroup) error {
	if nil != routerGroup {
		return fmt.Errorf("cannot initialize the web api twice")
	}

	routerGroup = apiRouterGroup

	return nil
}

// Router returns the root gin.RouterGroup that should
// be used to register new router groups and routes.
//
// This function should be used by components to create API endpoints.
func Router() *gin.RouterGroup {
	return routerGroup
}
