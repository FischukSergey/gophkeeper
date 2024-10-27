package command

// ICommand интерфейс для команд
type ICommand interface {
	Execute()
	Name() string
}	