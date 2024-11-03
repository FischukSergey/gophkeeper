package command

import (
	"fmt"
	"io"
	"os"

	"github.com/manifoldco/promptui"
)

const nameCommandExit = "Exit"

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
	confirmation := promptui.Prompt{
		Label: "Действительно хотите выйти? (y/n)",
	}
	response, err := confirmation.Run()
	if err != nil {
		fmt.Printf(errOutputMessage, err)
		return
	}
	if response != "y" && response != "Y" {
		return
	}
	os.Exit(0)
}

// Name возвращает имя команды.
func (c *commandExit) Name() string {
	return nameCommandExit
}
