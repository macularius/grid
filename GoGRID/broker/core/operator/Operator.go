package operator

import (
	"fmt"
	"grid/GoGRID/broker/core/settings"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

var answerChIn chan<- string // канал ответов

// Operator получатель и отправитель сообщений
type Operator struct {
	token []byte // идентификатор задачи в дистрибуторе
}

// Init регистрирует задачу в дистрибуторе
func (o *Operator) Init(taskCount int) (err error) {
	var (
		dHost = settings.Config.DistributorHost // хост дистрибутора
		dPort = settings.Config.DistributorPort // порт дистрибутора

		bHost = settings.Config.BrokerHost // хост дистрибутора
		bPort = settings.Config.BrokerPort // порт дистрибутора

		req  *http.Request
		resp *http.Response
		i    int
	)

	// формирование запроса
	req, err = http.NewRequest(http.MethodPost, "http://"+net.JoinHostPort(dHost, dPort)+"/broker/registration", nil)
	if err != nil {
		log.Printf("error Operator.SendTask : http.NewRequest, %v\n", err)
		return
	}
	req.PostForm.Add("task_count", strconv.Itoa(taskCount)) // токен задачи
	req.PostForm.Add("host", bHost)
	req.PostForm.Add("port", bPort)

	// формирование соединения
	client := &http.Client{}

	// отправка сообщения
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("error Operator.SendTask : client.Do, %v\n", err)
		return
	}

	i, err = resp.Body.Read(o.token)
	if err != nil {
		log.Printf("error Operator.SendTask : resp.Body.Read, %v\n", err)
		return
	}
	if i == 0 {
		err = fmt.Errorf("Длина токена равна нулю")
		return
	}

	return
}

// SendTask отправляет задачу в распределитель
func (o *Operator) SendTask(task string) (err error) {
	var (
		dHost = settings.Config.DistributorHost // хост дистрибутора
		dPort = settings.Config.DistributorPort // порт дистрибутора

		req  *http.Request  // запрос
		resp *http.Response // ответ

		b []byte // тело файла исполняемого кода
	)

	// получение файла исполняемого кода
	b, err = ioutil.ReadFile(settings.Config.WorkCodeFilePath)
	if err != nil {
		log.Printf("error Operator.SendTask : ioutil.ReadFile, %v\n", err)
		return
	}

	// формирование запроса
	req, err = http.NewRequest(http.MethodPost, "http://"+net.JoinHostPort(dHost, dPort)+"/broker/task", nil)
	if err != nil {
		log.Printf("error Operator.SendTask : http.NewRequest, %v\n", err)
		return
	}

	req.Form.Add("token", string(o.token))   // токен задачи
	req.PostForm.Add("task_body", task)      // тело задачи
	req.PostForm.Add("code_file", string(b)) // файл рабочего кода

	// формирование соединения
	client := &http.Client{}

	// отправка сообщения
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("error Operator.SendTask : client.Do, %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Статус не OK")
		return
	}

	return
}

// Listener ожидает ответы
func (o *Operator) Listener(answChIn chan<- string) (err error) {
	var (
		bHost = settings.Config.BrokerHost
		bPort = settings.Config.BrokerPort
	)

	answerChIn = answChIn

	http.HandleFunc("/distributor/solution", solution)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(bHost, bPort), nil))

	return
}
func solution(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error operator.solution : ioutil.ReadAll, %v\n", err)
		return
	}

	if len(b) > 0 {
		answer := string(b)
		answerChIn <- answer
	}

	if res := r.PostForm.Get("finish_sign"); res != "" && res == "finish" {
		answerChIn <- "finish"
		return
	}
}
