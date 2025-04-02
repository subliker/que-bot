package telebot

import "regexp"

var queryTextRegexp = regexp.MustCompile("[^A-Za-zА-Яа-яЁё0-9 ]+")
