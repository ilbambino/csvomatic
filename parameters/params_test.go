package parameters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCorrect(t *testing.T) {

	correct := `host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02`

	reader := strings.NewReader(correct)

	queries, err := ReadFromCSV(reader)

	assert.NoError(t, err, "We parsed the CSV without error")
	assert.Equal(t, len(queries), 2, "We have two entries")

	assert.Equal(t, queries[0].Hostname, "host_000008", "Hostname parsed")
	assert.Equal(t, queries[1].Hostname, "host_000001", "Hostname parsed")
}

func TestParseInCorrect(t *testing.T) {

	correct := `host_000008,2017-01-01 08,2017-01-01 09:59:22
host_000001,2017-01-02 14:02:02`

	reader := strings.NewReader(correct)

	queries, err := ReadFromCSV(reader)

	assert.Error(t, err, "We parsed the CSV with errors")
	assert.Equal(t, len(queries), 0, "We have no entries")

}

func TestParseEndBeforeStart(t *testing.T) {

	endBefore := `host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2016-01-02 14:02:02`

	reader := strings.NewReader(endBefore)

	queries, err := ReadFromCSV(reader)

	assert.Error(t, err, "We parsed the CSV with error")
	assert.Equal(t, len(queries), 1, "We have two entries")

}
func TestMissingFields(t *testing.T) {

	endBefore := `host_000008,2017-01-01 08:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02`

	reader := strings.NewReader(endBefore)

	queries, err := ReadFromCSV(reader)

	assert.Error(t, err, "We parsed the CSV with error")
	assert.Equal(t, len(queries), 1, "We have two entries")

}
