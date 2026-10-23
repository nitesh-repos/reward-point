package server

import (
	"net/http"
)

func (serv *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /members", serv.handleCreateMember)
	mux.HandleFunc("GET /members/{id}", serv.handleGetMember)
	mux.HandleFunc("POST /rewards", serv.handleCreateRewards)
	mux.HandleFunc("GET /members/{id}/rewards", serv.handleGetMemberRewards)
}

func (serv *Server) handleCreateMember(w http.ResponseWriter, r *http.Request) {

}

func (serv *Server) handleGetMember(w http.ResponseWriter, r *http.Request) {

}

func (serv *Server) handleCreateRewards(w http.ResponseWriter, r *http.Request) {

}

func (serv *Server) handleGetMemberRewards(w http.ResponseWriter, r *http.Request) {

}
