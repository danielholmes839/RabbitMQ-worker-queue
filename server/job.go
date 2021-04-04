package main

import (
	"encoding/json"
)

type Job struct {
	Id        int
	Name      string
	Completed bool
}

func (job *Job) serialize() []byte {
	data, _ := json.Marshal(job)
	return data
}

func NewJob(data []byte) *Job {
	job := &Job{}
	json.Unmarshal(data, job)
	return job
}
