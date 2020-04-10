package main

import (
	"grid/GoGRID/distributor/core/settings"
	"grid/GoGRID/distributor/core/solution_dispatcher"
	"grid/GoGRID/distributor/core/worker_dispatcher"
	"log"
)

func main() {
	var (
		configurator        *settings.ApplicationConfigurator       // экземпляр конфигуратора
		dispatcherWorkers   worker_dispatcher.IWorkerDispatcher     // экземпляр диспетчера воркеров
		dispatcherSolutions solution_dispatcher.ISolutionDispatcher // экземпляр диспетчера решений

		err error
	)

	// инициализация конфигуратора
	configurator = new(settings.ApplicationConfigurator)
	err = configurator.Init()
	if err != nil {
		log.Printf("error main.main : configurator.Init, %v\n", err)
		panic(err)
	}

	// инициализация диспетчера воркеров
	dispatcherWorkers = worker_dispatcher.GetWorkerDispatcher()
	go dispatcherWorkers.Run()

	// инициализация диспетчера решений
	dispatcherSolutions = solution_dispatcher.GetSolutionDispatcher()
	go dispatcherSolutions.Run()

	for {
	}
}
