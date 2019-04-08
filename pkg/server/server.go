package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ocrease/vboxanalyser/pkg/file"
	"github.com/ocrease/vboxanalyser/pkg/vbo"
)

type Server struct {
	router   *mux.Router
	fs       file.Service
	analyser *vbo.Analyser
	config   *vbo.Config
	port     int
}

func NewServer(config *vbo.Config) *Server {
	s := Server{router: mux.NewRouter(), fs: new(file.Explorer), analyser: new(vbo.Analyser), config: config}

	s.router.HandleFunc("/api/directory", s.directoryList).Methods("GET")
	s.router.HandleFunc("/api/analyse", s.analyseDirectory).Methods("GET")
	s.router.HandleFunc("/api/launch", s.launchInCT).Methods("POST")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))

	return &s
}

func (s *Server) Start() {
	log.Printf("Listening on port %v\n", s.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.config.Port), handlers.CORS()(s.router)))
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

func (s *Server) launchInCT(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = exec.Command(s.config.CTPath, string(b[:])).Run()
	if err != nil {
		fmt.Printf("Failed to launch Circuit Tools: %v\n", err)
	}
}
