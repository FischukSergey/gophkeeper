package command

import (
	"fmt"
	"io"
	"os"
)

const nameCommand = "exit"

// commandExit структура для команды выхода.
type commandExit struct {
	reader io.Reader
	writer io.Writer
}

// NewCommandExit создание команды выхода.
func NewCommandExit(reader io.Reader, writer io.Writer) *commandExit {
	return &commandExit{reader: reader, writer: writer}
}

// Execute выполнение команды выхода.
func (c *commandExit) Execute() {
	_, err := fmt.Fprint(c.writer, "Действительно хотите выйти? (y/n): ")
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	var response string
	_, err = fmt.Fscan(c.reader, &response)
	if err != nil {
		fmt.Printf(errReadMessage, err)
		return
	}
	if response != "y" && response != "Y" {
		fmt.Println("Exit отменен.")
		return
	}
	os.Exit(0)
}

// Name возвращает имя команды.
func (c *commandExit) Name() string {
	return nameCommand
}
