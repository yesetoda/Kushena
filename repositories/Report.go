package repositories

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yesetoda/kushena/models"

)
func getReportDir(period string) string {
	var dir string
	switch period {
	case "Daily", "daily":
		dir = filepath.Join("Reports", "daily")
	case "Weekly", "weekly":
		dir = filepath.Join("Reports", "weekly")
	case "Monthly", "monthly":
		dir = filepath.Join("Reports", "monthly")
	case "Yearly", "yearly":
		dir = filepath.Join("Reports", "yearly")
	default:
		dir = filepath.Join("Reports", "others")
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dir, err)
	}
	return dir
}

// Report generates all reports concurrently for daily, weekly, monthly, and yearly periods.
// It returns the combined JSON report as []byte along with any error.
func (repo *MongoRepository) Report(interval string) ([]byte, error) {
	endDate := time.Now()
	reportDir := getReportDir(interval)
	log.Println("📊 Starting Kushena Business Analytics Report generation...")

	var wg sync.WaitGroup
	wg.Add(1)

	// Variables to capture the result from the goroutine.
	var resultData []byte
	var resultErr error

	switch interval {
	case "daily", "Daily":
		beforeDay := endDate.AddDate(0, 0, -1)
		go func() {
			defer wg.Done()
			log.Println("Generating Daily Reports...")
			orderReport := generateOrderReport(repo.OrderCollection, beforeDay, endDate, "Daily")
			employeePerformanceReport := generateEmployeePerformanceReport(repo.AttendanceCollection, repo.OrderCollection, beforeDay, endDate, "Daily")
			operationalEfficiencyReport := generateOperationalEfficiencyReport(repo.AttendanceCollection, repo.OrderCollection, beforeDay, endDate, "Daily")
			revenueFinancialReport := generateRevenueFinancialReport(repo.OrderCollection, beforeDay, endDate, "Daily")

			combined := map[string]interface{}{
				"order_report":                  orderReport,
				"employee_performance_report":   employeePerformanceReport,
				"operational_efficiency_report": operationalEfficiencyReport,
				"revenue_financial_report":      revenueFinancialReport,
			}
			b, err := json.MarshalIndent(combined, "", "  ")
			if err != nil {
				resultErr = err
				return
			}
			resultData = b
			outfile := filepath.Join(reportDir, "combined_report_daily.json")
			if err := os.WriteFile(outfile, b, 0644); err != nil {
				log.Printf("Error writing Daily JSON report: %v", err)
			} else {
				log.Println("Daily Reports complete. Combined report saved to", outfile)
			}
		}()
	case "weekly", "Weekly":
		beforeWeek := endDate.AddDate(0, 0, -7)
		go func() {
			defer wg.Done()
			log.Println("Generating Weekly Reports...")
			orderReport := generateOrderReport(repo.OrderCollection, beforeWeek, endDate, "Weekly")
			employeePerformanceReport := generateEmployeePerformanceReport(repo.AttendanceCollection, repo.OrderCollection, beforeWeek, endDate, "Weekly")
			operationalEfficiencyReport := generateOperationalEfficiencyReport(repo.AttendanceCollection, repo.OrderCollection, beforeWeek, endDate, "Weekly")
			revenueFinancialReport := generateRevenueFinancialReport(repo.OrderCollection, beforeWeek, endDate, "Weekly")

			combined := map[string]interface{}{
				"order_report":                  orderReport,
				"employee_performance_report":   employeePerformanceReport,
				"operational_efficiency_report": operationalEfficiencyReport,
				"revenue_financial_report":      revenueFinancialReport,
			}
			b, err := json.MarshalIndent(combined, "", "  ")
			if err != nil {
				resultErr = err
				return
			}
			resultData = b
			outfile := filepath.Join(reportDir, "combined_report_weekly.json")
			if err := os.WriteFile(outfile, b, 0644); err != nil {
				log.Printf("Error writing Weekly JSON report: %v", err)
			} else {
				log.Println("Weekly Reports complete. Combined report saved to", outfile)
			}
		}()
	case "monthly", "Monthly":
		beforeMonth := endDate.AddDate(0, -1, 0)
		go func() {
			defer wg.Done()
			log.Println("Generating Monthly Reports...")
			orderReport := generateOrderReport(repo.OrderCollection, beforeMonth, endDate, "Monthly")
			employeePerformanceReport := generateEmployeePerformanceReport(repo.AttendanceCollection, repo.OrderCollection, beforeMonth, endDate, "Monthly")
			operationalEfficiencyReport := generateOperationalEfficiencyReport(repo.AttendanceCollection, repo.OrderCollection, beforeMonth, endDate, "Monthly")
			revenueFinancialReport := generateRevenueFinancialReport(repo.OrderCollection, beforeMonth, endDate, "Monthly")

			combined := map[string]interface{}{
				"order_report":                  orderReport,
				"employee_performance_report":   employeePerformanceReport,
				"operational_efficiency_report": operationalEfficiencyReport,
				"revenue_financial_report":      revenueFinancialReport,
			}
			b, err := json.MarshalIndent(combined, "", "  ")
			if err != nil {
				resultErr = err
				return
			}
			resultData = b
			outfile := filepath.Join(reportDir, "combined_report_monthly.json")
			if err := os.WriteFile(outfile, b, 0644); err != nil {
				log.Printf("Error writing Monthly JSON report: %v", err)
			} else {
				log.Println("Monthly Reports complete. Combined report saved to", outfile)
			}
		}()
	case "yearly", "Yearly":
		beforeYear := endDate.AddDate(-1, 0, 0)
		go func() {
			defer wg.Done()
			log.Println("Generating Yearly Reports...")
			orderReport := generateOrderReport(repo.OrderCollection, beforeYear, endDate, "Yearly")
			employeePerformanceReport := generateEmployeePerformanceReport(repo.AttendanceCollection, repo.OrderCollection, beforeYear, endDate, "Yearly")
			operationalEfficiencyReport := generateOperationalEfficiencyReport(repo.AttendanceCollection, repo.OrderCollection, beforeYear, endDate, "Yearly")
			revenueFinancialReport := generateRevenueFinancialReport(repo.OrderCollection, beforeYear, endDate, "Yearly")

			combined := map[string]interface{}{
				"order_report":                  orderReport,
				"employee_performance_report":   employeePerformanceReport,
				"operational_efficiency_report": operationalEfficiencyReport,
				"revenue_financial_report":      revenueFinancialReport,
			}
			b, err := json.MarshalIndent(combined, "", "  ")
			if err != nil {
				resultErr = err
				return
			}
			resultData = b
			outfile := filepath.Join(reportDir, "combined_report_yearly.json")
			if err := os.WriteFile(outfile, b, 0644); err != nil {
				log.Printf("Error writing Yearly JSON report: %v", err)
			} else {
				log.Println("Yearly Reports complete. Combined report saved to", outfile)
			}
		}()
	default:
		log.Printf("Interval %s not recognized.", interval)
		wg.Done()
	}
	wg.Wait()
	log.Println("✅ All reports have been generated.")
	return resultData, resultErr
}

