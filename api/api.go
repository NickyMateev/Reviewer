package api

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/api/project"
	"github.com/NickyMateev/Reviewer/api/pullrequest"
	"github.com/NickyMateev/Reviewer/api/user"
	"github.com/NickyMateev/Reviewer/web"
)

type defaultAPI struct {
	db *sql.DB
}

// Default returns an instance of the default API
func Default(db *sql.DB) web.API {
	return defaultAPI{
		db: db,
	}
}

// Controllers returns the default API's controllers
func (api defaultAPI) Controllers() []web.Controller {
	return []web.Controller{
		user.Controller(api.db),
		project.Controller(api.db),
		pullrequest.Controller(api.db),
	}
}
