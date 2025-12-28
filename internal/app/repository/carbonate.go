package repository

import (
	"time"
	"web_service_auth/internal/app/ds"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetDraftCarbonateID(creatorID uuid.UUID) uint {
	var carbonateID uint
	err := r.db.Model(&ds.Carbonate{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&carbonateID).Error
	if err != nil {
		return 0
	}
	return carbonateID
}

func (r *Repository) GetDraftCarbonate(creatorID uuid.UUID, create bool) (uint, error) {
	carbonateID := r.GetDraftCarbonateID(creatorID)
	if carbonateID == 0 && create {
		carbonate := ds.Carbonate{
			Status:     "черновик",
			DateCreate: time.Now(),
			CreatorID:  creatorID,
		}
		err := r.db.Create(&carbonate).Error
		if err != nil {
			return 0, err
		}
		carbonateID = carbonate.ID
	}
	return carbonateID, nil
}

func (r *Repository) GetCarbonatesWithFilter(status string, dateFrom, dateTo time.Time, creator uuid.UUID) ([]ds.Carbonate, error) {
	var carbonates []ds.Carbonate

	query := r.db.Model(&ds.Carbonate{}).Preload("Creator").Preload("Moderator")
	query = query.Where("status != ?", "удален")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if !dateFrom.IsZero() {
		query = query.Where("date_create >= ?", dateFrom)
	}

	if !dateTo.IsZero() {
		query = query.Where("date_create <= ?", dateTo)
	}

	if creator != uuid.Nil {
		query = query.Where("creator_id = ?", creator)
	}

	if err := query.Order("date_create DESC").Find(&carbonates).Error; err != nil {
		return nil, err
	}

	return carbonates, nil
}

func (r *Repository) GetCarbonateByID(id uint) (*ds.Carbonate, error) {
	var carbonate ds.Carbonate
	err := r.db.Preload("Creator").Preload("Moderator").Where("id = ? AND status != ?", id, "удален").First(&carbonate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &carbonate, nil
}

func (r *Repository) UpdateCarbonate(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.Carbonate{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) DeleteCarbonate(id uint) error {
	return r.db.Model(&ds.Carbonate{}).Where("id = ?", id).Update("status", "удален").Error
}
