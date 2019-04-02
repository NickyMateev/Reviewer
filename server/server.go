package server

import (
	"github.com/NickyMateev/Reviewer/api"
	"github.com/NickyMateev/Reviewer/job"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/google/go-github/github"
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
	Storage      storage.Storage
	Router       *mux.Router
	JobContainer job.Container
}

// Config consists of all server configuration settings
type Config struct {
	Port           int
	RequestTimeout time.Duration
}

// New creates a new Server instance
func New(cfg Config, storage storage.Storage, client *github.Client, slackConfig job.SlackConfig) *Server {
	defaultAPI := api.Default(storage)
	router := buildRouter(defaultAPI)

	return &Server{
		Config:       cfg,
		Storage:      storage,
		Router:       router,
		JobContainer: job.DefaultContainer(storage, client, slackConfig),
	}
}

// Run runs the application server
func (s *Server) Run() error {
	err := s.startJobs()
	if err != nil {
		return err
	}

	server := http.Server{
		Handler:      s.Router,
		Addr:         ":" + strconv.Itoa(s.Config.Port),
		ReadTimeout:  s.Config.RequestTimeout,
		WriteTimeout: s.Config.RequestTimeout,
	}

	log.Println("Server listening on port:", s.Config.Port)
	return server.ListenAndServe()
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

func (s *Server) startJobs() error {
	jobs := s.JobContainer.Jobs()

	scheduler := cron.New()
	for _, job := range jobs {
		err := scheduler.AddJob(job.Period(), job)
		if err != nil {
			return err
		}
	}

	jobNames := make([]string, 0)
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name())
	}
	log.Println("Scheduled jobs:", jobNames)

	scheduler.Start()
	return nil
}
