package main

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/api"
	"github.com/NickyMateev/Reviewer/job"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open(storage.DbType, storage.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = storage.UpdateSchema(db)
	if err != nil {
		panic(err)
	}
	log.Println("Database is up-to-date")

	router := buildRouter(api.Default(db))
	startJobs(db)
	log.Fatal(http.ListenAndServe(":8888", router))
}

func buildRouter(api web.API) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	controllers := api.Controllers()
	for _, controller := range controllers {
		routes := controller.Routes()
		for _, route := range routes {
			router.HandleFunc(route.Path, route.HandleFunc).Methods(route.Method)
		}
	}

	return router
}

func startJobs(db *sql.DB) {
	c := job.DefaultContainer(db)
	jobs := c.Jobs()

	scheduler := cron.New()
	for _, job := range jobs {
		scheduler.AddJob("@every " + job.Period(), job)
	}
	scheduler.Start()
}