package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Status  int    `json:"status"`
}

func Error(w http.ResponseWriter, err error, status int, msg ...string) {
	errResp := ErrorResponse{
		Message: msg[0],
		Error:   false,
		Status:  status,
	}

	if err != nil {
		log.Println(err)
		errResp.Error = true
	}

	err = json.NewEncoder(w).Encode(errResp)
	if err != nil {
		log.Println(err)
		return
	}
}
