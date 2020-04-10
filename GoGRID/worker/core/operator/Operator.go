package operator

import (
	"fmt"
	"grid/GoGRID/worker/core/settings"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
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

		wHost = settings.Config.WorkerHost // хост дистрибутора
		wPort = settings.Config.WorkerPort // порт дистрибутора

		resp *http.Response
	)
	// отправка сообщения
	vals := url.Values{}
	vals.Set("host", wHost)
	vals.Set("port", wPort)
	resp, err = http.PostForm("http://"+net.JoinHostPort(dHost, dPort)+"/worker/registration", vals)
	if err != nil {
		log.Printf("error Operator.Init : http.PostForm, %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Статус не OK")
		return
	}

	log.Printf("Отправлен запрос на регистрацию. Worker[%s:%s]\n", wHost, wPort)

	return
}

// Listener ожидает ответы
func (o *Operator) Listener() (err error) {

	var (
		wHost = settings.Config.WorkerHost
		wPort = settings.Config.WorkerPort
	)

	//прослушивание rest
	http.HandleFunc("/distributor/task", solution)
	log.Fatal(http.ListenAndServe(net.JoinHostPort(wHost, wPort), nil))

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
	defer r.Body.Close()

	r.ParseForm()

	//вытащить из запроса параметры
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

	//создать исполняемый файл и запустить его передав параметры:-token, -task_id
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
