package main

import (
	"grid/GoGRID/broker/core"
	"grid/GoGRID/broker/core/operator"
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

		taskCount int                 // оптимальное количество задач
		chTasks   chan core.PieOfBook // канал задач
		chAnswers chan string         // канал ответов

		configurator *settings.ApplicationConfigurator // экземпляр конфигуратора
		appOperator  *operator.Operator                // экземпляр оператора

		i int

		err error
	)

	// инициализация конфигуратора
	configurator = new(settings.ApplicationConfigurator)
	err = configurator.Init()
	if err != nil {
		log.Printf("error main.main : configurator.Init, %v\n", err)
		panic(err)
	}

	// инициализация оператора
	appOperator = new(operator.Operator)

	// получение книги и искомой строки
	book, substr, err = GetData()
	if err != nil {
		log.Printf("error main.main : GetData, %v\n", err)
		panic(err)
	}

	// получение оптимального количества задач
	taskCount = optimizer.Optimize3()

	// регистрируем оператор
	err = appOperator.Init(taskCount)
	if err != nil {
		log.Printf("error main.main : appOperator.Init, %v\n", err)
		panic(err)
	}

	// формирование задач
	chTasks = core.BrokeSync1(book, substr, taskCount)

	// получение результата
	chAnswers = make(chan string)
	go appOperator.Listener(chAnswers)

	// отправка задач в дистрибутор
	go func() {
		for {
			task := <-chTasks

			// задача в распределитель
			appOperator.SendTask(task)
		}
	}()

	log.Println("Ожидание результата работы...")

	// обработка результата
	result := 0
	for i = 0; i < taskCount; i++ {
		answer := <-chAnswers

		if answer == "finish" && i != taskCount-1 {
			log.Printf("Ошибка подсчета результата: обработано %v/%v решений. Результат = %v.\n", i, taskCount, result)
			return
		}

		if answer != "finish" {
			log.Printf("Получено решение строка %v, длина %v.\n", answer, len(answer))
		}

	}

	log.Printf("!!!CONGRATULATIONS!!!\n\tЗадач решено[%v/%v].\n", i, taskCount, result)
}

// GetData возвращает строку книги и искомой подстроки
func GetData() (book, substr string, err error) {
	var (
		f *os.File
		b []byte
	)

	//substr = settings.Config.Substr

	f, err = os.Open(settings.Config.Bookpath)
	if err != nil {
		log.Printf("error main.GetData : os.Open, %v\n", err)
		return
	}
	defer f.Close()

	b, err = ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error main.GetData : ioutil.ReadAll, %v\n", err)
		return
	}
	book = string(b)

	f, err = os.Open(settings.Config.Substr)
	if err != nil {
		log.Printf("error main.GetData : os.Open, %v\n", err)
		return
	}
	defer f.Close()

	b, err = ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error main.GetData : ioutil.ReadAll, %v\n", err)
		return
	}
	substr = string(b)

	return
}
