package version

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Get()); err != nil {
		slog.Default().Error("encoding version response", "error", err)
	}
}
