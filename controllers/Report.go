package controllers

import (
	"github.com/gin-gonic/gin"

)

func (controller *ControllerImplementation) DailyReport(c *gin.Context) {
	controller.Usecases.DailyReport()
	c.Status(200)
}

func (controller *ControllerImplementation) WeeklyReport(c *gin.Context) {
    controller.Usecases.WeeklyReport()
    c.Status(200)
}

func (controller *ControllerImplementation) MonthlyReport(c *gin.Context) {
    controller.Usecases.MonthlyReport()
    c.Status(200)
}

func (controller *ControllerImplementation) YearlyReport(c *gin.Context) {
	controller.Usecases.YearlyReport()
    c.Status(200)
}