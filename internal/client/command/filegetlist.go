package command

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const nameCommandFileGetList = "FileList"

// IFileGetListService интерфейс для получения списка файлов.
type IFileGetListService interface {
	GetFileList(ctx context.Context, token string) ([]models.File, error)
}

// CommandFileGetList структура для команды получения списка файлов.
type CommandFileGetList struct {
	fileGetListService IFileGetListService
	token              *grpcclient.Token
	reader             io.Reader
	writer             io.Writer
}

// NewCommandFileGetList создание новой команды получения списка файлов.
func NewCommandFileGetList(
	fileGetListService IFileGetListService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandFileGetList {
	return &CommandFileGetList{
		fileGetListService: fileGetListService,
		token:              token,
		reader:             reader,
		writer:             writer,
	}
}

// Name возвращает имя команды.
func (c *CommandFileGetList) Name() string {
	return nameCommandFileGetList
}

// Execute выполнение команды получения списка файлов.
func (c *CommandFileGetList) Execute() {
	//проверка наличия токена
	if c.token.Token == "" {
		fmt.Println("Вы не авторизованы. Авторизуйтесь с помощью команды login.")
		// ожидание нажатия клавиши
		fmt.Println("\nНажмите Enter для продолжения...")
		var input string
		fmt.Scanln(&input)
		return
	}
	// получение списка файлов
	files, err := c.fileGetListService.GetFileList(context.Background(), c.token.Token)
	if err != nil {
		fmt.Println("Ошибка при получении списка файлов:", err)
		return
	}

	// создаем новый tabwriter
	w := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', 0)

	// выводим заголовки таблицы
	fmt.Fprintln(w, "Имя файла\tРазмер\tДата создания\t")
	fmt.Fprintln(w, "----------\t------\t-------------\t")

	// выводим данные
	for _, file := range files {
		fmt.Fprintf(w, "%s\t%d kb\t%s\t\n",
			file.Filename,
			file.Size/1024,
			file.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	// применяем форматирование таблицы
	w.Flush()

	// ожидание нажатия клавиши
	fmt.Print("\nНажмите Enter для продолжения...")
	var input string
	fmt.Scanln(&input)
}
