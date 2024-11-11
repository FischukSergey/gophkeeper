package command

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

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
	// Проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return // Выходим из функции если токен невалидный
	}
	// получение списка файлов
	files, err := c.fileGetListService.GetFileList(context.Background(), c.token.Token)
	if err != nil {
		// проверка ошибки
		if strings.Contains(err.Error(), "токен не найден") {
			fmt.Println(errorAuth)
		} else {
			fmt.Println("Ошибка при получении списка файлов:", err)
		}
		return
	}

	// создаем новый tabwriter
	w := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', 0)

	// выводим заголовки таблицы
	_, err = fmt.Fprintln(w, "Имя файла\tРазмер\tДата создания")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
	}
	_, err = fmt.Fprintln(w, "----------\t------\t-------------")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
	}

	// выводим данные
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Filename) < strings.ToLower(files[j].Filename)
	})
	for _, file := range files {
		_, err = fmt.Fprintf(w, "%s\t%d kb\t%s\n",
			file.Filename,
			file.Size/1024,
			file.CreatedAt.Format(time.DateTime),
		)
		if err != nil {
			fmt.Printf(errOutputMessage, err)
		}
	}

	// применяем форматирование таблицы
	err = w.Flush()
	if err != nil {
		fmt.Printf(errOutputMessage, err)
	}

	// ожидание нажатия клавиши
	waitEnter(c.reader)
}
