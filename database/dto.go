package database

type TokenPaddedPan struct {
	//gorm.Model
	Token     string `gorm:"primaryKey"`
	PaddedPan string
}
