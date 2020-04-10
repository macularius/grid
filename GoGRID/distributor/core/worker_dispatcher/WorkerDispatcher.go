package worker_dispatcher

import (
	"grid/GoGRID/distributor/core/entities"
	"grid/GoGRID/distributor/core/settings"
	"log"
	"time"
)

// экземпляр диспетчера воркеров
var instance IWorkerDispatcher

// GetWorkerDispatcher возвращает экземпляр диспетчера воркеров
func GetWorkerDispatcher() IWorkerDispatcher {
	if instance == nil {
		instance = new(workerDispatcher)
	}

	return instance
}

// IWorkerDispatcher интерфейс диспетчера воркеров - синглтон
/*
	ведет учет воркеров
	предоставляет интерфейс отправки задач
*/
type IWorkerDispatcher interface {
	Init() error                  // инициализирует диспетчер
	Run()                         // запускает рабочий цикл диспетчера
	SendTask(entities.Task) error // отправляет задачу воркеру
}

// workerDispatcher диспетчера воркеров
type workerDispatcher struct {
	workersRegister map[*entities.Worker]entities.Priority // реестр воркеров
}

// Init инициализирует диспетчер
func (d *workerDispatcher) Init() (err error) {

	return
}

// Run запускает рабочий цикл диспетчера
func (d *workerDispatcher) Run() {

	return
}

// SendTask отправляет задачу воркеру
func (d *workerDispatcher) SendTask(task entities.Task) (err error) {
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
