# donezo

![Go](https://img.shields.io/badge/Go-1.25-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![CI](https://github.com/rhajizada/donezo/actions/workflows/ci.yml/badge.svg)](https://github.com/rhajizada/donezo/actions/workflows/ci.yml)
![coverage](https://signum.rhajizada.dev/api/badges/9fc8a30c-fcb9-4faa-904a-d62f59f4e9cf)

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
