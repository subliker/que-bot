package telebot

import tele "gopkg.in/telebot.v4"

type btn struct {
	tele.Btn
}

func newBtn(unique string) btn {
	return btn{
		Btn: tele.Btn{
			Unique: unique,
		},
	}
}

func (b *btn) tele() tele.Btn {
	return b.Btn
}
