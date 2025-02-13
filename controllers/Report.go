package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

)

func (controller *ControllerImplementation) DailyReport(c *gin.Context) {
	data, err := controller.Usecases.DailyReport()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, j)
}
func (controller *ControllerImplementation) WeeklyReport(c *gin.Context) {
	data, err := controller.Usecases.WeeklyReport()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, j)
}
func (controller *ControllerImplementation) MonthlyReport(c *gin.Context) {
	data, err := controller.Usecases.MonthlyReport()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, j)
}
func (controller *ControllerImplementation) YearlyReport(c *gin.Context) {
	data, err := controller.Usecases.YearlyReport()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var j map[string]interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, j)
}
