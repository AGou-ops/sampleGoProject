package main

import (
	"net/http"
)


func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// json := `{"status": "available", "environment": "%s", "port": %d}`
	// json = fmt.Sprintf(json, app.config.env, app.config.port)
	data := map[string]string {
		"status": "available",
		"environment": app.config.env,
		"version": version,
		"afterMarshal": "somethinghere",
	}
	if err := app.formatJson(w, http.StatusOK, data, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, "formatJson internal error", http.StatusInternalServerError)
		return
	}
}
