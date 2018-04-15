// set ts=2 shiftwidth=2 expandtab
package rest

import (
	"fmt"
	"log"
	"net/http"
	"path"
  "strconv"
	"text/template"

	"bitbucket.org/enlab/peak/models"
	"bitbucket.org/enlab/peak/controllers/datastore"
	"bitbucket.org/enlab/peak/utils"
	"github.com/emicklei/go-restful"
)

var rootdir = "./ui"

type IncludeLibs struct {
	WebixCSS string
	WebixJS  string
	LogicJS  string
}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// usr/pwd = admin/admin
	u, p, ok := req.Request.BasicAuth()
	if !ok || u != "admin" || p != "admin" {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
	chain.ProcessFilter(req, resp)
}

func simpleUI(req *restful.Request, resp *restful.Response) {
	p := &IncludeLibs{
		WebixCSS: "/static/assets/js/webix/skins/contrast.css",
		WebixJS:  "/static/assets/js/webix/webix.js",
		LogicJS:  "/static/assets/js/peak/main.js",
	}
	// you might want to cache compiled templates
	t, err := template.ParseFiles("ui/index.html")
	if err != nil {
		log.Fatalf("Template gave: %s", err)
	}
	t.Execute(resp.ResponseWriter, p)
}

type RestService struct {
	listen    string
	dataDir   string
	dataStore datastore.DataStorer
	container *restful.Container
}

func (service RestService) registerSimpleUIResource(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").Filter(basicAuthenticate).To(simpleUI))
	ws.Route(ws.GET("/static/{subpath:*}").To(staticFromPathParam))
	ws.Route(ws.GET("/static").To(staticFromQueryParam))
	container.Add(ws)
}

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	actual := path.Join(rootdir, req.PathParameter("subpath"))
	fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
	// respe.WriteEntity(req.PathParameter)
}

func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}

func (service RestService) registerApiResource(container *restful.Container) {
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	ws := new(restful.WebService)
	ws.
		Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/mountains").
		To(service.getMountains).
		Doc("Get all available mountains").
		Operation("getMountains").
		Returns(200, "OK", []models.Mountain{}))

	ws.Route(ws.PUT("/mountains").
		To(service.putMountain).
		Doc("Search for the mountain").
		Operation("putMountain").
		Reads(models.Mountain{}))

  ws.Route(ws.DELETE("/mountains/{mountainId}").
		To(service.deleteMountain).
		Doc("Delete mountain by ID").
		Operation("deleteMountain").
		Param(ws.PathParameter("mountainId", "identifier of the mountain").DataType("int")).
		Returns(200, "OK", models.Mountain{}))

	ws.Route(ws.GET("/mountains/{mountainId}").
		To(service.getMountains).
		Doc("Get mountain info").
		Operation("getMountain").
		Param(ws.PathParameter("mountainId", "identifier of the mountain").DataType("int")).
		Returns(200, "OK", models.Mountain{}))

	ws.Route(ws.GET("/about").
		To(service.getAbout).
		Doc("About author and project").
		Operation("getAbout").
		Returns(200, "OK", []models.DevInfo{}))

	ws.Route(ws.GET("/stat").
		To(service.getSummary).
		Doc("Get stat for catalog").
		Operation("getSummary").
		Returns(200, "OK", []models.Summary{}))

	container.Add(ws)
}


func (service RestService) putMountain(request *restful.Request, response *restful.Response) {
  mountain := models.Mountain{}
	err := request.ReadEntity(&mountain)
	if err == nil {
    err := service.dataStore.PutMountain(&mountain)
    if err == nil {
      response.WriteHeaderAndEntity(http.StatusCreated, mountain)
    }
	} else {
		response.WriteError(http.StatusInternalServerError, err)
  }
}

func (service RestService) deleteMountain(request *restful.Request, response *restful.Response) {
	mountainID, _ := strconv.ParseUint(request.PathParameter("mountainId"), 0, 32)
  service.dataStore.DeleteMountain(mountainID)
}

func (service RestService) getMountains(request *restful.Request, response *restful.Response) {
	search := models.Search{}
	request.ReadEntity(&search)
	mountainID, _ := strconv.Atoi(request.PathParameter("mountainId"))
	page := utils.ParseInt(request.QueryParameter("page"))
	per_page := utils.ParseInt(request.QueryParameter("per_page"))
	result, err := service.dataStore.GetMountains(mountainID, page, per_page)

	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Book wasn't found\n")
	}
}

func (service RestService) getAbout(request *restful.Request, response *restful.Response) {
	devinfo := models.DevInfo{}
	devinfo.Author = "Alexandr Mikhailenko a.k.a Alex M.A.K."
	devinfo.Email = "alex-m.a.k@yandex.kz"
	devinfo.Project.Name = "mOPDS"
	devinfo.Project.Version = "0.1.0"
	devinfo.Project.Link = "bitbucket.org/enlab/peak"
	devinfo.Project.Created = "24.03.18 22:59"
	response.WriteEntity(devinfo)
}

func (service RestService) getSummary(request *restful.Request, response *restful.Response) {
	result, err := service.dataStore.GetSummary()
	if err == nil {
		response.WriteEntity(result)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Book wasn't found\n")
	}
}

func (service RestService) StartListen() {
	log.Println("Start listening on ", service.listen)
	server := &http.Server{Addr: service.listen, Handler: service.container}
	log.Fatal(server.ListenAndServe())
}

func NewRestService(listen string, dataStore datastore.DataStorer, dataDir string) RestServer {
	service := new(RestService)
	service.listen = listen
	service.dataStore = dataStore
	service.dataDir = dataDir
	service.container = restful.NewContainer()
	service.container.Router(restful.CurlyRouter{})
	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-Total-Count"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Content-Length", "X-Total-Count=100"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      service.container}
	service.container.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	service.container.Filter(service.container.OPTIONSFilter)

	service.registerSimpleUIResource(service.container)
	service.registerApiResource(service.container)

	return service
}
