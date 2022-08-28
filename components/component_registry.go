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

package components

import (
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/lazybytez/jojo-discord-bot/api/components/bot_core"
	"github.com/lazybytez/jojo-discord-bot/api/components/bot_log"
	"github.com/lazybytez/jojo-discord-bot/components/pingpong"
	"github.com/lazybytez/jojo-discord-bot/components/statistics"
)

// Components contains all components that should be available.
//
// Enabled components should be registered here.
// When access to components is necessary use api.Components instead.
// Note that api.Components can only being accessed after the system has been initialized,
// which means the earliest point is in the LoadComponent lifecycle hooks.
var Components = []*api.Component{
	bot_core.C,
	bot_log.C,
	pingpong.C,
	statistics.C,
}
