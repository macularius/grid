package entities

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Solution тип решения
type Solution struct {
	Token     string            // идентификатор решения
	Broker    *Broker           // экземпляр брокера решения
	Tasks     map[int]*Task     // список задач
	TaskQueue map[int]time.Time // очередь задач на отправку и время, с которого считать задачу не решенной
}

// Task тип задачи
type Task struct {
	Token    string // идентификатор решения
	ID       int    // идентификатор задачи
	WorkCode []byte // исполняемый код
	Body     []byte // приложение к задаче
	Result   string // решение задачи
}

// Broker тип брокера
type Broker struct {
	Token     string // идентификатор решения
	Host      string
	Port      string
	TaskCount int
}

// Send отправляет результат в брокер
func (b *Broker) Send(res string) (err error) {
	var (
		resp *http.Response
		buf  *bytes.Buffer
		url  = "http://" + net.JoinHostPort(b.Host, b.Port) + "/distributor/solution"
	)

	// формирование запроса
	buf = new(bytes.Buffer)
	fmt.Fprint(buf, res)

	if b.TaskCount == 1 {
		url = url + "?finish_sign=finish"
	}

	// отправка сообщения
	resp, err = http.Post(url, "text/html", buf)
	if err != nil {
		log.Printf("error Worker.Send : http.PostForm, %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Не удалось отправить")
		return
	}

	return
}

// Worker тип воркера
type Worker struct {
	Host string
	Port string
}

// Send отправляет сообщение воркеру
func (w *Worker) Send(t *Task, token string) (err error) {
	var (
		resp *http.Response
	)

	// отправка сообщения
	resp, err = http.PostForm("http://"+net.JoinHostPort(w.Host, w.Port)+"/distributor/task", url.Values{"token": {token}, "task_id": {strconv.Itoa(t.ID)}, "task_body": {string(t.Body)}, "task_workcode": {string(t.WorkCode)}})
	if err != nil {
		log.Printf("error Worker.Send : http.PostForm, %v\n", err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Не удалось отправить")
		return
	}

	return
}

// Priority тип приоритета
type Priority uint8

// приоритеты воркеров, после 2ой неудачной попытки соединения с воркером, он удаляется
const (
	STABLE    Priority = 0 // стабильынй воркер
	UNSTABLE1 Priority = 1 // не ответил на 1 запрос
	UNSTABLE2 Priority = 2 // не ответил на 2 запроса
)
