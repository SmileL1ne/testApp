package handlers

import (
	"encoding/json"
	"net/http"
	"testApp/internal/entity"

	"github.com/julienschmidt/httprouter"
)

func (r *routes) getAllUsers(w http.ResponseWriter, req *http.Request) {
	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("pageSize")

	paginatedUserData, err := r.service.User.GetAll(req.Context(), pageStr, pageSizeStr)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&paginatedUserData); err != nil {
		r.serverError(w, req, err)
		return
	}
}

func (r *routes) saveUser(w http.ResponseWriter, req *http.Request) {
	var user entity.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		r.clientError(w, http.StatusBadRequest)
		return
	}

	status, err := r.service.User.Save(req.Context(), user)
	if err != nil {
		r.identifyStatus(w, req, status, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User has been saved successfully."))
}

func (r *routes) updateUser(w http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())
	idStr := params.ByName("id")

	var user entity.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		r.logger.Error("decoding updated user info", "error", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	status, err := r.service.User.Update(req.Context(), idStr, user)
	if err != nil {
		r.identifyStatus(w, req, status, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User has been updated successfully."))
}

func (r *routes) deleteUser(w http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())
	idStr := params.ByName("id")

	status, err := r.service.User.Delete(req.Context(), idStr)
	if err != nil {
		r.identifyStatus(w, req, status, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User has been deleted successfully."))
}
