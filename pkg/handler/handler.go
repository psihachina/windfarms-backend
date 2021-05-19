package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/psihachina/windfarms-backend/pkg/service"
	cors "github.com/rs/cors/wrapper/gin"
)

//Handler ...
type Handler struct {
	services *service.Service
}

//NewHandler ...
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

//InitRoutes ...
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.AllowAll())

	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.signUp)
		auth.POST("/sing-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		windfarms := api.Group("/windfarms")
		{
			windfarms.POST("/", h.createWindfarm)
			windfarms.GET("/", h.getAllWindfarms)
			windfarms.GET("/:id", h.getWindfarmByID)
			windfarms.PUT("/:id", h.updateWindfarm)
			windfarms.DELETE("/:id", h.deleteWindfarm)

			winds := windfarms.Group("/:id/winds")
			{
				winds.POST("/", h.createWinds)
				winds.GET("/", h.getAllWinds)
				winds.GET("/chart", h.getChartData)
				winds.GET("/table", h.getTableData)
			}

			models := windfarms.Group("/:id/models")
			{
				models.POST("/", h.createModel)
				models.GET("/", h.getAllModels)
				models.GET("/:model_id/map", h.getModelsMap)
				models.DELETE("/:model_id", h.deleteModel)
				//models.PUT("/:model_id", h.updateModel)
				models.GET("/:model_id", h.getModel)
			}
		}

		tubines := api.Group("/turbines")
		{
			tubines.POST("/", h.createTurbine)
			tubines.GET("/", h.getAllTurbines)
			tubines.GET("/:id", h.getTurbineID)
			tubines.PUT("/:id", h.updateTurbine)
			tubines.DELETE("/:id", h.deleteTurbine)
		}

	}

	return router
}
