package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/manifoldco/promptui"
)

const nameCommandFileUpload = "FileUpload"

// CommandFileUpload структура для команды загрузки файла.
type CommandFileUpload struct {
	fileUploadService IAuthService
	token             *grpcclient.Token
	reader            io.Reader
	writer            io.Writer
}

// NewCommandFileUpload создает новый экземпляр команды загрузки файла.
func NewCommandFileUpload(
	fileUploadService IAuthService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandFileUpload {
	return &CommandFileUpload{
		fileUploadService: fileUploadService,
		token:             token,
		reader:            reader,
		writer:            writer,
	}
}

// Name возвращает имя команды.
func (c *CommandFileUpload) Name() string {
	return nameCommandFileUpload
}

// Execute выполняет команду загрузки файла.
func (c *CommandFileUpload) Execute() {
	// Проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return // Выходим из функции если токен невалидный
	}
	//получение пути к файлу
	filePath, err := c.getFilePath()
	if err != nil {
		return
	}
	//чтение файла
	fileData, filename, err := c.getFileData(filePath)
	if err != nil {
		return
	}
	//загрузка файла на сервер
	if err := c.s3FileUpload(fileData, filename); err != nil {
		return
	}

	//ожидание нажатия клавиши
	waitEnter(c.reader)
}

// getFilePath получение пути к файлу.
func (c *CommandFileUpload) getFilePath() (string, error) {
	const maxAttempts = 2
	var filePath string
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		filePathPrompt := promptui.Prompt{
			Label: "Введите путь к файлу",
		}
		filePath, err = filePathPrompt.Run()
		if err == nil {
			// проверка, что файл существует
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				return filePath, nil
			}
		}

		if attempt < maxAttempts {
			fmt.Println("Ошибка при вводе пути к файлу. Попробуйте еще раз.")
			continue
		}
	}

	fmt.Println("Файл не найден. Будьте внимательны при вводе пути к файлу.")
	return "", err
}

// getFileData получение данных файла.
func (c *CommandFileUpload) getFileData(filePath string) ([]byte, string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return nil, "", err
	}
	//получение названия файла
	filename := filepath.Base(filePath)
	return fileData, filename, nil
}

// s3FileUpload загрузка файла на S3.
func (c *CommandFileUpload) s3FileUpload(fileData []byte, filename string) error {
	s3Filepath, err := c.fileUploadService.S3FileUpload(context.Background(), c.token.Token, fileData, filename)
	if err != nil {
		// проверка ошибки
		if strings.Contains(err.Error(), auth.ErrNotFound) ||
			strings.Contains(err.Error(), auth.ErrInvalid) ||
			strings.Contains(err.Error(), models.UserIDNotFound) {
			fmt.Println(errorAuth)
		} else {
			fmt.Printf(errOutputMessage, err)
		}
			return err
	}
	fmt.Println("Файл загружен на S3:", s3Filepath)
	return nil
}
