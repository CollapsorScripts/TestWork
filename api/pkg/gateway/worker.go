package gateway

import (
	"api/pkg/logger"
	"api/pkg/types"
	"api/pkg/utilities"
	"encoding/json"
	"github.com/IBM/sarama"
	"net/http"
	"time"
)

func (a *Api) Test(w http.ResponseWriter, r *http.Request) {
	logger.Info("Запрос прошел")
}

func (a *Api) responseWorker(requestID string, w http.ResponseWriter, hasBody bool) {
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

		logger.Info("Ответ который пришел: %s", utilities.ToJSON(response.Body))
		if hasBody {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(response.Body))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNoContent)
				logger.Error("%s", err.Error())
				return
			}
		} else {
			w.WriteHeader(http.StatusOK)
		}
	case <-time.After(15 * time.Second):
		a.mu.Lock()
		delete(a.responseChannels, requestID)
		a.mu.Unlock()
		http.Error(w, "Таймаут ожидания ответа", http.StatusGatewayTimeout)
	}
}
