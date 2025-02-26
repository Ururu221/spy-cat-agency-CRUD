package controllers

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"project1/db"
	"project1/models"
	"strconv"
	"strings"
)

type Breed struct {
	Name string `json:"name"`
}

func checkBreed(breedName string) (bool, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var breeds []Breed

	err = json.NewDecoder(resp.Body).Decode(&breeds)
	if err != nil {
		return false, err
	}

	for _, b := range breeds {
		if b.Name == breedName {
			return true, nil
		}
	}
	return false, nil
}

func CreateCat(c echo.Context) error {
	var cat models.Cat
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data format"})
	}

	valid, err := checkBreed(cat.Breed)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Breed verification error, error message: " + err.Error()})
	}
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Breed not found"})
	}

	if err := db.DB.Create(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error of saving to database, error message: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, cat)
}

func DeleteCatByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var cat models.Cat

	if err := db.DB.First(&cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if err := db.DB.Delete(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cat has been successfully removed"})
}

func UpdateSalaryCatByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	newSalary, err := strconv.ParseFloat(strings.TrimSpace(c.QueryParam("salary")), 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong salary"})
	}

	var cat models.Cat

	if err := db.DB.First(&cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	cat.Salary = newSalary

	if err := db.DB.Save(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update salary"})
	}

	return c.JSON(http.StatusOK, cat)
}

func GetAllCats(c echo.Context) error {
	var cats []models.Cat
	if err := db.DB.Find(&cats).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, cats)
}

func GetCatByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong ID"})
	}

	var cat models.Cat

	if err := db.DB.First(&cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, cat)
}
