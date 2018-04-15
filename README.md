# mPeak API
Need description...

## Install
```bash
mak@kaznii ~ $ go get bitbucket.org/enlab/mpeak
```

## Usage
```bash
  -about
    	About author and this project
  -catalog string
    	Directory of sqlite database (mandatory only for sqlite)
  -config string
    	Default configuration file (default "./conf/peak.conf")
  -dbtype string
    	Type used database: sqlite3, mysql or postgres (default "sqlite3")
  -del
    	Del mountain by ID
  -get
    	List all mountains
  -put
    	Put new mountain
  -id int
    	Id mountain (not mandory for -get)
  -page int
    	Pagination 1...n
  -per_page int
    	Limit results (-1 for no limit) (default 25)
  -database string
    	Database name connect to database
  -sslmode string
    	Whether to use ssl mode or not, here's the question: disable or enable
  -host string
    	IP address for connect to database (default "modps")
  -username string
    	Username for connect to database (default "peak")
  -password string
    	Password for connect to database (default "peak")
  -listen string
    	Set server listen address:port (default ":8000")
  -verbose
    	Verbose output
  -debug
    	Debug output
```

## Example configuration file
```bash
page 1
per_page 25
listen :8000
host localhost
username mopds
password mopds
database mopds
dbtype sqlite3
sslmode disable
debug false
```

#### Get all mountains
```bash
mak@denied ~ $ http GET ':8000/api/v1/mountains?page=1&per_page=25'
```
#### Put new mountain
```bash
mak@kaznii ~ $ http PUT :8000/api/v1/mountains title="TEST" web_id=222 height=2300 latitude="43.34232" longtitude="32.32232" type_link=2
```
#### Delete mountain by ID
```bash
mak@kaznii ~ $ http DELETE :8000/api/v1/mountains/1
```

### Dependencies
* github.com/namsreal/flag
* github.com/emicklei/go-restful
* github.com/jinzhu/gorm
