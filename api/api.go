package api

import (
	"github.com/NickyMateev/Reviewer/api/project"
	"github.com/NickyMateev/Reviewer/api/pullrequest"
	"github.com/NickyMateev/Reviewer/api/user"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
)

type defaultAPI struct {
	storage storage.Storage
}

// Default returns an instance of the default API
func Default(storage storage.Storage) web.API {
	return defaultAPI{
		storage: storage,
	}
}

// Controllers returns the default API's controllers
func (api defaultAPI) Controllers() []web.Controller {
	return []web.Controller{
		user.Controller(api.storage),
		project.Controller(api.storage),
		pullrequest.Controller(api.storage),
	}
}
