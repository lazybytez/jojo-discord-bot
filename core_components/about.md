## Do not put components here
The components located in the API package are core components
and only cover functionality that is deeply integrated in the bots API.
This includes core logging like when the bot joins/leaves guilds
or ensuring guilds are properly registered in the database.

Changing things here could seriously break other parts of the bot.

If you plan to extend the bot, use the components folder at the root of
the repository!
