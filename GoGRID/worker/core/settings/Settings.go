package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	WorkerPort string `ini:"workerPort"` // Порт воркера
	WorkerHost string `ini:"workerHost"` //Хост воркера
}
