package goquu

import (
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
)

type Job struct {
	Tags []string `json:"tags"`

	Command string   `json:"command"`
	Args    []string `json:"args"`
	Input   string   `json:"input"`
}

func (j *Job) Validate() (err error) {
	if len(j.Command) < 1 {
		return errors.New("Command must be given!")
	}
	if len(j.Tags) < 1 {
		return errors.New("Tags must contain at least one element")
	}
	for _, tag := range j.Tags {
		if len(tag) < 1 {
			return errors.New("Each tag should not be empty!")
		}
	}
	return
}

func (j *Job) Execute() (output string, errstr string, err error) {
	if err = j.Validate(); err != nil {
		return
	}
	command := j.Command
	args := j.Args
	input := j.Input
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return output, errstr, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return output, errstr, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return output, errstr, err
	}

	if err = cmd.Start(); err != nil {
		return output, errstr, err
	}

	io.WriteString(stdin, input)
	stdin.Close()

	tmp, err := ioutil.ReadAll(stdout)
	if err != nil {
		logger.Println(err)
	}
	output = string(tmp)

	tmp = []byte{}
	tmp, err = ioutil.ReadAll(stderr)
	if err != nil {
		logger.Println(err)
	}
	errstr = string(tmp)

	err = nil

	if err = cmd.Wait(); err != nil {
		return
	}
	return
}
