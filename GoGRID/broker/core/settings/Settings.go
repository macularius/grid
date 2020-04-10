package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	BrokerHost string `ini:"brokerHost"` // хост брокера
	BrokerPort string `ini:"brokerPort"` // порт брокера

	Bookpath         string // путь к книге
	Substr           string // подстрока
	WorkCodeFilePath string // путь к файлу исполняемого кода
}
