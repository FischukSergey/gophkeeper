package command

const (
	errOutputMessage = "Ошибка вывода сообщения: %s\n"
	errReadMessage   = "Ошибка чтения ответа: %s\n"
	errInputMessage  = "Ошибка ввода: %s\n"
	messageContinue  = "\nНажмите Enter для продолжения..."
)

// ICommand интерфейс для команд.
type ICommand interface {
	Execute()
	Name() string
}
