package operator

import (
	"fmt"
	"grid/GoGRID/broker/core/settings"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// Operator получатель и отправитель сообщений
type Operator struct {
	token []byte // идентификатор задачи в дистрибуторе
}

// Init регистрирует задачу в дистрибуторе
func (o *Operator) Init() (err error) {
	var (
		dHost = settings.Config.DistributorHost // хост дистрибутора
		dPort = settings.Config.DistributorPort // порт дистрибутора

		resp *http.Response
		i    int
	)

	resp, err = http.Get(net.JoinHostPort(dHost, dPort) + "/registration")
	if err != nil {
		log.Printf("error Operator.SendTask : http.Get, %v\n", err)
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
	req, err = http.NewRequest(http.MethodPost, net.JoinHostPort(dHost, dPort)+"/task", nil)
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
func (o *Operator) Listener(chAnswer chan<- int) (err error) {

	return
}
