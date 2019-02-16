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
	API          web.API
	JobContainer job.Container
}

// Config consists of all server configuration settings
type Config struct {
	Port           int
	RequestTimeout time.Duration
}

// New creates a new Server instance
func New(cfg Config, storage storage.Storage, client *github.Client, slackConfig job.SlackConfig) (*Server, error) {
	return &Server{
		Config:       cfg,
		Storage:      storage,
		API:          api.Default(storage),
		JobContainer: job.DefaultContainer(storage, client, slackConfig),
	}, nil
}

// Run runs the application server
func (s Server) Run() {
	defer s.Storage.Close()

	r := buildRouter(s.API)

	scheduler, err := scheduleJobs(s.JobContainer.Jobs())
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(s.Config.Port),
		ReadTimeout:  s.Config.RequestTimeout,
		WriteTimeout: s.Config.RequestTimeout,
	}

	log.Println("Starting scheduled jobs")
	scheduler.Start()

	log.Println("Server listening on port:", s.Config.Port)
	log.Fatal(server.ListenAndServe())
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

func scheduleJobs(jobs []job.Job) (*cron.Cron, error) {
	scheduler := cron.New()
	for _, job := range jobs {
		err := scheduler.AddJob(job.Period(), job)
		if err != nil {
			log.Printf("Unable to add job %q\n", job.Name())
			return nil, err
		}
	}

	jobNames := make([]string, 0)
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name())
	}
	log.Println("Scheduled jobs:", jobNames)

	return scheduler, nil
}
