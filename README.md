# Primitive Bot

![example](https://user-images.githubusercontent.com/49400499/120475779-78032880-c3b2-11eb-98ba-302ce8190f9d.gif)

## What is it

This is a Telegram Bot that uses the [fogleman/primitive](https://github.com/fogleman/primitive) library to reproduce
images sent to it using various geometric shapes.

Main features:

- Inline menu for setting desired options.
- Doesn't use a database. The queue can be restored from the logs.
- Sessions are stored in memory and cleared after some time of inactivity (30 minutes by default).

## Installation

If you have installation of Go on your machine:

```shell
go get github.com/lazy-void/primitive-bot/cmd/primitive-bot
```

Otherwise, you can download the precompiled binaries from
the [release](https://github.com/lazy-void/primitive-bot/releases) page.

## Usage

Basic Usage:

```shell
primitive-bot -token=$BOT_TOKEN
```

If you want to restore the queue from the logs, use the `-log` flag:

```shell
primitive-bot -token=$BOT_TOKEN -log=path/to/log.txt
```

Full list of options:

```commandline
  -i string
        Path to the directory where user-supplied images are stored. (default "inputs")
  -lang value
        Language of the bot (en, ru). (default "en")
  -limit int
        The number of operations that the user can add to the queue. (default 5)
  -log string
        Path to the previous log file. It is used to restore queue.
  -o string
        Path to the directory where resulting images are stored. (default "outputs")
  -size int
        The max value of image size that the user can specify. (default 3840)
  -steps int
        The max value of steps that the user can specify. (default 2000)
  -timeout duration
        The period of time that a session can be inactive before it's terminated. (default 30m0s)
  -token string
        The token for the Telegram Bot.
  -w int
        The number of parallel workers used to create a primitive image. (default to number of CPUs)
```
