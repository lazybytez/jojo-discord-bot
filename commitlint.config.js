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

// Configuration of commitlint to check commit message guidelines
module.exports = {
    extends: ['@commitlint/config-conventional'],
    parserPreset: 'conventional-changelog-conventionalcommits',
    rules: {
        'subject-max-length': [2, 'always', 50],
        'scope-enum': [2, 'always', [
            'deps', // Changes done on anything dependency related
            'devops', // Changes done on technical processes
            'api', // Changes to the public api
            'comp', // Changes to feature components
            'int', // Changes to internal stuff
            'serv', // Changes to the services that sit between internal and public api
            'core' // Changes on files in project root
        ]]
    }
};
