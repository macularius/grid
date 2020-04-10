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
	Listen()
	AttachListener(handleFunc http.HandlerFunc, url string)
}

// Operator получатель и отправитель сообщений
type operator struct {
}

// Init инициализирует оператор
func (o *operator) Init() (err error) {
	return
}

// Listen ожидает новые воркеры
func (o *operator) Listen() {
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))
}

// AttachListener обрабатывает запросы воркеров на регистрацию
func (o *operator) AttachListener(handleFunc http.HandlerFunc, url string) {
	// получение задач
	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	http.HandleFunc(url, handleFunc)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}
