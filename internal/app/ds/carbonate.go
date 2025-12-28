package ds

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Carbonate struct {
	ID          uint      `gorm:"primaryKey"`
	Status      string    `gorm:"type:varchar(15);not null"`
	DateCreate  time.Time `gorm:"not null"`
	DateUpdate  time.Time
	DateFinish  sql.NullTime `gorm:"default:null"`
	CreatorID   uuid.UUID    `gorm:"not null"`
	ModeratorID *uuid.UUID   `gorm:"default:null"`
	Mass        float32      `gorm:"default:null"`

	Creator   Users `gorm:"foreignKey:CreatorID"`
	Moderator Users `gorm:"foreignKey:ModeratorID"`
}
