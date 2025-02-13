package repositories

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yesetoda/kushena/infrastructures/analytics"
)

// Report2 generates all reports concurrently for daily, weekly, monthly, and yearly periods.
// It generates the extended report and saves multiple CSV files in a folder structure "Reports/orders/<period>".
func (repo *MongoRepository) Report2(interval string) {
	cols := map[string]mongo.Collection{
		"orders":     *repo.OrderCollection,
		"attendance": *repo.AttendanceCollection,
	}
	endDate := time.Now()

	log.Println("📊 Starting Kushena Business Analytics Report generation...")

	var wg sync.WaitGroup
	wg.Add(1)
	switch strings.ToLower(interval) {
	case "daily":
		before := endDate.AddDate(0, 0, -1)
		go func() {
			defer wg.Done()
			log.Println("Generating Daily Reports...")
			report, err := analytics.GenerateExtendedReport(cols, before, endDate)
			if err != nil {
				log.Printf("Error generating daily report: %v", err)
				return
			}
			if err := saveExtendedReportCSV(report, "daily"); err != nil {
				log.Printf("Error saving daily CSV: %v", err)
			}
			log.Println("Daily Reports complete.")
		}()
	case "weekly":
		before := endDate.AddDate(0, 0, -7)
		go func() {
			defer wg.Done()
			log.Println("Generating Weekly Reports...")
			report, err := analytics.GenerateExtendedReport(cols, before, endDate)
			if err != nil {
				log.Printf("Error generating weekly report: %v", err)
				return
			}
			if err := saveExtendedReportCSV(report, "weekly"); err != nil {
				log.Printf("Error saving weekly CSV: %v", err)
			}
			log.Println("Weekly Reports complete.")
		}()
	case "monthly":
		before := endDate.AddDate(0, -1, 0)
		go func() {
			defer wg.Done()
			log.Println("Generating Monthly Reports...")
			report, err := analytics.GenerateExtendedReport(cols, before, endDate)
			if err != nil {
				log.Printf("Error generating monthly report: %v", err)
				return
			}
			if err := saveExtendedReportCSV(report, "monthly"); err != nil {
				log.Printf("Error saving monthly CSV: %v", err)
			}
			log.Println("Monthly Reports complete.")
		}()
	case "yearly":
		before := endDate.AddDate(-1, 0, 0)
		go func() {
			defer wg.Done()
			log.Println("Generating Yearly Reports...")
			report, err := analytics.GenerateExtendedReport(cols, before, endDate)
			if err != nil {
				log.Printf("Error generating yearly report: %v", err)
				return
			}
			if err := saveExtendedReportCSV(report, "yearly"); err != nil {
				log.Printf("Error saving yearly CSV: %v", err)
			}
			log.Println("Yearly Reports complete.")
		}()
	default:
		log.Printf("Interval %s not recognized.", interval)
		wg.Done()
	}
	wg.Wait()
	log.Println("✅ All reports have been generated.")
}

func (repo *MongoRepository) DailyReport2() {
	repo.Report2("daily")
}

func (repo *MongoRepository) WeeklyReport2() {
	repo.Report2("weekly")
}

func (repo *MongoRepository) MonthlyReport2() {
	repo.Report2("monthly")
}

func (repo *MongoRepository) YearlyReport2() {
	repo.Report2("yearly")
}

// getReportDir2 returns the directory path for reports based on the period,
// creating a base folder "Reports/orders/<period>".
func getReportDir2(period string) string {
	baseDir := "Reports/"
	subDir := strings.ToLower(period)
	dir := filepath.Join(baseDir, subDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Failed to create directory %s: %v", dir, err)
	}
	return dir
}

