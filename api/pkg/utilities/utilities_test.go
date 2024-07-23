package utilities

import (
	"api/pkg/types"
	"encoding/json"
	_ "github.com/joho/godotenv"
	"testing"
)

func TestAny(t *testing.T) {
	request := types.DatabaseRequest{
		ID:      "123",
		Body:    string("111s"),
		AdminID: 1,
		Method:  "edit-category",
	}

	reqStr := ToJSON(request)
	t.Logf("request: %s", reqStr)
	request1 := types.DatabaseRequest{}

	err := json.Unmarshal([]byte(reqStr), &request1)
	if err != nil {
		t.Errorf("Ошибка: %v", err)
		return
	}

	t.Logf("Длина мапы: %d", len(request1.Params))

	t.Logf("Новый объект из строки: %+v", request1)

	str := "{\"ID\": \"7\"}"
	m1 := make(map[string]string)
	t.Logf("JSON строка: %s", str)
	if len(m1) == 0 {
		t.Logf("Мапа Nil")
	}

	err = json.Unmarshal([]byte(str), &m1)
	if err != nil {
		t.Errorf("Ошибка: %v", err)
		return
	}

	t.Logf("Мапа: %+v", m1)

	str2 := ToJSON(m1)

	t.Logf("Строка: %s", str2)
}
