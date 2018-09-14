package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hamcha/tg"
)

func webhook(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Re-encode request to ensure conformity
	var update tg.APIUpdate
	err := json.NewDecoder(req.Body).Decode(&update)
	if err != nil {
		log.Println("[webhook] Received incorrect request: " + err.Error())
		return
	}

	data, err := json.Marshal(tg.BrokerUpdate{
		Data:     update,
		Callback: nil,
	})
	if err != nil {
		log.Println("[webhook] Cannot re-encode json (??) : " + err.Error())
		return
	}

	broadcast(string(data))
}
