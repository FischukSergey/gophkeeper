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

// waitEnter ожидание нажатия клавиши.
func waitEnter(reader io.Reader) {
	fmt.Println(messageContinue)
	buffer := make([]byte, 1)
	_, err := reader.Read(buffer)
	if err != nil {
		fmt.Printf(errInputMessage, err)
	}
}

// checkToken проверка наличия токена.
func checkToken(token *grpcclient.Token, reader io.Reader) bool {
	if token.Token == "" {
		fmt.Println("Вы не авторизованы. Авторизуйтесь с помощью команды login.")
		waitEnter(reader)
		return false
	}
	return true
}

// validatePath проверяет существование и доступность указанного пути.
func validatePath(path string) error {
	// Проверяем существование директории.
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("путь '%s' не существует", path)
		}
		return fmt.Errorf("ошибка при проверке пути: %w", err)
	}

	// Проверяем, что это директория.
	if !fileInfo.IsDir() {
		return fmt.Errorf("путь '%s' не является директорией", path)
	}

	// Проверяем права на запись.
	tmpFile := path + "/.tmp_test"
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("нет прав на запись в директорию '%s': %w", path, err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии временного файла: %w", err)
	}
	err = os.Remove(tmpFile)
	if err != nil {
		return fmt.Errorf("ошибка при удалении временного файла: %w", err)
	}
	return nil
}

// checkFileExists проверка существования файла.
func checkFileExists(path, filename string) error {
	filepath := path + "/" + filename
	_, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("файл '%s' не существует", filepath)
	}
	return nil
}

// fprintln обертка над fmt.Fprintln, игнорирующая ошибку.
func fprintln(w io.Writer, a ...any) {
	_, _ = fmt.Fprintln(w, a...)
}

// fprint обертка над fmt.Fprint, игнорирующая ошибку.
func fprint(w io.Writer, a ...any) {
	_, _ = fmt.Fprint(w, a...)
}

// fscanln обертка над fmt.Fscanln, игнорирующая ошибку.
func fscanln(w io.Reader, a ...any) {
	_, _ = fmt.Fscanln(w, a...)
}

// fprintf обертка над fmt.Fprintf, игнорирующая ошибку.
func fprintf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}
