package operator

import (
	"fmt"
	"grid/GoGRID/worker/core/settings"
	"log"
	"net"
	"net/http"
)

// Operator получатель и отправитель сообщений
type Operator struct {
}

//Init Регистрирует воркер в дистрибуторе
func (o *Operator) Init() (err error) {
	var (
		dHost = settings.Config.DistributorHost // хост дистрибутора
		dPort = settings.Config.DistributorPort // порт дистрибутора

		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest(http.MethodPost, net.JoinHostPort(dHost, dPort)+"/registration", nil)
	if err != nil {
		log.Printf("error Operator.Init : http.NewRequest, %v\n", err)
		return
	}

	// формирование соединения
	client := &http.Client{}

	// отправка сообщения
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("error Operator.Init : client.Do, %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Статус не OK")
		return
	}
	return
}

// Listener ожидает ответы
func (o *Operator) Listener() (err error) {

	var (
		dHost = settings.Config.DistributorHost
		dPort = settings.Config.DistributorPort
	)

	//прослушивание rest
	http.HandleFunc("/distributor/task", solution)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(dHost, dPort), nil))

	return
}

func solution(w http.ResponseWriter, r *http.Request) {

	// r.PostForm.Get("token", string(token))
	// r.PostForm.Get("task_id", strconv.Itoa(t.ID))
	// r.PostForm.Get("task_body", string(t.Body))
	// r.PostForm.Get("task_workcode", string(t.WorkCode))

	return
}
