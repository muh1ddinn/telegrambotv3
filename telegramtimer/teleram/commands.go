package teleram

import tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Commands(bot *tgbot.BotAPI, update tgbot.Update) {
	switch update.Message.Command() {
	case "prayer_times":
		HandlePrayerTimes(bot, update)
		UserState.Store(update.Message.Chat.ID, "selecting_region")
	case "prayer_notifications":
		HandlePrayerNotifications(bot, update)
		UserState.Store(update.Message.Chat.ID, "awaiting_choice")
	case "start":
		HandleStart(bot, update)
	default:
		HandleUserResponse(bot, update)
	}
}
