package entities

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Solution тип решения
type Solution struct {
	Token     []byte            // идентификатор решения
	Broker    *Broker           // экземпляр брокера решения
	Tasks     map[int]*Task     // список задач
	TaskQueue map[int]time.Time // очередь задач на отправку и время, с которого считать задачу не решенной
}

// Task тип задачи
type Task struct {
	Token    []byte // идентификатор решения
	ID       int    // идентификатор задачи
	WorkCode []byte // исполняемый код
	Body     []byte // приложение к задаче
	Result   string // решение задачи
}

// Broker тип брокера
type Broker struct {
	Token     []byte // идентификатор решения
	Host      string
	Port      string
	TaskCount int
}

// Send отправляет результат в брокер
func (b *Broker) Send(res string) (err error) {
	var (
		req  *http.Request
		resp *http.Response
		buf  *bytes.Buffer
	)

	// формирование запроса
	fmt.Fprint(buf, res)
	req.Write(buf)
	if b.TaskCount == 0 {
		req.PostForm.Add("finish_sign", "finish")
	}

	// отправка сообщения
	resp, err = http.Post("http://"+net.JoinHostPort(b.Host, b.Port)+"/distributor/solution", "text/html", buf)
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
func (w *Worker) Send(t *Task, token []byte) (err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	// формирование запроса
	req, err = http.NewRequest(http.MethodPost, net.JoinHostPort(w.Host, w.Port)+"/distributor/task", nil)
	if err != nil {
		log.Printf("error Worker.Send : http.NewRequest, %v\n", err)
		return
	}
	req.PostForm.Add("token", string(token))
	req.PostForm.Add("task_id", strconv.Itoa(t.ID))
	req.PostForm.Add("task_body", string(t.Body))
	req.PostForm.Add("task_workcode", string(t.WorkCode))

	// формирование соединения
	client := &http.Client{}

	// отправка сообщения
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("error Worker.Send : client.Do, %v\n", err)
		return
	}

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
