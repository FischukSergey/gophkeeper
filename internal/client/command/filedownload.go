package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const nameCommandFileDownload = "FileDownload"

type IFileDownloadService interface {
	S3FileDownload(ctx context.Context, token string, filename string) ([]byte, error)
}

type CommandFileDownload struct {
	fileDownloadService IFileDownloadService
	token               *grpcclient.Token
	reader              io.Reader
	writer              io.Writer
}

func NewCommandFileDownload(
	fileDownloadService IFileDownloadService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandFileDownload {
	return &CommandFileDownload{
		fileDownloadService: fileDownloadService,
		token:               token,
		reader:              reader,
		writer:              writer,
	}
}

// Name возвращает имя команды.
func (c *CommandFileDownload) Name() string {
	return nameCommandFileDownload
}

// Execute выполняет команду загрузки файла.
func (c *CommandFileDownload) Execute() {
	//проверка наличия токена
	checkToken(c.token, c.reader)
	// запрос на загрузку файла
	fmt.Println("Введите имя файла из хранилища (внимание, чувствительно к регистру):")
	var filename string
	_, err := fmt.Scanln(&filename)
	if err != nil {
		fmt.Printf(errInputMessage, err)
		return
	}
	// введите путь для сохранения файла
	fmt.Println("Введите путь для сохранения файла:")
	var filepath string
	_, err = fmt.Scanln(&filepath)
	if err != nil {
		filepath = "."
	}
	// валидация пути
	err = validatePath(filepath)
	if err != nil {
		fmt.Printf(errInputMessage, err)
		return
	}
	// проверка существования файла
	err = checkFileExists(filepath, filename)
	if err == nil {
		fmt.Printf("Файл с таким именем уже существует\n")
		fmt.Println("Хотите перезаписать файл? (y/n)")
		var answer string
		_, err = fmt.Scanln(&answer)
		if err != nil {
			fmt.Printf(errInputMessage, err)
			return
		}
		if answer != "y" {
			return
		}
	}
	//TODO: проверить, что файл не занят другим процессом и можно ли его перезаписать

	// загрузка файла
	bytes, err := c.fileDownloadService.S3FileDownload(context.Background(), c.token.Token, filename)
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
			fmt.Printf("Ошибка при загрузке файла: %v\n", err)
		}
		// ожидание нажатия клавиши
		waitEnter(c.reader)
		return
	}
	// сохранение файла
	filename = filepath + "/" + filename
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	fmt.Println("Файл загружен на диск:", filename)
	// ожидание нажатия клавиши
	waitEnter(c.reader)
}
