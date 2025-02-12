package controllers

import (
	"github.com/gin-gonic/gin"

)

func (controller *ControllerImplementation) Report(c *gin.Context) {
	controller.Usecases.Report()
	c.Status(200)
}
