package pullrequest

import (
	"database/sql"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
)

type controller struct {
	storage storage.Storage
}

func (c *controller) getPullRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pullRequestID := vars["id"]
	log.Println("Getting pull request with id", pullRequestID)

	pullRequest, err := models.PullRequests(
		qm.Where("id = ?", pullRequestID),
		qm.Load("Author"),
		qm.Load("Approvers"),
		qm.Load("Commenters"),
		qm.Load("Idlers")).One(r.Context(), c.storage.Get())

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Missing pull request:", err)
			web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing pull request"})
		} else {
			log.Printf("Error getting pull request with id %v: %v\n", pullRequestID, err)
			web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		}
		return
	}

	showActivity := r.URL.Query().Get("details")
	if showActivity == "true" {
		log.Println("Enriching response with activity information for pull request with id", pullRequestID)
		enrichedPullRequest := struct {
			models.PullRequest
			Author     string   `json:"author"`
			Approvers  []string `json:"approvers"`
			Commenters []string `json:"commenters"`
			Idlers     []string `json:"idlers"`
		}{
			PullRequest: *pullRequest,
			Author:      pullRequest.R.Author.Username,
			Approvers:   getUserStrings(pullRequest.R.Approvers),
			Commenters:  getUserStrings(pullRequest.R.Commenters),
			Idlers:      getUserStrings(pullRequest.R.Idlers),
		}
		web.WriteResponse(w, http.StatusOK, enrichedPullRequest)
		return
	}

	web.WriteResponse(w, http.StatusOK, pullRequest)
}

func (c *controller) listPullRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all pull requests")
	pullRequests, err := models.PullRequests().All(r.Context(), c.storage.Get())
	if err != nil {
		log.Println("Error getting pull requests:", err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	if len(pullRequests) == 0 {
		pullRequests = []*models.PullRequest{}
	}

	web.WriteResponse(w, http.StatusOK, pullRequests)
}

func getUserStrings(users models.UserSlice) []string {
	result := make([]string, 0)
	for _, u := range users {
		result = append(result, u.Username)
	}
	return result
}
