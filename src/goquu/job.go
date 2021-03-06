package goquu

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
	"math"
)

type Job struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Input   string   `json:"input"`
	Timeout int	`json:"timeout"`

	Tags []string `json:"tags"`

	ErrorCallBack *Job `json:"errback"`
}

type JobResult struct {
	Status int    `json:"status"`
	ErrorMsg  string `json:"error_message"`
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`

	Job *Job `json:"job"`
}

func (j *Job) Validate() (err error) {
	if len(j.Command) < 1 {
		return errors.New("Command must be given!")
	}

	return
}

func (job *Job) executeCommand() (output string, errstr string, err error) {
	command := job.Command
	args := job.Args
	input := job.Input
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return output, errstr, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return output, errstr, err
	}

	var stdin io.WriteCloser
	if len(job.Input) > 0 {
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return output, errstr, err
		}
	}

	if err = cmd.Start(); err != nil {
		return output, errstr, err
	}

	ch := make(chan error)

	go func() {
		if stdin != nil {
			io.WriteString(stdin, input)
			stdin.Close()
		}

		tmp, read_err := ioutil.ReadAll(stdout)
		if read_err != nil {
			logger.Println(read_err)
		}
		output = string(tmp)

		tmp = []byte{}
		tmp, read_err = ioutil.ReadAll(stderr)
		if read_err != nil {
			logger.Println(read_err)
		}

		errstr = string(tmp)
		ch <- cmd.Wait()
	}()

	var duration time.Duration
	if job.Timeout > 0 {
		duration = time.Duration(job.Timeout) * time.Second
	} else {
		duration = math.MaxInt64
	}

	select {
	case <-time.After(duration):
		if err = cmd.Process.Kill(); err != nil {
		} else {
			err = errors.New(fmt.Sprintf("Process:\"%s\" is killed", strings.Join(cmd.Args, " ")))
		}
		return
	case err = <-ch:
		return
	}

	return
}

func (job *Job) Execute() (result *JobResult, err error) {
	if err = job.Validate(); err != nil {
		return
	}
	output, errstr, err := job.executeCommand()
	status := 0
	errmsg := ""
	if err != nil {
		status = 1 // TODO: get real status code from Process
		errmsg = err.Error()
	}
	result = &JobResult{
		Status: status,
		Stdout: output,
		Stderr: errstr,
		ErrorMsg: errmsg,
		Job: job,
	}
	return
}
