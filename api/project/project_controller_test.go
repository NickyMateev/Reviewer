package project

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage/storagefakes"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProjectsReturnsResponse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fakeStorage := storagefakes.FakeStorage{
		GetStub: func() *sql.DB {
			return db
		},
	}

	projects := buildDummyProjects(3)
	rows := sqlmock.NewRows([]string{"id", "name", "repo_name", "repo_owner"})
	for _, project := range projects {
		rows.AddRow(project.ID, project.Name, project.RepoName, project.RepoOwner)
	}
	mock.ExpectQuery("^SELECT (.+) FROM \"projects\";$").WillReturnRows(rows)

	c := controller{
		storage: &fakeStorage,
	}

	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/v1/projects").Expect().
		Status(http.StatusOK).
		JSON().Array().
		Equal(projects)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListProjectsReturnsEmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fakeStorage := storagefakes.FakeStorage{
		GetStub: func() *sql.DB {
			return db
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "repo_name", "repo_owner"})
	mock.ExpectQuery("^SELECT (.+) FROM \"projects\";$").WillReturnRows(rows)

	c := controller{
		storage: &fakeStorage,
	}

	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/v1/projects").Expect().
		Status(http.StatusOK).
		JSON().Array().
		Empty()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListProjectsReturnsInternalServerError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fakeStorage := storagefakes.FakeStorage{
		GetStub: func() *sql.DB {
			return db
		},
	}

	c := controller{
		storage: &fakeStorage,
	}

	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/v1/projects").Expect().
		Status(http.StatusInternalServerError)
}

func buildDummyProjects(count int64)  []models.Project {
	projects := make([]models.Project, 0)
	var i int64 = 0
	for ; i < count ; i++  {
		projects = append(projects, models.Project{
			ID: i,
			Name: fmt.Sprint("Name",i),
			RepoName: fmt.Sprint("RepoName",i),
			RepoOwner: fmt.Sprint("RepoOwner",i),
		})
	}
	return projects
}
