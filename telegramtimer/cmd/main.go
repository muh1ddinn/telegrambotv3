package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"timer/teleram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		botToken = "6765785166:AAGXaCcfpzpUR487nJ6OI1wEUub6XHyNeSg"
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Printf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil {
				teleram.Commands(bot, update)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Telegram Bot is running...")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
