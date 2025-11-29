package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type PatientInfo struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
	Gender      string `json:"gender"`
	Illness     string `json:"illness"`
	Email       string `json:"email"`
}

type PatientStore struct {
	mu       sync.RWMutex
	patients map[string]PatientInfo
}

var store *PatientStore

func init() {
	store = &PatientStore{
		patients: make(map[string]PatientInfo),
	}

	// Initialize with sample data
	samplePatients := []PatientInfo{
		{
			Name:        "Nobody Knows",
			DateOfBirth: "1985-03-15",
			Gender:      "Male",
			Illness:     "Hypertension",
			Email:       "nobody.knows@email.com",
		},
		{
			Name:        "Johnson Fake",
			DateOfBirth: "1990-07-22",
			Gender:      "Female",
			Illness:     "Type 2 Diabetes",
			Email:       "johnson.fake@email.com",
		},
		{
			Name:        "Michael Chen",
			DateOfBirth: "1978-11-08",
			Gender:      "Male",
			Illness:     "Asthma",
			Email:       "michael.chen@email.com",
		},
		{
			Name:        "Emily Lor",
			DateOfBirth: "1995-02-14",
			Gender:      "Female",
			Illness:     "Migraine",
			Email:       "emily.lor@email.com",
		},
	}

	for _, patient := range samplePatients {
		store.patients[strings.ToLower(patient.Name)] = patient
	}
}

func getPatientByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := strings.ToLower(vars["name"])

	store.mu.RLock()
	patient, exists := store.patients[name]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	var patient PatientInfo

	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if patient.Name == "" || patient.DateOfBirth == "" || patient.Email == "" {
		http.Error(w, "Name, date of birth, and email are required", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	store.patients[strings.ToLower(patient.Name)] = patient
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func listAllPatients(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	patients := make([]PatientInfo, 0, len(store.patients))
	for _, patient := range store.patients {
		patients = append(patients, patient)
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/patients", listAllPatients).Methods("GET")
	r.HandleFunc("/api/patients", createPatient).Methods("POST")
	r.HandleFunc("/api/patients/{name}", getPatientByName).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// CORS middleware
	handler := enableCORS(r)

	port := "8080"
	log.Printf("Starting patient service on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
