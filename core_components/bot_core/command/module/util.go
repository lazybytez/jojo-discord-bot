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

package module

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lazybytez/jojo-discord-bot/api/entities"
)

// findComponent tries to find a specific component by its code.
func findComponent(option *discordgo.ApplicationCommandInteractionDataOption) *entities.RegisteredComponent {
	var componentCode entities.ComponentCode

	switch v := option.Options[0].Value.(type) {
	case string:
		componentCode = entities.ComponentCode(v)
	default:
		return nil
	}

	for _, c := range C.EntityManager().RegisteredComponent().GetAvailable() {
		if c.Code == componentCode {
			return c
		}
	}

	return nil
}
