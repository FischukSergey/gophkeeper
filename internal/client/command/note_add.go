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
	noteAddService INoteService
	token          *grpcclient.Token
	reader         io.Reader
	writer         io.Writer
}

// NewCommandNoteAdd конструктор команды добавления заметки.
func NewCommandNoteAdd(
	noteAddService INoteService,
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
	// проверка токена
	if !checkToken(c.token, c.reader) {
		return
	}
	// ввод текста заметки
	note, err := c.inputNoteText()
	if err != nil {
		return
	}
	// ввод метаданных
	noteMetadata, err := c.inputNoteMetadata()
	if err != nil {
		return
	}
	//вызов сервиса
	if err := c.addNote(note, noteMetadata); err != nil {
		fprintln(c.writer, err.Error())
		return
	}

	fprintln(c.writer, "\nЗаметка добавлена")
	waitEnter(c.reader)
}

// inputNoteText ввод текста заметки.
func (c *CommandNoteAdd) inputNoteText() (string, error) {
	fprintln(c.writer, "Добавление заметки")
	fprintln(c.writer, "Введите текст заметки (для завершения нажмите Enter и затем Ctrl+D):")
	note := readMultilineString(c.reader) //читаем многострочную строку
	if note == "" {
		fprintln(c.writer, "Текст заметки не может быть пустым")
		return "", errors.New("текст заметки не может быть пустым")
	}
	return note, nil
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

// inputNoteMetadata добавляет метаданные в map.
func (c *CommandNoteAdd) inputNoteMetadata() (map[string]string, error) {
	fprintln(c.writer, "\nВведите метаданные для заметки:")
	// создание map для метаданных
	noteMetadata := make(map[string]string)
	for {
		fprintln(c.writer, "\nТекущие введенные метаданные:")
		if len(noteMetadata) == 0 {
			fprintln(c.writer, "Пока нет введенных метаданных")
		} else {
			for k, v := range noteMetadata {
				fprintf(c.writer, "%s: %s\n", k, v)
			}
		}
		// Спрашиваем, хочет ли пользователь добавить еще пару
		continuePrompt := promptui.Prompt{
			Label:   "Добавить новую пару key-value? (y/n)",
			Default: "y",
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
			fprintln(c.writer, "Ошибка при вводе ключа:", err)
			return nil, err
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
			fprintln(c.writer, "Ошибка при вводе значения:", err)
			return nil, err
		}
		//проверяем, нет ли такого ключа в map
		if _, ok := noteMetadata[key]; ok {
			fprintln(c.writer, "Такой ключ уже существует")
			continue
		}
			noteMetadata[key] = value //добавление пары ключ-значение в map
	}
	return noteMetadata, nil
}

// addNote добавление заметки.
func (c *CommandNoteAdd) addNote(note string, noteMetadata map[string]string) error {
	err := c.noteAddService.NoteAdd(context.Background(), note, noteMetadata, c.token.Token)
	if err != nil {
		return err
	}
	return nil
}
