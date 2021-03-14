package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/psihachina/windfarms-backend/models"
)

func (h *Handler) createTurbine(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	var inputTurbine models.Turbine

	var inputOutputs models.Outputs

	if err := c.ShouldBindBodyWith(&inputTurbine, binding.JSON); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ShouldBindBodyWith(&inputOutputs, binding.JSON); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Turbines.Create(userID, inputTurbine, inputOutputs)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllTurbines(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	winds, err := h.services.Turbines.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, winds)
}

func (h *Handler) getTurbineID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	turbineID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	item, err := h.services.Turbines.GetByID(userID, turbineID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateTurbine(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateTurbineInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Turbines.Update(userID, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteTurbine(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	turbineID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	err = h.services.Turbines.Delete(userID, turbineID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
