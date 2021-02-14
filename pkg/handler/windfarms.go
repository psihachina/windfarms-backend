package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/psihachina/windfarms-backend/models"
)

func (h *Handler) createWindfarm(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	var input models.Windfarm
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Windfarms.Create(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllWindfarms(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	items, err := h.services.Windfarms.GetAll(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) getWindfarmByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	item, err := h.services.Windfarms.GetByID(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateWindfarm(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateWindfarmInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Windfarms.Update(userID, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteWindfarm(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	windfarmID := c.Param("id")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	err = h.services.Windfarms.Delete(userID, windfarmID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
