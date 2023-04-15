<div align="center">

# JOJO Discord Bot

[![go-ref-badge][go-ref-badge]][go-ref]
[![gh-license-badge][gh-license-badge]][gh-license]
[![discord-badge][discord-badge]][discord]

[![codecov-badge][codecov-badge]][codecov]
[![gh-contributors-badge][gh-contributors-badge]][gh-contributors]
[![gh-stars-badge][gh-stars-badge]][gh-stars]

</div>

## Description

This is an open source Discord bot mainly developed by [Lazy Bytez][gh-team].  
If you want to take part in the development of the bot please check out
the [Contributing](https://github.com/lazybytez/jojo-discord-bot#contributing) section.

Open source doesn't mean everyone can do whatever they want with the bot so there is a
strict [LICENSE](https://github.com/lazybytez/jojo-discord-bot/blob/main/LICENSE) we want you to respect.

## Getting started

### Requirements

1. [Go 1.19](https://go.dev/doc/install)
2. Git
3. Docker
4. Make 

### Setup

Copies env and installs dependencies

```bash 
make setup
```

Copies env.example to .env

```bash 
make env
```

Installs go dependencies needed to run the bot (like discordgo)

```bash 
make install
```

### Running

Run your code to test and for development purposes.

```bash
make run 
```

Build executable for production usage.

```bash 
make build 
```

### QA

Shows test and codecov results.

```bash
make test
```

Local linting to assure code styling.

```bash
make lint
```

### Database services
To use the local PostgreSQL and Redis it is necesary to have `Docker` and `docker-compose`
installed locally.

Start local database and redis:
```bash
make services/start
```

Stop local database and redis:
```bash
make services/stop
```

Destroy local database and redis:
```bash
make services/start
```

## Contributing

If you want to take part in contribution, like fixing issues and contributing directly to the code base, please visit
the [How to Contribute][gh-contribute] document.

### Commit messages

Construct of a commit message:

```bash
prefix(scope): commit subject with max 50 chars
```

Example commit message:

```bash
feat(comp): add ping slash command
```

#### Scopes

Project specific scopes and what to use them for.

```bash
'deps', // Changes done on anything dependency related
'devops', // Changes done on technical processes
'api', // Changes to the public api
'comp', // Changes to feature components
'int', // Changes to internal stuff
'serv', // Changes to the services sit between internal and public api
'core' // Changes on files in project root
```

#### Prefixes:

Also see [CONTRIBUTING.md#commits](https://github.com/lazybytez/.github/blob/main/docs/CONTRIBUTING.md#commits)

```bash
'feat', // Some new feature that has been added
'fix', // Some fixes to an existing feature
'build', // Some change on how the project is built
'chore', // Some change that just has to be done (like updating dependencies)
'ci', // Some changes to the continues integration workflows
'docs', // Some changes to documentation located in the repo (either markdown files or code DocBlocks)
'perf', // Some performance improvements
'refactor', // Some code changes, that neither adds functionality or fixes a bug
'revert', // Some changes that revert already done changes
'style', // Some fixes regarding code style
'test', // Some automated tests that have been added
```

#### Branches:

| Branch     | Usage                                  |
|------------|----------------------------------------|
| main       | The default branch                     |
| feature/*  | For developing features                |
| fix/*      | For fixing bugs                        |

### Recommended IDEs

- [GoLand](https://www.jetbrains.com/de-de/go/) (paid)
- [Visual Studio Code](https://code.visualstudio.com/) (free)
  with [Go Language Extension](https://marketplace.visualstudio.com/items?itemName=golang.go) (free)

## Useful links

[License][gh-license] -
[Contributing][gh-contribute] -
[Code of conduct][gh-codeofconduct] -
[Issues][gh-issues] -
[Pull requests][gh-pulls]

<hr>  

###### Copyright (c) [Lazy Bytez][gh-team]. All rights reserved | Licensed under the AGPL-3.0 license.

<!-- Variables -->

[go-ref-badge]: https://img.shields.io/badge/godoc-reference-89dceb?style=for-the-badge&colorA=302D41&logo=go

[go-ref]: https://pkg.go.dev/github.com/lazybytez/jojo-discord-bot

[gh-license-badge]: https://img.shields.io/github/license/lazybytez/jojo-discord-bot?logo=gnu&style=for-the-badge&colorA=302D41&colorB=eba0ac

[gh-license]: https://github.com/lazybytez/jojo-discord-bot/blob/main/LICENSE

[discord-badge]: https://img.shields.io/discord/735171597362659328?label=Discord&logo=discord&style=for-the-badge&colorA=302D41&colorB=b4befe

[discord]: https://discord.gg/bcV6TN2k9V

[codecov-badge]: https://img.shields.io/codecov/c/github/lazybytez/jojo-discord-bot?style=for-the-badge&colorA=302D41

[codecov]: https://app.codecov.io/gh/lazybytez/jojo-discord-bot

[gh-contributors-badge]: https://img.shields.io/github/contributors/lazybytez/jojo-discord-bot?style=for-the-badge&colorA=302D41&colorB=cba6f7

[gh-contributors]: https://github.com/lazybytez/jojo-discord-bot/graphs/contributors

[gh-stars-badge]: https://img.shields.io/github/stars/lazybytez/jojo-discord-bot?colorA=302D41&colorB=f9e2af&style=for-the-badge

[gh-stars]: https://github.com/lazybytez/jojo-discord-bot/stargazers

[gh-contribute]: https://github.com/lazybytez/.github/blob/main/docs/CONTRIBUTING.md

[gh-codeofconduct]: https://github.com/lazybytez/.github/blob/main/docs/CODE_OF_CONDUCT.md

[gh-issues]: https://github.com/lazybytez/jojo-discord-bot/issues

[gh-pulls]: https://github.com/lazybytez/jojo-discord-bot/pulls

[gh-team]: https://github.com/lazybytez
