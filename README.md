# JOJO Discord Bot

![commit-info][commit-info]
![contributors-info][contributors-info]
![reposize-info][reposize-info]
![stars][stars]

## Description

This is an open source Discord bot mainly developed by Lazy Bytez.  
If you want to take part in the development of the bot please check out
the [Contributing](https://github.com/lazybytez/jojo-discord-bot#contributing) section.

Open source doesn't mean everyone can do whatever they want with the bot so there is a
strict [LICENSE](https://github.com/lazybytez/jojo-discord-bot/blob/main/LICENSE) we want you to respect.

## Getting started

### Requirements

1. [Go 1.18](https://go.dev/doc/install)
2. Git

### Create .env

```bash
cp .env.example .env
```

If you want to test, always insert the Token of your TestBot in the `.env`, because this file won't be committed to git.

### Install dependencies

```bash
go get
```

### Compile and run (for development usage)

```bash
go run .
```

### Build binary (for production usage)

```bash
go build .
```

## Contributing

If you want to take part in contribution, like fixing issues and contributing directly to the code base, please visit
the [How to Contribute][github-contribute] document.

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
'api', // Changes in /api/ directory
'comp', // Changes in /component/ directory
'int', // Changes in /internal/ directory
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

### Recommended IDEs

- [GoLand](https://www.jetbrains.com/de-de/go/) (paid)
- [Visual Studio Code](https://code.visualstudio.com/) (free)
  with [Go Language Extension](https://marketplace.visualstudio.com/items?itemName=golang.go) (free)

## Useful links

[License][github-license] -
[Contributing][github-contribute] -
[Code of conduct][github-codeofconduct] -
[Issues][github-issues] -
[Pull requests][github-pulls]

<hr>  

###### Copyright (c) [Lazy Bytez][github-team]. All rights reserved | Licensed under the AGPL-3.0 license.

<!-- Variables -->

[github-team]: https://github.com/lazybytez

[github-license]: https://github.com/lazybytez/general-template/blob/main/LICENSE

[github-contribute]: https://github.com/lazybytez/.github/blob/main/docs/CONTRIBUTING.md

[github-codeofconduct]: https://github.com/lazybytez/.github/blob/main/docs/CODE_OF_CONDUCT.md

[github-issues]: https://github.com/lazybytez/general-template/issues

[github-pulls]: https://github.com/lazybytez/general-template/pulls

[commit-info]: https://img.shields.io/github/last-commit/lazybytez/general-template?style=for-the-badge&colorA=302D41&colorB=b4befe

[contributors-info]: https://img.shields.io/github/contributors/lazybytez/general-template?style=for-the-badge&colorA=302D41&colorB=cba6f7

[reposize-info]: https://img.shields.io/github/repo-size/lazybytez/general-template?style=for-the-badge&colorA=302D41&colorB=89dceb

[stars]: https://img.shields.io/github/stars/lazybytez/jojo-discord-bot?colorA=302D41&colorB=f9e2af&style=for-the-badge
