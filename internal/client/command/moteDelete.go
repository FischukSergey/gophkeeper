package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/manifoldco/promptui"
)

const commandNoteDelete = "Delete"

// NoteDeleteCommand структура для удаления заметки.
type NoteDeleteCommand struct {
	noteDeleteService INoteDeleteService
	token             *grpcclient.Token
	reader            io.Reader
	writer            io.Writer
}

// INoteDeleteService интерфейс для сервиса удаления заметки.
type INoteDeleteService interface {
	NoteDeleteService(ctx context.Context, noteID int64, token string) error
}

// NewCommandNoteDelete создание новой команды удаления заметки.
func NewCommandNoteDelete(
	noteDeleteService INoteDeleteService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *NoteDeleteCommand {
	return &NoteDeleteCommand{
		noteDeleteService: noteDeleteService,
		token:             token,
		reader:            reader,
		writer:            writer,
	}
}

// Name возвращает имя команды.
func (c *NoteDeleteCommand) Name() string {
	return commandNoteDelete
}

// Execute выполнение команды удаления заметки.
func (c *NoteDeleteCommand) Execute() {
	// проверка токена
	if !checkToken(c.token, c.reader) {
		return
	}
	// ввод номера заметки
	fprint(c.writer, "Введите номер заметки для удаления: "+"\033[35m")
	var noteID int64
	_, err := fmt.Fscanln(c.reader, &noteID)
	if err != nil {
		fprintln(c.writer, "\033[0m"+"Неверный номер заметки")
		waitEnter(c.reader)
		return
	}
	// подтверждение удаления заметки
	continuePrompt := promptui.Prompt{
		Label:   "Вы уверены, что хотите удалить заметку? (y/n)",
		Default: "y",
		Validate: func(input string) error {
			if input != "y" && input != "n" {
				return errors.New("пожалуйста, введите 'y' или 'n'")
			}
			return nil
		},
	}
	confirm, err := continuePrompt.Run()
	if err != nil || confirm == "n" {
		return
	}
	// удаление заметки
	err = c.noteDeleteService.NoteDeleteService(context.Background(), noteID, c.token.Token)
	if err != nil {
		fprintln(c.writer, "\033[0m"+"Ошибка при удалении заметки:", err)
		waitEnter(c.reader)
		return
	}
	// вывод сообщения об удалении заметки
	fprintln(c.writer, "\033[0m"+"Заметка удалена")
	waitEnter(c.reader)
}
