package models

// ModelDB - структура БД
type ModelDB struct {
	ID        uint
	RequestID string
	Message   string
	Marked    bool `gorm:"default:false"`
}
