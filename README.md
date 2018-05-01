# Bare Minimum

[![CircleCI](https://circleci.com/gh/grzegorz-zur/bare-minimum.svg?style=svg)](https://circleci.com/gh/grzegorz-zur/bare-minimum)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b395ffed1b5c4a06a54f1416c08362b7)](https://www.codacy.com/app/grzegorz.zur/bare-minimum?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=grzegorz-zur/bare-minimum&amp;utm_campaign=Badge_Grade)
[![codebeat badge](https://codebeat.co/badges/8f64fe34-b32e-4ba5-a391-c02669f08b38)](https://codebeat.co/projects/github-com-grzegorz-zur-bare-minimum-master)
[![codecov](https://codecov.io/gh/grzegorz-zur/bare-minimum/branch/master/graph/badge.svg)](https://codecov.io/gh/grzegorz-zur/bare-minimum)

Minimalistic text editor for GNU/Linux.

## Goals

1. Effectively work with multiple files.
2. Effectively work with external tools that modify edited files.

## Usage

### Normal mode

![normal mode](keyboards/normal.svg "Normal mode")

### Commands

* ^q — quit editor
* ^w — close current file
* ^e — write current file
* ^z — stop the editor (background)
* ^d — switch to previous file
* ^f — switch to next file

### Input mode

Type to input text in the current file.

### Switch mode

Type to filter files. Navigate with cursor keys and press enter to select a file.
