package new_rest

import (
  "fmt"
  "log"
  "time"
  "strconv"
	"net/http"
	"github.com/jinzhu/gorm"
	"bitbucket.org/enlab/peak/models"
	"github.com/ant0ine/go-json-rest/rest"
	"bitbucket.org/enlab/peak/controllers/paginate"
  "github.com/StephanDollberg/go-json-rest-middleware-jwt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var rootdir = "./ui"
const (
	DEFAULT_PAGE_SIZE int = 20
	MAX_PAGE_SIZE     int = 1000
)

type dbStore struct {
	DB *gorm.DB
}

func handle_auth(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"authed": r.Env["REMOTE_USER"].(string)})
}

func NewRestService(config *models.DBConfig) {
  i := dbStore{}
	i.NewDBStore(config)

  jwt_middleware := &jwt.JWTMiddleware{
		Key:        []byte("secret key"),
		Realm:      "jwt auth",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return userId == "admin" && password == "admin"
    },
  }

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
  // we use the IfMiddleware to remove certain paths from needing authentication
	/*api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: jwt_middleware,
  })*/
  api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == "http://localhost"
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
  })
	router, err := rest.MakeRouter(
    rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/auth_test", handle_auth),
    rest.Get("/refresh_token", jwt_middleware.RefreshHandler),
		rest.Get("/mountains", i.GetMountains),
		rest.Put("/mountains", i.PutMountain),
		rest.Get("/mountains/:id", i.GetMountains),
		rest.Delete("/mountains/:id", i.DeleteMountain),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
  http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.MakeHandler()))
  http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(rootdir))))
  log.Fatal(http.ListenAndServe(":8000", nil))
}

func (i *dbStore) PutMountain(w rest.ResponseWriter, r *rest.Request) {
	mountain := models.Mountain{}
	if err := r.DecodeJsonPayload(&mountain); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  fmt.Println(mountain)
	if err := i.DB.Save(&mountain).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&mountain)
}

func (i *dbStore) GetMountains(w rest.ResponseWriter, r *rest.Request) {
	count := 0
	result := []models.Mountain{}

  mountainID, _ := strconv.Atoi(r.PathParam("id"))
  page, _ := strconv.Atoi(r.URL.Query().Get("page"))
  per_page, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	m := i.DB.Select("mountains.*").Table("mountains")
  if mountainID != 0 {
    m = m.Where("id = ?", mountainID)
  }
	if per_page <= 0 {
		per_page = DEFAULT_PAGE_SIZE
	}
	if per_page > MAX_PAGE_SIZE {
		per_page = MAX_PAGE_SIZE
	}
	if page == 0 {
		page = 1
	}
	m.Count(&count)

  if count > 1 {
    p := paginate.NewPaginatedList(page, per_page, count)
    m = m.Limit(p.Limit())
    m = m.Offset(p.Offset())
    m.Find(&result)
    p.Items = result
    w.WriteJson(&p)
  } else if count == 1 {
    m.Find(&result)
    w.WriteJson(&result[0])
  } else {
    rest.Error(w, "mountains not found", 400)
  }
}

func (store *dbStore) DeleteMountain(w rest.ResponseWriter, r *rest.Request) {
  id, _ := strconv.Atoi(r.PathParam("id"))
	mountain := models.Mountain{}
	if store.DB.First(&mountain, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := store.DB.Delete(&mountain).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  w.WriteHeader(http.StatusOK)
}

func (store *dbStore) GetSummary(w rest.ResponseWriter, r *rest.Request) {
	summary := models.Summary{}
	mountains_count := 0
	store.DB.Table("mountains").Count(&mountains_count)
	summary.Mountains = mountains_count

  w.WriteJson(&summary)
}

// NewDBStore creates new instance of datastorer
func (i *dbStore) NewDBStore(config *models.DBConfig) {
  var err error
	i.DB, err = gorm.Open(config.DBType, config.DBParams)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
  i.DB.LogMode(config.DBLog)
}
