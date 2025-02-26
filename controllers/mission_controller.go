package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"project1/db"
	"project1/models"
	"strconv"
)

func CreateMissionAndTargets(c echo.Context) error {
	req := struct {
		Targets []struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		} `json:"targets"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if len(req.Targets) < 1 || len(req.Targets) > 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Targets must be 1-3"})
	}

	mission := models.Mission{
		Complete: false,
	}
	if err := db.DB.Create(&mission).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error of saving to database, error message: " + err.Error()})
	}

	for _, t := range req.Targets {
		target := models.Target{
			Model:     gorm.Model{},
			Name:      t.Name,
			Country:   t.Country,
			Notes:     "",
			Complete:  false,
			MissionID: &mission.ID,
		}
		if err := db.DB.Create(&target).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error of saving to database, error message: " + err.Error()})
		}
	}

	var loadedMission models.Mission
	if err := db.DB.Preload("Targets").First(&loadedMission, mission.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error loading mission: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, loadedMission)
}

func DeleteMissionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var mission models.Mission

	if err := db.DB.First(&mission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if mission.CatID != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete mission with assigned cat"})
	}

	if err := db.DB.Delete(&mission).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Mission has been successfully removed"})
}

func MarkMissionAsCompleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var mission models.Mission

	if err := db.DB.Preload("Targets").Preload("Cat").First(&mission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	mission.Complete = true

	if err := db.DB.Save(&mission).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update mission"})
	}
	return c.JSON(http.StatusOK, mission)
}

func AssignCatToMissionByID(c echo.Context) error {
	catID, err := strconv.Atoi(c.QueryParam("cat_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong cat ID"})
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

	var cat models.Cat

	if err := db.DB.First(&cat, catID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if mission.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot assign cat to completed mission"})
	}

	var existingMission models.Mission
	if err := db.DB.Where("cat_id = ?", catID).First(&existingMission).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cat is already assigned to another mission"})
	}

	catIDabs := uint(catID)
	mission.CatID = &catIDabs

	if err := db.DB.Save(&mission).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update mission"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Cat has been successfully assigned to mission"})
}

func GetAllMissions(c echo.Context) error {
	var missions []models.Mission
	if err := db.DB.Preload("Targets").Preload("Cat").Find(&missions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, missions)
}

func GetMissionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var mission models.Mission
	if err := db.DB.Preload("Targets").Preload("Cat").First(&mission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, mission)
}
