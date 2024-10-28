package command

const (
	errOutputMessage = "Ошибка вывода сообщения: %s\n"
	errReadMessage   = "Ошибка чтения ответа: %s\n"
)

// ICommand интерфейс для команд.
type ICommand interface {
	Execute()
	Name() string
}
