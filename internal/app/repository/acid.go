package repository

import (
	"spa_InDb/internal/app/ds"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetAllAcids() ([]ds.Acid, error) {
	var acids []ds.Acid
	err := r.db.Where("is_delete = false").Find(&acids).Error
	if err != nil {
		return nil, err
	}
	return acids, nil
}

func (r *Repository) GetAcidByID(id int) (*ds.Acid, error) {
	acid := &ds.Acid{}
	err := r.db.Where("is_delete = false and id = $1", id).First(&acid).Error
	if err != nil {
		return nil, err
	}
	return acid, nil
}

func (r *Repository) SearchAcidsByName(name string) ([]ds.Acid, error) {
	var acids []ds.Acid
	err := r.db.Where("Name_ext ILIKE ? and is_delete = ?", "%"+name+"%", false).Find(&acids).Error
	if err != nil {
		return nil, err
	}
	return acids, nil
}

func (r *Repository) GetAcidCount() int64 {
	var carbonateID uint
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.Carbonate{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&carbonateID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CarbonateAcid{}).Where("carbonate_id = ?", carbonateID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in carbonate_acids:", err)
	}

	return count
}

func (r *Repository) GetDraftCarbonateID(creatorID uint) uint {
	var carbonateID uint
	err := r.db.Model(&ds.Carbonate{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&carbonateID).Error
	if err != nil {
		return 0
	}
	return carbonateID
}

func (r *Repository) GetDraftCarbonate(creatorID uint) (uint, error) {
	carbonateID := r.GetDraftCarbonateID(creatorID)
	if carbonateID == 0 {
		carbonate := ds.Carbonate{
			Status:      "черновик",
			DateCreate:  time.Now(),
			CreatorID:   creatorID,
			ModeratorID: nil,
		}
		err := r.db.Create(&carbonate).Error
		if err != nil {
			return 0, err
		}
		carbonateID = carbonate.ID
	}
	return carbonateID, nil
}

func (r *Repository) AddToCarbonate(acidId, carbonateID uint) error {
	ca := ds.CarbonateAcid{
		AcidID:      acidId,
		CarbonateID: carbonateID,
		Mass:        0,
	}

	result := r.db.Model(&ds.CarbonateAcid{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&ca)

	return result.Error
}

func (r *Repository) GetCarbonate(carbonateID uint) ([]ds.CarbonateAcid, error) {
	var rows []ds.CarbonateAcid
	err := r.db.Where("carbonate_id = ?", carbonateID).Preload("Acid").Preload("Carbonate").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) GetCarbonateAmount(carbonateID uint) float32 {
	var res float32
	err := r.db.Model(&ds.Carbonate{}).Where("id = ?", carbonateID).Select("mass").First(&res).Error
	if err != nil {
		return 0
	}
	return res
}

func (r *Repository) UpdateCarbonateStatus(carbonateID uint, status string) error {
	query := "UPDATE carbonates SET status = ? WHERE id = ?"
	if err := r.db.Exec(query, status, carbonateID).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetCarbonateStatus(carbonateID uint) (string, error) {
	var carbonate ds.Carbonate
	if err := r.db.Where("id = ?", carbonateID).First(&carbonate).Error; err != nil {
		return "", err
	}

	return carbonate.Status, nil
}
