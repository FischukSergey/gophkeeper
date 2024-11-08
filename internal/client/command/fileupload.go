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
	"github.com/manifoldco/promptui"
)

const nameCommandFileUpload = "FileUpload"

// IFileUploadService интерфейс для загрузки файла в хранилище.
type IFileUploadService interface {
	S3FileUpload(ctx context.Context, token string, fileData []byte, filename string) (string, error)
}

// CommandFileUpload структура для команды загрузки файла.
type CommandFileUpload struct {
	fileUploadService IFileUploadService
	token             *grpcclient.Token
	reader            io.Reader
	writer            io.Writer
}

// NewCommandFileUpload создает новый экземпляр команды загрузки файла.
func NewCommandFileUpload(
	fileUploadService IFileUploadService,
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
	//проверка наличия токена
	if c.token.Token == "" {
		fmt.Println("Вы не авторизованы. Авторизуйтесь с помощью команды login.")
		// ожидание нажатия клавиши
		fmt.Println(messageContinue)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf(errInputMessage, err)
		}
		return
	}
	//ввод пути к файлу
	filePathPrompt := promptui.Prompt{
		Label: "Введите путь к файлу",
	}
	filePath, err := filePathPrompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе пути к файлу:", err)
		return
	}
	//проверка, что файл существует
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Файл не найден:", err)
		return
	}
	//чтение файла
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	//получение названия файла
	filename := filepath.Base(filePath)
	//загрузка файла на сервер
	s3Filepath, err := c.fileUploadService.S3FileUpload(context.Background(), c.token.Token, fileData, filename)
	if err != nil {
		// проверка ошибки
		if strings.Contains(err.Error(), auth.ErrNotFound) ||
			strings.Contains(err.Error(), auth.ErrInvalid) ||
			strings.Contains(err.Error(), "user ID not found in context") {
			fmt.Println("Ошибка авторизации. Пожалуйста, войдите в систему заново")
		} else {
			fmt.Printf(errOutputMessage, err)
		}
		return
	}
	fmt.Println("Файл загружен на S3:", s3Filepath)
	//ожидание нажатия клавиши
	fmt.Println(messageContinue)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Printf(errInputMessage, err)
	}
}
