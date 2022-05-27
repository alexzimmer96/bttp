package main

import (
	"bttp"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/example", bttp.Handle(exampleHandler))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := bttp.ListenGracefully(srv)
	if err != nil {
		log.Printf("error while serving: %s", err)
	}
}

type ExampleRequest struct {
	Name string `json:"name"`
}

func exampleHandler(r *http.Request) bttp.Response {
	var req ExampleRequest

	// bttp.DecodeBody returns a bttp.Body indicating a BadRequest to the user
	// if the request could not have been parsed successfully.
	if ok, resp := bttp.DecodeBody(r, &req); !ok {
		return *resp
	}

	return bttp.Ok(map[string]string{
		"message": fmt.Sprintf("Hello %s", req.Name),
	})
}
