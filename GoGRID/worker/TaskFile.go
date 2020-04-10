package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	flag.StringVar(&taskID, "task_id", "", "Идентификатор задачи")
	flag.StringVar(&token, "token", "", "Идентификатор глобальной задачи")
	flag.StringVar(&url, "URL", "", "Адрес для ответа")

	flag.Parse()

	if taskID != "" {
		panic(2)
	}
	if token != "" {
		panic(2)
	}
	if url != "" {
		panic(2)
	}

	// формирование объекта задачи
	task := new(task)
	b, err := ioutil.ReadFile("task.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &task)
	if err != nil {
		panic(err)
	}

	// выполнение задачи
	res := 0
	for i := 0; i < len(task.Str)-len(task.Substr); i++ {
		strPiece := task.Str[i : i+len(task.Substr)]
		if strPiece == task.Substr {
			res++
		}
	}

	// отправка результата
	var (
		req  *http.Request
		resp *http.Response
		buf  *bytes.Buffer
	)

	// формирование запроса
	buf = new(bytes.Buffer)
	fmt.Fprint(buf, res)
	req.Write(buf)

	// отправка сообщения
	resp, err = http.Post(url, "text/html", buf)
	if err != nil {
		log.Printf("error main.main : http.PostForm, %v\n", err)
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Не удалось отправить")
		return
	}

	return
}

type task struct {
	Str    string `json:"str"`
	Substr string `json:"substr"`
}
