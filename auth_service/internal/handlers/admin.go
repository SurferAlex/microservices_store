package handlers

import (
	"auth_service/internal/middleware"
	"auth_service/internal/repository/psql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type assignRoleReq struct {
	Role string `json:"role"`
}

func AssignRole(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	has, err := psql.UserHasPermission(userID, "manage_system")
	if err != nil || !has {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := mux.Vars(r)["id"]
	userID, err = strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req assignRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Role == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := psql.AssignRoleToUser(userID, req.Role); err != nil {
		http.Error(w, "assign failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET обработчик для проерки ролей
func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Проверка, что только админ может смотреть роли
	has, err := psql.UserHasPermission(userID, "manage_system")
	if err != nil || !has {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	targetUserID, err := strconv.Atoi(idStr)
	if err != nil || targetUserID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	roles, err := psql.GetUserRoles(targetUserID)
	if err != nil {
		http.Error(w, "failed to get roles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": targetUserID,
		"roles":   roles,
	})
}