// generateOrderReport produces an order report (business analytics & sales trends).
func generateOrderReport(collection *mongo.Collection, startDate, endDate time.Time, period string) []models.Order {
	log.Printf("Generating %s Order Report...\n", period)
	filter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalf("Error fetching orders for %s report: %v", period, err)
	}
	defer cursor.Close(context.TODO())

	var orders []models.Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		log.Fatalf("Error decoding orders for %s report: %v", period, err)
	}

	totalRevenue := 0.0
	totalOrders := len(orders)
	peakHours := make(map[string]int)
	itemSales := make(map[string]int)
	var orderValues []float64
	foodSales := make(map[string]int)
	drinkSales := make(map[string]int)

	for _, order := range orders {
		totalRevenue += order.TotalPrice
		orderValues = append(orderValues, order.TotalPrice)
		hourStr := fmt.Sprintf("%02d", order.CreatedAt.Hour())
		peakHours[hourStr]++
		for _, food := range order.Foods {
			key := food.FoodId.Hex()
			itemSales[key]++
			foodSales[key]++
		}
		for _, drink := range order.Drinks {
			key := drink.DrinkId.Hex()
			itemSales[key]++
			drinkSales[key]++
		}
	}

	var minVal, maxVal, avgVal float64
	if totalOrders > 0 {
		minVal = orderValues[0]
		maxVal = orderValues[0]
		sum := 0.0
		for _, v := range orderValues {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
			sum += v
		}
		avgVal = sum / float64(totalOrders)
	}

	mostSelling, leastSelling := getMostAndLeast(itemSales)
	topFoods := getTopN(foodSales, 3)
	topDrinks := getTopN(drinkSales, 3)
	var monthlyOrders map[string]int
	if period == "Yearly" {
		monthlyOrders = getMonthlyOrderCounts(orders)
	}

	fmt.Printf("\n📌 %s Order Report\n", period)
	fmt.Printf("Total Orders: %d | Total Revenue: $%.2f\n", totalOrders, totalRevenue)
	fmt.Printf("Order Value Distribution: Min=$%.2f, Max=$%.2f, Avg=$%.2f\n", minVal, maxVal, avgVal)
	fmt.Println("Peak Sales Hours:")
	var hrs []string
	for h := range peakHours {
		hrs = append(hrs, h)
	}
	sort.Strings(hrs)
	for _, h := range hrs {
		fmt.Printf("  %s:00 - %d orders\n", h, peakHours[h])
	}
	fmt.Printf("Most Selling Item: %s (%d sales)\n", mostSelling.Key, mostSelling.Value)
	fmt.Printf("Least Selling Item: %s (%d sales)\n", leastSelling.Key, leastSelling.Value)
	fmt.Println("Top Ordered Foods:", topFoods)
	fmt.Println("Top Ordered Drinks:", topDrinks)
	if period == "Yearly" {
		fmt.Println("Monthly Orders (Seasonal Demand):", monthlyOrders)
	}

	reportDir := getReportDir(period)
	saveCSV(filepath.Join(reportDir, "order_report_"+period+".csv"), []string{"Order ID", "Total Price", "Created At"}, orders)
	return orders
}

