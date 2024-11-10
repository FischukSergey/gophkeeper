package command

import (
	"fmt"
	"io"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
)

const (
	errOutputMessage = "Ошибка вывода сообщения: %s\n"
	errReadMessage   = "Ошибка чтения ответа: %s\n"
	errInputMessage  = "Ошибка ввода: %s\n"
	messageContinue  = "\nНажмите Enter для продолжения..."
)

// ICommand интерфейс для команд.
type ICommand interface {
	Execute()
	Name() string
}

// Ожидание нажатия клавиши
func waitEnter(reader io.Reader) {
	fmt.Println(messageContinue)
	buffer := make([]byte, 1)
	_, err := reader.Read(buffer)
	if err != nil {
		fmt.Printf(errInputMessage, err)
	}
}

// проверка наличия токена
func checkToken(token *grpcclient.Token, reader io.Reader) {
	if token.Token == "" {
		fmt.Println("Вы не авторизованы. Авторизуйтесь с помощью команды login.")
		waitEnter(reader)
	}
}

// validatePath проверяет существование и доступность указанного пути
func validatePath(path string) error {
	// Проверяем существование директории
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Путь '%s' не существует\n", path)
			return err
		}
		fmt.Printf("Ошибка при проверке пути: %v\n", err)
		return err
	}

	// Проверяем, что это директория
	if !fileInfo.IsDir() {
		fmt.Printf("Путь '%s' не является директорией\n", path)
		return err
	}

	// Проверяем права на запись
	tmpFile := path + "/.tmp_test"
	f, err := os.Create(tmpFile)
	if err != nil {
		fmt.Printf("Нет прав на запись в директорию '%s': %v\n", path, err)
		return err
	}
	f.Close()
	os.Remove(tmpFile)
	return nil
}

// проверка существования файла
func checkFileExists(path, filename string) error {
	filepath := path + "/" + filename
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	return nil
}
