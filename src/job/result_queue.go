package job

import (
	"encoding/json"
	"../queue"
	"fmt"
)

type JobResult struct {
	Status int    `json:"status"`
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`

	Job *Job `json:"job"`
}

type ResultQueue struct {
	*queue.Queue
}

func NewResultQueue(path string) (resultQueue *ResultQueue, err error) {
	q, err := queue.New(path, "resultQueue-")
	if err != nil {
		return
	}
	resultQueue = &ResultQueue{q}
	return
}

func (resultQueue *ResultQueue) PushResult(result *JobResult)(err error) {
	str, err := json.Marshal(*result)
	if err != nil {
		return
	}
	return resultQueue.Push(str)
}

func (resultQueue *ResultQueue) ListResults()(results []JobResult) {
	for _, bytes := range resultQueue.List() {
		var result JobResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			fmt.Println(err)
		} else {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		results = make([]JobResult, 0)
	}

	return
}
