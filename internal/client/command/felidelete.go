package command

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const nameCommandFileDelete = "FileDelete"

// CommandFileDelete структура для команды удаления файла.
type CommandFileDelete struct {
	fileDeleteService IFileDeleteService
	token             *grpcclient.Token
	reader            io.Reader
	writer            io.Writer
}

// IFileDeleteService интерфейс для удаления файла.
type IFileDeleteService interface {
	S3FileDelete(ctx context.Context, token string, filename string) error
}

// NewCommandFileDelete функция для создания команды удаления файла.
func NewCommandFileDelete(
	fileDeleteService IFileDeleteService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandFileDelete {
	return &CommandFileDelete{
		fileDeleteService: fileDeleteService,
		token:             token,
		reader:            reader,
		writer:            writer,
	}
}

// Name возвращает имя команды удаления файла.
func (c *CommandFileDelete) Name() string {
	return nameCommandFileDelete
}

// Execute выполнение команды удаления файла.
func (c *CommandFileDelete) Execute() {
	//проверка наличия токена
	checkToken(c.token, c.reader)
	// запрос на удаление файла
	fmt.Println("Введите имя файла для удаления (внимание, чувствительно к регистру):")
	var filename string
	_, err := fmt.Scanln(&filename)
	if err != nil {
		fmt.Printf(errInputMessage, err)
	}
	// удаление файла
	err = c.fileDeleteService.S3FileDelete(context.Background(), c.token.Token, filename)
	// проверка ошибки
	if err != nil {
		switch {
		case strings.Contains(err.Error(), models.ErrFileNotExist.Error()):
			fmt.Println("Файл не найден")
		case strings.Contains(err.Error(), auth.ErrNotFound) ||
			strings.Contains(err.Error(), auth.ErrInvalid) ||
			strings.Contains(err.Error(), auth.ErrExpired):
			fmt.Println("Ошибка авторизации. Пожалуйста, войдите в систему заново")
		default:
			fmt.Printf("Ошибка при удалении файла: %v\n", err)
		}
		// ожидание нажатия клавиши
		waitEnter(c.reader)
		return
	}

	fmt.Println("Файл успешно удален")
	waitEnter(c.reader)
}
