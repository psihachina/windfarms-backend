package handler

import (
	"net/http"
	"strconv"

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

func (h *Handler) getChartData(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	from := c.Query("from")
	to := c.Query("to")
	height, err := strconv.Atoi(c.Query("height"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	winds, err := h.services.Winds.GetWindForChart(userID, windfarmID, from, to, height)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, winds)
}

func (h *Handler) getTableData(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if windfarmID == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid windfarm id param")
	}

	winds, err := h.services.Winds.GetWindForTable(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, winds)
}
