# Tool to get some stats from TimescaleDB

Small command line tool to measure how fast queries are run against a Postgres (TimescaleDB) database.

The queries are run against a database with a table called `cpu_usage` which needs to have the following schema:

```
  ts    TIMESTAMPTZ,
  host  TEXT,
  usage DOUBLE PRECISION
```

The queries that are run against the database come from a CSV file, that you need to pass to the tool (using the `-input` option). The CSV needs to have the following format:

```
hostname,start_time,end_time
host_name,2017-01-01 08:59:22,2017-01-01 09:59:22
…
```

If everything is working properly at the end you should get some stats like:

```
Tool Time: 571.636873ms
Total Queries: 200
Total Query Time: 974.882794ms
Query Times:
	Max:		21.198181ms
	Min:		3.308055ms
	Median:		4.747431ms
	Average:	4.874413ms
```

## Run it

To run it you need to have Go >= 1.11 as it uses Go modules. So pull the repo into some path that is not in your GOPATH (if you have it). Then just run `go build` to get the binary. Dependencies are not vendored so you need a working Internet connection to get started. 

Once you have the binary, all you need is to run `csvomatic -workers 2 -input pathto.csv`  to run it.

### Configure the connection to the database

To connect to a database you need to pass the needed configuration. To do so you do it with environment variables. If you don't have those set there will be some defaults in use.

| Information | Env Var | Default |
| ------------|---------| --------|
| Hostname | DB_HOST | localhost |
| Port | DB_PORT | 5432 |
| Username | DB_USER | postgres |
| Password | DB_PASSWORD | password |
| Database | DB_NAME | homework |

So far no SSL connections supported.

If you want to run it with env vars you can also do something like:

```
$ DB_HOST="somehost" DB_PASSWORD="mypass" DB_NAME="myDB" ./csvomatic -input myfile.csv
```

