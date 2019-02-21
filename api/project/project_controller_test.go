package project

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/storage/storagefakes"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

var projectColumns = []string{models.ProjectColumns.ID, models.ProjectColumns.Name, models.ProjectColumns.RepoName, models.ProjectColumns.RepoOwner}
const insertQuery = "INSERT INTO \"projects\""
const selectQuery = "SELECT (.+) FROM \"projects\""
const deleteQuery = "DELETE FROM \"projects\""

func TestCreateProjectReturnsResponse(t *testing.T) {
	return
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	project := buildDummyProjectWithID(1)
	sqlMock.ExpectExec(insertQuery).WithArgs(project.Name, project.RepoName, project.RepoOwner).WillReturnResult(sqlmock.NewResult(1,1))
	sqlMock.ExpectCommit()

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.createProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.POST(web.ProjectsURL).
		WithJSON(struct {
		Name string `json:"name"`
		RepoName string `json:"repo_name"`
		RepoOwner string `json:"repo_owner"`
	}{
		Name: project.Name,
		RepoName: project.RepoName,
		RepoOwner: project.RepoOwner,
	}).
		Expect().
		Status(http.StatusCreated)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProjectReturnsResponse(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	projectToFind := buildDummyProjectWithID(0)
	rowToFind := sqlmock.NewRows(projectColumns)
	rowToFind.AddRow(projectToFind.ID, projectToFind.Name, projectToFind.RepoName, projectToFind.RepoOwner)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rowToFind)

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.getProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", projectToFind.ID).Expect().
		Status(http.StatusOK).
		JSON().
		Equal(projectToFind)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProjectReturnsNotFound(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows(projectColumns)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rows)

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.getProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusNotFound)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProjectReturnsInternalServerError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf(""))

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.getProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusInternalServerError)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListProjectsReturnsResponse(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	projects := buildDummyProjects(3)
	rows := sqlmock.NewRows(projectColumns)
	for _, project := range projects {
		rows.AddRow(project.ID, project.Name, project.RepoName, project.RepoOwner)
	}
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rows)

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL).Expect().
		Status(http.StatusOK).
		JSON().Array().
		Equal(projects)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListProjectsReturnsEmptyResult(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows(projectColumns)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rows)

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL).Expect().
		Status(http.StatusOK).
		JSON().Array().
		Empty()

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListProjectsReturnsInternalServerError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf(""))

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.listProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL).Expect().
		Status(http.StatusInternalServerError)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProjectReturnsResponse(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	projectToFind := buildDummyProjectWithID(0)
	rowToFind := sqlmock.NewRows(projectColumns)
	rowToFind.AddRow(projectToFind.ID, projectToFind.Name, projectToFind.RepoName, projectToFind.RepoOwner)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rowToFind)
	sqlMock.ExpectExec(deleteQuery).WillReturnResult(sqlmock.NewResult(1,1))
	sqlMock.ExpectCommit()

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusNoContent)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProjectReturnsNotFound(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	rowToFind := sqlmock.NewRows(projectColumns)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rowToFind)
	sqlMock.ExpectRollback()

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusNotFound)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProjectWhenTxCannotBeginReturnsInternalServerError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusInternalServerError)
}

func TestDeleteProjectOnGetQueryWhenTxCannotRollbackReturnsInternalServerError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf(""))

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusInternalServerError)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProjectOnDeleteQueryWhenTxCannotRollbackReturnsInternalServerError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	projectToFind := buildDummyProjectWithID(0)
	rowToFind := sqlmock.NewRows(projectColumns)
	rowToFind.AddRow(projectToFind.ID, projectToFind.Name, projectToFind.RepoName, projectToFind.RepoOwner)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rowToFind)

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusInternalServerError)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteProjectWhenCommitFailsReturnsInternalServerError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	projectToFind := buildDummyProjectWithID(0)
	rowToFind := sqlmock.NewRows(projectColumns)
	rowToFind.AddRow(projectToFind.ID, projectToFind.Name, projectToFind.RepoName, projectToFind.RepoOwner)
	sqlMock.ExpectQuery(selectQuery).WillReturnRows(rowToFind)
	sqlMock.ExpectExec(deleteQuery).WillReturnResult(sqlmock.NewResult(1,1))

	c := controller{
		storage: buildFakeStorage(db),
	}
	server := httptest.NewServer(web.HandleFunc(c.deleteProject))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET(web.ProjectsURL + "/{id}").WithPath("id", 1).Expect().
		Status(http.StatusInternalServerError)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateProject(t *testing.T)  {
	project := buildDummyProjectWithID(0)
	err := validateProject(&project)
	if err != nil {
		t.Fatal("an error occurred when it wasn't supposed to for valid project")
	}
}

func TestValidateProjectWithInvalidName(t *testing.T)  {
	project := buildDummyProjectWithID(0)
	project.Name = ""
	err := validateProject(&project)
	if err == nil {
		t.Fatal("an error should have occurred project with invalid name")
	}
}

func TestValidateProjectWithInvalidRepoName(t *testing.T)  {
	project := buildDummyProjectWithID(0)
	project.RepoName = ""
	err := validateProject(&project)
	if err == nil {
		t.Fatal("an error should have occurred project with invalid repo_name")
	}
}

func TestValidateProjectWithInvalidRepoOwner(t *testing.T)  {
	project := buildDummyProjectWithID(0)
	project.RepoOwner = ""
	err := validateProject(&project)
	if err == nil {
		t.Fatal("an error should have occurred project with invalid repo_owner")
	}
}

func buildFakeStorage(db *sql.DB) storage.Storage {
	return &storagefakes.FakeStorage{
		GetStub: func() *sql.DB {
			return db
		},
	}
}

func buildDummyProjects(count int64)  []models.Project {
	projects := make([]models.Project, 0)
	var i int64 = 0
	for ; i < count ; i++  {
		projects = append(projects, buildDummyProjectWithID(i))
	}
	return projects
}

func buildDummyProjectWithID(id int64)  models.Project {
	return models.Project{
		ID: id,
		Name:      fmt.Sprint(projectColumns[1], id),
		RepoName:  fmt.Sprint(projectColumns[2], id),
		RepoOwner: fmt.Sprint(projectColumns[3], id),
	}
}
