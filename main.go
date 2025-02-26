package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"project1/controllers"
	"project1/db"
)

func main() {
	db.InitDB("host=localhost user=postgres password=2413050505 dbname=Spy_Cat_Agency port=1234 sslmode=disable")

	e := echo.New()

	e.Use(middleware.Logger())

	cats := e.Group("/cats")
	cats.POST("", controllers.CreateCat)
	cats.DELETE("/:id", controllers.DeleteCatByID)
	cats.PUT("/:id", controllers.UpdateSalaryCatByID)
	cats.GET("", controllers.GetAllCats)
	cats.GET("/:id", controllers.GetCatByID)

	missions := e.Group("/missions")
	missions.POST("", controllers.CreateMissionAndTargets)
	missions.DELETE("/:id", controllers.DeleteMissionByID)
	missions.PUT("/complete/:id", controllers.MarkMissionAsCompleteByID)
	missions.GET("/assign-cat", controllers.AssignCatToMissionByID)
	missions.GET("", controllers.GetAllMissions)
	missions.GET("/:id", controllers.GetMissionByID)

	targets := e.Group("/targets")
	targets.PUT("/complete/:id", controllers.MarkTargetAsCompleteByID)
	targets.PUT("/update-note/:id", controllers.UpdateNotesTargetByID)
	targets.PUT("/delete-from-mission/:id", controllers.DeleteTargetFromMissionByID)
	targets.PUT("/add-to-mission/", controllers.AddTargetToMissionByID)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
