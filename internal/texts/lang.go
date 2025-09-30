package texts

import (
	"context"

	"github.com/Traliaa/KineticVPN-Bot/internal/texts/main_setting"
	"github.com/Traliaa/KineticVPN-Bot/internal/texts/router_setting"
)

const (
	notTranslate = "Translation for this message is not available. Please contact @traliaa on Telegram for assistance."
)

type Lang int8

const (
	RU Lang = iota
	EN
)

type TextID int8

const (
	WelcomeMain TextID = iota
	RouterLogin
	RouterPassword
	RouterAddress
)

var ruLang = map[TextID]string{
	WelcomeMain: main_setting.WelcomeMainRU,

	RouterLogin:    router_setting.RouterLoginRU,
	RouterPassword: router_setting.RouterPasswordRU,
	RouterAddress:  router_setting.RouterAddressRU,
}

var enLang = map[TextID]string{
	WelcomeMain: main_setting.WelcomeMainEN,

	RouterLogin:    router_setting.RouterLoginEN,
	RouterPassword: router_setting.RouterPasswordEN,
	RouterAddress:  router_setting.RouterAddressEN,
}

func GetTranslate(ctx context.Context, Lang Lang, TextID TextID) string {
	switch Lang {
	case RU:
		text, ok := ruLang[TextID]
		if ok {
			return text
		}
	case EN:
		text, ok := enLang[TextID]
		if ok {
			return text
		}
	default:
		return notTranslate
	}
	return notTranslate
}
