package repository

import (
	"web_service_auth/internal/app/ds"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Repository) GetCarbonateAcids(carbonateID uint) ([]ds.CarbonateAcid, error) {
	var acids []ds.CarbonateAcid
	err := r.db.Where("carbonate_id = ?", carbonateID).Preload("Acid").Preload("Carbonate").Find(&acids).Error
	return acids, err
}

func (r *Repository) GetAcidCount(userID uuid.UUID) int64 {
	var carbonateID uint
	var count int64

	err := r.db.Model(&ds.Carbonate{}).Where("creator_id = ? AND status = ?", userID, "черновик").Select("id").First(&carbonateID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CarbonateAcid{}).Where("carbonate_id = ?", carbonateID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in carbonate_acids:", err)
	}

	return count
}

func (r *Repository) UpdateCarbonateAcidAmount(carbonateID, acidID uint, mass float32) error {
	return r.db.Model(&ds.CarbonateAcid{}).
		Where("carbonate_id = ? AND acid_id = ?", carbonateID, acidID).
		Update("mass", mass).Error
}

func (r *Repository) UpdateCarbonateAcidResult(carbonateAcidID uint, result float32) error {
	return r.db.Model(&ds.CarbonateAcid{}).
		Where("id = ?", carbonateAcidID).
		Update("result", result).Error
}

func (r *Repository) RemoveAcidFromCarbonate(carbonateID, acidID uint) error {
	return r.db.Where("carbonate_id = ? AND acid_id = ?", carbonateID, acidID).
		Delete(&ds.CarbonateAcid{}).Error
}

func (r *Repository) GetCalculated(id uint) int64 {
	var count int64
	err := r.db.Model(&ds.CarbonateAcid{}).Where("carbonate_id = ? AND result != 0", id).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}