// saveExtendedReportCSV writes multiple CSV files (one per section) into the folder "Reports/orders/<period>".
func saveExtendedReportCSV(report *analytics.ExtendedReport, period string) error {
	dir := getReportDir2(period)

	// 1. Order Metrics
	orderFileName := filepath.Join(dir, fmt.Sprintf("order_metrics_%s.csv", period))
	orderFile, err := os.Create(orderFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", orderFileName, err)
	}
	defer orderFile.Close()
	orderWriter := csv.NewWriter(orderFile)
	defer orderWriter.Flush()
	orderHeader := []string{"TotalOrders", "TotalRevenue", "AverageOrderValue", "MedianOrderValue", "MinOrderValue", "MaxOrderValue", "OrderCompletionRate"}
	if err := orderWriter.Write(orderHeader); err != nil {
		return fmt.Errorf("failed to write order header: %v", err)
	}
	om := report.OrderMetrics
	orderRow := []string{
		fmt.Sprintf("%d", om.TotalOrders),
		fmt.Sprintf("%.2f", om.TotalRevenue),
		fmt.Sprintf("%.2f", om.AverageOrderValue),
		fmt.Sprintf("%.2f", om.MedianOrderValue),
		fmt.Sprintf("%.2f", om.MinOrderValue),
		fmt.Sprintf("%.2f", om.MaxOrderValue),
		fmt.Sprintf("%.2f", om.OrderCompletionRate),
	}
	if err := orderWriter.Write(orderRow); err != nil {
		return fmt.Errorf("failed to write order row: %v", err)
	}

	// 2. Employee Metrics
	empFileName := filepath.Join(dir, fmt.Sprintf("employee_metrics_%s.csv", period))
	empFile, err := os.Create(empFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", empFileName, err)
	}
	defer empFile.Close()
	empWriter := csv.NewWriter(empFile)
	defer empWriter.Flush()
	empHeader := []string{"EmployeeID", "OrdersProcessed", "TotalSales", "AverageOrderValue", "AverageProcessingTime", "ProcessingTimeStdDev", "TotalWorkDuration"}
	if err := empWriter.Write(empHeader); err != nil {
		return fmt.Errorf("failed to write employee header: %v", err)
	}
	for _, emp := range report.EmployeeMetrics {
		empRow := []string{
			emp.EmployeeID,
			fmt.Sprintf("%d", emp.OrdersProcessed),
			fmt.Sprintf("%.2f", emp.TotalSales),
			fmt.Sprintf("%.2f", emp.AverageOrderValue),
			emp.AverageProcessingTime,
			emp.ProcessingTimeStdDev,
			emp.TotalWorkDuration,
		}
		if err := empWriter.Write(empRow); err != nil {
			return fmt.Errorf("failed to write employee row: %v", err)
		}
	}

	// 3. Item Metrics
	itemFileName := filepath.Join(dir, fmt.Sprintf("item_metrics_%s.csv", period))
	itemFile, err := os.Create(itemFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", itemFileName, err)
	}
	defer itemFile.Close()
	itemWriter := csv.NewWriter(itemFile)
	defer itemWriter.Flush()
	itemHeader := []string{"Type", "ItemID", "TotalQuantity"}
	if err := itemWriter.Write(itemHeader); err != nil {
		return fmt.Errorf("failed to write item header: %v", err)
	}
	// Foods
	for _, item := range report.ItemMetrics.BestSellingFoods {
		row := []string{"Food", item.ItemId, fmt.Sprintf("%.2f", item.TotalQuantity)}
		if err := itemWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write food row: %v", err)
		}
	}
	// Drinks
	for _, item := range report.ItemMetrics.BestSellingDrinks {
		row := []string{"Drink", item.ItemId, fmt.Sprintf("%.2f", item.TotalQuantity)}
		if err := itemWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write drink row: %v", err)
		}
	}

	// 4. Sales Trend
	salesFileName := filepath.Join(dir, fmt.Sprintf("sales_trend_%s.csv", period))
	salesFile, err := os.Create(salesFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", salesFileName, err)
	}
	defer salesFile.Close()
	salesWriter := csv.NewWriter(salesFile)
	defer salesWriter.Flush()
	salesHeader := []string{"Month", "Revenue", "Orders"}
	if err := salesWriter.Write(salesHeader); err != nil {
		return fmt.Errorf("failed to write sales header: %v", err)
	}
	months := make([]string, 0, len(report.SalesTrend.MonthlyRevenue))
	for m := range report.SalesTrend.MonthlyRevenue {
		months = append(months, m)
	}
	sort.Strings(months)
	for _, m := range months {
		row := []string{
			m,
			fmt.Sprintf("%.2f", report.SalesTrend.MonthlyRevenue[m]),
			fmt.Sprintf("%d", report.SalesTrend.MonthlyOrders[m]),
		}
		if err := salesWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write sales row: %v", err)
		}
	}
	if err := salesWriter.Write([]string{"GrowthRate", fmt.Sprintf("%.2f", report.SalesTrend.GrowthRate), ""}); err != nil {
		return fmt.Errorf("failed to write growth rate: %v", err)
	}

	// 5. Cohort Analysis
	cohortFileName := filepath.Join(dir, fmt.Sprintf("cohort_analysis_%s.csv", period))
	cohortFile, err := os.Create(cohortFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", cohortFileName, err)
	}
	defer cohortFile.Close()
	cohortWriter := csv.NewWriter(cohortFile)
	defer cohortWriter.Flush()
	cohortHeader := []string{"CohortLabel", "OrderCount", "TotalRevenue", "AverageOrderValue"}
	if err := cohortWriter.Write(cohortHeader); err != nil {
		return fmt.Errorf("failed to write cohort header: %v", err)
	}
	for _, c := range report.CohortAnalysis {
		row := []string{
			c.CohortLabel,
			fmt.Sprintf("%d", c.OrderCount),
			fmt.Sprintf("%.2f", c.TotalRevenue),
			fmt.Sprintf("%.2f", c.AverageOrderValue),
		}
		if err := cohortWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write cohort row: %v", err)
		}
	}

	// 6. Employee Efficiency Ranking
	effFileName := filepath.Join(dir, fmt.Sprintf("employee_efficiency_ranking_%s.csv", period))
	effFile, err := os.Create(effFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", effFileName, err)
	}
	defer effFile.Close()
	effWriter := csv.NewWriter(effFile)
	defer effWriter.Flush()
	effHeader := []string{"EmployeeID", "EfficiencyRatio"}
	if err := effWriter.Write(effHeader); err != nil {
		return fmt.Errorf("failed to write efficiency header: %v", err)
	}
	for _, eff := range report.EmployeeEfficiencyRanking {
		row := []string{
			eff.EmployeeID,
			fmt.Sprintf("%.2f", eff.EfficiencyRatio),
		}
		if err := effWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write efficiency row: %v", err)
		}
	}

	// 7. Employee Revenue Ranking
	revFileName := filepath.Join(dir, fmt.Sprintf("employee_revenue_ranking_%s.csv", period))
	revFile, err := os.Create(revFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", revFileName, err)
	}
	defer revFile.Close()
	revWriter := csv.NewWriter(revFile)
	defer revWriter.Flush()
	revHeader := []string{"EmployeeID", "TotalSales", "OrdersProcessed", "AverageOrderValue"}
	if err := revWriter.Write(revHeader); err != nil {
		return fmt.Errorf("failed to write revenue header: %v", err)
	}
	for _, emp := range report.EmployeeRevenueRanking {
		row := []string{
			emp.EmployeeID,
			fmt.Sprintf("%.2f", emp.TotalSales),
			fmt.Sprintf("%d", emp.OrdersProcessed),
			fmt.Sprintf("%.2f", emp.AverageOrderValue),
		}
		if err := revWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write revenue row: %v", err)
		}
	}

	log.Printf("CSV reports saved in %s:\n  %s\n  %s\n  %s\n  %s\n  %s\n  %s\n  %s",
		dir, orderFileName, empFileName, itemFileName, salesFileName, cohortFileName, effFileName, revFileName)
	return nil
}
