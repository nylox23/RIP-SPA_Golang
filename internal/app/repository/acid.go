package repository

import (
	"web_service_auth/internal/app/ds"

	"gorm.io/gorm/clause"
)

func (r *Repository) GetAcidsWithFilter(search string) ([]ds.Acid, error) {
	var acids []ds.Acid

	query := r.db.Model(&ds.Acid{}).Where("is_delete = false")

	if search != "" {
		query = query.Where("name_ext ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&acids).Error; err != nil {
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

func (r *Repository) CreateAcid(acid *ds.Acid) error {
	return r.db.Create(acid).Error
}

func (r *Repository) UpdateAcid(id int, updates map[string]interface{}) error {
	return r.db.Model(&ds.Acid{}).Where("id = ? AND is_delete = false", id).Updates(updates).Error
}

func (r *Repository) DeleteAcid(id int) error {
	return r.db.Model(&ds.Acid{}).Where("id = ?", id).Update("is_delete", true).Error
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
