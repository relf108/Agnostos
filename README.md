# Agnostos

_A Language agnostic Env manager_

<!--toc:start-->

- [Agnostos](#agnostos)
  - [Goals](#goals)
  - [Features](#features)
  - [Future Features](#future-features)
  <!--toc:end-->

## Goals

- Replace language specific env managers like `Conda`, `NVM`, etc.
- Improve development -> production environment parity.

## Features

- [ ] Uses full docker environment allowing us to set an env for bash and get a full env specific cli.
- [ ] Support for major bash alternatives (zsh, fish, etc).
- [ ] bashrc/zshrc type file for setting up default (or per-environment) cli.
  - [ ] Point to local path containing dotfiles to set default (or per-environment) dotfiles.
  - [ ] Point to local path containing rc files to set default (or per-environment) rc files.
- [ ] Export env to dockerfile.
- [ ] Support remote dev via VScode/Jetbrains IDEs

## Future Features

- [ ] Optional remote dev env setup (minimal) or full dev env (mount local dirs)
