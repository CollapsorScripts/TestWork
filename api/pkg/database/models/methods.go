package models

import (
	"api/pkg/database"
	"api/pkg/logger"
)

// FindOne - Поиск записи по id (возвращает все поля)
func (m *ModelDB) FindOne() error {
	db := database.GetDB()
	query := "SELECT * FROM model_dbs WHERE id = ?"
	result := db.Raw(query, m.ID).Scan(&m)

	return result.Error
}

// FindOne (не методом) - Поиск записи по id (возвращает все поля)
func FindOne(id uint) (*ModelDB, error) {
	db := database.GetDB()

	model := &ModelDB{}
	query := "SELECT * FROM model_dbs WHERE id = ?"
	result := db.Raw(query, id).Scan(&model)

	return model, result.Error
}

// Create - создание новой записи
func (m *ModelDB) Create() error {
	db := database.GetDB()
	var ID uint
	query := "INSERT INTO model_dbs (request_id, message) VALUES (?, ?) RETURNING id"
	result := db.Raw(query, m.RequestID, m.Message).Scan(&ID)

	if result.Error != nil {
		return result.Error
	}

	m, err := FindOne(ID)

	return err
}

// AllModels - все записи
func AllModels() ([]*ModelDB, error) {
	db := database.GetDB()
	var models []*ModelDB
	query := "SELECT * FROM model_dbs"
	result := db.Raw(query)
	if result.Error != nil {
		return nil, result.Error
	}

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var model ModelDB
		err := db.ScanRows(rows, &model)
		if err != nil {
			return nil, err
		}
		models = append(models, &model)
	}

	return models, nil
}

// UpdateMark - обновлении записи
func (m *ModelDB) UpdateMark() error {
	db := database.GetDB()
	query := "UPDATE model_dbs SET marked = ? WHERE request_id = ?"
	logger.Info("ID обновления: %d", m.RequestID)
	result := db.Exec(query, !m.Marked, m.RequestID)
	return result.Error
}

// Delete - удаление записи
func (m *ModelDB) Delete() error {
	db := database.GetDB()
	query := "DELETE FROM model_dbs WHERE id = ?"
	result := db.Exec(query, m.ID)

	return result.Error
}

// DeleteById - удаление записи по ID
func DeleteById(id uint) error {
	db := database.GetDB()
	query := "DELETE FROM model_dbs WHERE id = ?"
	result := db.Exec(query, id)

	return result.Error
}

// CountModels - количество записей
func CountModels() (int64, error) {
	db := database.GetDB()
	query := "SELECT COUNT(*) FROM model_dbs"
	var count int64
	result := db.Raw(query).Scan(&count)

	return count, result.Error
}
