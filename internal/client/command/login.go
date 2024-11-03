package command

import (
	"context"
	"fmt"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
	"github.com/manifoldco/promptui"
)

const nameCommandLogin = "login"

// IAuthService интерфейс для сервиса авторизации.
type IAuthService interface {
	Authorization(ctx context.Context, login, password string) (string, error)
}

// CommandLogin структура для команды авторизации.
type CommandLogin struct {
	authService IAuthService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

func NewCommandLogin(
	authService IAuthService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandLogin {
	return &CommandLogin{
		authService: authService,
		token:       token,
		reader:      reader,
		writer:      writer,
	}
}

// Name возвращает имя команды.
func (c *CommandLogin) Name() string {
	return nameCommandLogin
}

// Execute выполнение команды авторизации.
func (c *CommandLogin) Execute() {
	//ввод логина
	loginPrompt := promptui.Prompt{
		Label: "Введите логин: ",
	}
	login, err := loginPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе логина:", err)
		return
	}
	//валидация логина
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		fmt.Println(err)
		return
	}
	//ввод пароля
	passwordPrompt := promptui.Prompt{
		Label: "Введите пароль: ",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return
	}
	//валидация пароля
	err = modelsclient.ValidatePassword(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	//авторизация
	token, err := c.authService.Authorization(context.Background(), login, password)
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}
	c.token.Token = token
	fmt.Println("\nАвторизация прошла успешно")
	// Ожидание нажатия клавиши
	fmt.Print("\nНажмите Enter для продолжения...")
	var input string
	fmt.Scanln(&input)
}
