package api

import (
	"encoding/json"
	"net/http"

	"github.com/paul-schwendenman/magic-log-ui/internal/config"
)

var configPath string

func init() {
	configPath = config.GetConfigPath()
}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cfg)
}

func SaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	var newCfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := newCfg.Validate(); err != nil {
		http.Error(w, "Invalid config: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.SaveToFile(configPath, &newCfg); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
