package repository

import (
	"web_service_auth/internal/app/ds"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetUserByLogin(login string) (*ds.Users, error) {
	var user ds.Users
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uuid.UUID) (*ds.Users, error) {
	var user ds.Users
	err := r.db.Where("uuid = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *ds.Users) error {
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(id uuid.UUID, updates map[string]interface{}) error {
	return r.db.Model(&ds.Users{}).Where("uuid = ?", id).Updates(updates).Error
}
