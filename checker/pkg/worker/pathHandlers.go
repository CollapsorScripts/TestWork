package worker

import (
	"checker/pkg/logger"
	"checker/pkg/types"
	"net/http"
)

// parsedCmd - структура для хранения данных после парса запроса
type parsedCmd struct {
	requestID string
	body      string
	method    string
}

// commandHandler - типа данных для удобства работы с хэндлерами
type commandHandler func(*parsedCmd) (*types.Response, error)

// funcHandlers - глобальная мапа для хранения хэндлеров и путей к ним
var funcHandlers map[string]commandHandler

// funcHandlersBeforeInit - мапа для инициализации funcHandlers
var funcHandlersBeforeInit = map[string]commandHandler{
	//Сообщения
	"send": handleSend,
}

//--------------------------HANDLERS----------------------------------------

func handleUnimplemented(cmd *parsedCmd) (*types.Response, error) {
	resp := new(types.Response)
	resp.ID = cmd.requestID
	resp.ErrCode = http.StatusNotImplemented
	resp.ErrString = "Данный метод не реализован"
	return resp, nil
}

//--------------------------------------------------------------------------

// funcUnimplemented - нереализованные хэндлеры
var funcUnimplemented = map[string]commandHandler{
	"unimplemented": handleUnimplemented,
}

// InitHandlers - инициализация хэндлеров
func InitHandlers() {
	funcHandlers = funcHandlersBeforeInit
}

func parseCmd(req *types.Request) (*parsedCmd, error) {
	var err error = nil

	cmd := &parsedCmd{
		requestID: req.ID,
		body:      req.Body,
		method:    req.Method,
	}

	return cmd, err
}

// processRequest - обработка запроса
func (w *KafkaWorker) processRequest(request *types.Request) *types.Response {
	cmd, err := parseCmd(request)
	if err != nil {
		logger.Error("Ошибка при парсе запроса: %v", err)
		return nil
	}

	resp, err := standartCmdResult(cmd)
	if err != nil {
		logger.Error("Ошибка: %v", err)
		resp.ErrCode = http.StatusInternalServerError
		resp.ErrString = err.Error()
	}

	return resp
}

func standartCmdResult(cmd *parsedCmd) (*types.Response, error) {
	handler, ok := funcHandlers[cmd.method]
	if ok {
		goto handled
	} else {
		handler = funcUnimplemented["unimplemented"]
		goto handled
	}

handled:
	return handler(cmd)
}
