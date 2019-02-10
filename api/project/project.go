package project

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/web"
	"net/http"
)

// Controller returns an instance of the Project controller
func Controller(db *sql.DB) controller {
	return controller{
		db: db,
	}
}

// Routes returns all Project routes
func (c controller) Routes() []web.Route {
	return []web.Route{
		{
			Path:       web.ProjectsURL,
			Method:     http.MethodPost,
			HandleFunc: c.createProject,
		},
		{
			Path:       web.ProjectsURL,
			Method:     http.MethodGet,
			HandleFunc: c.listProject,
		},
		{
			Path:       web.ProjectsURL + "/{id}",
			Method:     http.MethodGet,
			HandleFunc: c.getProject,
		},
		{
			Path:       web.ProjectsURL + "/{id}",
			Method:     http.MethodDelete,
			HandleFunc: c.deleteProject,
		},
	}
}
