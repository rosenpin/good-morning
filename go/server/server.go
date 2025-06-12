package server

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.com/rosenpin/good-morning/provider"
)

// Server serves the cached image every day
type Server struct {
	provider       provider.ImageProvider
	maxDailyReload uint8
	dailyCounter   uint8
}

// New creates a new instance of the server
func New(provider provider.ImageProvider, maxDailyReload int) Server {
	fmt.Println("mdr", maxDailyReload)
	return Server{provider, uint8(maxDailyReload), 0}
}

// Start starts the server
func (s Server) Start() {
	go func() {
		for true {
			fmt.Println("Resetting daily counter")
			s.dailyCounter = 0
			time.Sleep(24 * time.Hour)
		}
	}()

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		if s.dailyCounter >= s.maxDailyReload {
			w.WriteHeader(403)
			w.Write([]byte("You can no longer request reloads today. Daily quota exceeded"))
			return
		}

		s.dailyCounter++

		_, err := s.provider.ForceReload()
		if err != nil {
			w.WriteHeader(501)
			w.Write([]byte(fmt.Sprint("Error reload image", err)))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Image reloaded successfully"))
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		image, err := s.provider.Provide()
		if err != nil {
			w.WriteHeader(501)
			w.Write([]byte(err.Error()))
			return
		}
		defer image.Close()

		w.WriteHeader(200)
		io.Copy(w, image)
		return
	})

	http.ListenAndServe(":8181", nil)
}
