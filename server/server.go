package server

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/api"
	"github.com/NickyMateev/Reviewer/job"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Server represents the application's server
type Server struct {
	Config       Config
	DB           *sql.DB
	API          web.API
	JobContainer job.Container
}

// Config consists of all server configuration settings
type Config struct {
	Port           int
	RequestTimeout time.Duration
}

// New creates a new Server instance
func New(cfg Config, db *sql.DB) (*Server, error) {
	return &Server{
		Config:       cfg,
		DB:           db,
		API:          api.Default(db),
		JobContainer: job.DefaultContainer(db),
	}, nil
}

// Run runs the application server
func (s Server) Run() {
	defer s.DB.Close()

	r := s.buildRouter()

	err := s.startJobs()
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(s.Config.Port),
		ReadTimeout:  s.Config.RequestTimeout,
		WriteTimeout: s.Config.RequestTimeout,
	}

	log.Println("Server listening on port:", s.Config.Port)
	log.Fatal(server.ListenAndServe())
}

func (s Server) buildRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	controllers := s.API.Controllers()
	for _, controller := range controllers {
		routes := controller.Routes()
		for _, route := range routes {
			router.HandleFunc(route.Path, route.HandleFunc).Methods(route.Method)
		}
	}

	return router
}

func (s Server) startJobs() error {
	jobs := s.JobContainer.Jobs()

	scheduler := cron.New()
	for _, job := range jobs {
		err := scheduler.AddJob("@every "+job.Period(), job)
		if err != nil {
			return err
		}
	}
	scheduler.Start()
	return nil
}
