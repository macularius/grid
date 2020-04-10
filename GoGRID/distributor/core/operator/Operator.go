package operator

import (
	"grid/GoGRID/distributor/core/entities"
	"grid/GoGRID/distributor/core/settings"
	"log"
	"net"
	"net/http"
)

var (
	instance    *operator
	workersChIn chan<- *entities.Worker
)

// GetOperator возвращает экземпляр оператора
func GetOperator() IOperator {
	if instance == nil {
		instance = new(operator)
		instance.Init()
	}

	return instance
}

// IOperator интерфейс оператора - синглтон
type IOperator interface {
	WorkerListener(chan<- *entities.Worker) error // ожидает новые воркеры
}

// Operator получатель и отправитель сообщений
type operator struct {
}

// Init инициализирует оператор
func (o *operator) Init() (err error) {
	return
}

// WorkerListener ожидает новые воркеры
func (o *operator) WorkerListener(workerChIn chan<- *entities.Worker) (err error) {
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	workersChIn = workerChIn

	http.HandleFunc("/worker/registration", workersHandler)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}
func workersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	host := r.PostForm.Get("host")
	port := r.PostForm.Get("port")

	worker := &entities.Worker{
		Host: host,
		Port: port,
	}

	log.Printf("Синициирована регистрация воркера %+v\n", worker)

	workersChIn <- worker
}

// BrokersListener обрабатывает запросы брокеров на регистрацию
func (o *operator) BrokersListener(handleFunc http.HandlerFunc) {
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	http.HandleFunc("/broker/registration", handleFunc)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}

// WorkersListener обрабатывает запросы воркеров на регистрацию
func (o *operator) WorkersListener(handleFunc http.HandlerFunc) {
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	http.HandleFunc("/worker/registration", handleFunc)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}

// SolutionsListener обрабатывает запросы воркеров с решением
func (o *operator) SolutionsListener(handleFunc http.HandlerFunc) {
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	http.HandleFunc("/worker/solution", handleFunc)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}
