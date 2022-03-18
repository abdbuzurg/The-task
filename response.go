package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// General way of how api will give response
func RespondWithJson(w http.ResponseWriter, data interface{}, success bool) {
	dataFormat := Response{
		Success: success,
		Data:    data,
	}

	response, err := json.Marshal(dataFormat)
	if err != nil {
		fmt.Println(err.Error())
		panic("Internal server error while forming the response")
	}

	w.Write(response)
}
