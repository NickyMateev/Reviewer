package project

import (
	"github.com/NickyMateev/Reviewer/storage/storagefakes"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProjects(t *testing.T) {
	c := controller{
		storage: &storagefakes.FakeStorage{},
	}

	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/v1/projects").Expect().Status(http.StatusOK)
}
