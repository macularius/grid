package operator

import (
	"fmt"
	"grid/GoGRID/worker/core/settings"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
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
	req, err = http.NewRequest(http.MethodPost, "http://"+net.JoinHostPort(dHost, dPort)+"/registration", nil)
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

	var (
		err           error
		token         string
		task_id       string
		task_body     string
		task_workcode string
		URL           = "http://" + net.JoinHostPort(settings.Config.DistributorHost, settings.Config.DistributorPort) + "/solution"
	)
	//вытащить из запроса параметры
	//создать исполняемый файл и запустить его передав параметры:-token, -task_id

	if token = r.PostForm.Get("token"); token == "" {
		err = fmt.Errorf("Токен не может быть пустым")
		log.Printf("error operator.solution : r.PostForm.Get, %v\n", err)
		return
	}

	if task_id = r.PostForm.Get("task_id"); task_id == "" {
		err = fmt.Errorf("Идентификационный номер подзадачи отсутствует")
		log.Printf("error operator.solution : r.PostForm.Get, %v\n", err)
		return
	}

	if task_body = r.PostForm.Get("task_body"); task_body == "" {
		err = fmt.Errorf("Отсутствует тело задачи")
		log.Printf("error operator.solution : r.PostForm.Get, %v\n", err)
		return
	}

	if task_workcode = r.PostForm.Get("task_workcode"); task_workcode == "" {
		err = fmt.Errorf("Отсутствует исполняемый код")
		log.Printf("error operator.solution : r.PostForm.Get, %v\n", err)
		return
	}

	ioutil.WriteFile("TaskFile.go", []byte(task_body), 0777)
	if err != nil {
		log.Printf("error Operator.solution : ioutil.WriteFile, %v\n", err)
		return
	}

	cmd := exec.Command("go", "run", "TaskFile.go", "-task_id", task_id, "-token", token, "-URL", URL)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))

	return
}
