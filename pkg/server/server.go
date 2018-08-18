package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ocrease/vboxanalyser/pkg/file"
	"github.com/ocrease/vboxanalyser/pkg/vbo"
)

type Server struct {
	router   *mux.Router
	fs       file.Service
	analyser *vbo.Analyser
	port     int
}

func NewServer(port int) *Server {
	s := Server{router: mux.NewRouter(), fs: new(file.Explorer), analyser: new(vbo.Analyser), port: port}

	s.router.HandleFunc("/api/directory", s.directoryList).Methods("GET")
	s.router.HandleFunc("/api/analyse", s.analyseDirectory).Methods("GET")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))

	return &s
}

func (s *Server) Start() {
	log.Printf("Listening on port %v\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), handlers.CORS()(s.router)))
}

func (s *Server) directoryList(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	json.NewEncoder(w).Encode(s.fs.GetDirectory(path))
}

func (s *Server) analyseDirectory(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	summaries := make([]vbo.FileSummary, 0)
	s.analyser.AnalyseDirectory(path, func(fs vbo.FileSummary) {
		summaries = append(summaries, fs)
	})
	json.NewEncoder(w).Encode(summaries)
}
