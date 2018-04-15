package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"bitbucket.org/enlab/peak/models"
	"bitbucket.org/enlab/peak/controllers/new_rest"
	"bitbucket.org/enlab/peak/controllers/datastore"
	"bitbucket.org/enlab/peak/utils"
	"github.com/namsral/flag"
)

func setDatabaseConfig(Catalog string, username string, password string, host string, dbname string, dbtype string, dblog bool, sslmode string) *models.DBConfig {
	result := new(models.DBConfig)

	result.DBType = dbtype
	result.DBLog = dblog
	if dbtype == "sqlite3" {
		fileData, err := ioutil.ReadFile(filepath.Join(Catalog, "mountains.db"))
		if err == nil {
			err = json.Unmarshal(fileData, result)
		}

		if err != nil { // fallback to sqlite
			result.DBParams = filepath.Join(Catalog, "mountains.db")
		}
	} else if dbtype == "postgres" {
		if sslmode == "" {
			sslmode = "disable"
		}
		result.DBParams = fmt.Sprintf("user=%s password=%s DB.name=%s sslmode=%s", username, password, dbname, sslmode)
	}

	return result
}

func main() {
	var (
    Catalog           string
		config            string
		Page              int
		PerPage           int
		Listen            string
		Verbose           bool
		Debug             bool
		About             bool
		Host              string
		Username          string
		Password          string
		DBName            string
		DBType            string
		DBLog             bool
		SSLMode           string

    GetMountains         bool
    DelMountain         bool
    PutMountain         bool
    IdMountain              int
	)

	flag.StringVar(&config, "config", "./conf/peak.conf", "Default configuration file")
	flag.StringVar(&Catalog, "catalog", "", "Directory of sqlite database (mandatory only for sqlite)")
	flag.IntVar(&Page, "page", 0, "Pagination 1...n")
	flag.IntVar(&PerPage, "per_page", 25, "Limit results (-1 for no limit)")
	flag.StringVar(&Listen, "listen", ":8000", "Set server listen address:port")
	flag.BoolVar(&Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&Debug, "debug", false, "Debug output")
	flag.BoolVar(&About, "about", false, "About author and this project")
	flag.StringVar(&Host, "host", "modps", "IP address for connect to database")
	flag.StringVar(&Username, "username", "peak", "Username for connect to database")
	flag.StringVar(&Password, "password", "peak", "Password for connect to database")
	flag.StringVar(&DBType, "dbtype", "sqlite3", "Type used database: sqlite3, mysql or postgres")
	flag.StringVar(&DBName, "database", "", "Database name connect to database")
	flag.StringVar(&SSLMode, "sslmode", "", "Whether to use ssl mode or not, here's the question: disable or enable")

	flag.IntVar(&IdMountain, "id", 0, "Id mountain")
	flag.BoolVar(&GetMountains, "get", false, "List all mountains")
	flag.BoolVar(&DelMountain, "del", false, "Del mountain by ID")
	flag.BoolVar(&PutMountain, "put", false, "Put new mountain")
	flag.Parse()

	DBLog = Debug
	conf := setDatabaseConfig(Catalog, Username, Password, Host, DBName, DBType, DBLog, SSLMode)
	store, err := datastore.NewDBStore(conf)
	if err != nil {
		log.Fatalln("Failed to open database")
	}
	defer store.Close()

  if GetMountains {
    result, err := store.GetMountains(IdMountain, Page, PerPage)
    if err == nil {
      utils.PrintJson(result, true)
    } else {
      log.Println("Nothing found")
    }
  } else if PutMountain {
    fmt.Println("not ready")
  } else if DelMountain {
    fmt.Println("not ready")
  } else if About {
		devinfo := models.DevInfo{}
		devinfo.Author = "Alexandr Mikhailenko a.k.a Alex M.A.K."
		devinfo.Email = "alex-m.a.k@yandex.kz"
		devinfo.Project.Name = "mOPDS"
		devinfo.Project.Version = "0.1.0"
		devinfo.Project.Link = "bitbucket.org/enlab/peak"
		devinfo.Project.Created = "24.03.18 22:59"

		utils.PrintJson(devinfo, true)
	} else {
		//rest.NewRestService(Listen, store, Catalog).StartListen()
		new_rest.NewRestService(conf)

  }
}
