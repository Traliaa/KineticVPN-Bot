package telgram_bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleCommand обрабатывает команды бота
func (b *BotService) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()
	chatID := message.Chat.ID

	var responseText string

	switch command {
	case "start":
		responseText = "🚀 Добро пожаловать! Это тестовый бот.\n\n" +
			"Доступные команды:\n" +
			"/start - начать работу\n" +
			"/help - получить помощь\n" +
			"/settings - настройки бота"

	case "help":
		responseText = "📖 Помощь по боту:\n\n" +
			"Это демонстрационный бот с тремя командами:\n\n" +
			"• /start - приветственное сообщение\n" +
			"• /help - показывает эту справку\n" +
			"• /settings - настройки бота\n\n" +
			"Для использования просто отправьте одну из команд."

	case "settings":
		responseText = "⚙️ Настройки бота:\n\n" +
			"Здесь будут настройки вашего бота.\n" +
			"В будущем здесь можно будет:\n" +
			"• Изменить язык\n" +
			"• Настроить уведомления\n" +
			"• Выбрать тему оформления\n\n" +
			"Функционал настроек в разработке."

	default:
		responseText = "❌ Неизвестная команда. Используйте /help для списка команд."
	}

	msg := tgbotapi.NewMessage(chatID, responseText)
	msg.ReplyToMessageID = message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
