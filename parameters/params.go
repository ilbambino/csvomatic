package parameters

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"time"
)

// QueryParams is the structure that contains all parameters read from a CSV
// and used to make a query to the database
type QueryParams struct {
	Hostname string
	Start    time.Time
	End      time.Time
}

// the date format found in the CSV files
const dateFormat = "2006-01-02 15:04:05"

// ReadFromCSV tries to read all the parameters from a CSV and returns them in
// a list (otherwise an error)
// if there are errors parsing the file, it will print them to stdout.
// TODO limit the amount of errors to some sane amount
func ReadFromCSV(reader io.Reader) ([]QueryParams, error) {

	var queryList []QueryParams
	someError := false
	csvReader := csv.NewReader(reader)
	lineNumber := 0
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil && lineNumber > 0 {
			log.Println(err)
			someError = true
			continue
		}

		if len(line) != 3 { //Make sure we have the needed fields
			log.Println("Not the expected fields, line:", lineNumber)
			someError = true
			continue
		}

		start, err := time.Parse(dateFormat, line[1])
		if err != nil && lineNumber > 0 {
			log.Println(err)
			someError = true
			continue
		}

		end, err := time.Parse(dateFormat, line[2])
		if err != nil && lineNumber > 0 {
			log.Println(err)
			someError = true
			continue
		}

		if end.Before(start) {
			log.Println("End time before start time, line:", lineNumber)
			someError = true
			continue
		}

		queryList = append(queryList, QueryParams{
			Hostname: line[0],
			Start:    start,
			End:      end,
		})
	}

	if someError {
		return queryList, fmt.Errorf("File contains errors")
	}

	return queryList, nil
}
