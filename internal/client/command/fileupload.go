package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/manifoldco/promptui"
)

const nameCommandFileUpload = "fileupload"

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
	//ввод пути к файлу
	filePathPrompt := promptui.Prompt{
		Label: "Введите путь к файлу: ",
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
		fmt.Println("Ошибка при загрузке файла:", err)
		return
	}
	fmt.Println("Файл загружен на S3:", s3Filepath)
	//ожидание нажатия клавиши
	fmt.Print("\nНажмите Enter для продолжения...")
	var input string
	fmt.Scanln(&input)
}
