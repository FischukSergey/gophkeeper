package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const (
	errOutputMessage   = "Ошибка вывода сообщения: %s\n"
	errReadMessage     = "Ошибка чтения ответа: %s\n"
	errInputMessage    = "Ошибка ввода: %s\n"
	errLoginMessage    = "Ошибка при вводе логина: "
	errPasswordMessage = "Ошибка при вводе пароля: "
	errGetCardsMessage = "Ошибка получения списка карт: "
	messageContinue    = "\nНажмите Enter для продолжения..."
)

// ICommand интерфейс для команд.
type ICommand interface {
	Execute()
	Name() string
}

// IAuthService интерфейс для сервиса авторизации.
type IAuthService interface {
	Register(ctx context.Context, login, password string) (string, error)
	Authorization(ctx context.Context, login, password string) (string, error)
	Check(ctx context.Context) error
	S3FileUpload(ctx context.Context, token string, fileData []byte, filename string) (string, error)
	S3FileDownload(ctx context.Context, token string, filename string) ([]byte, error)
	S3FileDelete(ctx context.Context, token string, filename string) error
	GetFileList(ctx context.Context, token string) ([]models.File, error)
}

// INoteService интерфейс для сервиса заметок.
type INoteService interface {
	NoteAdd(ctx context.Context, note string, metadata map[string]string, token string) error
	NoteDelete(ctx context.Context, noteID int64, token string) error
	NoteGetList(ctx context.Context, token string) ([]models.Note, error)
}

// ICardService интерфейс для сервиса карт.
type ICardService interface {
	CardAdd(ctx context.Context, card models.Card, token string) error
	AddCardMetadata(ctx context.Context, cardID int64, metadata map[string]string, token string) error
	DeleteCard(ctx context.Context, cardID int64, token string) error
	GetCardList(ctx context.Context, token string) ([]models.Card, error)
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

	// Проверяем права на запись используя os.CreateTemp
	f, err := os.CreateTemp(path, "gophkeeper-test-*")
	if err != nil {
		return fmt.Errorf("нет прав на запись в директорию '%s': %w", path, err)
	}
	tmpName := f.Name()
	if err = f.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии временного файла: %w", err)
	}
	if err = os.Remove(tmpName); err != nil {
		return fmt.Errorf("ошибка при удалении временного файла: %w", err)
	}
	return nil
}

// checkFileExists проверка существования файла.
func checkFileExists(path, filename string) error {
	filepath := filepath.Join(path, filename)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return fmt.Errorf("файл '%s' не существует", filepath)
	}
	if err != nil {
		return fmt.Errorf("ошибка при проверке файла: %w", err)
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
