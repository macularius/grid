package worker_dispatcher

import (
	"grid/GoGRID/distributor/core/entities"
	"grid/GoGRID/distributor/core/operator"
	"grid/GoGRID/distributor/core/settings"
	"log"
	"net/http"
	"time"
)

// экземпляр диспетчера воркеров
var instance *workerDispatcher

// GetWorkerDispatcher возвращает экземпляр диспетчера воркеров
func GetWorkerDispatcher() IWorkerDispatcher {
	if instance == nil {
		instance = new(workerDispatcher)
		instance.Init()
	}

	return instance
}

// IWorkerDispatcher интерфейс диспетчера воркеров - синглтон
/*
	ведет учет воркеров
	предоставляет интерфейс отправки задач
*/
type IWorkerDispatcher interface {
	Init() error                           // инициализирует диспетчер
	Run()                                  // запускает рабочий цикл диспетчера
	SendTask(*entities.Task, string) error // отправляет задачу воркеру
}

// workerDispatcher диспетчера воркеров
type workerDispatcher struct {
	workersRegister map[*entities.Worker]entities.Priority // реестр воркеров
	newWorkersCh    chan *entities.Worker                  // канал новых воркеров
	appOperator     operator.IOperator                     // оператор
}

// Init инициализирует диспетчер
func (d *workerDispatcher) Init() (err error) {
	d.appOperator = operator.GetOperator()
	d.newWorkersCh = make(chan *entities.Worker)
	d.workersRegister = make(map[*entities.Worker]entities.Priority)

	d.appOperator.AttachListener(workersHandler, "/worker/registration")

	return
}

// Run запускает рабочий цикл диспетчера
func (d *workerDispatcher) Run() {
	for {
	}
}

// SendTask отправляет задачу воркеру
func (d *workerDispatcher) SendTask(task *entities.Task, token string) (err error) {
	var (
		curPriority     = entities.STABLE                                                     // текущий искомый приоритет
		unstableAllFlag bool                                                                  // признак отсутствия стабильных источников
		dur             = time.Second * time.Duration(settings.Config.WaitNewWorkersDuration) // длительность ожидания
	)

	for {
		if unstableAllFlag {
			log.Printf("Нет рабочих воркеров следующая попытка в %v.\n", time.Now().Add(dur).Format("15:04:05"))
			time.Sleep(dur)
		}

		worker := d.getWorker(curPriority)
		if worker != nil {
			// попытка отправить сообщение
			err = worker.Send(task, token)
			if err != nil {
				d.workerNextPriority(worker, curPriority)
			} else {
				delete(d.workersRegister, worker)
				return
			}
		}

		// переходим на следующий приоритет
		switch curPriority {
		case entities.STABLE:
			curPriority = entities.UNSTABLE1
		case entities.UNSTABLE1:
			curPriority = entities.UNSTABLE2
		case entities.UNSTABLE2:
			unstableAllFlag = true
		}
	}
}
func workersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	r.ParseForm()

	host := r.PostForm.Get("host")
	port := r.PostForm.Get("port")

	worker := &entities.Worker{
		Host: host,
		Port: port,
	}

	log.Printf("Синициирована регистрация воркера %+v\n", worker)

	instance.workersRegister[worker] = entities.STABLE

	log.Printf("Зарегистрирован воркер %+v\n", worker)
}

// getWorker возвращает воркер
func (d *workerDispatcher) getWorker(curPriority entities.Priority) (worker *entities.Worker) {
	var priority entities.Priority

	for worker, priority = range d.workersRegister {
		if priority == curPriority {
			return
		}
	}

	return
}

// workerNextPriority устанавливает следующий приоритет
func (d *workerDispatcher) workerNextPriority(worker *entities.Worker, curPriority entities.Priority) {
	switch curPriority {
	case entities.STABLE:
		d.workersRegister[worker] = entities.UNSTABLE1
	case entities.UNSTABLE1:
		d.workersRegister[worker] = entities.UNSTABLE2
	case entities.UNSTABLE2:
		delete(d.workersRegister, worker)
	}
}
