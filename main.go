package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/ilbambino/csvomatic/parameters"
	"github.com/ilbambino/csvomatic/stats"
	"github.com/ilbambino/csvomatic/worker"
	_ "github.com/lib/pq"
)

var numberWorkers int
var csvFile string

func setupArgs() {
	flag.IntVar(&numberWorkers, "workers", 2, "Number of workers")
	flag.StringVar(&csvFile, "input", "query_params.csv", "Path to csv with query arguments")
	flag.Parse()

}

// postgresConf holds the configuration to connect to the DB
type postgresConf struct {
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
	Database string `env:"DB_NAME" envDefault:"homework"`
}

// getDBConnectionFromEnv constructs a Postgres connection string from env vars or from defaults
func getDBConnectionFromEnv() (string, error) {

	cfg := postgresConf{}
	err := env.Parse(&cfg)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database), nil
}

func main() {

	commandStart := time.Now()
	setupArgs()

	connStr, err := getDBConnectionFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fileReader, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	queries, err := parameters.ReadFromCSV(fileReader)
	if err != nil {
		log.Fatal(err)
	}

	if len(queries) == 0 {
		log.Println("No valid queries in source file")
		return
	}

	pool := worker.CreatePool(numberWorkers, len(queries), db)

	for _, query := range queries {
		pool.Queue(worker.Job{QueryParams: query})
	}
	measures := pool.WaitUntilDone()

	fmt.Println("Tool Time:", time.Since(commandStart))

	metrics := stats.ProcessDurations(measures)
	tmpl, err := template.New("metrics").Parse(outTemplate)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(os.Stdout, metrics)
	if err != nil {
		log.Fatal(err)
	}

}

var outTemplate = `Total Queries: {{.Queries}}
Total Query Time: {{.Total}}
Query Times:
	Max:		{{.Max}}
	Min:		{{.Min}}
	Median:		{{.Median}}
	Average:	{{.Average}}
`
