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

const nameCommandRegister = "register"

// IRegisterService интерфейс для сервиса регистрации.
type IRegisterService interface {
	Register(ctx context.Context, login string, password string) (string, error)
}

// commandRegister структура для команды регистрации.
type commandRegister struct {
	registerService IRegisterService
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
	//повторный ввод пароля
	passwordConfirmPrompt := promptui.Prompt{
		Label: "Введите пароль повторно: ",
		Mask:  '*',
	}
	passwordConfirm, err := passwordConfirmPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return
	}
	if password != passwordConfirm {
		fmt.Println("Пароли не совпадают")
		return
	}
	//вызываем регистрацию	
	token, err := c.registerService.Register(context.Background(), strings.TrimSpace(login), strings.TrimSpace(password))
	if err != nil {
		fmt.Println("Ошибка при регистрации:", err)
		return
	}
	c.token.Token = token
	fmt.Println("\nРегистрация прошла успешно")
	// Ожидание нажатия клавиши
	fmt.Print("\nНажмите Enter для продолжения...")
	var input string
	fmt.Scanln(&input)
}

// NewCommandRegister создание команды регистрации.
func NewCommandRegister(
	registerService IRegisterService,
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
