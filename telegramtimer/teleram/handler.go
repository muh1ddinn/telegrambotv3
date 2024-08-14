package teleram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"timer/model"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbot.BotAPI, update tgbot.Update) {
	keyboard := CmdKeyboard()
	msg := tgbot.NewMessage(update.Message.Chat.ID, "Welcome! Please choose one:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
	UserState.Store(update.Message.Chat.ID, "awaiting_choice")
}

func HandleUserResponse(bot *tgbot.BotAPI, update tgbot.Update) {
	userID := update.Message.Chat.ID
	state, _ := UserState.Load(userID)

	switch state {
	case "awaiting_choice":
		choice := strings.ToLower(update.Message.Text)
		if choice == "prayer_times" {
			HandlePrayerTimes(bot, update)
		} else if choice == "prayer_notifications" {
			HandlePrayerNotifications(bot, update)
		} else {
			msg := tgbot.NewMessage(userID, "Invalid choice. Please select either /prayer_times or /prayer_notifications.")
			bot.Send(msg)
		}
	case "selecting_region":
		HandleRegionSelection(bot, update)
	default:
		msg := tgbot.NewMessage(userID, "Unknown command. Please start with /start.")
		bot.Send(msg)
	}
}

// HandlePrayerTimes shows the list of available regions
func HandlePrayerTimes(bot *tgbot.BotAPI, update tgbot.Update) {
	regions := "Available regions:\nAndijon\nBuxoro\nJizzax\nQarshi\nNavoiy\nNamangan\nSamarqand\nGuliston\nTermiz\nToshkent\nQo'qon\nUrganch"
	msg := tgbot.NewMessage(update.Message.Chat.ID, regions)
	bot.Send(msg)
	UserState.Store(update.Message.Chat.ID, "selecting_region")
}

func HandleRegionSelection(bot *tgbot.BotAPI, update tgbot.Update) {
	region := update.Message.Text
	nomoztime, err := Text(region)
	if err != nil {
		msg := tgbot.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error: %s", err))
		bot.Send(msg)
		return
	}

	response := fmt.Sprintf(
		"🕘 Namoz time for %s on %s\n🗓 Hafta: %s\n🗓 Hijri Date: Month: %s \n🌌 Bomdod: %s\n🌅 Quyosh: %s\n🏞 Peshin: %s\n🌇 Asr: %s\n🏙 Shom Iftor: %s\n🌃 Xufton: %s",
		nomoztime.Region,
		nomoztime.Date,
		nomoztime.Weekday,
		nomoztime.HijriDate.Month,
		nomoztime.DailyTimee.Tong_saharlik,
		nomoztime.DailyTimee.Quyosh,
		nomoztime.DailyTimee.Peshin,
		nomoztime.DailyTimee.Asr,
		nomoztime.DailyTimee.Shom_iftor,
		nomoztime.DailyTimee.Hufton,
	)
	msg := tgbot.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)

	schedulePrayerTimes(bot, update.Message.Chat.ID, nomoztime)

	UserState.Store(update.Message.Chat.ID, "awaiting_choice")
}

func HandlePrayerNotifications(bot *tgbot.BotAPI, update tgbot.Update) {
	msg := tgbot.NewMessage(update.Message.Chat.ID, "You will receive notifications🕘 for prayer times. If you haven't used the /prayer_times command yet, please do so.")
	bot.Send(msg)
	UserState.Store(update.Message.Chat.ID, "awaiting_choice")
}

func schedulePrayerTimes(bot *tgbot.BotAPI, chatID int64, nomoztime model.Nomoztime) {
	prayerTimes := map[string]string{
		"🌌 Bomdod":     nomoztime.DailyTimee.Tong_saharlik,
		"🏞 Peshin":     nomoztime.DailyTimee.Peshin,
		"🌇 Asr:":       nomoztime.DailyTimee.Asr,
		"🏙 Shom Iftor": nomoztime.DailyTimee.Shom_iftor,
		"🌃 Xufton":     nomoztime.DailyTimee.Hufton,
	}

	now := time.Now()
	location := now.Location()

	for prayer, timeStr := range prayerTimes {
		prayerTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			fmt.Printf("Failed to parse prayer time %s: %v\n", timeStr, err)
			continue
		}

		// Use the local time zone
		prayerTimeToday := time.Date(
			now.Year(),          // Year
			now.Month(),         // Month
			now.Day(),           // Day
			prayerTime.Hour(),   // Hour
			prayerTime.Minute(), // Minute
			0,                   // Second
			0,                   // Nanosecond
			location,            // Location
		)

		if prayerTimeToday.After(now) {
			duration := time.Until(prayerTimeToday)

			time.AfterFunc(duration, func() {
				msg := tgbot.NewMessage(chatID, fmt.Sprintf("🕘 It's time for %s,%s", prayer, timeStr))
				bot.Send(msg)
			})
		}
	}
}

// Text fetches prayer times for a given region
func Text(region string) (model.Nomoztime, error) {
	if region == "" {
		region = "Toshkent"
	}

	url := fmt.Sprintf("https://islomapi.uz/api/present/day?region=%s", region)
	res, err := http.Get(url)
	if err != nil {
		return model.Nomoztime{}, fmt.Errorf("error fetching data from API: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Nomoztime{}, fmt.Errorf("error reading response body: %v", err)
	}

	var nomoztime model.Nomoztime
	if err := json.Unmarshal(body, &nomoztime); err != nil {
		return model.Nomoztime{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return nomoztime, nil
}

// CmdKeyboard creates and returns a reply keyboard
func CmdKeyboard() tgbot.ReplyKeyboardMarkup {
	return tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("/prayer_times"),
			tgbot.NewKeyboardButton("/prayer_notifications"),
		),
	)
}
