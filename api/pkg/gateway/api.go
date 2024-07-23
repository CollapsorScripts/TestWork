package gateway

import (
	"api/cmd/config"
	"api/pkg/logger"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"sync"
	"time"
)

const (
	databaseResponseTopic = "databaseResponse"
	databaseRequestTopic  = "databaseRequest"
)

type Api struct {
	r                     *mux.Router
	responseChannels      map[string]chan *sarama.ConsumerMessage
	mu                    sync.Mutex
	producer              *sarama.SyncProducer
	consumer              *sarama.Consumer
	gatewayConsumer       *sarama.PartitionConsumer
	databaseResponseTopic string
	databaseRequestTopic  string
}

func New() (*Api, error) {
	r := mux.NewRouter()

	// Создание продюсера Kafka
	configKafka := sarama.NewConfig()
	configKafka.Producer.MaxMessageBytes = 1073741824
	configKafka.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Cfg.Kafka, configKafka)
	if err != nil {
		return nil, err
	}

	// Создание консьюмера Kafka
	consumer, err := sarama.NewConsumer(config.Cfg.Kafka, nil)
	if err != nil {
		producer.Close()
		return nil, err
	}

	// Подписка на партицию "authResponse" в Kafka
	gatewayConsumer, err := consumer.ConsumePartition(databaseResponseTopic, 0, sarama.OffsetNewest)
	if err != nil {
		producer.Close()
		consumer.Close()
		return nil, err
	}

	api := &Api{
		r:                     r,
		responseChannels:      make(map[string]chan *sarama.ConsumerMessage),
		mu:                    sync.Mutex{},
		producer:              &producer,
		consumer:              &consumer,
		gatewayConsumer:       &gatewayConsumer,
		databaseResponseTopic: databaseResponseTopic,
		databaseRequestTopic:  databaseRequestTopic,
	}

	go api.kafkaPoll()

	return api, nil
}

//------------------Подзагрузка маршрутов-----------------------------------

func (a *Api) PreloadRoutes() *http.Server {
	r := mux.NewRouter()
	r.Use(cors.Default().Handler, mux.CORSMethodMiddleware(r))
	//"github.com/google/uuid"

	//REST
	{
		//Отправка сообщения
		r.HandleFunc("/send", a.Send).Methods(http.MethodPost,
			http.MethodOptions)

		//Получени статистики
		r.HandleFunc("/statistic", a.Statistic).Methods(http.MethodGet,
			http.MethodOptions)
	}

	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := crs.Handler(r)

	srv := &http.Server{
		Addr:         ":8070",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.AllowAll().Handler(handler),
	}

	return srv
}

// kafkaPoll - прослушивает сообщения от kafka
func (a *Api) kafkaPoll() {
	for {
		select {
		case msg, ok := <-(*a.gatewayConsumer).Messages():
			{
				if !ok {
					logger.Warn("Канал закрыт, выход из горутины.")
					return
				}
				responseID := string(msg.Key)
				a.mu.Lock()
				ch, exists := a.responseChannels[responseID]
				if exists {
					ch <- msg
					delete(a.responseChannels, responseID)
				}
				a.mu.Unlock()
			}
		}
	}
}
