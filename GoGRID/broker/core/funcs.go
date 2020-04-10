package core

import "log"

// BrokeAsync возвращает канал, в который по мере разбиения будут поступать задачи
func BrokeAsync(book, substr string, taskCount int) (chTasks chan string) {
	chTasks = make(chan string)

	go func(chTasks chan string, book, substr string, taskCount int) {
		log.Printf("Начато асинхронное формирование задач\n\tкнига[%v], substr[%v], taskCount[%v]\n", len(book), len(substr), taskCount)

		var (
			x      int // сдвиг
			piecln int // отрывок книги(задача)
			i      int
		)

		if taskCount == 1 {
			chTasks <- book
			return
		}

		substrln := len(substr)
		bookln := len(book)

		x = (bookln - substrln) / taskCount
		piecln = substrln + x

		for ; i < bookln-piecln; i += 1 + x {
			chTasks <- book[i : i+piecln]
		}

		// если остаток книги после задачи меньше задачи, то включаем ее в последнюю задачу
		if i < bookln {
			chTasks <- book[i:]
		}
	}(chTasks, book, substr, taskCount)

	return
}

// BrokeSync возвращает массив задач
func BrokeSync(book, substr string, taskCount int) (chTasks chan string) {
	chTasks = make(chan string, taskCount)

	var (
		x      int // сдвиг
		piecln int // отрывок книги(задача)
		i      int
	)

	if taskCount == 1 {
		chTasks <- book
		return
	}

	substrln := len(substr)
	bookln := len(book)

	x = (bookln - substrln) / taskCount
	piecln = substrln + x

	for ; i < bookln-piecln; i += 1 + x {
		chTasks <- book[i : i+piecln]
	}
	// если остаток книги после задачи меньше задачи, то включаем ее в последнюю задачу
	if i < bookln {
		chTasks <- book[i:]
	}

	return
}

// AnswerCount принимает канал ответов и количество задач. Суммирует решени
func AnswerCount(chAnswer <-chan int, taskCount int) (result int) {
	for i := 0; i < taskCount; i++ {
		result += <-chAnswer
	}

	return
}
