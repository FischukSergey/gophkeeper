package command

import (
	"context"
	"fmt"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
	"github.com/manifoldco/promptui"
)

const nameCommandLogin = "Login"

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
	//получение логина
	login, err := c.getLogin()
	if err != nil {
		fmt.Println("Ошибка при вводе логина:", err)
		return
	}
	//получение пароля
	password, err := c.getPassword()
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return
	}
	//авторизация
	token, err := c.authorization(login, password)
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}
	c.token.Token = token
	fmt.Println("\nАвторизация прошла успешно")

	// Ожидание нажатия клавиши
	waitEnter(c.reader)
}

// getLogin получение логина.
func (c *CommandLogin) getLogin() (string, error) {
	loginPrompt := promptui.Prompt{
		Label: "Введите логин",
	}
	login, err := loginPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе логина:", err)
		return "", err
	}
	//валидация логина
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return login, nil
}

// getPassword получение пароля.
func (c *CommandLogin) getPassword() (string, error) {
	passwordPrompt := promptui.Prompt{
		Label: "Введите пароль",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return "", err
	}
	//валидация пароля
	err = modelsclient.ValidatePassword(password)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return password, nil
}

// authorization выполнение авторизации.
func (c *CommandLogin) authorization(login, password string) (string, error) {
	token, err := c.authService.Authorization(context.Background(), login, password)
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return "", err
	}
	return token, nil
}
