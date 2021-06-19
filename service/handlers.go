package service

import (
	"github.com/emanonsm/software_project/service/handlers/patient"
	"github.com/emanonsm/software_project/service/handlers/root"
	"github.com/emanonsm/software_project/service/handlers/search"
	"github.com/emanonsm/software_project/service/handlers/visit"
	"net/http"
)

// ShowRoot - обертка функции ShowRoot
func (s *Service) ShowRoot(w http.ResponseWriter, r *http.Request) {
	root.ShowRoot(w, r, s.db)
}

// CreatePatient - обертка функции CreatePatient
func (s *Service) CreatePatient(w http.ResponseWriter, r *http.Request) {
	patient.CreatePatient(w, r, s.db)
}

// CreateVisit - обертка функции CreatePatient
func (s *Service) CreateVisit(w http.ResponseWriter, r *http.Request) {
	visit.CreateVisit(w, r, s.db)
}

// Search - обертка функции Search
func (s *Service) Search(w http.ResponseWriter, r *http.Request) {
	search.Search(w, r, s.db)
}

// Visits - обертка функции Visits
func (s *Service) Visits(w http.ResponseWriter, r *http.Request) {
	search.Visits(w, r, s.db)
}

// ShowVisit - обертка функции ShowVisit
func (s *Service) ShowVisit(w http.ResponseWriter, r *http.Request) {
	search.ShowVisit(w, r, s.db)
}


