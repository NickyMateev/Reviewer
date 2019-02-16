package server

import (
	"github.com/NickyMateev/Reviewer/web"
	"github.com/NickyMateev/Reviewer/web/webfakes"
	"net/http"
	"reflect"
	"testing"
)

const fakeURL = "/v1/fake"

func fakeCreate(w http.ResponseWriter, r *http.Request) {}
func fakeGetAll(w http.ResponseWriter, r *http.Request) {}
func fakeGet(w http.ResponseWriter, r *http.Request)    {}
func fakePatch(w http.ResponseWriter, r *http.Request)  {}
func fakeDelete(w http.ResponseWriter, r *http.Request) {}

func TestBuildRouter(t *testing.T) {
	fakeCtrl := &webfakes.FakeController{}
	fakeRoutes := []web.Route{
		{
			Path:       fakeURL,
			Method:     http.MethodPost,
			HandleFunc: fakeCreate,
		},
		{
			Path:       fakeURL,
			Method:     http.MethodGet,
			HandleFunc: fakeGetAll,
		},
		{
			Path:       fakeURL + "/{id}",
			Method:     http.MethodGet,
			HandleFunc: fakeGet,
		},
		{
			Path:       fakeURL + "/{id}",
			Method:     http.MethodPatch,
			HandleFunc: fakePatch,
		},
		{
			Path:       fakeURL + "/{id}",
			Method:     http.MethodDelete,
			HandleFunc: fakeDelete,
		},
	}
	fakeCtrl.RoutesReturns(fakeRoutes)

	fakeAPI := &webfakes.FakeAPI{}
	fakeAPI.ControllersReturns([]web.Controller{fakeCtrl})

	router := buildRouter(fakeAPI)

	reflectRouter := reflect.ValueOf(*router)
	routes := reflectRouter.FieldByName("routes")

	if routes.Len() != len(fakeRoutes) {
		t.Errorf("route number mismatch; actual=%v; expected=%v", routes.Len(), len(fakeRoutes))
	}
}
