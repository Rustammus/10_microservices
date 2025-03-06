package rest

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"microservices/task_5/server/internal/models"
	"microservices/task_5/server/internal/schemas"
	"net/http"
	"strconv"
)

func (s *API) initUserController(r *httprouter.Router) {
	r.POST("/users", s.userRegister)
	r.GET("/users/:id", s.userFind)
	r.PUT("/users/:id", s.userUpdate)
	r.DELETE("/users/:id", s.userDelete)
}

func (s *API) userRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userReg := schemas.UserRegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read body err")
		return
	}

	id, err := s.s.User.Register(r.Context(), models.User{
		Name:     userReg.Name,
		Email:    userReg.Email,
		Password: userReg.Password,
	})
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "service err")
		return
	}

	WriteResponse(w, 201, schemas.UserRegisterResponse{ID: id}, "created")
}

func (s *API) userFind(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getIdFromParams(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read param err")
		return
	}

	user, err := s.s.User.FindById(r.Context(), id)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "service err")
		return
	}

	WriteResponse(w, 200, schemas.UserFindByIdResponse{Name: user.Name, Email: user.Email}, "found")
}

func (s *API) userUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getIdFromParams(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read param err")
		return
	}

	userUpdate := schemas.UserUpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read body err")
		return
	}

	err = s.s.User.Update(r.Context(), models.User{
		ID:       id,
		Name:     userUpdate.Name,
		Email:    userUpdate.Email,
		Password: userUpdate.Password,
	})
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "service err")
		return
	}

	WriteResponse(w, 200, struct{}{}, "updated")
}

func (s *API) userDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := getIdFromParams(ps)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "read param err")
		return
	}

	err = s.s.User.DeleteById(r.Context(), id)
	if err != nil {
		WriteResponseErr(w, http.StatusBadRequest, err, "service err")
		return
	}

	WriteResponse(w, 200, struct{}{}, "deleted")
}

func getIdFromParams(ps httprouter.Params) (int64, error) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	return id, err
}
