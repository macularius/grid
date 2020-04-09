package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	f, _ := os.Open("book1.txt")
	b, _ := ioutil.ReadAll(f)
	book := string(string(b))
	// book := "Мама_мыла_рамуМама_мыла_раму"

	substr := `Болконский`
	// substr := `Билибин был человек лет тридцати пяти, холостой, одного общества с князем Андреем. Они были знакомы еще в Петербурге, но еще ближе познакомились в последний приезд князя Андрея в Вену вместе с Кутузовым. Как князь Андрей был молодой человек, обещающий пойти далеко на военном поприще, так, и еще более, обещал Билибин на дипломатическом. Он был еще молодой человек, но уже немолодой дипломат, так как он начал служить с шестнадцати лет, был в Париже, в Копенгагене и теперь в Вене занимал довольно значительное место. И канцлер и наш посланник в Вене знали его и дорожили им. Он был не из того большого количества дипломатов, которые обязаны иметь только отрицательные достоинства, не делать известных вещей и говорить по-французски для того, чтобы быть очень хорошими дипломатами; он был один из тех дипломатов, которые любят и умеют работать, и, несмотря на свою лень, он иногда проводил ночи за письменным столом. Он работал одинаково хорошо, в чем бы ни состояла сущность работы. Его интересовал не вопрос "зачем?", а вопрос "как?". В чем состояло дипломатическое дело, ему было всё равно; но составить искусно, метко и изящно циркуляр, меморандум или донесение - в этом он находил большое удовольствие. Заслуги Билибина ценились, кроме письменных работ, еще и по его искусству обращаться и говорить в высших сферах.`
	// substr := `фывфывфывфывфывфывфывфывфывфывфaasdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddsdadddddddddddddddывфывфы`

	// i := 1
	// li := 1
	// for i < len(book)-len(substr)+1 {
	// 	fmt.Printf("Кол-во задач %v, li %v\n", i, li)
	// 	fmt.Printf(" Отношение задачи к кинге %f\n", float64(len(substr)+((len(book)-len(substr))/i))/float64((len(book))))

	// 	fmt.Printf("\tФормирование\n")
	// 	tasks := CreateTasks1(book, substr, i)

	// 	fmt.Printf("\tТаски\n")
	// 	handle(tasks, substr)

	// 	t := i
	// 	i += li
	// 	li = t
	// }

	fmt.Println(Optimise(book, substr))
}

// CreateTasks1 ##@@
func CreateTasks1(book string, entryStr string, n int) (tasks []string) {
	var (
		x      int // сдвиг
		piecln int // отрывок книги(подзадача)
		tStart time.Time
	)
	tStart = time.Now()
	defer func() {
		fmt.Printf("\tВремя формирования: %v\n", time.Now().Sub(tStart))
	}()

	if n == 1 {
		tasks = append(tasks, book)
		return
	}

	entryStrLength := len(entryStr)
	bookLength := len(book)

	x = (bookLength - entryStrLength) / n
	piecln = entryStrLength + x

	i := 0
	for ; i < bookLength-piecln; i += 1 + x {
		// if i+1+x >= bookLength-piecln {
		// 	tasks = append(tasks, book[i:])
		// 	break
		// }

		tasks = append(tasks, book[i:i+piecln])
	}
	if i < bookLength {
		tasks = append(tasks, book[i:])
	}

	return
}

func handle(tasks []string, entryStr string) {
	var (
		n       int
		tStart  time.Time
		average time.Duration
	)

	tStart = time.Now()
	defer func() {
		fmt.Printf("\tСреднее время выполнения одной задачи: %v\nКоличество вхождений: %v\n", average, n)
	}()

	for _, task := range tasks {
		tStart = time.Now()

		str := strings.ToLower(task)
		substr := strings.ToLower(entryStr)

		for i := 0; i < len(str)-len(substr); i++ {
			strPiece := str[i : i+len(substr)]
			if strPiece == substr {
				n++
			}
		}

		average = (average + time.Now().Sub(tStart)) / 2
	}
}

func Optimise(book, substr string) (taskCount int) {
	var (
		bookln           = len(book)
		substrln         = len(substr)
		i        float32 = 0.0015 // нижняя граница оптимальной длины задачи

		x int
	)

	// когда длина подстроки выходит за рамки оптимальной длины задачи, брать за длину задачи - длину подстроки
	if float32(substrln)/float32(bookln) < 0.000152 {
		i = 0.000152
	} else {
		i = float32(substrln) / float32(bookln)
	}

	for x <= 0 {
		x = int(float32(bookln)*i) - substrln // сдвиг
		if x < 1 {
			i += 0.0001
			continue
		}

		taskCount = (bookln - substrln) / x

		fmt.Printf("Task count %v, i = %v\n", taskCount, i)
		fmt.Printf("float32(bookln)*%v = %v\n", i, int(float32(bookln)*i))
		fmt.Printf("x = %v\n", x)
		fmt.Printf("substrln = %v\n\n", substrln)

		i += 0.0001
	}

	return
}
