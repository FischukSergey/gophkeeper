package command

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/manifoldco/promptui"
)
const nameCommandNoteAdd = "NoteAdd"

// CommandNoteAdd структура команды добавления заметки.	
type CommandNoteAdd struct {
	noteAddService INoteAddService
	token          *grpcclient.Token
	reader         io.Reader
	writer         io.Writer
}

// INoteAddService интерфейс для сервиса добавления заметки.
type INoteAddService interface {
	NoteAdd(ctx context.Context, note string, metadata map[string]string, token string) error
}

// NewCommandNoteAdd конструктор команды добавления заметки.
func NewCommandNoteAdd(
	noteAddService INoteAddService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandNoteAdd {
	return &CommandNoteAdd{
		noteAddService: noteAddService,
		token:          token,
		reader:         reader,
		writer:         writer,
	}
}

// Name возвращает название команды.
func (c *CommandNoteAdd) Name() string {
	return nameCommandNoteAdd
}

// Execute выполнение команды.
func (c *CommandNoteAdd) Execute() {
	if !checkToken(c.token, c.reader) {
		return
	}
	fprintln(c.writer, "Добавление заметки")
	fprintln(c.writer, "Введите текст заметки (для завершения нажмите Enter и затем Ctrl+D):")
	note := readMultilineString(c.reader) //читаем многострочную строку
	if note == "" {
		fprintln(c.writer, "Текст заметки не может быть пустым")
		return
	}

	fprintln(c.writer, "\nВведите метаданные для заметки:")
	noteMetadata := make(map[string]string)
	addMetadata(c.writer, &noteMetadata) //добавление метаданных

	//вызов сервиса
	if err := c.noteAddService.NoteAdd(context.Background(), note, noteMetadata, c.token.Token); err != nil {
		fprintln(c.writer, err.Error())
		return
	}

	fprintln(c.writer, "\nЗаметка добавлена")
	waitEnter(c.reader)
}

// readMultilineString читает многострочную строку из io.Reader.
func readMultilineString(reader io.Reader) string {
	scanner := bufio.NewScanner(reader)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n")
}

// addMetadata добавляет метаданные в map.
func addMetadata(writer io.Writer, metadata *map[string]string) {
	for {
		fprintln(writer, "\nТекущие введенные метаданные:")
		if len(*metadata) == 0 {
			fprintln(writer, "Пока нет введенных метаданных")
		} else {
			for k, v := range *metadata {
				fprintf(writer, "%s: %s\n", k, v)
			}
		}
		// Спрашиваем, хочет ли пользователь добавить еще пару
		continuePrompt := promptui.Prompt{
			Label:     "Добавить новую пару key-value? (y/n)",
			Default:   "y",
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return errors.New("пожалуйста, введите 'y' или 'n'")
				}
				return nil
			},
		}

		shouldContinue, err := continuePrompt.Run()
		if err != nil || shouldContinue == "n" {
			break //выход из цикла, если пользователь не хочет добавлять больше метаданных
		}

		// Ввод ключа
		keyPrompt := promptui.Prompt{
			Label: "Введите ключ",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("ключ не может быть пустым")
				}
				return nil
			},
		}

		key, err := keyPrompt.Run()
		if err != nil {
			fprintln(writer, "Ошибка при вводе ключа:", err)
			return
		}

		// Ввод значения
		valuePrompt := promptui.Prompt{
			Label: "Введите значение",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("значение не может быть пустым")
				}
				return nil
			},
		}

		value, err := valuePrompt.Run()
		if err != nil {
			fprintln(writer, "Ошибка при вводе значения:", err)
			return
		}
		//проверяем, нет ли такого ключа в map
		if _, ok := (*metadata)[key]; ok {
			fprintln(writer, "Такой ключ уже существует")
			continue
		}
		(*metadata)[key] = value //добавление пары ключ-значение в map
	}
}
