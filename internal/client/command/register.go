package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"syscall"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
	"golang.org/x/term"
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
	_, err := fmt.Fprint(c.writer, "Введите логин: ")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	scanner := bufio.NewScanner(c.reader)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при вводе логина:", err)
		return
	}
	login := scanner.Text()
	//валидация логина
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprint(c.writer, "Введите пароль: ")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return
	}
	password := string(bytePassword)
	//валидация пароля
	err = modelsclient.ValidatePassword(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	//повторный ввод пароля
	_, err = fmt.Fprint(c.writer, "\nВведите пароль повторно: ")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	bytePassword, err = term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println("Ошибка при вводе пароля:", err)
		return
	}
	passwordConfirm := string(bytePassword)
	if password != passwordConfirm {
		fmt.Println("Пароли не совпадают")
		return
	}

	//вызываем регистрацию
	token, err := c.registerService.Register(context.Background(), strings.TrimSpace(login), strings.TrimSpace(password))
	if err != nil {
		//проверяем текст ошибки
		if strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "SQLSTATE 23505") {
			fmt.Println("\nПользователь с таким логином уже зарегистрирован")
			return
		}
		fmt.Println(err)
		return
	}
	c.token.Token = token
	fmt.Println("\nРегистрация прошла успешно")
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
