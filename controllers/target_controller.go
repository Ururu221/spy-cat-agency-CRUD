package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"project1/db"
	"project1/models"
	"strconv"
	"strings"
)

func MarkTargetAsCompleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var target models.Target

	if err := db.DB.First(&target, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	target.Complete = true

	if err := db.DB.Save(&target).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update target"})
	}

	var mission models.Mission
	if err := db.DB.Preload("Targets").First(&mission, target.MissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target is not in any mission"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, mission)
}

func UpdateNotesTargetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	newNotes := strings.TrimSpace(c.QueryParam("notes"))
	if newNotes == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong notes"})
	}

	var target models.Target

	if err := db.DB.First(&target, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var mission models.Mission

	if err := db.DB.First(&mission, target.MissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target is not in any mission"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if mission.Complete || target.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot update notes for completed target or mission"})
	}

	sb := new(strings.Builder)
	sb.WriteString(target.Notes)
	sb.WriteString("\n")
	sb.WriteString(newNotes)

	target.Notes = sb.String()

	if err := db.DB.Save(&target).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update target"})
	}
	return c.JSON(http.StatusOK, mission)
}

func DeleteTargetFromMissionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var target models.Target

	if err := db.DB.First(&target, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if target.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete completed target"})
	}

	var mission models.Mission
	if err := db.DB.First(&mission, target.MissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target is not in any mission"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if mission.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete target from completed mission"})
	}
	if len(mission.Targets) == 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete last target from mission"})
	}

	target.MissionID = nil

	if err := db.DB.Save(&target).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update target"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Target has been successfully removed from mission"})
}

func AddTargetToMissionByID(c echo.Context) error {
	targetID, err := strconv.Atoi(c.QueryParam("target_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong target ID"})
	}

	missionID, err := strconv.Atoi(c.QueryParam("mission_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong mission ID"})
	}

	var mission models.Mission

	if err := db.DB.First(&mission, missionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if len(mission.Targets) == 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "The mission already has the maximum number of targets (3)"})
	}

	var target models.Target

	if err := db.DB.First(&target, targetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if mission.Complete || target.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot add target to completed mission of completed target"})
	}

	if target.MissionID != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Target is already in a mission"})
	}

	missionIDabs := uint(missionID)
	target.MissionID = &missionIDabs

	if err := db.DB.Save(&target).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update target"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Target has been successfully added to mission"})
}
