package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/yesetoda/Kushena/controllers"
	"github.com/yesetoda/Kushena/infrastructures/auth_services"
	"github.com/yesetoda/Kushena/repositories"
	"github.com/yesetoda/Kushena/router"
	"github.com/yesetoda/Kushena/usecases"
)

func scheduleReports(scheduler *gocron.Scheduler, repo repositories.RepositoryInterface, reportTime, reportDay, reportMonth string) {
	// Daily Report
	scheduler.Every(1).Day().At(reportTime).Do(repo.DailyReport)

	// Weekly Report (on specified day)
	if reportDay != "" {
		weekDays := map[string]time.Weekday{
			"Sunday": time.Sunday, "Monday": time.Monday, "Tuesday": time.Tuesday, "Wednesday": time.Wednesday,
			"Thursday": time.Thursday, "Friday": time.Friday, "Saturday": time.Saturday,
		}

		weekDay, exists := weekDays[reportDay]
		if exists {
			_, err := scheduler.Every(1).Week().Weekday(weekDay).At(reportTime).Do(repo.WeeklyReport)
			if err != nil {
				fmt.Println("Error scheduling weekly report:", err)
			}
		} else {
			fmt.Println("Invalid week day:", reportDay)
		}
	}

	// Monthly Report (on 1st of each month)
	scheduler.Every(1).Month(1).At(reportTime).Do(repo.MonthlyReport)

	// Yearly Report (on January 1st)
	scheduler.Every(1).Day().At(reportTime).Do(func() {
		now := time.Now()
		if now.Month().String() == reportMonth && now.Day() == 1 {
			fmt.Println("📆 Running Yearly Report...")
			repo.YearlyReport()
		}
	})
}

func main() {
	fmt.Println("Hello, Kushena!")
	fmt.Println("Welcome to our restaurant!")

	// Initialize Dependencies
	repo := repositories.NewRepo()
	scheduler := gocron.NewScheduler(time.Local)

	after1min := time.Now().Add(1 * time.Minute).Format("15:04")
	fmt.Println("Scheduling report at:", after1min)

	// Schedule Reports
	scheduleReports(scheduler, repo, after1min, "Wednesday", "January")

	// Start the Scheduler in a separate goroutine
	scheduler.StartAsync()

	// Start Gin Router (Blocking Call)
	usecase := usecases.NewUsecase(repo)
	controller := controllers.NewController(usecase)
	auth := auth_services.NewAuthController(usecase)
	router := router.NewGinRoute(controller, auth)

	router.Run() // This will keep running indefinitely
}
