package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"twitter_clone/data"
)

type KeyProduct struct{}

type ProfileHandler struct {
	logger *log.Logger
	// NoSQL: injecting product repository
	repo *data.ProfileRepo
}

//Injecting the logger makes this code much more testable.
func NewProfilesHandler(l *log.Logger, r *data.ProfileRepo) *ProfileHandler {
	return &ProfileHandler{l, r}
}

func (p *ProfileHandler) GetAllProfiles(rw http.ResponseWriter, h *http.Request) {
	profiles, err := p.repo.GetAll()
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if profiles == nil {
		return
	}

	err = profiles.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *ProfileHandler) GetProfileById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	profile, err := p.repo.GetById(id)
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if profile == nil {
		http.Error(rw, "Patient with given id not found", http.StatusNotFound)
		p.logger.Printf("Patient with id: '%s' not found", id)
		return
	}

	err = profile.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *ProfileHandler) GetProfilesByName(rw http.ResponseWriter, h *http.Request) {
	name := h.URL.Query().Get("name")

	profiles, err := p.repo.GetByName(name)
	if err != nil {
		p.logger.Print("Database exception: ", err)
	}

	if profiles == nil {
		return
	}

	err = profiles.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		p.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (p *ProfileHandler) PostProfile(rw http.ResponseWriter, h *http.Request) {
	profile := h.Context().Value(KeyProduct{}).(*data.Profile)
	p.repo.Insert(profile)
	rw.WriteHeader(http.StatusCreated)
}

func (p *ProfileHandler) PatchProfile(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	profile := h.Context().Value(KeyProduct{}).(*data.Profile)

	p.repo.Update(id, profile)
	rw.WriteHeader(http.StatusOK)
}

func (p *ProfileHandler) DeleteProfile(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	p.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (p *ProfileHandler) MiddlewarePatientDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		profile := &data.Profiles{}
		err := profile.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, profile)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (p *ProfileHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
