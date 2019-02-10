package project

import (
	"database/sql"
	"encoding/json"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
)

type controller struct {
	db *sql.DB
}

func (c controller) createProject(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	project := models.Project{}
	err := decoder.Decode(&project)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	project.Insert(r.Context(), c.db, boil.Infer())

	json.NewEncoder(w).Encode(project)
}

func (c controller) getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	project, err := models.Projects(qm.Where("id = ?", projectID)).One(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(project)
}

func (c controller) listProject(w http.ResponseWriter, r *http.Request) {
	projects, err := models.Projects().All(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(projects)
}

func (c controller) deleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["id"]

	project, err := models.Projects(qm.Where("id = ?", projectID)).DeleteAll(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(project)
}