// generateEmployeePerformanceReport produces an employee performance report (attendance, sales & efficiency).
func generateEmployeePerformanceReport(attendanceCollection, orderCollection *mongo.Collection, startDate, endDate time.Time, period string) map[string]interface{} {
	log.Printf("Generating %s Employee Performance Report...\n", period)

	attendanceFilter := bson.M{"time": bson.M{"$gte": startDate, "$lt": endDate}}
	attendanceCursor, err := attendanceCollection.Find(context.TODO(), attendanceFilter)
	if err != nil {
		log.Fatalf("Error fetching attendance for %s report: %v", period, err)
	}
	defer attendanceCursor.Close(context.TODO())
	var attendances []models.Attendance
	if err := attendanceCursor.All(context.TODO(), &attendances); err != nil {
		log.Fatalf("Error decoding attendance for %s report: %v", period, err)
	}

	orderFilter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	orderCursor, err := orderCollection.Find(context.TODO(), orderFilter)
	if err != nil {
		log.Fatalf("Error fetching orders for %s report: %v", period, err)
	}
	defer orderCursor.Close(context.TODO())
	var orders []models.Order
	if err := orderCursor.All(context.TODO(), &orders); err != nil {
		log.Fatalf("Error decoding orders for %s report: %v", period, err)
	}

	employeeWorkHours := make(map[string]time.Duration)
	employeeLateCount := make(map[string]int)
	employeeAttendanceCount := make(map[string]int)
	checkInTimes := make(map[string]time.Time)
	lateThresholdHour := 9

	for _, rec := range attendances {
		empID := rec.EmployeeID.Hex()
		employeeAttendanceCount[empID]++
		if rec.Type == "in" {
			if rec.Time.Hour() >= lateThresholdHour {
				employeeLateCount[empID]++
			}
			if _, exists := checkInTimes[empID]; !exists {
				checkInTimes[empID] = rec.Time
			}
		} else if rec.Type == "out" {
			if inTime, exists := checkInTimes[empID]; exists {
				duration := rec.Time.Sub(inTime)
				employeeWorkHours[empID] += duration
				delete(checkInTimes, empID)
			}
		}
	}
	if len(employeeAttendanceCount) == 0 {
		log.Printf("No attendance records found for %s period\n", period)
	}
	absenteeThreshold := 1
	absentees := []string{}
	for empID, count := range employeeAttendanceCount {
		if count <= absenteeThreshold {
			absentees = append(absentees, empID)
		}
	}

	employeeOrders := make(map[string]int)
	employeeRevenue := make(map[string]float64)
	employeeOrderTimestamps := make(map[string][]time.Time)
	for _, order := range orders {
		empID := order.EmployeeId.Hex()
		employeeOrders[empID]++
		employeeRevenue[empID] += order.TotalPrice
		employeeOrderTimestamps[empID] = append(employeeOrderTimestamps[empID], order.CreatedAt)
	}
	employeeAvgOrderValue := make(map[string]float64)
	for empID, count := range employeeOrders {
		if count > 0 {
			employeeAvgOrderValue[empID] = employeeRevenue[empID] / float64(count)
		}
	}

	employeeProcessingTime := make(map[string]time.Duration)
	for empID, timestamps := range employeeOrderTimestamps {
		if len(timestamps) > 1 {
			sort.Slice(timestamps, func(i, j int) bool {
				return timestamps[i].Before(timestamps[j])
			})
			var totalGap time.Duration
			for i := 1; i < len(timestamps); i++ {
				totalGap += timestamps[i].Sub(timestamps[i-1])
			}
			employeeProcessingTime[empID] = totalGap / time.Duration(len(timestamps)-1)
		}
	}

	fmt.Printf("\n📊 %s Employee Performance Report\n", period)
	fmt.Printf("%-24s %-8s %-12s %-14s %-16s %-20s\n", "Employee ID", "Orders", "Revenue($)", "Avg Order($)", "Work Hours", "Late Checkins")
	for empID := range employeeAttendanceCount {
		fmt.Printf("%-24s %-8d %-12.2f %-14.2f %-16s %-20d\n",
			empID,
			employeeOrders[empID],
			employeeRevenue[empID],
			employeeAvgOrderValue[empID],
			employeeWorkHours[empID].Round(time.Minute).String(),
			employeeLateCount[empID])
	}
	fmt.Println("Absenteeism (low attendance):", absentees)
	fmt.Println("Average Order Processing Time per Employee:")
	for empID, procTime := range employeeProcessingTime {
		fmt.Printf("  %s: %s\n", empID, procTime.Round(time.Second).String())
	}

	reportDir := getReportDir(period)
	csvHeaders := []string{"Employee ID", "Orders", "Revenue", "Avg Order Value", "Work Hours", "Late Checkins"}
	var csvData [][]string
	for empID := range employeeAttendanceCount {
		row := []string{
			empID,
			strconv.Itoa(employeeOrders[empID]),
			fmt.Sprintf("%.2f", employeeRevenue[empID]),
			fmt.Sprintf("%.2f", employeeAvgOrderValue[empID]),
			employeeWorkHours[empID].Round(time.Minute).String(),
			strconv.Itoa(employeeLateCount[empID]),
		}
		csvData = append(csvData, row)
	}
	saveCSVTable(filepath.Join(reportDir, "employee_report_"+period+".csv"), csvHeaders, csvData)
	return map[string]interface{}{
		"orders":          employeeOrders,
		"revenue":         employeeRevenue,
		"avg_order":       employeeAvgOrderValue,
		"work_hours":      employeeWorkHours,
		"late_checkins":   employeeLateCount,
		"processing_time": employeeProcessingTime,
		"absentees":       absentees,
	}
}

