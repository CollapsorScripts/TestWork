package models

import (
	"api/cmd/config"
	"api/pkg/logger"
	"testing"
)

func TestAllModels(t *testing.T) {
	config.Init()
	logger.New()
	tests := []struct {
		name string
	}{
		{
			name: "Все записи",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AllModels()
			if err != nil {
				t.Errorf("AllModels() error = %v", err)
				return
			}

			t.Logf("Записи: %+v", got)
		})
	}
}

func TestDeleteById(t *testing.T) {
	config.Init()
	type args struct {
		id uint
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Удаление записи",
			args: args{id: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteById(tt.args.id); err != nil {
				t.Errorf("DeleteById() error = %v", err)
			}
		})
	}
}

func TestFindOne(t *testing.T) {
	config.Init()
	type args struct {
		id uint
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Поиск записи по ID",
			args: args{id: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindOne(tt.args.id)
			if err != nil {
				t.Errorf("FindOne() error = %v", err)
				return
			}

			t.Logf("Запись: %+v", got)
		})
	}
}

func TestModelDB_Create(t *testing.T) {
	config.Init()
	type fields struct {
		ID        uint
		RequestID string
		Message   string
		Marked    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Создание записи",
			fields: fields{
				RequestID: "1",
				Message:   "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelDB{
				RequestID: tt.fields.RequestID,
				Message:   tt.fields.Message,
			}
			t.Logf("Модель: %+v", m)
			if err := m.Create(); err != nil {
				t.Errorf("Create() error = %v", err)
			}
		})
	}
}

func TestModelDB_Delete(t *testing.T) {
	config.Init()
	type fields struct {
		ID        uint
		RequestID string
		Message   string
		Marked    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Удаление записи методом",
			fields: fields{
				ID: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelDB{
				ID: tt.fields.ID,
			}
			if err := m.Delete(); err != nil {
				t.Errorf("Delete() error = %v", err)
			}
		})
	}
}

func TestModelDB_Update(t *testing.T) {
	config.Init()
	type fields struct {
		ID        uint
		RequestID string
		Message   string
		Marked    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Обновление маркера",
			fields: fields{
				ID:     3,
				Marked: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelDB{
				ID:        tt.fields.ID,
				RequestID: tt.fields.RequestID,
				Message:   tt.fields.Message,
				Marked:    tt.fields.Marked,
			}
			if err := m.UpdateMark(); err != nil {
				t.Errorf("Update() error = %v", err)
			}
		})
	}
}
