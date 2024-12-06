// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/handlers/handlers.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	reflect "reflect"

	models "github.com/FischukSergey/gophkeeper/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockProtoKeeperSaver is a mock of ProtoKeeperSaver interface.
type MockProtoKeeperSaver struct {
	ctrl     *gomock.Controller
	recorder *MockProtoKeeperSaverMockRecorder
}

// MockProtoKeeperSaverMockRecorder is the mock recorder for MockProtoKeeperSaver.
type MockProtoKeeperSaverMockRecorder struct {
	mock *MockProtoKeeperSaver
}

// NewMockProtoKeeperSaver creates a new mock instance.
func NewMockProtoKeeperSaver(ctrl *gomock.Controller) *MockProtoKeeperSaver {
	mock := &MockProtoKeeperSaver{ctrl: ctrl}
	mock.recorder = &MockProtoKeeperSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProtoKeeperSaver) EXPECT() *MockProtoKeeperSaverMockRecorder {
	return m.recorder
}

// Authorization mocks base method.
func (m *MockProtoKeeperSaver) Authorization(ctx context.Context, login, password string) (models.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorization", ctx, login, password)
	ret0, _ := ret[0].(models.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorization indicates an expected call of Authorization.
func (mr *MockProtoKeeperSaverMockRecorder) Authorization(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorization", reflect.TypeOf((*MockProtoKeeperSaver)(nil).Authorization), ctx, login, password)
}

// FileDeleteFromS3 mocks base method.
func (m *MockProtoKeeperSaver) FileDeleteFromS3(ctx context.Context, userID int64, filename string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileDeleteFromS3", ctx, userID, filename)
	ret0, _ := ret[0].(error)
	return ret0
}

// FileDeleteFromS3 indicates an expected call of FileDeleteFromS3.
func (mr *MockProtoKeeperSaverMockRecorder) FileDeleteFromS3(ctx, userID, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileDeleteFromS3", reflect.TypeOf((*MockProtoKeeperSaver)(nil).FileDeleteFromS3), ctx, userID, filename)
}

// FileDownloadFromS3 mocks base method.
func (m *MockProtoKeeperSaver) FileDownloadFromS3(ctx context.Context, userID int64, filename string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileDownloadFromS3", ctx, userID, filename)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileDownloadFromS3 indicates an expected call of FileDownloadFromS3.
func (mr *MockProtoKeeperSaverMockRecorder) FileDownloadFromS3(ctx, userID, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileDownloadFromS3", reflect.TypeOf((*MockProtoKeeperSaver)(nil).FileDownloadFromS3), ctx, userID, filename)
}

// FileGetListFromS3 mocks base method.
func (m *MockProtoKeeperSaver) FileGetListFromS3(ctx context.Context, userID int64) ([]models.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileGetListFromS3", ctx, userID)
	ret0, _ := ret[0].([]models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileGetListFromS3 indicates an expected call of FileGetListFromS3.
func (mr *MockProtoKeeperSaverMockRecorder) FileGetListFromS3(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileGetListFromS3", reflect.TypeOf((*MockProtoKeeperSaver)(nil).FileGetListFromS3), ctx, userID)
}

// FileUploadToS3 mocks base method.
func (m *MockProtoKeeperSaver) FileUploadToS3(ctx context.Context, fileData io.Reader, filename string, userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileUploadToS3", ctx, fileData, filename, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileUploadToS3 indicates an expected call of FileUploadToS3.
func (mr *MockProtoKeeperSaverMockRecorder) FileUploadToS3(ctx, fileData, filename, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileUploadToS3", reflect.TypeOf((*MockProtoKeeperSaver)(nil).FileUploadToS3), ctx, fileData, filename, userID)
}

// Ping mocks base method.
func (m *MockProtoKeeperSaver) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockProtoKeeperSaverMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockProtoKeeperSaver)(nil).Ping), ctx)
}

// RegisterUser mocks base method.
func (m *MockProtoKeeperSaver) RegisterUser(ctx context.Context, login, password string) (models.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, login, password)
	ret0, _ := ret[0].(models.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockProtoKeeperSaverMockRecorder) RegisterUser(ctx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockProtoKeeperSaver)(nil).RegisterUser), ctx, login, password)
}