// generateOperationalEfficiencyReport produces an operational efficiency report (utilization, idle time, etc.).
func generateOperationalEfficiencyReport(attendanceCollection, orderCollection *mongo.Collection, startDate, endDate time.Time, period string) map[string]interface{} {
	log.Printf("Generating %s Operational Efficiency Report...\n", period)
	attendanceFilter := bson.M{"time": bson.M{"$gte": startDate, "$lt": endDate}}
	attendanceCursor, err := attendanceCollection.Find(context.TODO(), attendanceFilter)
	if err != nil {
		log.Fatalf("Error fetching attendance for %s operational report: %v", period, err)
	}
	defer attendanceCursor.Close(context.TODO())
	var attendances []models.Attendance
	if err := attendanceCursor.All(context.TODO(), &attendances); err != nil {
		log.Fatalf("Error decoding attendance for %s operational report: %v", period, err)
	}

	orderFilter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	orderCursor, err := orderCollection.Find(context.TODO(), orderFilter)
	if err != nil {
		log.Fatalf("Error fetching orders for %s operational report: %v", period, err)
	}
	defer orderCursor.Close(context.TODO())
	var orders []models.Order
	if err := orderCursor.All(context.TODO(), &orders); err != nil {
		log.Fatalf("Error decoding orders for %s operational report: %v", period, err)
	}

	employeeOrders := make(map[string]int)
	employeeAttendanceCount := make(map[string]int)
	for _, order := range orders {
		empID := order.EmployeeId.Hex()
		employeeOrders[empID]++
	}
	for _, rec := range attendances {
		empID := rec.EmployeeID.Hex()
		employeeAttendanceCount[empID]++
	}

	orderToAttendance := make(map[string]float64)
	for empID, attCount := range employeeAttendanceCount {
		ord := employeeOrders[empID]
		if attCount > 0 {
			orderToAttendance[empID] = float64(ord) / float64(attCount)
		}
	}

	idleEmployees := []string{}
	for empID, attCount := range employeeAttendanceCount {
		if employeeOrders[empID] == 0 && attCount > 0 {
			idleEmployees = append(idleEmployees, empID)
		}
	}

	abnormalEmployees := []string{}
	for empID, ratio := range orderToAttendance {
		if ratio < 0.5 || ratio > 2.0 {
			abnormalEmployees = append(abnormalEmployees, empID)
		}
	}

	fmt.Printf("\n⚙️ %s Operational Efficiency Report\n", period)
	fmt.Printf("%-24s %-18s %-18s\n", "Employee ID", "Attendance Count", "Order/Attendance Ratio")
	for empID, attCount := range employeeAttendanceCount {
		fmt.Printf("%-24s %-18d %-18.2f\n", empID, attCount, orderToAttendance[empID])
	}
	fmt.Println("Idle Employees:", idleEmployees)
	fmt.Println("Abnormal Order Patterns:", abnormalEmployees)

	reportDir := getReportDir(period)
	csvHeaders := []string{"Employee ID", "Attendance Count", "Order/Attendance Ratio"}
	var csvData [][]string
	for empID, attCount := range employeeAttendanceCount {
		row := []string{
			empID,
			strconv.Itoa(attCount),
			fmt.Sprintf("%.2f", orderToAttendance[empID]),
		}
		csvData = append(csvData, row)
	}
	saveCSVTable(filepath.Join(reportDir, "operational_efficiency_"+period+".csv"), csvHeaders, csvData)

	result := map[string]interface{}{
		"employee_attendance_count": employeeAttendanceCount,
		"order_to_attendance_ratio": orderToAttendance,
		"idle_employees":            idleEmployees,
		"abnormal_employees":        abnormalEmployees,
	}
	return result
}

