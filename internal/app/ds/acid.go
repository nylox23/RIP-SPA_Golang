package ds

type Acid struct {
	ID        int     `gorm:"primaryKey"`
	IsDelete  bool    `gorm:"type:boolean not null;default:false"`
	Img       string  `gorm:"type:varchar(100)"`
	NameExt   string  `gorm:"type:varchar(25);not null"`
	Info      string  `gorm:"type:varchar(500)"`
	Name      string  `gorm:"type:varchar(15);not null"`
	Hplus     int     `gorm:"type:integer;not null"`
	MolarMass float32 `gorm:"type:numeric;not null"`
}
