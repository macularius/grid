package settings

// Settings конфигурация
type Settings struct {
	DistributorHost string `ini:"distributorHost"` // хост распределителя
	DistributorPort string `ini:"distributorPort"` // порт распределителя

	Bookpath string // путь к книге
	Substr   string // подстрока
}