// generateRevenueFinancialReport produces a revenue & financial report.
func generateRevenueFinancialReport(orderCollection *mongo.Collection, startDate, endDate time.Time, period string) map[string]interface{} {
	log.Printf("Generating %s Revenue & Financial Report...\n", period)
	filter := bson.M{"created_at": bson.M{"$gte": startDate, "$lt": endDate}}
	cursor, err := orderCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalf("Error fetching orders for %s revenue report: %v", period, err)
	}
	defer cursor.Close(context.TODO())
	var orders []models.Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		log.Fatalf("Error decoding orders for %s revenue report: %v", period, err)
	}

	totalRevenue := 0.0
	orderCount := len(orders)
	var orderValues []float64
	totalFoodItems := 0
	totalDrinkItems := 0

	for _, order := range orders {
		totalRevenue += order.TotalPrice
		orderValues = append(orderValues, order.TotalPrice)
		totalFoodItems += len(order.Foods)
		totalDrinkItems += len(order.Drinks)
	}

	var minVal, maxVal, avgVal float64
	if orderCount > 0 {
		minVal = orderValues[0]
		maxVal = orderValues[0]
		sum := 0.0
		for _, val := range orderValues {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
			sum += val
		}
		avgVal = sum / float64(orderCount)
	} else {
		log.Printf("No order records found for %s Revenue Report\n", period)
	}

	totalItems := totalFoodItems + totalDrinkItems
	var foodShare, drinkShare float64
	if totalItems > 0 {
		foodShare = float64(totalFoodItems) / float64(totalItems) * 100
		drinkShare = float64(totalDrinkItems) / float64(totalItems) * 100
	}
	avgItemsPerOrder := 0.0
	if orderCount > 0 {
		avgItemsPerOrder = float64(totalItems) / float64(orderCount)
	}

	fmt.Printf("\n💰 %s Revenue & Financial Report\n", period)
	fmt.Printf("Total Revenue: $%.2f | Total Orders: %d\n", totalRevenue, orderCount)
	fmt.Printf("Order Value: Min=$%.2f, Max=$%.2f, Avg=$%.2f\n", minVal, maxVal, avgVal)
	fmt.Printf("Food vs. Drink Share: Food=%.2f%%, Drink=%.2f%%\n", foodShare, drinkShare)
	fmt.Printf("Average Items Per Order: %.2f\n", avgItemsPerOrder)

	reportDir := getReportDir(period)
	csvHeaders := []string{"Metric", "Value"}
	csvData := [][]string{
		{"Total Revenue", fmt.Sprintf("%.2f", totalRevenue)},
		{"Total Orders", strconv.Itoa(orderCount)},
		{"Min Order Value", fmt.Sprintf("%.2f", minVal)},
		{"Max Order Value", fmt.Sprintf("%.2f", maxVal)},
		{"Avg Order Value", fmt.Sprintf("%.2f", avgVal)},
		{"Food Share (%)", fmt.Sprintf("%.2f", foodShare)},
		{"Drink Share (%)", fmt.Sprintf("%.2f", drinkShare)},
		{"Avg Items Per Order", fmt.Sprintf("%.2f", avgItemsPerOrder)},
	}
	saveCSVTable(filepath.Join(reportDir, "revenue_financial_"+period+".csv"), csvHeaders, csvData)
	return map[string]interface{}{
		"total_revenue":       totalRevenue,
		"total_orders":        orderCount,
		"min_order_value":     minVal,
		"max_order_value":     maxVal,
		"avg_order_value":     avgVal,
		"food_share_percent":  foodShare,
		"drink_share_percent": drinkShare,
		"avg_items_per_order": avgItemsPerOrder,
	}
}

