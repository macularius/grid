package operator

// Operator получатель и отправитель сообщений
type Operator struct {
}

// SendTasks отправляет задачи в распределитель
func (o *Operator) SendTasks(tasks []string) (err error) {

	// формирование сообщения

	// установление соединения

	// отправка сообщения

	return
}

// Listener ожидает ответы
func (o *Operator) Listener(chAnswer chan<- int) (err error) {

	return
}
