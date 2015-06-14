package goquu

import (
	"encoding/json"
	"net/http"
	"errors"
)

type Server struct {
	jobQueue    *JobQueue
	resultQueue *ResultQueue
	channels    []chan bool
	config *Config
}

func (server *Server) worker(ch chan bool) {
	for {
		<-ch
	inner:
		for {

			job, err := server.jobQueue.PopJob()
			if err != nil {
				break inner
			}
			for job != nil {
				result, err := job.Execute()
				server.resultQueue.PushResult(result)
				if err != nil && job.ErrorCallBack != nil {
					job = job.ErrorCallBack
				} else {
					job = nil
				}
			}
		}
	}
}

func NewServer() (server *Server, err error) {
	config, err := loadConfig()
	if err != nil {
		return
	}

	err = SetLoggerFromFile(config.LogFile, config.LogFlag)
	if err != nil {
		return
	}

	jq, err := NewJobQueue(config.DBDirectory)
	if err != nil {
		return
	}

	rq, err := NewResultQueue(config.DBDirectory)
	if err != nil {
		return
	}

	if config.WorkerSize < 1 {
		err = errors.New("Worker size must be greater than zero!")
		return
	}

	return &Server{jobQueue: jq, resultQueue: rq, channels: make([]chan bool, config.WorkerSize), config: config}, err
}

func (server *Server) ResultsAPIHandler(w http.ResponseWriter, r *http.Request) {
	results := server.resultQueue.ListResults()
	str, err := json.Marshal(results)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(str)
	}
}

func (server *Server) createNewJobHandler(w http.ResponseWriter, r *http.Request) {
	var job Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
	} else if err := job.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(err)
	} else {
		server.jobQueue.PushJob(&job)
		for _, ch := range server.channels {
			ch <- true
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (server *Server) getListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs := server.jobQueue.ListJobs()
	str, err := json.Marshal(jobs)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(str)
	}
}

func (server *Server) JobsAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "POST":
		server.createNewJobHandler(w, r)
	case "GET":
		server.getListJobsHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (server *Server) Run() (err error) {
	for i:=0; i< server.config.WorkerSize; i++ {
		server.channels[i] = make(chan bool)
		go server.worker(server.channels[i])
	}
	for _, ch := range server.channels {
		ch <- true
	}
	http.HandleFunc("/jobs", server.JobsAPIHandler)
	http.HandleFunc("/results", server.ResultsAPIHandler)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		return
	}
	return
}
