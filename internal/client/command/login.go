package command

import (
	"context"
	"fmt"
	"io"
	"strings"
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
		fmt.Println(err)
		return
	}
	//получение пароля
	password, err := c.getPassword()
	if err != nil {
		fmt.Println(err)
		return
	}
	//получение названия приложения
	applicationName, err := c.getApplicationName()
	if err != nil {
		fmt.Println(err)
		return
	}
	//авторизация
	token, err := c.authorization(login, password, applicationName)
	if err != nil {
		fmt.Println(err)
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
		return "", fmt.Errorf("%s: %w", errLoginMessage, err)
	}
	//валидация логина
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errLoginMessage, err)
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
		return "", fmt.Errorf("%s: %w", errPasswordMessage, err)
	}
	//валидация пароля
	err = modelsclient.ValidatePassword(password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errPasswordMessage, err)
	}
	return password, nil
}
// getApplicationName получение названия приложения.
func (c *CommandLogin) getApplicationName() (string, error) {
	applicationNamePrompt := promptui.Prompt{
		Label: "Введите название приложения (не обязательно)",
	}
	applicationName, err := applicationNamePrompt.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %w", errApplicationNameMessage, err)
	}
	applicationName = strings.ToLower(strings.TrimSpace(applicationName))
	return applicationName, nil
}

// authorization выполнение авторизации.
func (c *CommandLogin) authorization(login, password, applicationName string) (string, error) {
	token, err := c.authService.Authorization(context.Background(), login, password, applicationName)
	if err != nil {
		return "", fmt.Errorf("ошибка при авторизации: %w", err)
	}
	return token, nil
}
