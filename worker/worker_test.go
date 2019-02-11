package worker

import (
	"testing"
	"time"

	"github.com/ilbambino/csvomatic/parameters"
	"github.com/stretchr/testify/assert"
)

func TestJobWorkerID(t *testing.T) {

	job1 := Job{QueryParams: parameters.QueryParams{Hostname: "myHost1", Start: time.Now(), End: time.Now()}}

	assert.Equal(t, job1.GetWorkerID(10), uint(4), "We get the expected ID")
	assert.Equal(t, job1.GetWorkerID(7), uint(3), "We get the expected ID different")

	job1.Hostname = "myHost2"
	assert.Equal(t, job1.GetWorkerID(10), uint(2), "Host changed, ID changed")

	job1.Hostname = "myHost1"
	assert.Equal(t, job1.GetWorkerID(10), uint(4), "Original Name, Original ID")

}
