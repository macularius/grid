package main

import (
	"grid/GoGRID/broker/core"
	"grid/GoGRID/broker/core/settings"
	"grid/GoGRID/broker/optimizer"
	"io/ioutil"
	"log"
	"os"
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
	err = configurator.Init()
	if err != nil {
		log.Printf("error main.main : configurator.Init, %v\n", err)
		return
	}

	book, substr = GetData()

	taskCount = optimizer.Optimize2(12391238, 20)

	chTasks = core.BrokeAsync(book, substr, taskCount)
	chTasks = chTasks

<<<<<<< HEAD
	// for {
	// 	select {
	// 	case task := <-chTasks:
	// 		// задача в распределитель
	// 		task = task
	// 	}
	// }
=======
	for {
		select {
		case task := <-chTasks:
			// задача в распределитель
<<<<<<< HEAD

=======
			task = task
>>>>>>> Init on my worflow
		}
	}
>>>>>>> Init on my worflow
}

// GetData возвращает строку книги и искомой подстроки
func GetData(book, substr string) (err error) {
	var (
		f *os.File
		b []byte
	)

	substr = settings.Config.Substr

	f, err = os.Open(settings.Config.Bookpath)
	if err != nil {
		log.Printf("error main.GetData : os.Open, %v\n", err)
		return
	}

	b, err = ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error main.GetData : ioutil.ReadAll, %v\n", err)
		return
	}

	return
}
