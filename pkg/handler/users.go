package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllUsers(c *gin.Context) {

	users, err := h.services.Users.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) deleteUser(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid email param")
	}

	err := h.services.Users.Delete(email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) confirmUser(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid email param")
	}

	err := h.services.Users.Confirm(email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
