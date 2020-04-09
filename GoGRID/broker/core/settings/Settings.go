package settings

// Settings конфигурация
type Settings struct {
	distributorHost string `ini:"distributorHost"` // хост распределителя
	distributorPort string `ini:"distributorPort"` // порт распределителя

	bookpath string // путь к книге
	substr   string // подстрока
}
