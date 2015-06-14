package goquu

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	jobQueue    *JobQueue
	resultQueue *ResultQueue
	channels    []chan bool
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
			output, errstr, err := job.Execute()
			status := 0
			if err != nil {
				logger.Println(err)
				status = 1 // TODO: get real status code from Process
			}

			server.resultQueue.PushResult(&JobResult{
				Status: status,
				Stdout: output,
				Stderr: errstr,
				Job:    job,
			})
		}
	}
}

func NewServer() (server *Server, err error) {
	jq, err := NewJobQueue("./queue.db")
	if err != nil {
		return
	}
	rq, err := NewResultQueue("./queue.db")
	if err != nil {
		return
	}
	return &Server{jobQueue: jq, resultQueue: rq, channels: make([]chan bool, 0)}, err
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
	server.channels = append(server.channels, make(chan bool))
	go server.worker(server.channels[0])
	http.HandleFunc("/jobs", server.JobsAPIHandler)
	http.HandleFunc("/results", server.ResultsAPIHandler)
	if err = http.ListenAndServe(":8080", nil); err != nil {
		return
	}
	return
}
