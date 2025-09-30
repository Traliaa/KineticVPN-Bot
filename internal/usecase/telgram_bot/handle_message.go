package telgram_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//// HandleMessage обрабатывает обычные текстовые сообщения
//func (b *BotService) HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	chatID := message.Chat.ID
//	text := strings.ToLower(message.Text)
//
//	var responseText string
//
//	// Простой пример обработки текстовых сообщений
//	switch {
//	case strings.Contains(text, "привет") || strings.Contains(text, "hello"):
//		responseText = "Привет! 👋 Используйте команды для взаимодействия с ботом."
//	case strings.Contains(text, "спасибо") || strings.Contains(text, "thanks"):
//		responseText = "Пожалуйста! 😊 Если нужна помощь - используйте /help"
//	default:
//		responseText = "Я понимаю только команды. Используйте /help для списка доступных команд."
//	}
//
//	msg := tgbotapi.NewMessage(chatID, responseText)
//	if _, err := bot.Send(msg); err != nil {
//		log.Printf("Error sending message: %v", err)
//	}
//}

// handleMessage обрабатывает текстовые сообщения и команды
func (b *BotService) HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Инициализируем сессию пользователя если её нет
	if _, exists := UserSession[chatID]; !exists {
		UserSession[chatID] = &UserData{Step: "start"}
	}

	user := UserSession[chatID]

	// Обрабатываем команды
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendWelcomeMessage(bot, chatID, user)
		case "help":
			sendHelpMessage(bot, chatID)
		case "reset":
			delete(UserSession, chatID)
			UserSession[chatID] = &UserData{Step: "start"}
			sendWelcomeMessage(bot, chatID, UserSession[chatID])
		}
		return
	}

	// Обрабатываем шаги настройки
	switch user.Step {
	case "ask_setup":
		b.handleSetupResponse(bot, message, user)
	case "router_url":
		handleRouterURL(bot, message, user)
	case "username":
		handleUsername(bot, message, user)
	case "password":
		handlePassword(bot, message, user)
	case "auth_code":
		handleAuthCode(bot, message, user)
	default:
		sendMainMenu(bot, chatID, "Используйте кнопки для навигации:")
	}
}

// UserData хранит данные пользователя
type UserData struct {
	RouterURL    string   `json:"router_url"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	AuthCode     string   `json:"auth_code"`
	SelectedApps []string `json:"selected_apps"`
	Step         string   `json:"step"`
}

// UserSession хранит сессии пользователей
var UserSession = make(map[int64]*UserData)
