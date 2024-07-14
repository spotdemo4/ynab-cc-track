package models

type Item struct {
	ID   int `gorm:"primaryKey"`
	Name string
}