// ─── UTILITY FUNCTIONS ─────────────────────────────────────────────
func saveCSVTable(filename string, headers []string, rows [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create CSV file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		log.Fatalf("Failed to write CSV headers: %v", err)
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write CSV row: %v", err)
		}
	}
}

func saveCSV(filename string, headers []string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create CSV file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(headers)
	switch v := data.(type) {
	case []models.Order:
		for _, order := range v {
			writer.Write([]string{
				order.Id.Hex(),
				fmt.Sprintf("%.2f", order.TotalPrice),
				order.CreatedAt.String(),
			})
		}
	case map[string]int:
		for key, value := range v {
			writer.Write([]string{key, strconv.Itoa(value)})
		}
	case map[string]float64:
		for key, value := range v {
			writer.Write([]string{key, fmt.Sprintf("%.2f", value)})
		}
	}
}

type kv struct {
	Key   string
	Value int
}

func getMostAndLeast(m map[string]int) (most, least kv) {
	var items []kv
	for k, v := range m {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})
	if len(items) > 0 {
		most = items[0]
		least = items[len(items)-1]
	}
	return
}

func getTopN(m map[string]int, n int) []kv {
	var items []kv
	for k, v := range m {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})
	if len(items) > n {
		return items[:n]
	}
	return items
}

func getMonthlyOrderCounts(orders []models.Order) map[string]int {
	monthly := make(map[string]int)
	for _, order := range orders {
		month := order.CreatedAt.Format("Jan")
		monthly[month]++
	}
	return monthly
}

func (repo *MongoRepository) DailyReport() ([]byte, error) {
	return repo.Report("daily")
}

func (repo *MongoRepository) WeeklyReport() ([]byte, error) {
	return repo.Report("weekly")
}

func (repo *MongoRepository) MonthlyReport() ([]byte, error) {
	return repo.Report("monthly")
}

func (repo *MongoRepository) YearlyReport() ([]byte, error) {
	return repo.Report("yearly")
}
