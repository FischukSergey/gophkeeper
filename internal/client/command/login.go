package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"golang.org/x/term"
	"syscall"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
)

const nameCommandLogin = "login"

// IAuthService интерфейс для сервиса авторизации
type IAuthService interface {
	Authorization(ctx context.Context, login, password string) (string, error)
}

// CommandLogin структура для команды авторизации
type CommandLogin struct {
	authService IAuthService
	token *grpcclient.Token
	reader io.Reader
	writer io.Writer	
}
	
func NewCommandLogin(
	authService IAuthService, 
	token *grpcclient.Token, 
	reader io.Reader, 
	writer io.Writer,
) *CommandLogin {
	return &CommandLogin{
		authService: authService, 
		token: token, 
		reader: reader, 
		writer: writer,
	}
}	

func (c *CommandLogin) Name() string {
	return nameCommandLogin
}

func (c *CommandLogin) Execute() {
	//ввод логина
	fmt.Fprint(c.writer, "Введите логин: ")	
	scanner := bufio.NewScanner(c.reader)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при вводе логина:", err)
		return
	}	
	login := scanner.Text()
	//валидация логина
	err := modelsclient.ValidateLogin(login)
	if err != nil {
		fmt.Println(err)
		return
	}

	//ввод пароля без отображения в терминале
	fmt.Fprint(c.writer, "Введите пароль: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
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

	//авторизация
	token, err := c.authService.Authorization(context.Background(), login, password)
	if err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}
	c.token.Token = token
	fmt.Println("Авторизация успешна")
}	
