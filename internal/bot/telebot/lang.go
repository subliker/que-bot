package telebot

import (
	"regexp"

	"github.com/subliker/que-bot/internal/lang"
)

func (c *controller) telegramLangCodeToI18N(tlc string) string {
	switch tlc {
	case "ru":
		return "ru-RU"
	}
	return "ru-RU"
}

func (c *controller) langBundle(langCode string) lang.Messages {
	return lang.MessagesForOrDefault(c.telegramLangCodeToI18N(langCode))
}

var inputText = regexp.MustCompile("[^A-Za-zА-Яа-яЁё0-9 ]+")
