package service

import (
	"errors"
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ─── GetUserByID ────────────────────────────────────────────────────────────

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

// ─── CreateUser ─────────────────────────────────────────────────────────────

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

// ─── RegisterUser ───────────────────────────────────────────────────────────

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 2, Name: "Alice", Email: "alice@mail.com"}
	mockRepo.EXPECT().GetByEmail("alice@mail.com").Return(existing, nil)

	err := svc.RegisterUser(&repository.User{Name: "Alice2"}, "alice@mail.com")
	assert.EqualError(t, err, "user with this email already exists")
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 3, Name: "Bob", Email: "bob@mail.com"}
	mockRepo.EXPECT().GetByEmail("bob@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := svc.RegisterUser(newUser, "bob@mail.com")
	assert.NoError(t, err)
}

func TestRegisterUser_RepoErrorOnCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 4, Name: "Charlie"}
	mockRepo.EXPECT().GetByEmail("charlie@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(errors.New("db error"))

	err := svc.RegisterUser(newUser, "charlie@mail.com")
	assert.EqualError(t, err, "db error")
}

// ─── UpdateUserName ─────────────────────────────────────────────────────────

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.UpdateUserName(1, "")
	assert.EqualError(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("user not found"))

	err := svc.UpdateUserName(99, "NewName")
	assert.EqualError(t, err, "user not found")
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
		// Verify name was actually changed before update
		assert.Equal(t, "NewName", u.Name)
		return nil
	})

	err := svc.UpdateUserName(2, "NewName")
	assert.NoError(t, err)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))

	err := svc.UpdateUserName(2, "NewName")
	assert.EqualError(t, err, "update failed")
}

// ─── DeleteUser ─────────────────────────────────────────────────────────────

func TestDeleteUser_AttemptDeleteAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.DeleteUser(1)
	assert.EqualError(t, err, "it is not allowed to delete admin user")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(2).Return(nil)

	err := svc.DeleteUser(2)
	assert.NoError(t, err)
}

func TestDeleteUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(5).Return(errors.New("db error"))

	err := svc.DeleteUser(5)
	assert.EqualError(t, err, "db error")
}
