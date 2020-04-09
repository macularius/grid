package main

import (
	"fmt"
	"grid/GoGRID/broker/core"
	"grid/GoGRID/broker/core/settings"
	"grid/GoGRID/broker/optimizer"
)

func main() {
	var (
		book   string
		substr string

		taskCount int

		chTasks chan string

		configurator *settings.ApplicationConfigurator

		err error
	)

	configurator = new(settings.ApplicationConfigurator)
	err = configurator.ReadConfig()
	if err != nil {
		fmt.Printf("error main.main : configurator.ReadConfig, %v\n", err)
		return
	}

	book, substr = GetData()

	taskCount = optimizer.Optimize2(12391238, 20)

	chTasks = core.BrokeAsync(book, substr, taskCount)

	for {
		select {
		case task := <-chTasks:
			// задача в распределитель

		}
	}
}

// GetData возвращает строку книги и искомой подстроки
func GetData() (book, substr string) {

	return
}
