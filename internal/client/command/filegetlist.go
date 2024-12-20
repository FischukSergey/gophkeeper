package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/modelsclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const (
	nameCommandFileGetList = "FileList"
	kb                     = 1024
)

// CommandFileGetList структура для команды получения списка файлов.
type CommandFileGetList struct {
	fileGetListService IAuthService
	token              *grpcclient.Token
	reader             io.Reader
	writer             io.Writer
}

// NewCommandFileGetList создание новой команды получения списка файлов.
func NewCommandFileGetList(
	fileGetListService IAuthService,
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
	files, err := c.getFileList()
	if err != nil {
		c.handleError(err)
		waitEnter(c.reader)
		return
	}
	// вывод списка файлов
	c.displayFileList(files)

	// ожидание нажатия клавиши
	waitEnter(c.reader)
}

// getFileList получение списка файлов.
func (c *CommandFileGetList) getFileList() ([]models.File, error) {
	files, err := c.fileGetListService.GetFileList(context.Background(), c.token.Token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка файлов: %w", err)
	}
	if len(files) == 0 {
		return nil, errors.New("список файлов пуст")
	}
	return files, nil
}

// handleError обрабатывает ошибки.
func (c *CommandFileGetList) handleError(err error) {
	if strings.Contains(err.Error(), modelsclient.ErrTokenNotFound) {
		fmt.Println(errorAuth)
	} else {
		fmt.Println(err)
	}
}

// displayFileList выводит список файлов.
func (c *CommandFileGetList) displayFileList(files []models.File) {
	w := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', 0)

	// сортировка файлов по имени
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Filename) < strings.ToLower(files[j].Filename)
	})

	// вывод заголовков
	c.displayHeader(w)
	// вывод файлов
	for _, file := range files {
		c.displayFile(w, file)
	}

	// вывод данных
	err := w.Flush()
	if err != nil {
		fmt.Println(err)
	}
}

// displayHeader выводит заголовки таблицы.
func (c *CommandFileGetList) displayHeader(w *tabwriter.Writer) {
	headers := []string{
		"Имя файла\tРазмер\tДата создания",
		"----------\t------\t-------------",
	}
	for _, header := range headers {
		if _, err := fmt.Fprintln(w, header); err != nil {
			fmt.Printf(errOutputMessage, err)
		}
	}
}

// displayFile выводит информацию о файле.
func (c *CommandFileGetList) displayFile(w *tabwriter.Writer, file models.File) {
	_, err := fmt.Fprintf(w, "%s\t%d kb\t%s\n",
		file.Filename,
		file.Size/kb,
		file.CreatedAt.Format(time.DateTime),
	)
	if err != nil {
		fmt.Printf(errOutputMessage, err)
	}
}
