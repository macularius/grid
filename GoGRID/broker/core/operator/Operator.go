package operator

import (
	"encoding/json"
	"fmt"
	"grid/GoGRID/broker/core/settings"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

var answerChIn chan<- string // канал ответов

// Operator получатель и отправитель сообщений
type Operator struct {
	token        []byte // идентификатор задачи в дистрибуторе
	workCodeFile []byte
}

// Init регистрирует задачу в дистрибуторе
func (o *Operator) Init(taskCount int) (err error) {
	var (
		dHost = settings.Config.DistributorHost // хост дистрибутора
		dPort = settings.Config.DistributorPort // порт дистрибутора

		bHost = settings.Config.BrokerHost // хост дистрибутора
		bPort = settings.Config.BrokerPort // порт дистрибутора

		resp *http.Response
	)

	// получение файла исполняемого кода
	o.workCodeFile, err = ioutil.ReadFile(settings.Config.WorkCodeFilePath)
	if err != nil {
		log.Printf("error Operator.SendTask : ioutil.ReadFile, %v\n", err)
		return
	}

	// формирование запроса
	resp, err = http.PostForm("http://"+net.JoinHostPort(dHost, dPort)+"/broker/registration", url.Values{"task_count": {strconv.Itoa(taskCount)}, "host": {bHost}, "port": {bPort}})
	if err != nil {
		log.Printf("error Operator.SendTask : http.PostForm, %v\n", err)
		return
	}
	defer resp.Body.Close()

	// log.Printf("Получен ответ на запрос регистрации брокера в дистрибуторе.\nresp: %+v\n\n", resp)

	o.token, _ = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error Operator.SendTask : ioutil.ReadAll, %v\n", err)
		return
	}
	if len(o.token) == 0 {
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

		resp *http.Response // ответ
	)

	taskBody := TaskFile{
		Str:    task,
		Substr: settings.Config.Substr,
	}
	taskb, err := json.Marshal(taskBody)
	if err != nil {
		log.Printf("error Operator.SendTask : json.Marshal, %v\n", err)
		return
	}

	// отправка сообщения
	resp, err = http.PostForm("http://"+net.JoinHostPort(dHost, dPort)+"/broker/task", url.Values{"token": {string(o.token)}, "task_body": {string(taskb)}, "code_file": {string(o.workCodeFile)}})
	if err != nil {
		log.Printf("error Operator.SendTask : http.PostForm, %v\n", err)
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

	r.ParseForm()

	if res := r.PostForm.Get("finish_sign"); res != "" && res == "finish" {
		answerChIn <- "finish"
		return
	}
}

// TaskFile тип файла задачи
type TaskFile struct {
	Str    string `json:"str"`
	Substr string `json:"substr"`
}
