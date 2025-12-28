package ds

import (
	"web_service_auth/internal/app/role"

	"github.com/google/uuid"
)

type Users struct {
	UUID     uuid.UUID `gorm:"primary_key;type:uuid" json:"id"`
	Login    string    `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password string    `gorm:"type:varchar(100);not null" json:"password"`
	Role     role.Role `sql:"type:string;"`
}
