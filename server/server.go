package server

import (
	"fmt"
	"io"
	"net/http"

	"gitlab.com/rosenpin/good-morning/provider"
)

// Server serves the cached image every day
type Server struct {
	provider provider.ImageProvider
}

// New creates a new winstace of the server
func New(provider provider.ImageProvider) Server {
	return Server{provider}
}

// Start starts the server
func (s Server) Start() {
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		image, err := s.provider.Provide()
		if err != nil {
			w.WriteHeader(501)
			w.Write([]byte(err.Error()))
			return
		}
		defer image.Close()

		w.WriteHeader(200)
		io.Copy(w, image)
	})
	http.ListenAndServe(":8080", nil)
}
