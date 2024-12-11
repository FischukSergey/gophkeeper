package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
	"github.com/manifoldco/promptui"
)

const nameCommandRegister = "Registration"

// commandRegister структура для команды регистрации.
type commandRegister struct {
	registerService IAuthService
	token           *grpcclient.Token
	reader          io.Reader
	writer          io.Writer
}

// Name возвращает имя команды.
func (c *commandRegister) Name() string {
	return nameCommandRegister
}

// Execute выполняет команду регистрации.
func (c *commandRegister) Execute() {
	//получение логина
	login, err := c.promptLogin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//получение пароля
	password, err := c.promptPassword()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//вызываем регистрацию
	token, err := c.registerUser(login, password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.token.Token = token
	fmt.Println("\nРегистрация прошла успешно")
	waitEnter(c.reader)
}

// NewCommandRegister создание команды регистрации.
func NewCommandRegister(
	registerService IAuthService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *commandRegister {
	return &commandRegister{
		registerService: registerService,
		token:           token,
		reader:          reader,
		writer:          writer,
	}
}

// promptLogin ввод логина.
func (c *commandRegister) promptLogin() (string, error) {
	loginPrompt := promptui.Prompt{
		Label: "Введите логин",
	}
	login, err := loginPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %w", errLoginMessage, err)
	}
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", errLoginMessage, err)
	}
	return login, nil
}

// promptPassword ввод пароля.
func (c *commandRegister) promptPassword() (string, error) {
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
	//повторный ввод пароля
	passwordConfirmPrompt := promptui.Prompt{
		Label: "Введите пароль повторно",
		Mask:  '*',
	}
	passwordConfirm, err := passwordConfirmPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %w", errPasswordMessage, err)
	}
	if password != passwordConfirm {
		return "", errors.New("пароли не совпадают")
	}
	return password, nil
}

// registerUser регистрация пользователя.
func (c *commandRegister) registerUser(login string, password string) (string, error) {
	token, err := c.registerService.Register(context.Background(), login, password)
	if err != nil {
		return "", fmt.Errorf("ошибка при регистрации: %w", err)
	}
	return token, nil
}
