package settings

import (
	"log"
	"os"

	"github.com/go-ini/ini"
)

// Config экземпляр конфигурации
var Config Settings

// ApplicationConfigurator Структура для работы с конфигурацией
type ApplicationConfigurator struct {
}

// Init инициализацияконфигурации
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

	return
}
