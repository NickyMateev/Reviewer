package user

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/web"
	"net/http"
)

// Controller returns an instance of the User controller
func Controller(db *sql.DB) *controller {
	return &controller{
		db: db,
	}
}

// Routes returns all User routes
func (c *controller) Routes() []web.Route {
	return []web.Route{
		{
			Path:       web.UsersURL,
			Method:     http.MethodGet,
			HandleFunc: c.listUsers,
		},
		{
			Path:       web.UsersURL + "/{id}",
			Method:     http.MethodGet,
			HandleFunc: c.getUser,
		},
		{
			Path:       web.UsersURL + "/{id}",
			Method:     http.MethodPatch,
			HandleFunc: c.patchUser,
		},
		{
			Path:       web.UsersURL + "/{id}",
			Method:     http.MethodDelete,
			HandleFunc: c.deleteUser,
		},
	}
}
