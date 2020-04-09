package main

import (
	"fmt"
	"grid/GoGRID/broker/core"
	"grid/GoGRID/broker/core/settings"
	"grid/GoGRID/broker/optimizer"
)

func main() {
	var (
		book   string // книга, получать по пути из аргумента консоли
		substr string // подстрока, получать из аргумента консоли

		taskCount int         // оптимальное количество задач
		chTasks   chan string // канал задач

		configurator *settings.ApplicationConfigurator // экземпляр конфигуратора

		err error
	)

	configurator = new(settings.ApplicationConfigurator)
	err = configurator.ReadConfig()
	if err != nil {
		fmt.Printf("error main.main : configurator.ReadConfig, %v\n", err)
		return
	}
	err = configurator.GetArgs()
	if err != nil {
		fmt.Printf("error main.main : configurator.GetArgs, %v\n", err)
		return
	}

	book, substr = GetData()

	taskCount = optimizer.Optimize2(12391238, 20)

	chTasks = core.BrokeAsync(book, substr, taskCount)
	chTasks = chTasks

	// for {
	// 	select {
	// 	case task := <-chTasks:
	// 		// задача в распределитель
	// 		task = task
	// 	}
	// }
}

// GetData возвращает строку книги и искомой подстроки
func GetData() (book, substr string) {

	return
}
