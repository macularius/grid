package operator

var answerChIn chan<- string // канал ответов

// Operator получатель и отправитель сообщений
type Operator struct {
}

// Init регистрирует задачу в дистрибуторе
func (o *Operator) Init(taskCount int) (err error) {
}

// SendTask отправляет задачу в распределитель
func (o *Operator) SendTask(task string) (err error) {

	return
}

// Listener ожидает ответы
func (o *Operator) Listener(answChIn chan<- string) (err error) {

	return
}
