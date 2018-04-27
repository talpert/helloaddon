package api

import (
	"net/http"

	"github.com/InVisionApp/rye"
)

func (a *API) homeHandler(rw http.ResponseWriter, r *http.Request) {
	rye.WriteJSONStatus(rw, "Oh, hello there!", "Refer to README.md for talpert/helloaddon API usage", http.StatusOK)
}

func (a *API) versionHandler(rw http.ResponseWriter, r *http.Request) {
	rye.WriteJSONStatus(rw, "version", "talpert/helloaddon "+a.Version, http.StatusOK)
}
