package utils

import (
	"encoding/json"
	"log"
)

func LogReqBody(req interface{}, tag string) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("[%s] %s", tag, string(jsonBytes))
	}
}
