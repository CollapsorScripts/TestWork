package worker

import (
	"checker/pkg/types"
	"checker/pkg/utilities"
	"net/http"
)

// handleSend - обработчик эндпоинта /send, искусственно в 20% случаях не обрабатывает сообщения,
// что бы формировать статистику
func handleSend(cmd *parsedCmd) (*types.Response, error) {
	resp := new(types.Response)
	resp.ID = cmd.requestID

	if utilities.RandInt(0, 100) <= 20 {
		resp.ErrCode = http.StatusBadRequest
		resp.ErrString = "Сообщение не обработано"
		return resp, nil
	}

	resp.Body = cmd.body
	resp.ErrCode = http.StatusOK
	return resp, nil
}
