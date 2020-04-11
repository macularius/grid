package optimizer

import "fmt"

// Optimize возвращает оптимальное количество задач для указанных книги и подстроки
func Optimize(book, substr string) (taskCount int) {
	var (
		bookln   = len(book)
		substrln = len(substr)
		i        float32 // нижняя граница оптимальной длины задачи

		x int
	)

	// если книга сильно больше подстроки, то выставляем нижнюю границу оптимально промежутка длины задачи
	// иначе отношение подстроки к книге
	if float32(substrln)/float32(bookln) < 0.00152 {
		i = 0.00152
	} else {
		i = float32(substrln) / float32(bookln)
	}

	for x <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		if x < 1 {
			i += 0.0005
			continue
		}

		taskCount = (bookln - substrln) / x

		fmt.Printf("Task count %v, i = %v\n", taskCount, i)

		i += 0.0005
	}

	return
}

func Optimize3() (taskCount int) {
	return 200
}

// Optimize2 возвращает оптимальное количество задач для указанных книги и подстроки
func Optimize2(book, substr string) (taskCount int) {
	var (
		bookln   = len(book)
		substrln = len(substr)
		i        float32 // нижняя граница оптимальной длины задачи

		x int
	)

	// если книга сильно больше подстроки, то выставляем нижнюю границу оптимально промежутка длины задачи
	// иначе отношение подстроки к книге
	if float32(substrln)/float32(bookln) < 0.00452 {
		i = 0.00452
	} else {
		i = float32(substrln) / float32(bookln)
	}

	for x <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		if x < 1 {
			i += 0.0005
			continue
		}

		taskCount = (bookln - substrln) / x

		fmt.Printf("Task count %v, i = %v\n", taskCount, i)

		i += 0.0005
	}

	return
}
