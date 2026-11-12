package handler

import (
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"time"
)

func GetMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil data profil user berdasarkan user.ID dari database
	profile, err := lmodel.GetUserByID(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Tangkap input JSON dari request body
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		// Tambahkan field lain jika diperlukan sesuai struct model user kalian
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: Update data user ke database berdasarkan user.ID
	if err := lmodel.UpdateUser(user.ID, req.Username, req.Email); err != nil {
		http.Error(w, "Failed to update user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data updated successfully",
		"username": req.Username,
	})
}

func GetMyAuditLog(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil riwayat audit log khusus milik user tersebut dari database
	logs, err := lmodel.GetAuditLogsByUserID(user.ID) // Atau berdasarkan user.Username tergantung skema
	if err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user.Username,
		"total": len(logs),
		"logs":  logs,
	})
}

func ExportMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Ambil profil, proyek, dan audit log user untuk portabilitas data UU PDP
	profile, _ := lmodel.GetUserByID(user.ID)
	projects, _ := lmodel.GetProjectsByUserID(user.ID)
	logs, _ := lmodel.GetAuditLogsByUserID(user.ID)

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"profile":     profile,
		"projects":    projects,
		"audit_logs":  logs,
	}

	w.Header().Set("Content-Disposition", "attachment; filename=my_data_export.json")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(export)
}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// TODO: Hapus akun dan seluruh data terkait secara permanen (Right to be Forgotten)
	if err := lmodel.DeleteUserAccount(user.ID); err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your account and all data have been permanently deleted",
	})
}
