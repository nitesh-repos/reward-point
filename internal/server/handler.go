package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reward-points/internal/models"
	"strconv"
)

func (serv *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /members", serv.handleCreateMember)
	mux.HandleFunc("GET /members/{id}", serv.handleGetMember)
	mux.HandleFunc("POST /rewards", serv.handleCreateRewards)
	mux.HandleFunc("GET /members/{id}/rewards", serv.handleGetMemberRewards)
}

func (serv *Server) handleCreateMember(w http.ResponseWriter, r *http.Request) {
	var mr MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		serv.logger.Error("could not decode payload", "error", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON body"})
		return
	}

	member, err := serv.store.CreateMember(models.Member{
		Name:  mr.Name,
		Email: mr.Email,
	})
	if err != nil {
		serv.logger.Error("create member failed", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "could not create member"})
		return
	}

	serv.logger.Info("member created", "id", member.MemberID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(member)

}

func (serv *Server) handleGetMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		serv.logger.Error("invalid member id", "id", idStr)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid member id in path parameter"})
		return
	}

	memberInfo, err := serv.store.GetMember(id)
	if err != nil {
		serv.logger.Error("member fetch error", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("could not fetch member for id %d", id)})
		return
	}

	serv.logger.Info("member fetched successfully", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberInfo)
}

func (serv *Server) handleCreateRewards(w http.ResponseWriter, r *http.Request) {
	var rr RewardRequest
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		serv.logger.Error("could not decode payload", "error", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON body"})
		return
	}

	reward, err := serv.store.CreateMemberReward(models.Reward{
		MemberID:    rr.MemberID,
		PointTypeID: rr.PointTypeID,
		Points:      rr.Points,
		Description: rr.Description,
	})
	if err != nil {
		serv.logger.Error("create member reward failed", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "could not create member reward"})
		return
	}

	serv.logger.Info("member reward created", "id", reward.RewardID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reward)
}

func (serv *Server) handleGetMemberRewards(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		serv.logger.Error("invalid member id", "id", idStr)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid member id in path parameter"})
		return
	}

	rewards, err := serv.store.GetMemberReward(id)
	if err != nil {
		serv.logger.Error("member rewards fetch error", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("could not fetch member rewards for id %d", id)})
		return
	}

	serv.logger.Info("member rewards fetched successfully", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rewards)
}
