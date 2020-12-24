package master_worker_job

type Worker struct {
	Id          int
	WorkerQueue chan chan Job
	JobChan     chan Job
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan Job) *Worker {
	return &Worker{
		Id:          id,
		WorkerQueue: workerQueue,
		JobChan:     make(chan Job),
		QuitChan:    make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.JobChan
			select {
			case job := <-w.JobChan:
				job.Do()
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
