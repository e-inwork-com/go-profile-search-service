package api

import (
	"io"
	"net/http"
)

func (app *Application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(
		"GET",
		app.Config.SolrURL+"/api/collections/"+
			app.Config.SolrProfile+"/select?"+
			r.URL.Query().Encode(),
		nil)

	client := &http.Client{}
	res, _ := client.Do(req)

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
	w.Write([]byte(body))
}
