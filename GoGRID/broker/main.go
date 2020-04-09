package main

import (
	"fmt"
	"grid/GoGRID/broker/optimizer"
)

func main() {
	var (
		// book   string
		// substr string

		taskCount int

		// chTasks chan string
	)

	// book, substr = GetData()

	taskCount = optimizer.Optimize2(12391238, 20)

	fmt.Println(taskCount)
	return

	// chTasks = core.BrokeAsync(book, substr, taskCount)

	for {
		select {
		// case task := <-chTasks:
		// задача в распределитель

		}
	}
}

// GetData возвращает строку книги и искомой подстроки
func GetData() (book, substr string) {

	return
}
