package httpj

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func RenderErr(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
}
