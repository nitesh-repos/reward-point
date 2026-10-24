package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reward-points/internal/database"
	"reward-points/internal/models"
	"strconv"
)

func (serv *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /members", serv.handleCreateMember)
	mux.HandleFunc("GET /members/{id}", serv.handleGetMember)
	mux.HandleFunc("POST /rewards", serv.handleCreateRewards)
	mux.HandleFunc("GET /members/{id}/rewards", serv.handleGetMemberRewards)
}

var ErrBadRequest = errors.New("bad request")

func (serv *Server) handleCreateMember(w http.ResponseWriter, r *http.Request) {
	var mr MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		serv.writeError(w, fmt.Errorf("invalid JSON body: %w", ErrBadRequest))
		return
	}

	member, err := serv.store.CreateMember(models.Member{
		Name:  mr.Name,
		Email: mr.Email,
	})
	if err != nil {
		serv.writeError(w, err)
		return
	}

	serv.logger.Info("member created", "id", member.MemberID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)

}

func (serv *Server) handleGetMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		serv.writeError(w, fmt.Errorf("invalid member id \"%s\" in path: %w", idStr, ErrBadRequest))
		return
	}

	memberInfo, err := serv.store.GetMember(id)
	if err != nil {
		serv.writeError(w, err)
		return
	}

	serv.logger.Info("member fetched successfully", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberInfo)
}

func (serv *Server) handleCreateRewards(w http.ResponseWriter, r *http.Request) {
	var rr RewardRequest
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		serv.writeError(w, fmt.Errorf("invalid JSON body: %w", ErrBadRequest))
		return
	}

	reward, err := serv.store.CreateMemberReward(models.Reward{
		MemberID:    rr.MemberID,
		PointTypeID: rr.PointTypeID,
		Points:      rr.Points,
		Description: rr.Description,
	})
	if err != nil {
		serv.writeError(w, err)
		return
	}

	serv.logger.Info("member reward created", "id", reward.RewardID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reward)
}

func (serv *Server) handleGetMemberRewards(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		serv.writeError(w, fmt.Errorf("invalid member id \"%s\" in path: %w", idStr, ErrBadRequest))
		return
	}

	rewards, err := serv.store.GetMemberReward(id)
	if err != nil {
		serv.writeError(w, err)
		return
	}

	serv.logger.Info("member rewards fetched successfully", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rewards)
}

func (serv *Server) writeError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal server error"

	switch {
	case errors.Is(err, ErrBadRequest):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, database.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, database.ErrConflict):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, database.ErrInvalid):
		status, msg = http.StatusUnprocessableEntity, err.Error()
	}

	if status >= 500 {
		serv.logger.Error("request failed", "status", status, "error", err)
	} else {
		serv.logger.Warn("request rejected", "status", status, "error", err)
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
