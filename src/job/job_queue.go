package job

import (
	"encoding/json"
	"../queue"
	"fmt"
)

type JobQueue struct {
	*queue.Queue
}

func NewJobQueue(path string) (jobQueue *JobQueue, err error) {
	q, err := queue.New(path, "jobQueue-")
	if err != nil {
		return
	}
	jobQueue = &JobQueue{q}
	return
}


func (jobQueue *JobQueue) PopJob() (j *Job, err error) {
	jsonStr, err := jobQueue.Pop()
	if err != nil {
		return
	}
	var result Job
	err = json.Unmarshal(jsonStr, &result)
	j = &result
	if err != nil {
		return
	}
	return
}

func (jobQueue *JobQueue) PushJob(j *Job)(err error) {
	str, err := json.Marshal(*j)
	if err != nil {
		return
	}
	return jobQueue.Push(str)
}

func (jobQueue *JobQueue) ListJobs()(jobs []Job) {
	for _, bytes := range jobQueue.List() {
		var j Job
		if err := json.Unmarshal(bytes, &j); err != nil {
			fmt.Println(err)
		} else {
			jobs = append(jobs, j)
		}
	}

	if len(jobs) == 0 {
		jobs = make([]Job, 0)
	}

	return
}
