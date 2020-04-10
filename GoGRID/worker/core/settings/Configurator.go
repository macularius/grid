package settings

import (
	"flag"
	"log"
	"os"

	"github.com/go-ini/ini"
)

// Config экземпляр конфигурации
var Config Settings

// ApplicationConfigurator Структура для работы с конфигурацией
type ApplicationConfigurator struct {
}

//Init для инициализации
func (c *ApplicationConfigurator) Init() (err error) {
	c.ReadConfig()
	if err != nil {
		log.Printf("error ApplicationConfigurator.Init : configurator.Init, %v\n", err)
		return
	}

	c.GetArgs()
	if err != nil {
		log.Printf("error ApplicationConfigurator.Init : configurator.Init, %v\n", err)
		return
	}

	return
}

// ReadConfig Зачитывает конфигурацию из ini файла /../conf/conf.ini
func (c *ApplicationConfigurator) ReadConfig() (err error) {
	cfg := new(ini.File)

	cfg, err = ini.Load("conf.ini")
	if err != nil {
		log.Printf("Fail to read file: %v\n", err)
		os.Exit(1)
	}

	err = cfg.MapTo(&Config)
	if err != nil {
		log.Printf("Fail to map file: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Readed config: %+v\n", Config)

	return
}

// GetArgs обработка аргументов консольной строки
func (c *ApplicationConfigurator) GetArgs() (err error) {
	var (
		workerHost      string
		workerPort      string
		distributorHost string
		distributorPort string
	)

	// внесение агрументов в объект конфига
	flag.StringVar(&workerHost, "workerhost", "", "Хост, который слушает воркер")
	flag.StringVar(&workerPort, "workerport", "", "Порт, который слушает воркер")
	flag.StringVar(&distributorHost, "distributorHost", "", "хост распределителя")
	flag.StringVar(&distributorPort, "distributorPort", "", "Порт распределителя")

	// after all flag definitions you must call
	flag.Parse()

	// then we can access our values
	log.Printf("Value of workerhost is: %s\n", workerHost)
	log.Printf("Value of workerport is: %s\n", workerPort)
	log.Printf("Value of distributorHost is: %s\n", distributorHost)
	log.Printf("Value of distributorPort is: %s\n", distributorPort)

	if workerHost != "" {
		Config.WorkerHost = workerHost
	}
	if workerPort != "" {
		Config.WorkerPort = workerPort
	}
	if distributorHost != "" {
		Config.DistributorHost = distributorHost
	}
	if distributorPort != "" {
		Config.DistributorPort = distributorPort
	}
	return
}
