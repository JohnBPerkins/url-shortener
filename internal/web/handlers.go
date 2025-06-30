package web

import (
	"encoding/json"
	"log"
	"net/http"

	pb "github.com/JohnBPerkins/url-shortener/gen"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Code string `json:"code"`
}

type ResolveRequest struct {
	Code string `json:"code"`
}

type ResolveResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewShrinkHandler(svc pb.ShortenerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("→  HTTP %s %s\n", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
			return
		}

		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		if req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "URL is required"})
			return
		}

        grpcReq := &pb.ShortenRequest{Url: req.URL}
        grpcResp, err := svc.Shrink(r.Context(), grpcReq)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(ShortenResponse{
            Code: grpcResp.GetCode(),
        })
	}
}

func NewResolveHandler(svc pb.ShortenerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}