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
			if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"}); encodeErr != nil {
				log.Printf("handlers.go: failed to write 405 JSON: %v", encodeErr)
			}
			return
		}

		var req ShortenRequest
		if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"}); encodeErr != nil {
				log.Printf("handlers.go: failed to write 400 JSON: %v", encodeErr)
			}
			return
		}

		if req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{Error: "URL is required"}); encodeErr != nil {
				log.Printf("handlers.go: failed to write 400 JSON: %v", encodeErr)
			}
			return
		}

        grpcReq := &pb.ShortenRequest{Url: req.URL}
        grpcResp, err := svc.Shorten(r.Context(), grpcReq)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()}); encodeErr != nil {
				log.Printf("handlers.go: failed to write 500 JSON: %v", encodeErr)
			}
            return
        }

        w.WriteHeader(http.StatusOK)
        encodeErr := json.NewEncoder(w).Encode(ShortenResponse{
            Code: grpcResp.GetCode(),
        })
		if encodeErr != nil {
			log.Printf("handlers.go: failed to write 200 JSON: %v", encodeErr)
		}
	}
}

func NewResolveHandler(svc pb.ShortenerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}