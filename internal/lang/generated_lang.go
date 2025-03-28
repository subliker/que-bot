/** Code generated using https://github.com/MrNemo64/go-n-i18n 
 * Any changes to this file will be lost on the next tool run */

package lang

import (
    "fmt"
    "strings"
)

func MessagesFor(tag string) (Messages, bool) {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "ru-RU":
        return ru_RU_Messages{}, true
    }
    return nil, false
}

func MessagesForMust(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "ru-RU":
        return ru_RU_Messages{}
    }
    panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "ru-RU":
        return ru_RU_Messages{}
    }
    return ru_RU_Messages{}
}

type Messages interface{
    StartMessage() startMessage
    Query() query
    Callback() callback
}
type startMessage interface{
    Head(user_name string) string
    Main(bot_user_name string) string
}
type query interface{
    Main() querymain
    Btns() querybtns
}
type querymain interface{
    Title(queue_name string) string
    Description() string
    Text(queue_name string) string
    TextNoGroup() string
}
type querybtns interface{
    New() string
}
type callback interface{
    QueueNew() callbackqueueNew
    Queue() callbackqueue
    Btns() callbackbtns
}
type callbackqueueNew interface{
    Main(queue_name string) string
}
type callbackqueue interface{
    Head(queue_name string) string
    Member(num int, first_name string, last_name string, user_name string) string
}
type callbackbtns interface{
    SubmitFirst() string
    Submit(num int) string
}

type ru_RU_Messages struct{}
func (ru_RU_Messages) StartMessage() startMessage {
    return ru_RU_startMessage{}
}
type ru_RU_startMessage struct{}
func (ru_RU_startMessage) Head(user_name string) string {
    if user_name == "" {
        return "✨ Добро пожаловать! ✨"
    } else {
        return fmt.Sprintf("✨ Добро пожаловать, %s! ✨", user_name)
    }
}
func (ru_RU_startMessage) Main(bot_user_name string) string {
    return "<i>Я упрощаю процесс занятия очереди при помощи одной кнопки!</i> 🔲" + "\n" +
        "" + "\n" +
        "<i>Знакомо, что весь чат в разнобой занимает места?</i> 🥲" + "\n" +
        "<b>Подобный процесс можно автоматизировать при помощи меня</b> 🤩" + "\n" +
        "" + "\n" +
        "Бот работает в режиме inline. Это означает, что вы можете взаимодействовать с ним из текстового поля. 💬" + "\n" +
        fmt.Sprintf(" 1. Введите в строке @%s", bot_user_name) + "\n" +
        " 2. Нажмите появившуюся кнопку для создания очереди или введите через пробел опциональные параметры(присутствует подсказка-placeholder)" + "\n" +
        "" + "\n" +
        "Добавьте меня в групповой чат, увидимся там!"
}
func (ru_RU_Messages) Query() query {
    return ru_RU_query{}
}
type ru_RU_query struct{}
func (ru_RU_query) Main() querymain {
    return ru_RU_querymain{}
}
type ru_RU_querymain struct{}
func (ru_RU_querymain) Title(queue_name string) string {
    if queue_name != "" {
        return fmt.Sprintf("Создать очередь %s 🔨", queue_name)
    } else {
        return "Создать очередь 🔨"
    }
}
func (ru_RU_querymain) Description() string {
    return "Стандартная очередь по нажатию кнопки"
}
func (ru_RU_querymain) Text(queue_name string) string {
    if queue_name != "" {
        return fmt.Sprintf("Чтобы создать стандартную очередь *%s*, нажми на кнопку ниже 🚀", queue_name)
    } else {
        return "Чтобы создать стандартную очередь, нажми на кнопку ниже 🚀"
    }
}
func (ru_RU_querymain) TextNoGroup() string {
    return "Очередь можно создать только в группе! 🫥"
}
func (ru_RU_query) Btns() querybtns {
    return ru_RU_querybtns{}
}
type ru_RU_querybtns struct{}
func (ru_RU_querybtns) New() string {
    return "🔨 Создать очередь 🔨"
}
func (ru_RU_Messages) Callback() callback {
    return ru_RU_callback{}
}
type ru_RU_callback struct{}
func (ru_RU_callback) QueueNew() callbackqueueNew {
    return ru_RU_callbackqueueNew{}
}
type ru_RU_callbackqueueNew struct{}
func (ru_RU_callbackqueueNew) Main(queue_name string) string {
    return fmt.Sprintf("Очередь на *%s* на данный момент пуста! Стань первым 😎", queue_name)
}
func (ru_RU_callback) Queue() callbackqueue {
    return ru_RU_callbackqueue{}
}
type ru_RU_callbackqueue struct{}
func (ru_RU_callbackqueue) Head(queue_name string) string {
    return fmt.Sprintf("На данный момент очередь на *%s* выглядит следующим образом 📒:", queue_name)
}
func (ru_RU_callbackqueue) Member(num int, first_name string, last_name string, user_name string) string {
    return fmt.Sprintf("%d. [%s %s](https://t.me/%s)", num, first_name, last_name, user_name)
}
func (ru_RU_callback) Btns() callbackbtns {
    return ru_RU_callbackbtns{}
}
type ru_RU_callbackbtns struct{}
func (ru_RU_callbackbtns) SubmitFirst() string {
    return "Я первый! 🥇"
}
func (ru_RU_callbackbtns) Submit(num int) string {
    if num >= 30 {
        return "Занять место 🫡"
    } else if num >= 25 {
        return "Занять место 🫠"
    } else if num >= 20 {
        return "Занять место 😱"
    } else if num >= 15 {
        return "Занять место 😨"
    } else if num >= 10 {
        return "Занять место 🥲"
    } else if num >= 5 {
        return "Занять место 🤠"
    } else if num < 5 {
        return "Занять место 😁"
    } else {
        return fmt.Sprintf("%d", num)
    }
}


