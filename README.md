# Telegram_C2

![Language](https://img.shields.io/badge/language-Go-00ADD8?style=flat-square&logo=go) ![Platform](https://img.shields.io/badge/platform-Windows-0078D6?style=flat-square&logo=windows) ![Purpose](https://img.shields.io/badge/purpose-Educational-yellow?style=flat-square)

## ⚠️ Disclaimer

> **This project is intended strictly for educational and research purposes.**
> Use only in controlled lab environments with explicit permission. The author assumes no responsibility for misuse of this code.

## Description

A simple Command & Control (C2) proof-of-concept written in Go that uses a Telegram bot as its command server. The target agent registers itself as a system service, awaits instructions sent via Telegram messages, and reports results back to the operator through the same channel. It was built as a study exercise to understand how threat actors leverage legitimate communication platforms for C2 infrastructure.

## Features

- **Telegram bot integration** — uses a Telegram bot as the C2 communication channel
- **Remote command execution** — runs arbitrary shell/system commands on the target machine and returns output
- **Screenshot capture** — takes screenshots of the target's desktop and sends them via Telegram
- **Service installation** — registers itself as a Windows system service for persistence
- **Charset normalization** — handles Windows codepage output encoding to avoid garbled text
- **Environment-based configuration** — bot token and chat ID loaded from a `.env` file

## Tech Stack

| Technology | Purpose |
|---|---|
| [Go](https://go.dev/) | Core language |
| [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api) | Telegram bot communication |
| [kardianos/service](https://github.com/kardianos/service) | Cross-platform service installation |
| [vova616/screenshot](https://github.com/vova616/screenshot) | Desktop screenshot capture |
| [joho/godotenv](https://github.com/joho/godotenv) | `.env` file loading |
| [go-errors/errors](https://github.com/go-errors/errors) | Enhanced error handling |
| [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) | IANA charset encoding/decoding |

## Prerequisites

- [Go](https://go.dev/dl/) 1.18 or later
- A Telegram bot token — create one via [@BotFather](https://t.me/BotFather)
- Your Telegram chat ID (obtainable via [@userinfobot](https://t.me/userinfobot))
- For cross-compilation to Windows from Linux/macOS:
  - `mingw-w64` toolchain (`apt install mingw-w64` on Debian/Ubuntu)

## Getting Started

**1. Clone the repository**

```bash
git clone https://github.com/MaximilianoMeyer/Telegram_C2.git
cd Telegram_C2
```

**2. Install Go dependencies**

```bash
go mod tidy
```

**3. Configure environment variables**

Create a `.env` file in the project root:

```bash
cp .env.example .env   # if available, otherwise create manually
```

```ini
TELEGRAM_TOKEN=your_bot_token_here
CHAT_ID=your_chat_id_here
```

**4. Build for Windows (native)**

```bash
go build -o c2_agent.exe c2.go
```

**5. Build for Windows from Linux/macOS (cross-compile)**

```bash
chmod +x compile_windows.sh
./compile_windows.sh
```

## Usage

**Deploy and run directly on a Windows target (for lab testing):**

```powershell
.\c2_agent.exe
```

**Install as a Windows service (requires elevated privileges):**

```powershell
.\c2_agent.exe install
.\c2_agent.exe start
```

**Send commands via Telegram:**

Once the agent is running, open your Telegram bot and send commands directly:

| Command | Description |
|---|---|
| `whoami` | Returns the current user |
| `dir` / `ls` | Lists directory contents |
| `screenshot` | Captures and sends a desktop screenshot |
| Any shell command | Executes and returns stdout output |

Output is sent back to your configured Telegram chat automatically.

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `TELEGRAM_TOKEN` | Bot token obtained from @BotFather | ✅ Yes |
| `CHAT_ID` | Telegram chat ID of the operator | ✅ Yes |

## Project Structure

```
Telegram_C2/
├── c2.go                 # Main agent — service registration, bot listener, command execution
├── compile_windows.sh    # Shell script for cross-compiling to Windows via mingw-w64
├── .gitignore            # Git ignore rules
└── README.md             # Project documentation
```

## Contributing

This is a personal study project, but feedback and educational contributions are welcome.

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Commit your changes: `git commit -m 'Add your feature'`
4. Push to the branch: `git push origin feat/your-feature`
5. Open a Pull Request

Please ensure any contribution remains within the scope of **educational use only**.

## License

No license is currently specified for this project. All rights reserved by the author unless otherwise stated. Contact the repository owner before reusing or distributing this code.

---

<p align="center">
  <a href="https://readmebot.dev?utm_source=readme_badge&utm_medium=github">
    <img src="https://readmebot.dev/badge.svg" alt="Docs by ReadmeBot" height="20">
  </a>
</p>
