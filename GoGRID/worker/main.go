package main

import (
	"grid/GoGRID/worker/core/operator"
	"grid/GoGRID/worker/core/settings"
	"log"
)

func main() {

	var (
		configurator *settings.ApplicationConfigurator

		err error
	)

	// инит конфигурации
	configurator = new(settings.ApplicationConfigurator)
	err = configurator.Init()
	if err != nil {
		log.Printf("error main.main : configurator.Init, %v\n", err)
		return
	}

	appOperator := new(operator.Operator)
	appOperator.Init()
	err = appOperator.Listener()

	// запуск слушателя порта

	/* при получении задачи и файла приложения:
	создать рабочую директорию
	поместить туда файл исполняемого кода и файл приложения(json)
	запустить исполняемый файл, при этом в аргументы передаем хост и порт дистрибутора и путь к файлу приложению
	*/
}
