package user

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

func (c controller) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	err := decoder.Decode(&user)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	user.Insert(r.Context(), c.db, boil.Infer())
	json.NewEncoder(w).Encode(user)
}

func (c controller) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := models.Users(qm.Where("id = ?", userID)).One(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (c controller) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.Users().All(r.Context(), c.db)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (c controller) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	_, err := models.Users(qm.Where("id = ?", userID)).DeleteAll(r.Context(), c.db)
	if err != nil {
		return
	}
}
