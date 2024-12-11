package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/manifoldco/promptui"
)

const commandNoteDelete = "Delete"

// NoteDeleteCommand структура для удаления заметки.
type NoteDeleteCommand struct {
	noteDeleteService INoteService
	token             *grpcclient.Token
	reader            io.Reader
	writer            io.Writer
}

// NewCommandNoteDelete создание новой команды удаления заметки.
func NewCommandNoteDelete(
	noteDeleteService INoteService,
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
	// получение списка заметок
	notes, err := c.fetchNotes()
	if err != nil {
		fprintln(c.writer, err.Error())
		waitEnter(c.reader)
		return
	}
	// ввод номера заметки
	noteID, err := c.inputNoteID()
	if err != nil {
		fprintln(c.writer, err.Error())
		waitEnter(c.reader)
		return
	}
	//проверяем, что есть такая заметка перебором списка заметок
	var exist bool
	for _, note := range notes {
		if note.NoteID == noteID {
			exist = true
			break
		}
	}
	if !exist {
		fprintln(c.writer, "Заметка с таким ID не найдена")
		waitEnter(c.reader)
		return
	}
	// подтверждение удаления заметки
	confirm, err := c.confirmDelete()
	if err != nil {
		return
	}
	if !confirm {
		return
	}
	// удаление заметки
	err = c.deleteNote(noteID)
	if err != nil {
		return
	}
	// вывод сообщения об удалении заметки
	fprintln(c.writer, "\033[0m"+"Заметка удалена")
	waitEnter(c.reader)
}

// fetchNotes получение списка заметок.
func (c *NoteDeleteCommand) fetchNotes() ([]models.Note, error) {
	notes, err := c.noteDeleteService.NoteGetList(context.Background(), c.token.Token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка заметок: %w", err)
	}
	if len(notes) == 0 {
		return nil, errors.New("список заметок пуст")
	}
	return notes, nil
}

// inputNoteID ввод номера заметки.
func (c *NoteDeleteCommand) inputNoteID() (int64, error) {
	// ввод номера заметки
	fprint(c.writer, "Введите номер заметки для удаления:")
	var noteID int64
	_, err := fmt.Fscanln(c.reader, &noteID)
	if err != nil {
		return 0, fmt.Errorf("\033[0m"+"Неверный номер заметки: %w", err)
	}
	return noteID, nil
}

// confirmDelete подтверждение удаления заметки.
func (c *NoteDeleteCommand) confirmDelete() (bool, error) {
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
		return false, fmt.Errorf("ошибка при подтверждении удаления заметки: %w", err)
	}
	return true, nil
}

// deleteNote удаление заметки.
func (c *NoteDeleteCommand) deleteNote(noteID int64) error {
	err := c.noteDeleteService.NoteDelete(context.Background(), noteID, c.token.Token)
	if err != nil {
		return fmt.Errorf("ошибка при удалении заметки: %w", err)
	}
	return nil
}
