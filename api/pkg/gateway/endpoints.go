package gateway

import (
	"api/pkg/database/models"
	"api/pkg/logger"
	"api/pkg/types"
	"api/pkg/utilities"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// Send - отправляет сообщение
func (a *Api) Send(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()

	var model *models.ModelDB

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	model.RequestID = requestID

	if err := model.Create(); err != nil {
		logger.Error("Ошибка при создании записи: %v", err)
		return
	}

	request := types.Request{
		ID:     requestID,
		Body:   utilities.ToJSON(model),
		Method: "send",
	}

	logger.Info("Что будем отправлять: %s", utilities.ToJSON(request))

	msg := &sarama.ProducerMessage{
		Topic: a.databaseRequestTopic,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.ByteEncoder(utilities.ToJSON(request)),
	}

	// отправка сообщения в Kafka
	_, _, err := (*a.producer).SendMessage(msg)
	if err != nil {
		logger.Error("Ошибка при отправлении сообщения в kafka: %v", err)
		http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	responseCh := make(chan *sarama.ConsumerMessage)
	a.mu.Lock()
	a.responseChannels[requestID] = responseCh
	a.mu.Unlock()

	select {
	case responseMsg := <-responseCh:
		response := types.Response{}
		if err := json.Unmarshal(responseMsg.Value, &response); err != nil {
			logger.Error("Ошибка при десериализации ответа: %v", err)
			http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
			return
		}

		if response.ErrCode != 200 && len(response.ErrString) != 0 {
			logger.Error("Код ошибки: %d", response.ErrCode)
			http.Error(w, response.ErrString, response.ErrCode)
			return
		}

		err := json.Unmarshal([]byte(response.Body), &model)
		if err != nil {
			logger.Error("Ошибка при десериализации ответа: %v", err)
			http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
			return
		}

		err = model.UpdateMark()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	case <-time.After(15 * time.Second):
		a.mu.Lock()
		delete(a.responseChannels, requestID)
		a.mu.Unlock()
		http.Error(w, "Таймаут ожидания ответа", http.StatusGatewayTimeout)
	}
}

// Statistic - Статистика по сообщениям
func (a *Api) Statistic(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Mark   []models.ModelDB
		Unmark []models.ModelDB
		Count  int64
	}{}

	response.Mark = []models.ModelDB{}
	response.Unmark = []models.ModelDB{}
	response.Count = 0

	count, err := models.CountModels()
	if err != nil {
		http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	response.Count = count

	allModels, err := models.AllModels()
	if err != nil {
		http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	for _, model := range allModels {
		if model.Marked {
			response.Mark = append(response.Mark, *model)
		} else {
			response.Unmark = append(response.Unmark, *model)
		}
	}

	_, err = w.Write([]byte(utilities.ToJSON(response)))
	if err != nil {
		http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}
}
