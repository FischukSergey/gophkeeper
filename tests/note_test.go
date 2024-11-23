package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestNoteAdd(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	//добавляем заметку
	tests := []struct {
		name     string
		note     string
		metaData map[string]string
		token    string
		wantErr  bool
	}{
		{name: "empty token", note: "", metaData: nil, token: "", wantErr: true},
		{name: "empty note", note: "", metaData: nil, token: token, wantErr: true},
		{name: "success", note: gofakeit.Sentence(10), metaData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		}, token: token, wantErr: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := noteService.NoteAdd(context.Background(), test.note, test.metaData, test.token)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// тест на получение заметки.
func TestNoteGet(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)

	noteText := gofakeit.Sentence(10)
	//добавляем заметку
	tests := []struct {
		name     string
		note     string
		metaData map[string]string
		token    string
		wantErr  bool
	}{
		{name: "empty token", note: "", metaData: nil, token: "", wantErr: true},
		{name: "success", note: noteText, metaData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		}, token: token, wantErr: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := noteService.NoteAdd(context.Background(), test.note, test.metaData, test.token)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
	//получаем заметку
	note, err := noteService.NoteGetList(context.Background(), token)
	require.NoError(t, err)
	require.NotNil(t, note)
	require.NotEmpty(t, note)
	//проверяем, что заметка соответствует добавленной
	require.Contains(t, note[0].NoteText, noteText)
}

// тест на удаление заметки.
func TestNoteDelete(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	//добавляем заметку
	err = noteService.NoteAdd(context.Background(), gofakeit.Sentence(10), map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, token)
	require.NoError(t, err)
	//получаем заметку
	note, err := noteService.NoteGetList(context.Background(), token)
	require.NoError(t, err)
	require.NotNil(t, note)
	require.NotEmpty(t, note)
	//удаляем заметку
	err = noteService.NoteDeleteService(context.Background(), note[0].NoteID, token)
	require.NoError(t, err)
	//получаем заметку
	note, err = noteService.NoteGetList(context.Background(), token)
	require.NoError(t, err)
	require.Empty(t, note)
}
