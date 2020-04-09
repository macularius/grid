package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	WorkerPort string // Порт который слушает воркер
}
