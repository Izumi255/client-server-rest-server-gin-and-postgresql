package handlers

import (
	"net/http"

	"incident-tracker/database"
	"incident-tracker/dto"
	"incident-tracker/models"

	"github.com/gin-gonic/gin"
)

// GetIncidentsHandler godoc
// @Summary Отримання списку всіх інцидентів
// @Tags incidents
// @Produce json
// @Success 200 {array} models.Incident
// @Failure 500 {object} map[string]string
// @Router /incidents [get]
func GetIncidentsHandler(c *gin.Context) {
	var incidents []models.Incident
	if err := database.DB.Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch incidents"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

// GetIncidentByIDHandler godoc
// @Summary Отримання конкретного інциденту за ID
// @Tags incidents
// @Produce json
// @Param id path int true "ID інциденту"
// @Success 200 {object} models.Incident
// @Failure 404 {object} map[string]string
// @Router /incidents/{id} [get]
func GetIncidentByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var incident models.Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// CreateIncidentHandler godoc
// @Summary Створення нового технічного інциденту
// @Tags incidents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateIncidentDTO true "Дані інциденту"
// @Success 201 {object} models.Incident
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /incidents [post]
func CreateIncidentHandler(c *gin.Context) {
	var input dto.CreateIncidentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID uint
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uint); ok {
			userID = uid
		}
	}

	incident := models.Incident{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
		UserID:      userID,
	}

	if err := database.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incident"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// UpdateIncidentStatusHandler godoc
// @Summary Оновлення статусу інциденту
// @Tags incidents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID інциденту"
// @Param input body dto.UpdateStatusDTO true "Новий статус"
// @Success 200 {object} models.Incident
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /incidents/{id}/status [patch]
func UpdateIncidentStatusHandler(c *gin.Context) {
	id := c.Param("id")
	var incident models.Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	var input dto.UpdateStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&incident).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update incident status"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// DeleteIncidentHandler godoc
// @Summary Видалення інциденту за ID
// @Tags incidents
// @Security BearerAuth
// @Param id path int true "ID інциденту"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /incidents/{id} [delete]
func DeleteIncidentHandler(c *gin.Context) {
	id := c.Param("id")
	var incident models.Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	if err := database.DB.Delete(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incident"})
		return
	}

	c.Status(http.StatusNoContent)
}
