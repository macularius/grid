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
	if float32(substrln)/float32(bookln) < 0.000152 {
		i = 0.000152
	} else {
		i = float32(substrln) / float32(bookln)
	}

	for x <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		if x < 1 {
			i += 0.00005
			continue
		}

		taskCount = (bookln - substrln) / x

		fmt.Printf("Task count %v, i = %v\n", taskCount, i)

		i += 0.00005
	}

	return
}

// Optimize2 возвращает оптимальное количество задач для указанных книги и подстроки
func Optimize2(bookln, substrln int) (taskCount int) {
	var (
		i float32 // нижняя граница оптимальной длины задачи

		x int
	)

	// если книга сильно больше подстроки, то выставляем нижнюю границу оптимально промежутка длины задачи
	// иначе отношение подстроки к книге
	if float32(substrln)/float32(bookln) < 0.000152 {
		i = 0.000152
	} else {
		i = float32(substrln) / float32(bookln)
	}

	for x <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		if x < 1 {
			i += 0.00005
			continue
		}

		taskCount = (bookln - substrln) / x

		fmt.Printf("Task count %v, i = %v\n", taskCount, i)
		// fmt.Printf("float32(bookln)*%v = %v\n", i, int(float32(bookln)*i))
		fmt.Printf("x = %v\n", x)
		// fmt.Printf("substrln = %v\n\n", substrln)

		i += 0.00005
	}

	return
}
