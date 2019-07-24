// Package worker contains all the logic or workers to execute queries to a DB
// It makes queries to a DB and returns how long it took to perform the operation
// The jobs are assigned deterministic based on the query (Hostname)
package worker

import (
	"database/sql"
	"hash/crc32"
	"log"
	"time"

	"github.com/ilbambino/csvomatic/parameters"
)

// Job holds the information needed to execute a query. The query read from the file
// and the DB connection
type Job struct {
	parameters.QueryParams
	db *sql.DB
}

// Result of the execution of a job
type Result struct {
	OK       bool
	Duration time.Duration
}

type jobQueue chan Job
type resultQueue chan Result

// Pool contains a pool of worker goroutines that execute the jobs.
type Pool struct {
	jobQueues []jobQueue
	results   resultQueue
	size      int
	jobs      int
	db        *sql.DB
}

// CreatePool creates a pool of workers of the given size and also the channels needed
// We need to pass the expected amount of results, so we can simplify the code as we will
// reserve a channel of the size of the results, hence we won't have a dead lock
func CreatePool(size int, expectedResults int, db *sql.DB) Pool {

	pool := Pool{size: size, db: db}
	pool.jobQueues = make([]jobQueue, size)
	pool.results = make(resultQueue, expectedResults)

	for j := 0; j < size; j++ {
		pool.jobQueues[j] = make(jobQueue, expectedResults/size)

		go worker(j, pool.jobQueues[j], pool.results)
	}

	return pool
}

// Queue gets a job and puts it in a queue to be processed
func (p *Pool) Queue(job Job) {

	job.db = p.db
	p.jobQueues[job.GetWorkerID(p.size)] <- job
	p.jobs++
}

// WaitUntilDone blocks and waits until all jobs have been executed and returns a slice of
// Duration of all the successful results
// This method closes the input channels, so the worker goroutines will die when done.
func (p Pool) WaitUntilDone() []time.Duration {

	for j := 0; j < p.size; j++ {
		close(p.jobQueues[j]) //close all job queues, no more jobs accepted
	}

	measures := make([]time.Duration, 0, p.size) //we know the max size we expect

	// wait until all jobs are processed
	for i := 0; i < p.jobs; i++ {
		res := <-p.results
		if res.OK {
			measures = append(measures, res.Duration)
		}
	}
	return measures
}

// GetWorkerID generates a number [0, numberWorkers-1] based on the host name
// the same hostname always gives the same worker (if the number of workers is the same)
func (j Job) GetWorkerID(numberWorkers int) uint {

	hashed := crc32.Checksum([]byte(j.Hostname), crc32.IEEETable)
	return uint(int(hashed) % numberWorkers) // we don't care if there is overflow
}

// Execute makes the query to the DB with the job params. Returns the Duration taken
func (j Job) Execute() (time.Duration, error) {

	queryStart := time.Now()
	rows, err := j.db.Query("SELECT MAX(usage) FROM cpu_usage WHERE host = $1 AND ts >= $2 AND ts < $3", j.Hostname, j.Start, j.End)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var max float32
	for rows.Next() {
		err := rows.Scan(&max) //Force fetching from DB
		if err != nil {
			return 0, err
		}
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	lambda := time.Since(queryStart)
	return lambda, nil
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		// fmt.Println("worker", id, "started  job", j)
		lambda, err := j.Execute()
		// fmt.Println("worker", id, "finished job", j)
		results <- Result{OK: err == nil, Duration: lambda}
	}
}
