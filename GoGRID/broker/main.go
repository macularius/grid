package main

import (
	"grid/GoGRID/broker/core"
	"grid/GoGRID/broker/core/operator"
	"grid/GoGRID/broker/core/settings"
	"grid/GoGRID/broker/optimizer"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	var (
		book   string // книга, получать по пути из аргумента консоли
		substr string // подстрока, получать из аргумента консоли

		taskCount int         // оптимальное количество задач
		chTasks   chan string // канал задач
		chAnswers chan string // канал ответов

		configurator *settings.ApplicationConfigurator // экземпляр конфигуратора
		appOperator  *operator.Operator                // экземпляр оператора

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
	taskCount = optimizer.Optimize(book, substr)

	// регистрируем оператор
	err = appOperator.Init(taskCount)
	if err != nil {
		log.Printf("error main.main : appOperator.Init, %v\n", err)
		panic(err)
	}

	// формирование задач
	chTasks = core.BrokeAsync(book, substr, taskCount)

	// получение результата
	chAnswers = make(chan string)
	err = appOperator.Listener(chAnswers)
	if err != nil {
		log.Printf("error main.main : appOperator.Listener, %v\n", err)
		panic(err)
	}

	// отправка задач в дистрибутор
	go func() {
		for {
			select {
			case task := <-chTasks:
				log.Printf("Сформирована задача^ %+v\n", task)

				// задача в распределитель
				appOperator.SendTask(task)
			}
		}
	}()

	log.Println("Ожидание результата работы...")

	// обработка результата
	result := 0
	for i := 0; i < taskCount; i++ {
		var res int

		answer := <-chAnswers

		if answer == "finish" && i != taskCount-1 {
			log.Printf("Ошибка подсчета результата: обработано %v/%v решений. Результат = %v.\n", i, taskCount, result)
			return
		}

		res, err = strconv.Atoi(answer)
		if err != nil {
			log.Printf("error main.main : strconv.Atoi, %v\n", err)
			continue
		}

		result += res

		log.Printf("Получено решение %v. Текущий результат = %v.\n", res, result)
	}

	log.Printf("!!!CONGRATULATIONS!!!\n\tЗадача решена, результат = %v.\n", result)
}

// GetData возвращает строку книги и искомой подстроки
func GetData() (book, substr string, err error) {
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
	defer f.Close()

	b, err = ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error main.GetData : ioutil.ReadAll, %v\n", err)
		return
	}
	book = string(b)

	return
}
