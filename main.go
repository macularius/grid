package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	// book := "aaaa_aa_aa_aa_a"
	f, _ := os.Open("book.txt")
	b, _ := ioutil.ReadAll(f)

	book := string(b)
	substr := "Колобок"

	CreateTasks1(book, substr)

	fmt.Println(strings.Count(book, "Колобок"))
	fmt.Println(strings.Count(book, "колобок"))
}

// CreateTasks1 ##@@
func CreateTasks1(book string, entryStr string) {
	entryStrLength := len(entryStr)
	bookLength := len(book)

	fmt.Println("Начало построения задач")

	// смещение
	for i := 1; entryStrLength+i < bookLength/2; i++ {
		var tasks []string

		// обход книги - формирование задач
		for c := 0; c+i-1+entryStrLength <= bookLength-entryStrLength; c += i {
			tasks = append(tasks, book[c:c+i-1+entryStrLength])

			if c+i-1+entryStrLength <= bookLength-entryStrLength && bookLength%(i-1+entryStrLength) > 0 {
				// fmt.Println("Есть остаток")

				// все что отсалось после формирования полных задач
				tasks = append(tasks, book[bookLength-bookLength/(i-1+entryStrLength)*(i-1+entryStrLength):])
			}
		}

		fmt.Printf("\nСмещение: %v, booklen: %v, количество задач: %v\n", i-1, bookLength, len(tasks))
		// fmt.Printf("tasks: %v\n", tasks[:5])

		// решить задачи
		handle(tasks, entryStr)
	}
}

func handle(tasks []string, substr string) {
	var (
		n int

		tStart   time.Time     // время начала
		tFinish  time.Time     // время начала
		tAverage time.Duration // среднее время решения задачи
	)

	tStart = time.Now()

	for _, s := range tasks {

		// начало решения задачи
		t1 := time.Now()

		// решение
		for {
			i := strings.Index(s, substr)
			if i == -1 {
				break
			}
			n++
			s = s[i+1:]
		}

		// конец решения задачи
		t2 := time.Now()
		tAverage = (tAverage + t2.Sub(t1)) / 2
	}

	tFinish = time.Now()
	fmt.Printf("\tВремя решения всех задач: %v\n\tСреднее время выполнения: %v\n\tКоличество вхождений: %v\n", tFinish.Sub(tStart), tAverage, n)
}
