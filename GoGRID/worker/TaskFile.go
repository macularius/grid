package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// получение аргументов консольной строки: путь к json файлу задачи, адрес получателя результата
	var (
		taskID string
		token  string
		url    string
	)

	// внесение агрументов в объект конфига
	flag.StringVar(&taskID, "id", "", "Идентификатор задачи")
	flag.StringVar(&token, "token", "", "Идентификатор глобальной задачи")
	flag.StringVar(&url, "URL", "", "Адрес для ответа")

	flag.Parse()

	if taskID == "" {
		fmt.Println("taskID is nil")
		return
	}
	if token == "" {
		fmt.Println("token is nil")
		return
	}
	if url == "" {
		fmt.Println("url is nil")
		return
	}

	fmt.Printf("taskID[%v], token[%v], url[%v]\n", taskID, token, url)
	// return

	// формирование объекта задачи
	task := new(task)
	b, err := ioutil.ReadFile("task.json")
	if err != nil {
		fmt.Printf("error ioutil.ReadFile, %v\n", err)
		return
	}
	json.Unmarshal(b, &task)
	if err != nil {
		fmt.Printf("error json.Unmarshal, %v\n", err)
		return
	}

	// выполнение задачи
	var (
		book1  string
		book2  string
		x      = 1
		xLast  = 1
		res    []string
		resStr = ""
	)

	book1 = task.Str
	book2 = task.Substr

	fmt.Println(len(book1))
	fmt.Println(len(book2))

	for i := 0; i < len(book1); i++ {
		fmt.Println("1")
		if string(book1[i]) == string(book2[i]) {
			x++
			resStr += string(book1[i])

			continue
		}

		if x > xLast {
			fmt.Println("2")
			res = []string{resStr}
			xLast = x
		} else {

			if x == xLast {

				res = append(res, resStr)
			}
		}
		resStr = ""
		x = 1

	}

	fmt.Printf("Результат: %v\n", res)

	// отправка результата
	var (
		resp *http.Response
		buf  *bytes.Buffer
	)

	// формирование запроса
	buf = new(bytes.Buffer)
	fmt.Fprint(buf, res)

	// отправка сообщения
	resp, err = http.Post(url+"?token="+token+"&task_id="+taskID, "text/html", buf)
	if err != nil {
		fmt.Printf("error http.Post, %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Не удалось отправить")
		fmt.Printf("error http.Post, %v\n", err)
		return
	}

	return
}

type task struct {
	Str    string `json:"str"`
	Substr string `json:"substr"`
}
