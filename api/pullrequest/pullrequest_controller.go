package pullrequest

import (
	"database/sql"
	"encoding/json"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
)

type controller struct {
	db *sql.DB
}

func (c controller) getPullRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pullRequestID := vars["id"]

	pullRequest, err := models.PullRequests(qm.Where("id = ?", pullRequestID)).One(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(pullRequest)
}

func (c controller) listPullRequest(w http.ResponseWriter, r *http.Request) {
	pullRequests, err := models.PullRequests().All(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(pullRequests)
}
