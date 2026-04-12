package repo

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"

	"github.com/google/uuid"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	if err := u.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
	var userFromDB entity.User
	if err := u.PG.Conn.Where("username = ?", user.Username).First(&userFromDB).Error; err != nil {
		return nil, fmt.Errorf("username not found: %v", err)
	}
	return &userFromDB, nil
}

// GetUserByID fetches a single user by their UUID primary key.
func (u *UserRepo) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}

// PromoteUser sets role = "admin" for the given user and returns the updated record.
func (u *UserRepo) PromoteUser(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if err := u.PG.Conn.Model(&user).Update("role", "admin").Error; err != nil {
		return nil, fmt.Errorf("promote user: %v", err)
	}
	return &user, nil
}
