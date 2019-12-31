# Bare Minimum

Minimalistic text editor.

## Status

[![Actions](https://github.com/grzegorz-zur/bm/workflows/Test/badge.svg)](https://github.com/grzegorz-zur/bm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/grzegorz-zur/bm)](https://goreportcard.com/report/github.com/grzegorz-zur/bm)
[![codecov](https://codecov.io/gh/grzegorz-zur/bm/branch/master/graph/badge.svg)](https://codecov.io/gh/grzegorz-zur/bm)

## Goals

1.  Effectively work with multiple files.
2.  Effectively work with external tools that modify edited files.

## Installation

To install or update run the following command.

```sh
go get -u github.com/grzegorz-zur/bm
```

## Usage

```sh
bm
```

### Command mode

Execute commands.

![command mode](command.svg "Command mode")

### Input mode

Type to input text in the current file.

![input mode](input.svg "Input mode")

### Select mode

Select text and performs operation on the selection

![select mode](select.svg "Select mode")

### Switch mode

Type to filter files. Navigate with cursor keys and press enter to select a file.

![switch mode](switch.svg "Switch mode")

