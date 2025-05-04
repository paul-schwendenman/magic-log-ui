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

	if errs := newCfg.Validate(); len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string][]string{
			"errors": errorStrings(errs),
		})
		return
	}

	if err := config.SaveToFile(configPath, &newCfg); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func errorStrings(errs []error) []string {
	out := make([]string, len(errs))
	for i, e := range errs {
		out[i] = e.Error()
	}
	return out
}
