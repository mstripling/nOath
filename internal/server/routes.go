package server

import (
    "encoding/json"
    "log"
    "net/http"
    "fmt"
    "io"

    "github.com/mstripling/nOauth/util"
)

func (s *Server) RegisterRoutes() http.Handler {

    mux := http.NewServeMux()
    mux.HandleFunc("/", s.HelloWorldHandler)
    mux.HandleFunc("/post", s.PostHandler)

    return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
    var rawPayload util.InboundPayload

    err := json.NewDecoder(r.Body).Decode(&rawPayload)
    if err != nil {
        http.Error(w, fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest)
        return 
    }

    refreshResp, err := util.RefreshAuthToken(r, &rawPayload)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error processing request", err), http.StatusBadRequest)
        return
    }

    //post to api
    resp, err := util.Post(r, refreshResp, &rawPayload)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshalling request to JSON", err), http.StatusInternalServerError)
    }
    
    w.WriteHeader(resp.StatusCode)
    _, err = io.Copy(w, resp.Body)
    
    if err != nil {
        http.Error(w, fmt.Sprintf("Error copying response body: %s", err), http.StatusInternalServerError)
        return
    }
    
    defer resp.Body.Close()
}
