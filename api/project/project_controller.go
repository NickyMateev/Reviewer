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

	db := c.storage.Get()
	tx, txErr := db.BeginTx(r.Context(), nil)
	if txErr != nil {
		log.Printf("Unable to begin create transaction")
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	err = project.Insert(r.Context(), tx, boil.Infer())
	if err != nil {
		txErr = tx.Rollback()
		if txErr != nil {
			log.Printf("Could not rollback create transaction properly")
		}
		log.Println("Error creating new project:", err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Could not commit create transaction properly")
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


	db := c.storage.Get()
	tx, txErr := db.BeginTx(r.Context(), nil)
	if txErr != nil {
		log.Printf("Unable to begin delete transaction")
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	project, err := models.Projects(qm.Where("id = ?", projectID)).One(r.Context(), tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("Could not rollback delete transaction properly")
		}

		if txErr == nil && err == sql.ErrNoRows {
			log.Println("Missing project:", err)
			web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing project"})
		} else {
			log.Printf("Error getting project with id %v: %v\n", projectID, err)
			web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		}
		return
	}

	_, err = project.Delete(r.Context(), tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			log.Printf("Could not rollback delete transaction properly")
		}
		log.Printf("Error deleting project with id %v: %v\n", projectID, err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	txErr = tx.Commit()
	if txErr != nil {
		log.Printf("Could not commit delete transaction properly")
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
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
