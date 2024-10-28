package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"syscall"

	"golang.org/x/term"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
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
	// валидация логина
	err = modelsclient.ValidateLogin(login)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ввод пароля без отображения в терминале
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

	// авторизация
	token, err := c.authService.Authorization(context.Background(), login, password)
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}
	c.token.Token = token
	fmt.Println("\nАвторизация успешна")
}
