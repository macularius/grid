package solution_dispatcher

import (
	"grid/GoGRID/distributor/core/entities"
	"grid/GoGRID/distributor/core/operator"
	"grid/GoGRID/distributor/core/worker_dispatcher"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	instance *solutionDispatcher
)

// GetSolutionDispatcher возвращает экземпляр дисптчера решений
func GetSolutionDispatcher() ISolutionDispatcher {
	if instance == nil {
		instance = new(solutionDispatcher)
		instance.Init()
	}

	return instance
}

// ISolutionDispatcher интерфейс диспетчера решений - синглтон
/*
	обрабатывает входящие запросы на регистрацию брокеров(решений)
	ведет массив решений
	ведет учет задач
	отправляет задачи в воркеры
*/
type ISolutionDispatcher interface {
	Init() error                // инициализирует диспетчер решений
	Run()                       // запускает рабочий цикл диспетчера решений
	Resolve(*entities.Solution) // запускает решение
}

type solutionDispatcher struct {
	solutions        map[string]*entities.Solution
	workerDispatcher worker_dispatcher.IWorkerDispatcher
	appOperator      operator.IOperator

	newBrokersCh chan *entities.Broker
	newTasksCh   chan *entities.Task
}

// Init инициализирует диспетчер решений
func (d *solutionDispatcher) Init() (err error) {
	d.appOperator = operator.GetOperator()
	d.workerDispatcher = worker_dispatcher.GetWorkerDispatcher()
	d.solutions = make(map[string]*entities.Solution)

	d.newBrokersCh = make(chan *entities.Broker)
	d.newTasksCh = make(chan *entities.Task)

	d.appOperator.AttachListener(brokersListener, "/broker/registration")
	d.appOperator.AttachListener(solutionsListener, "/worker/solution")
	d.appOperator.AttachListener(tasksListener, "/broker/task")

	return
}

// Run запускает рабочий цикл диспетчера решений
func (d *solutionDispatcher) Run() {
	for {
		select {
		// слушать новые брокеры
		case broker := <-d.newBrokersCh:
			tasks := make(map[int]*entities.Task)
			taskQueue := make(map[int]time.Time)

			s := &entities.Solution{
				Broker:    broker,
				TaskQueue: taskQueue,
				Tasks:     tasks,
				Token:     broker.Token,
			}
			d.solutions[string(broker.Token)] = s

			go d.Resolve(s)
		// слушать новые задачи
		case task := <-d.newTasksCh:
			s := d.solutions[string(task.Token)]
			id := len(s.Tasks) + 1

			task.Result = "-1"
			s.Tasks[id] = task
		}

	}
}
func brokersListener(w http.ResponseWriter, r *http.Request) {
	// получение новых брокеров
	var (
		taskCountStr string
	)
	defer r.Body.Close()

	r.ParseForm()

	// формирование объекта брокера
	taskCountStr = r.PostForm.Get("task_count")

	broker := &entities.Broker{}

	broker.Host = r.PostForm.Get("host")
	broker.Port = r.PostForm.Get("port")
	broker.TaskCount, _ = strconv.Atoi(taskCountStr)

	// формирование токена
	broker.Token, _ = uuid.New().MarshalBinary()

	// запись токена в тело ответа
	w.Write(broker.Token)
}
func tasksListener(w http.ResponseWriter, r *http.Request) {
	// получение новых задач
	var (
		tokenStr    string
		bodyStr     string
		workcodeStr string
	)
	defer r.Body.Close()

	r.ParseForm()

	// формирование объекта задачи
	tokenStr = r.Form.Get("token")            // токен задачи
	bodyStr = r.PostForm.Get("task_body")     // тело задачи
	workcodeStr = r.PostForm.Get("code_file") // файл рабочего кода

	task := &entities.Task{}
	task.Token = []byte(tokenStr)
	task.Body = []byte(bodyStr)
	task.WorkCode = []byte(workcodeStr)
	task.Result = "-1"

	// присвоение id задаче
	task.ID = len(instance.solutions[tokenStr].Tasks) + 1

	// отправка задачи в канал
	instance.newTasksCh <- task

}
func solutionsListener(w http.ResponseWriter, r *http.Request) {
	// получение решений
	defer r.Body.Close()

	r.ParseForm()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error solution_dispatcher.solutionsListener : ioutil.ReadAll, %v\n", err)
		return
	}
	result := string(b)
	taskIDStr := r.PostForm.Get("task_id")
	taskID, _ := strconv.Atoi(taskIDStr)
	token := r.PostForm.Get("token")

	var (
		solution = instance.solutions[token]
	)

	// зафиксировать задачу решенной
	solution.Tasks[taskID].Result = result
	delete(solution.TaskQueue, taskID)
	solution.Broker.TaskCount--

	// отправить решение в брокер
	for {
		err = solution.Broker.Send(result)
		if err != nil {
			log.Printf("error solution_dispatcher.solutionsListener : solution.Broker.Send, %v\n", err)
			time.Sleep(time.Second * 10)
			continue
		}

		break
	}

	return
}

// Resolve запускает решение
func (d *solutionDispatcher) Resolve(s *entities.Solution) {
	for s.Broker.TaskCount > 0 {
		// получить задачу без решения
		for _, task := range s.Tasks {
			if _, ok := s.TaskQueue[task.ID]; !ok && task.Result == "-1" {
				// отправить задачу в воркер
				err := d.workerDispatcher.SendTask(task, s.Token)
				if err != nil {
					log.Printf("error operator.Resolve : ioutil.ReadAll, %v\n", err)
					continue
				}

				// TODO вынести delay в конфиг
				// поместить задачу в список отправленных задач
				s.TaskQueue[task.ID] = time.Now().Add(time.Second * 60)
			}
		}

		// очистить список от просроченных задач
		for id, t := range s.TaskQueue {
			if time.Now().After(t) {
				delete(s.TaskQueue, id)
			}
		}
	}
}
