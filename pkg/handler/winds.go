package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createWinds(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	err = h.services.Winds.Create(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"Status": "OK",
	})
}

func (h *Handler) getAllWinds(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	winds, err := h.services.Winds.GetAll(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, winds)
}
