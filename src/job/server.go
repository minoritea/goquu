package job
import (
	"net/http"
	"encoding/json"
	"fmt"
)

type Server struct {
	jobQueue *JobQueue
	resultQueue *ResultQueue
}

func worker(ch chan bool, server *Server) {
	for {
		<- ch
		inner:
		for {
			j, err := server.jobQueue.PopJob()
			if err != nil {
				break inner
			}
			output, errstr, err := j.Execute()
			status := 0
			if err != nil {
				fmt.Println(err)
				status = 1 // TODO: get real status code from Process
			}

			server.resultQueue.PushResult(&JobResult{
				Status: status,
				Stdout: output,
				Stderr: errstr,
				Job: j,
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
	return &Server{jobQueue: jq, resultQueue: rq}, err
}

func (server *Server) Run() (err error) {
	ch := make(chan bool)
	q := server.jobQueue
	go worker(ch, server)
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "PUT" || r.Method == "POST" {
			var j Job
			err := json.NewDecoder(r.Body).Decode(&j)
			if err != nil {
		                w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
			} else if err := j.Validate(); err != nil {
		                w.WriteHeader(http.StatusBadRequest)
				fmt.Println(err)
			} else {
				q.PushJob(&j)
				ch <- true
		                w.WriteHeader(http.StatusOK)
			}
		} else {
			jobs := q.ListJobs()
			str, err := json.Marshal(jobs)
			if err != nil {
				fmt.Println(err)
		                w.WriteHeader(http.StatusInternalServerError)
			} else {
		                w.WriteHeader(http.StatusOK)
				w.Write(str)
			}
		}
	})
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request){
		results := server.resultQueue.ListResults()
		str, err := json.Marshal(results)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(str)
		}
	})
	http.ListenAndServe(":8080", nil)
	return
}
