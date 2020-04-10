package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	BrokerHost string `ini:"brokerHost"` // хост брокера
	BrokerPort string `ini:"brokerPort"` // порт брокера

	Bookpath         string `ini:"bookpath"` // путь к книге
	Substr           string `ini:"substr"` // подстрока
	WorkCodeFilePath string `ini:"workCodeFilePath"` // путь к файлу исполняемого кода
}
