package crud

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jexroid/gopi/api"
	"github.com/jexroid/gopi/pkg/models"
	"gorm.io/gorm"
)

// Read godoc
// @Summary Get a user by ID
// @Description Retrieve a single user by their UUID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   uuid path string true "User UUID"
// @Success 200 {object} models.User "User found"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{uuid} [get]
func (db Database) Read(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	var car models.User
	result := db.DB.First(&car, "id = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(car)
}

// All godoc
// @Summary Get all users with pagination and search
// @Description Retrieve a paginated list of users with optional phone number search
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   page query int false "Page number (default: 1)"
// @Param   limit query int false "Number of items per page (default: 10, max: 100)"
// @Param   phone query string false "Phone number to search for (partial match)"
// @Success 200 {object} UsersResponse "List of users with pagination info"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /users [get]
func (db Database) All(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	phone := r.URL.Query().Get("phone")

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 100 {
				limit = 100
			} else {
				limit = l
			}
		}
	}

	offset := (page - 1) * limit

	var users []models.User
	query := db.DB

	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	var total int64
	countQuery := db.DB.Model(&models.User{})
	if phone != "" {
		countQuery = countQuery.Where("phone LIKE ?", "%"+phone+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := query.Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	response := api.UsersResponse{
		Users: users,
		Pagination: api.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update godoc
// @Summary Update a user
// @Description Update an existing user's information
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   uuid path string true "User UUID"
// @Param   user body models.User true "User data to update"
// @Success 200 {object} models.User "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{uuid} [put]
func (db Database) Update(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	var updateCar models.User
	if err := json.NewDecoder(r.Body).Decode(&updateCar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var car models.User
	result := db.DB.First(&car, "id = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	db.DB.Model(&car).Updates(updateCar)
	json.NewEncoder(w).Encode(car)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by their UUID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   uuid path string true "User UUID"
// @Success 204 "User deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{uuid} [delete]
func (db Database) Delete(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	result := db.DB.Delete(&models.User{}, "id = ?", uuid)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
