package optimizer

import "fmt"

// Optimize возвращает оптимальное количество задач для указанных книги и подстроки
func Optimize(book, substr string) (taskCount int) {
	var (
		bookln           = len(book)
		substrln         = len(substr)
		i        float32 = 0.0015 // нижняя граница оптимальной длины задачи

		x int
	)

	// когда длина подстроки выходит за рамки оптимальной длины задачи, брать за длину задачи - длину подстроки
	if float32(substrln)/float32(bookln) > 0.001017 {
		fmt.Println("Длина подстроки выходит за рамки оптимальной длины задачи. float32(substrln) / float32(bookln) = ", float32(substrln)/float32(bookln))

		taskCount = 1
		return
	}

	for taskCount <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		taskCount = (bookln - substrln) / x

		i += 0.00001
	}

	return
}
