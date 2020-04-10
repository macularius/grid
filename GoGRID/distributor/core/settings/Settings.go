package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	BrokerHost string `ini:"brokerHost"` // хост брокера
	BrokerPort string `ini:"brokerPort"` // порт брокера

	WaitNewWorkersDuration int `ini:"waitNewWorkersDuration"` // время ожидания новых воркеров, если отсутствуют рабочие воркеры(в секундах)
}
