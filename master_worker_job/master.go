package master_worker_job

const (
	MaxWorkerNum = 2048
)

type Master struct {
	WorkerNum   int
	WorkerQueue chan chan Job
	JobQueue    chan Job
	WorkerList  []*Worker
}

func NewMaster(workerNum int, jobNum int) *Master {
	if workerNum <= 0 || workerNum > MaxWorkerNum {
		workerNum = MaxWorkerNum
	}

	if jobNum > 2*workerNum {
		jobNum = 2 * workerNum
	}

	return &Master{
		WorkerNum:   workerNum,
		WorkerQueue: make(chan chan Job, workerNum),
		JobQueue:    make(chan Job, jobNum),
	}
}

func (m *Master) Run() {
	for i := 0; i < m.WorkerNum; i++ {
		worker := NewWorker(i, m.WorkerQueue)
		m.WorkerList = append(m.WorkerList, worker)
		worker.Start()
	}

	go m.Dispatch()
}

func (m *Master) Dispatch() {
	for {
		select {
		case job := <-m.JobQueue:
			go func() {
				jobChan := <-m.WorkerQueue
				jobChan <- job
			}()
		}
	}
}

func (m *Master) AddJob(job Job) {
	m.JobQueue <- job
}

func (m *Master) Stop() {
	for _, w := range m.WorkerList {
		w.Stop()
	}
}
