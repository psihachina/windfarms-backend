package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/psihachina/windfarms-backend/models"
)

func (h *Handler) getAllModels(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	models, err := h.services.Models.GetAll(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, models)
}

func (h *Handler) getModel(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	model, err := h.services.Models.GetByID(userID, windfarmID, modelID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) getModelsMap(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	model, err := h.services.Models.GetMapData(userID, windfarmID, modelID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) generateModel(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	var inputModel models.Model

	if err := c.BindJSON(&inputModel); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	model, err := h.services.Models.GenerateModel(userID, windfarmID, modelID, inputModel)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) createNewModel(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	var inputModel models.Model

	if err := c.BindJSON(&inputModel); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	model, err := h.services.Models.CreateModel(userID, windfarmID, inputModel)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model)
}

func (h *Handler) updateModel(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	var input models.UpdateModelInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Models.Update(userID, windfarmID, modelID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteModel(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	err = h.services.Models.Delete(userID, windfarmID, modelID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteModelTurbine(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	turbineModelID := c.Param("turbine_id")
	if turbineModelID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid model id param")
	}

	err := h.services.Models.DeleteTurbine(modelID, turbineModelID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
