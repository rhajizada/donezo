# donezo

![Go](https://img.shields.io/badge/Go-1.22-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![ci](https://github.com/rhajizada/donezo/actions/workflows/ci.yml/badge.svg)

**donezo** is a simple TUI to-do app in written Go using
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and SQLite.

## Features

- Task Management: Organize tasks into boards, each with individual items.
- TUI Interface: An interactive command-line UI built with Bubble Tea.
- SQLite Database: Data is stored locally in an SQLite database.
- Boards and Items: Create, update, delete, and list boards and items, with
  support for toggling item completion status.
- Tags: Tag, un-tag items, view items by tags.

## Installation

```bash
git clone https://github.com/rhajizada/donezo.git
cd donezo
make install
```

## Roadmap

- 1.0
  - Implement comprehensive test coverage for all modules.
  - Set up CI using GitHub Actions.
  - Automate releases and publish them via GitHub Actions.
  - Add support for custom app styling.
  - Introduce configuration options for app customization.
