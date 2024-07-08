# kbot

This project implements a Telegram bot using Go and the gopkg.in/telebot.v3 library.

## Description

This bot provides weather information based on user location and sends random cat pictures upon command.

## Bot exmample

[t.me/alex_m_devops_kbot](t.me/alex_m_devops_kbot)

## Installation

### Prerequisites
- Go installed on your system.
- Obtain a Telegram bot token from the BotFather on Telegram.

1. Set up environment variables:

```bash
export TELE_TOKEN=<your-telegram-bot-token>
```

2. Build and run the bot:

```bash
go build -ldflags "-X="github.com/alex-moiseenko/kbot/cmd.appVersion=v1.0.2
./kbot
```

## Usage

Once the bot is running, interact with it on Telegram:

- Start a chat with bot.
- Send /cat to get a random cat picture.
- Send your location to get weather information.
