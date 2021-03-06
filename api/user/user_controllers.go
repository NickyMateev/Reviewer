package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/NickyMateev/Reviewer/models"
	"github.com/NickyMateev/Reviewer/storage"
	"github.com/NickyMateev/Reviewer/web"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
)

type controller struct {
	storage storage.Storage
}

func (c *controller) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Println("Getting user with id", userID)

	user, err := models.Users(
		qm.Where("id = ?", userID),
		qm.Load("ApprovedPullRequests"),
		qm.Load("CommentedPullRequests"),
		qm.Load("IdledPullRequests")).One(r.Context(), c.storage.Get())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Missing user:", err)
			web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing user"})
		} else {
			log.Printf("Error getting user with id %v: %v\n", userID, err)
			web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		}
		return
	}

	showDetails := r.URL.Query().Get("details")
	if showDetails == "true" {
		log.Println("Enriching response with details information for user with id", userID)
		enrichedUser := struct {
			models.User
			ApprovedPullRequests  models.PullRequestSlice `json:"approved_pull_requests,omitempty"`
			CommentedPullRequests models.PullRequestSlice `json:"commented_pull_requests,omitempty"`
			IdledPullRequests     models.PullRequestSlice `json:"idled_pull_requests,omitempty"`
		}{
			User:                  *user,
			ApprovedPullRequests:  user.R.ApprovedPullRequests,
			CommentedPullRequests: user.R.CommentedPullRequests,
			IdledPullRequests:     user.R.IdledPullRequests,
		}
		web.WriteResponse(w, http.StatusOK, enrichedUser)
		return
	}

	web.WriteResponse(w, http.StatusOK, user)
}

func (c *controller) listUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all users")

	result := make([]*models.User, 0)
	users, err := models.Users().All(r.Context(), c.storage.Get())
	if err != nil {
		log.Println("Error getting users:", err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	if len(users) > 0 {
		result = users
	}

	web.WriteResponse(w, http.StatusOK, result)
}

func (c *controller) patchUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Println("Updating user with id", userID)

	decoder := json.NewDecoder(r.Body)
	reqUser := models.User{}
	err := decoder.Decode(&reqUser)
	if err != nil {
		log.Println("Error decoding user payload:", err)
		web.WriteResponse(w, http.StatusBadRequest, web.ErrorResponse{Error: "decoding error"})
		return
	}

	err = validateUser(&reqUser)
	if err != nil {
		log.Println("Validation error:", err)
		web.WriteResponse(w, http.StatusBadRequest, web.ErrorResponse{Error: err.Error()})
		return
	}

	var user *models.User
	txErr := c.storage.Transaction(r.Context(), func(context context.Context, tx *sql.Tx) error {
		var err error
		user, err = models.Users(qm.Where("id = ?", userID)).One(context, tx)
		if err != nil {
			return err
		}
		user.Metadata = reqUser.Metadata
		_, err = user.Update(context, tx, boil.Infer())
		return err
	})

	if txErr == sql.ErrNoRows {
		log.Println("Missing user:", err)
		web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing user"})
	} else if txErr != nil {
		log.Printf("Error updating user with id %v: %v\n", userID, err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
	}

	web.WriteResponse(w, http.StatusOK, user)
}

func (c *controller) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	log.Println("Deleting user with id", userID)

	rows, err := models.Users(qm.Where("id = ?", userID)).DeleteAll(r.Context(), c.storage.Get())
	if err != nil {
		log.Printf("Error deleting user with id %v: %v", userID, err)
		web.WriteResponse(w, http.StatusInternalServerError, struct{}{})
		return
	}

	if rows == 0 {
		log.Println("Missing user:", err)
		web.WriteResponse(w, http.StatusNotFound, web.ErrorResponse{Error: "missing user"})
		return
	}

	web.WriteResponse(w, http.StatusNoContent, struct{}{})
}

func validateUser(user *models.User) error {
	if user.Username != "" || user.ID != 0 || user.GithubID != 0 {
		return errors.New("only user metadata property is updatable")
	}
	return nil
}
