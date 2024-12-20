package command

import (
	"context"
	"io"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
)

const nameCommandGetList = "Getlist"

type NoteGetListCommand struct {
	noteGetListService INoteService
	token              *grpcclient.Token
	reader             io.Reader
	writer             io.Writer
}

// NewCommandNoteGetList создает команду для получения списка заметок.
func NewCommandNoteGetList(
	noteGetListService INoteService, // сервис для получения списка заметок
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer) *NoteGetListCommand {
	return &NoteGetListCommand{
		noteGetListService: noteGetListService,
		token:              token,
		reader:             reader,
		writer:             writer,
	}
}

// Name возвращает имя команды.
func (c *NoteGetListCommand) Name() string {
	return nameCommandGetList
}

// Execute выполняет команду.
func (c *NoteGetListCommand) Execute() {
	if !checkToken(c.token, c.reader) {
		return
	}
	fprintln(c.writer, "Список заметок:")
	//вызываем сервис получения списка заметок
	notes, err := c.noteGetListService.NoteGetList(context.Background(), c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при получении списка заметок:", err)
		return
	}
	if len(notes) == 0 {
		fprintln(c.writer, "Список заметок пуст")
		waitEnter(c.reader)
		return
	}
	//выводим список заметок
	for _, note := range notes {
		//выводим номер заметки
		fprintln(c.writer, "\nНомер заметки: "+"\033[35m"+strconv.FormatInt(note.NoteID, 10)+"\033[0m")
		//выводим текст заметки
		fprintln(c.writer, "Текст заметки:")
		fprintln(c.writer, "\033[32m"+note.NoteText+"\033[0m")
		//выводим метаданные заметки построчно
		fprintln(c.writer, "Метаданные заметки:")
		for _, meta := range note.Metadata {
			fprintln(c.writer, "\033[36m"+meta.Key+"\033[0m: \033[32m"+meta.Value+"\033[0m")
		}
	}
	waitEnter(c.reader)
}
