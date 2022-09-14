package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// json := `{"status": "available", "environment": "%s", "port": %d}`
	// json = fmt.Sprintf(json, app.config.env, app.config.port)
	// data := map[string]string {
	// 	"status": "available",
	// 	"environment": app.config.env,
	// 	"version": version,
	// 	"afterMarshal": "somethinghere",
	// }
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	if err := app.writeJson(w, *r, http.StatusOK, data, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
}
