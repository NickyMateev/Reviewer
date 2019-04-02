package project

import (
	"database/sql"
	"encoding/json"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
)

type controller struct {
	storage storage.Storage
}

func (c *controller) createProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new project")

	decoder := json.NewDecoder(r.Body)
	project := models.Project{}
	err := decoder.Decode(&project)
	if err != nil {
		log.Println("Error decoding project payload:", err)
		web.WriteResponse(w, http.StatusBadRequest, web.ErrorResponse{Error: "decoding error"})
		return
	}

	err = validateProject(&project)
	if err != nil {
		log.Println("Validation error:", err)
		web.WriteResponse(w, http.StatusBadRequest, web.ErrorResponse{Error: err.Error()})
		return
	}

	err = project.Insert(r.Context(), c.storage.Get(), boil.Infer())
	if err != nil {
		log.Println("Error creating new project:", err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	web.WriteResponse(w, http.StatusCreated, project)
}

func (c *controller) getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]
	log.Println("Getting project with id", projectID)

	project, err := models.Projects(qm.Where("id = ?", projectID)).One(r.Context(), c.storage.Get())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Missing project:", err)
			web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing project"})
		} else {
			log.Printf("Error getting project with id %v: %v\n", projectID, err)
			web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		}
		return
	}
	web.WriteResponse(w, http.StatusOK, project)
}

func (c *controller) listProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all projects")
	projects, err := models.Projects().All(r.Context(), c.storage.Get())
	if err != nil {
		log.Println("Error getting projects:", err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	if len(projects) == 0 {
		projects = []*models.Project{}
	}

	web.WriteResponse(w, http.StatusOK, projects)
}

func (c *controller) deleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]
	log.Println("Deleting project with id", projectID)

	rows, err := models.Projects(qm.Where("id = ?", projectID)).DeleteAll(r.Context(), c.storage.Get())
	if err != nil {
		log.Printf("Error deleting project with id %v: %v", projectID, err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	if rows == 0 {
		log.Println("Missing project:", err)
		web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing project"})
		return
	}

	web.WriteResponse(w, http.StatusNoContent, struct{}{})
}

func validateProject(project *models.Project) error {
	if project.Name == "" {
		return errors.New("name is missing")
	}
	if project.RepoName == "" {
		return errors.New("repo_name is missing")
	}
	if project.RepoOwner == "" {
		return errors.New("repo_owner is missing")
	}
	return nil
}
