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

package api

import "github.com/bwmarrin/discordgo"

// SyncApplicationComponentGlobalCommands ensures that the available discordgo.ApplicationCommand
// are synced for the given component globally.
//
// This means that disabled commands are enabled and enabled commands are disabled
// depending on their global enable state.
//
// Also orphaned commands are cleaned up.
// This is executed whenever a guild is joined or a component is toggled.
//
// Sync is a four-step process:
//   - remove orphaned commands
//   - remove disabled commands
//   - add new commands
//   - update existing commands
func (c *SlashCommandManager) SyncApplicationComponentGlobalCommands(
	session *discordgo.Session,
) {
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, "")
	if nil != err {
		slashCommandManagerLogger.Err(err, "Failed to handle global slash-command sync!")

		return
	}

	slashCommandManagerLogger.Info("Syncing slash-commands globally...")
	registeredCommands = c.removeOrphanedCommands(session, "", registeredCommands)
	registeredCommands = c.removeCommandsByComponentState(session, "", registeredCommands)
	registeredCommands = c.addCommandsByComponentState(session, "", registeredCommands)
	_ = c.updateRegisteredCommands(session, "", registeredCommands)

	slashCommandManagerLogger.Info(
		"Finished syncing slash-commands globally...")
}
