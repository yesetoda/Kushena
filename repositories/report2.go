package repositories

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"time"
//
//
//
//
//
//
//
//

// 	"github.com/yesetoda/kushena/infrastructures/analytics"

// )

// // Report2 generates all reports concurrently for daily, weekly, monthly, and yearly periods.
// // It generates the extended report and saves a CSV file in a folder structure "Records/<period>".
// func (repo *MongoRepository) Report2(interval string) {
// 	// cols := map[string]*mongo.Collection{
// 	// 	"orders":     repo.OrderCollection,
// 	// 	"attendance": repo.AttendanceCollection,
// 	// }
// 	endDate := time.Now()

// 	log.Println("📊 Starting Kushena Business Analytics Report generation...")

// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	switch strings.ToLower(interval) {
// 	case "daily":
// 		before := endDate.AddDate(0, 0, -1)
// 		go func() {
// 			defer wg.Done()
// 			log.Println("Generating Daily Reports...")
// 			report, err := analytics.GenerateEnhancedReport(nil, before, endDate) // pass appropriate db/cols as needed
// 			if err != nil {
// 				log.Printf("Error generating daily report: %v", err)
// 				return
// 			}
// 			if err := saveExtendedReportCSV(report, "daily"); err != nil {
// 				log.Printf("Error saving daily CSV: %v", err)
// 			}
// 			log.Println("Daily Reports complete.")
// 			log.Println(report)
// 		}()
// 	case "weekly":
// 		before := endDate.AddDate(0, 0, -7)
// 		go func() {
// 			defer wg.Done()
// 			log.Println("Generating Weekly Reports...")
// 			report, err := analytics.GenerateEnhancedReport(nil, before, endDate)
// 			if err != nil {
// 				log.Printf("Error generating weekly report: %v", err)
// 				return
// 			}
// 			if err := saveExtendedReportCSV(report, "weekly"); err != nil {
// 				log.Printf("Error saving weekly CSV: %v", err)
// 			}
// 			log.Println("Weekly Reports complete.")
// 			log.Println(report)
// 		}()
// 	case "monthly":
// 		before := endDate.AddDate(0, -1, 0)
// 		go func() {
// 			defer wg.Done()
// 			log.Println("Generating Monthly Reports...")
// 			report, err := analytics.GenerateEnhancedReport(nil, before, endDate)
// 			if err != nil {
// 				log.Printf("Error generating monthly report: %v", err)
// 				return
// 			}
// 			if err := saveExtendedReportCSV(report, "monthly"); err != nil {
// 				log.Printf("Error saving monthly CSV: %v", err)
// 			}
// 			log.Println("Monthly Reports complete.")
// 			log.Println(report)
// 		}()
// 	case "yearly":
// 		before := endDate.AddDate(-1, 0, 0)
// 		go func() {
// 			defer wg.Done()
// 			log.Println("Generating Yearly Reports...")
// 			report, err := analytics.GenerateEnhancedReport(nil, before, endDate)
// 			if err != nil {
// 				log.Printf("Error generating yearly report: %v", err)
// 				return
// 			}
// 			if err := saveExtendedReportCSV(report, "yearly"); err != nil {
// 				log.Printf("Error saving yearly CSV: %v", err)
// 			}
// 			log.Println("Yearly Reports complete.")
// 			log.Println(report)
// 		}()
// 	default:
// 		log.Printf("Interval %s not recognized.", interval)
// 		wg.Done()
// 	}
// 	wg.Wait()
// 	log.Println("✅ All reports have been generated.")
// }

// func (repo *MongoRepository) DailyReport2() {
// 	repo.Report2("daily")
// }

// func (repo *MongoRepository) WeeklyReport2() {
// 	repo.Report2("weekly")
// }

// func (repo *MongoRepository) MonthlyReport2() {
// 	repo.Report2("monthly")
// }

// func (repo *MongoRepository) YearlyReport2() {
// 	repo.Report2("yearly")
// }

// // getReportDir2 returns the directory path for reports based on the period,
// // creating a base folder "Records/<period>".
// func getReportDir2(period string) string {
// 	baseDir := "Records"
// 	subDir := strings.ToLower(period)
// 	dir := filepath.Join(baseDir, subDir)
// 	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
// 		log.Printf("Failed to create directory %s: %v", dir, err)
// 	}
// 	return dir
// }

// // saveExtendedReportCSV writes a CSV file with a summary of the extended report
// // into the folder "Records/<period>" with a timestamped filename.
// func saveExtendedReportCSV(report *analytics.ExtendedReport, period string) error {
// 	dir := getReportDir2(period)
// 	filename := filepath.Join(dir, fmt.Sprintf("extended_report_%s_%s.csv", period, time.Now().Format("20060102")))
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file %s: %v", filename, err)
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Write header
// 	header := []string{"Section", "Metric", "Value"}
// 	if err := writer.Write(header); err != nil {
// 		return fmt.Errorf("failed to write header: %v", err)
// 	}

// 	// Example: Write Attendance Summary Metrics
// 	att := report.Attendance.Summary
// 	rows := [][]string{
// 		{"Attendance", "TotalWorkHours", fmt.Sprintf("%.2f", att.TotalWorkHours)},
// 		{"Attendance", "AvgDailyAttendance", fmt.Sprintf("%.2f", att.AvgDailyAttendance)},
// 		{"Attendance", "LateArrivals", fmt.Sprintf("%d", att.LateArrivals)},
// 		{"Attendance", "EarlyDepartures", fmt.Sprintf("%d", att.EarlyDepartures)},
// 		{"Attendance", "AbsenteeismRate", fmt.Sprintf("%.2f", att.AbsenteeismRate)},
// 	}

// 	// (Extend with other sections such as Productivity, Sales, Operations, etc.)
// 	for _, row := range rows {
// 		if err := writer.Write(row); err != nil {
// 			return fmt.Errorf("failed to write row: %v", err)
// 		}
// 	}

// 	log.Printf("CSV report saved: %s", filename)
// 	return nil
// }
