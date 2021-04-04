package main

import (
	"sort"
	"sync"
)

type Database struct {
	// Mock in-mememory database
	*sync.Mutex
	counter int
	jobs    map[int]*Job
}

func NewDatabase() *Database {
	return &Database{&sync.Mutex{}, 0, make(map[int]*Job)}
}

func (db *Database) createJob(name string) *Job {
	// Create a new job
	db.Lock()
	defer db.Unlock()

	db.counter++
	db.jobs[db.counter] = &Job{db.counter, name, false}
	return db.jobs[db.counter]
}

func (db *Database) updateJob(id int) {
	// "Update" a job as completed
	db.Lock()
	defer db.Unlock()

	db.jobs[id].Completed = true
}

func (db *Database) read() []*Job {
	// Read jobs from the db
	jobs := make([]*Job, 0, len(db.jobs))

	for _, job := range db.jobs {
		jobs = append(jobs, job)
	}

	sort.SliceStable(jobs, func(i, j int) bool {
		return jobs[i].Id < jobs[j].Id
	})

	return jobs
}
