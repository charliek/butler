package main

import (
	"github.com/charliek/butler/service"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	log "github.com/ngmoco/timber"
	"net/http"
)

func setupServices() *service.ServiceRegistry {
	registry := service.NewRegistry()
	registry.AddService(&service.ButlerService{
		Name:        "grails-blog",
		Display:     "grails-blog",
		ServiceType: "upstart",
		Port:        8080,
	})
	registry.AddService(&service.ButlerService{
		Name:        "blog-service",
		Display:     "blog-service",
		ServiceType: "upstart",
		Port:        5678,
	})
	return registry
}

func setupLogging() {
	log.AddLogger(log.ConfigLogger{
		LogWriter: new(log.ConsoleWriter),
		Level:     log.DEBUG,
		Formatter: log.NewPatFormatter("[%D %T] [%L] %s %M"),
	})
}

func main() {
	setupLogging()
	registry := setupServices()

	m := martini.Classic()

	// setup middleware
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", registry.List())
	})

	m.Get("/api/services", func(r render.Render) {
		r.JSON(200, registry.List())
	})

	m.Get("/api/task/stop/:name", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		service, ok := registry.GetByName(params["name"])
		if ok {
			log.Info("Stopping service %s", service.Name)
			service.Stop()
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			http.NotFound(res, req)
		}
	})

	m.Get("/api/task/start/:name", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		service, ok := registry.GetByName(params["name"])
		if ok {
			log.Info("Starting service %s", service.Name)
			service.Start()
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			http.NotFound(res, req)
		}
	})

	m.Get("/api/task/local/:name", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		service, ok := registry.GetByName(params["name"])
		if ok {
			service.RunLocal()
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			http.NotFound(res, req)
		}
	})

	m.Get("/api/task/vagrant/:name", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		service, ok := registry.GetByName(params["name"])
		if ok {
			service.RunVagrant()
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			http.NotFound(res, req)
		}
	})

	m.Run()
}
