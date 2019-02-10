package pullrequest

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/web"
	"net/http"
)

// Controller returns an instance of the Pull Request controller
func Controller(db *sql.DB) controller {
	return controller{
		db: db,
	}
}

// Routes returns all Pull Request routes
func (c controller) Routes() []web.Route {
	return []web.Route{
		{
			Path:       web.PullRequestsURL,
			Method:     http.MethodGet,
			HandleFunc: c.listPullRequest,
		},
		{
			Path:       web.PullRequestsURL + "/{id}",
			Method:     http.MethodGet,
			HandleFunc: c.getPullRequest,
		},
	}
}
