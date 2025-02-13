package analytics

// import (
// 	"context"
// 	"encoding/csv"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"sort"
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
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"

// )

// // ──────────────────────────────────────────
// // DATA STRUCTURES
// // ──────────────────────────────────────────

// type ExtendedReport struct {
// 	Period           DateRange             `json:"period"`
// 	Attendance       AttendanceAnalytics   `json:"attendance"`
// 	Productivity     ProductivityAnalytics `json:"productivity"`
// 	Sales            SalesAnalytics        `json:"sales"`
// 	Operations       OperationalAnalytics  `json:"operations"`
// 	FraudIndicators  FraudAnalysis         `json:"fraud_indicators"`
// 	SeasonalTrends   []SeasonalTrend       `json:"seasonal_trends"`
// 	PredictedMetrics PredictiveMetrics     `json:"predicted_metrics"`
// }

// type DateRange struct {
// 	Start time.Time `json:"start"`
// 	End   time.Time `json:"end"`
// }

// type AttendanceAnalytics struct {
// 	Summary         AttendanceSummary          `json:"summary"`
// 	DailyTrends     []DailyAttendance          `json:"daily_trends"`
// 	ShiftAnalysis   []ShiftMetrics             `json:"shift_analysis"`
// 	EmployeeDetails []EmployeeAttendanceDetail `json:"employee_details"`
// }

// type AttendanceSummary struct {
// 	TotalWorkHours     float64 `json:"total_work_hours"`
// 	AvgDailyAttendance float64 `json:"avg_daily_attendance"`
// 	LateArrivals       int     `json:"late_arrivals"`
// 	EarlyDepartures    int     `json:"early_departures"`
// 	AbsenteeismRate    float64 `json:"absenteeism_rate"`
// }

// type DailyAttendance struct {
// 	Date              string  `json:"date"`
// 	PresentEmployees  int     `json:"present_employees"`
// 	AvgWorkHours      float64 `json:"avg_work_hours"`
// 	PeakArrivalTime   string  `json:"peak_arrival_time"`
// 	PeakDepartureTime string  `json:"peak_departure_time"`
// }

// type ShiftMetrics struct {
// 	ShiftName       string  `json:"shift_name"`
// 	AvgHoursWorked  float64 `json:"avg_hours_worked"`
// 	TotalOrders     int     `json:"total_orders"`
// 	TotalRevenue    float64 `json:"total_revenue"`
// 	EfficiencyRatio float64 `json:"efficiency_ratio"`
// }

// type EmployeeAttendanceDetail struct {
// 	EmployeeID      string   `json:"employee_id"`
// 	TotalHours      float64  `json:"total_hours"`
// 	AvgDailyHours   float64  `json:"avg_daily_hours"`
// 	LateArrivals    int      `json:"late_arrivals"`
// 	EarlyDepartures int      `json:"early_departures"`
// 	AttendanceDays  []string `json:"attendance_days"`
// 	OvertimeHours   float64  `json:"overtime_hours"`
// }

// type ProductivityAnalytics struct {
// 	EmployeeRankings []EmployeeProductivity `json:"employee_rankings"`
// 	OrderEfficiency  OrderEfficiencyMetrics `json:"order_efficiency"`
// 	WorkloadBalance  WorkloadDistribution   `json:"workload_balance"`
// 	AccuracyAnalysis AccuracyMetrics        `json:"accuracy_analysis"`
// }

// type EmployeeProductivity struct {
// 	EmployeeID       string  `json:"employee_id"`
// 	OrdersProcessed  int     `json:"orders_processed"`
// 	RevenueGenerated float64 `json:"revenue_generated"`
// 	EfficiencyScore  float64 `json:"efficiency_score"`
// 	AccuracyRate     float64 `json:"accuracy_rate"`
// 	AvgOrderTime     float64 `json:"avg_order_time"`
// }

// type OrderEfficiencyMetrics struct {
// 	AvgProcessingTime  float64 `json:"avg_processing_time"`
// 	PeakEfficiencyHour string  `json:"peak_efficiency_hour"`
// 	SlowestHour        string  `json:"slowest_hour"`
// 	StdDeviation       float64 `json:"std_deviation"`
// }

// type WorkloadDistribution struct {
// 	MostLoadedEmployee  string  `json:"most_loaded_employee"`
// 	LeastLoadedEmployee string  `json:"least_loaded_employee"`
// 	WorkloadDisparity   float64 `json:"workload_disparity"`
// }

// type AccuracyMetrics struct {
// 	TotalCompleted   int      `json:"total_completed"`
// 	TotalCancelled   int      `json:"total_cancelled"`
// 	CancellationRate float64  `json:"cancellation_rate"`
// 	TopErrorSources  []string `json:"top_error_sources"`
// }

// type SalesAnalytics struct {
// 	RevenueTrends      RevenueAnalysis  `json:"revenue_trends"`
// 	ProductPerformance ProductMetrics   `json:"product_performance"`
// 	CustomerBehavior   CustomerInsights `json:"customer_behavior"`
// 	ShiftComparison    []ShiftSales     `json:"shift_comparison"`
// }

// type RevenueAnalysis struct {
// 	TotalRevenue     float64            `json:"total_revenue"`
// 	RevenueByDayPart map[string]float64 `json:"revenue_by_day_part"`
// 	GrowthRate       float64            `json:"growth_rate"`
// 	RecurringRevenue float64            `json:"recurring_revenue"`
// }

// type ProductMetrics struct {
// 	TopSellingItems   []ProductRanking   `json:"top_selling_items"`
// 	WorstSellingItems []ProductRanking   `json:"worst_selling_items"`
// 	MarginAnalysis    map[string]float64 `json:"margin_analysis"`
// 	WasteAnalysis     map[string]float64 `json:"waste_analysis"`
// }

// type ProductRanking struct {
// 	ProductID   string  `json:"product_id"`
// 	ProductName string  `json:"product_name"`
// 	TotalSales  float64 `json:"total_sales"`
// 	TotalUnits  int     `json:"total_units"`
// }

// type CustomerInsights struct {
// 	AvgOrderValue      float64        `json:"avg_order_value"`
// 	PeakOrderTimes     []HourlySales  `json:"peak_order_times"`
// 	LoyaltySegments    map[string]int `json:"loyalty_segments"`
// 	RepeatCustomerRate float64        `json:"repeat_customer_rate"`
// }

// type HourlySales struct {
// 	Hour   string  `json:"hour"`
// 	Sales  float64 `json:"sales"`
// 	Orders int     `json:"orders"`
// }

// type ShiftSales struct {
// 	ShiftName       string  `json:"shift_name"`
// 	SalesTotal      float64 `json:"sales_total"`
// 	AvgOrderValue   float64 `json:"avg_order_value"`
// 	CustomersServed int     `json:"customers_served"`
// }

// type OperationalAnalytics struct {
// 	Efficiency    OperationalEfficiency `json:"efficiency"`
// 	ResourceUsage ResourceUtilization   `json:"resource_usage"`
// 	CostAnalysis  CostMetrics           `json:"cost_analysis"`
// }

// type OperationalEfficiency struct {
// 	LaborCostPerOrder   float64 `json:"labor_cost_per_order"`
// 	RevenuePerLaborHour float64 `json:"revenue_per_labor_hour"`
// 	BreakEvenRatio      float64 `json:"break_even_ratio"`
// }

// type ResourceUtilization struct {
// 	EquipmentUsage   map[string]float64 `json:"equipment_usage"`
// 	SpaceUtilization map[string]float64 `json:"space_utilization"`
// }

// type CostMetrics struct {
// 	FoodCostPercentage  float64 `json:"food_cost_percentage"`
// 	LaborCostPercentage float64 `json:"labor_cost_percentage"`
// 	PrimeCost           float64 `json:"prime_cost"`
// }

// type FraudAnalysis struct {
// 	SuspiciousActivities []SuspiciousActivity `json:"suspicious_activities"`
// 	AnomalySummary       AnomalySummary       `json:"anomaly_summary"`
// }

// type SuspiciousActivity struct {
// 	EmployeeID   string    `json:"employee_id"`
// 	ActivityType string    `json:"activity_type"`
// 	Occurrences  int       `json:"occurrences"`
// 	LastOccurred time.Time `json:"last_occurred"`
// }

// type AnomalySummary struct {
// 	TotalFraudCases int     `json:"total_fraud_cases"`
// 	PotentialLoss   float64 `json:"potential_loss"`
// 	MostCommonType  string  `json:"most_common_type"`
// }

// type SeasonalTrend struct {
// 	Period         string             `json:"period"`
// 	TrendIndicator string             `json:"trend_indicator"`
// 	ImpactMetrics  map[string]float64 `json:"impact_metrics"`
// }

// type PredictiveMetrics struct {
// 	NextWeekForecast        ForecastPrediction       `json:"next_week_forecast"`
// 	StaffingRecommendations []StaffingRecommendation `json:"staffing_recommendations"`
// }

// type ForecastPrediction struct {
// 	ExpectedOrders  int     `json:"expected_orders"`
// 	ExpectedRevenue float64 `json:"expected_revenue"`
// 	ConfidenceLevel float64 `json:"confidence_level"`
// }

// type StaffingRecommendation struct {
// 	ShiftName        string `json:"shift_name"`
// 	RecommendedStaff int    `json:"recommended_staff"`
// 	ExpectedDemand   int    `json:"expected_demand"`
// }

// // ──────────────────────────────────────────
// // MAIN REPORT GENERATION FUNCTION
// // ──────────────────────────────────────────

// // GenerateEnhancedReport concurrently fetches data from multiple collections,
// // computes derived metrics and predictions, and returns an ExtendedReport.
// func GenerateEnhancedReport(db map[string]*mongo.Collection, start, end time.Time) (*ExtendedReport, error) {
// 	var wg sync.WaitGroup
// 	errChan := make(chan error, 5)
// 	var report ExtendedReport
// 	ctx := context.TODO()

// 	report.Period = DateRange{Start: start, End: end}

// 	wg.Add(5)
// 	go fetchAttendanceData(ctx, db, start, end, &report, errChan, &wg)
// 	go fetchOrderData(ctx, db, start, end, &report, errChan, &wg)
// 	go fetchInventoryData(ctx, db, start, end, &report, errChan, &wg)
// 	go fetchFraudData(ctx, db, start, end, &report, errChan, &wg)
// 	go fetchHistoricalData(ctx, db, start, end, &report, errChan, &wg)

// 	go func() {
// 		wg.Wait()
// 		close(errChan)
// 	}()

// 	for err := range errChan {
// 		if err != nil {
// 			return nil, fmt.Errorf("error generating report: %v", err)
// 		}
// 	}

// 	calculateDerivedMetrics(&report)
// 	generatePredictions(&report)
// 	applySortingRankings(&report)

// 	return &report, nil
// }

// // ──────────────────────────────────────────
// // HELPER FUNCTIONS
// // ──────────────────────────────────────────

// // fetchAttendanceData retrieves and aggregates attendance records.
// func fetchAttendanceData(ctx context.Context, db map[string]*mongo.Collection, start, end time.Time, report *ExtendedReport, errChan chan error, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	if db == nil {
// 		errChan <- fmt.Errorf("database collections map is nil")
// 		return
// 	}
// 	attendanceCol, ok := db["attendance"]
// 	if !ok || attendanceCol == nil {
// 		errChan <- fmt.Errorf("attendance collection is nil")
// 		return
// 	}

// 	pipeline := []bson.M{
// 		{"$match": bson.M{"time": bson.M{"$gte": start, "$lte": end}}},
// 		{"$group": bson.M{
// 			"_id":              "$employee_id",
// 			"total_hours":      bson.M{"$sum": "$hours_worked"},
// 			"late_arrivals":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$gt": []interface{}{"$check_in_time", 9}}, 1, 0}}},
// 			"early_departures": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$lt": []interface{}{"$check_out_time", 17}}, 1, 0}}},
// 		}},
// 	}

// 	cursor, err := attendanceCol.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		errChan <- fmt.Errorf("failed to fetch attendance data: %v", err)
// 		return
// 	}
// 	defer cursor.Close(ctx)

// 	var results []bson.M
// 	if err = cursor.All(ctx, &results); err != nil {
// 		errChan <- fmt.Errorf("failed to decode attendance data: %v", err)
// 		return
// 	}

// 	var totalHours float64
// 	var totalLate, totalEarly int
// 	var employeeDetails []EmployeeAttendanceDetail

// 	for _, res := range results {
// 		hours, _ := res["total_hours"].(float64)
// 		late, _ := res["late_arrivals"].(int32)
// 		early, _ := res["early_departures"].(int32)

// 		totalHours += hours
// 		totalLate += int(late)
// 		totalEarly += int(early)

// 		employeeDetails = append(employeeDetails, EmployeeAttendanceDetail{
// 			EmployeeID:      fmt.Sprintf("%v", res["_id"]),
// 			TotalHours:      hours,
// 			LateArrivals:    int(late),
// 			EarlyDepartures: int(early),
// 			AttendanceDays:  []string{},  // Dummy value
// 			AvgDailyHours:   hours / 5.0, // Assume 5 working days
// 			OvertimeHours:   0,           // Dummy value
// 		})
// 	}

// 	report.Attendance.Summary = AttendanceSummary{
// 		TotalWorkHours:     totalHours,
// 		AvgDailyAttendance: float64(len(results)) / 5.0, // Dummy calculation
// 		LateArrivals:       totalLate,
// 		EarlyDepartures:    totalEarly,
// 		AbsenteeismRate:    0.05, // Dummy value
// 	}
// 	report.Attendance.EmployeeDetails = employeeDetails
// 	report.Attendance.DailyTrends = []DailyAttendance{} // Not implemented
// 	report.Attendance.ShiftAnalysis = []ShiftMetrics{}  // Not implemented
// }

// // fetchOrderData retrieves order data and aggregates productivity and sales metrics.
// func fetchOrderData(ctx context.Context, db map[string]*mongo.Collection, start, end time.Time, report *ExtendedReport, errChan chan error, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	if db == nil {
// 		errChan <- fmt.Errorf("database collections map is nil")
// 		return
// 	}
// 	ordersCol, ok := db["orders"]
// 	if !ok || ordersCol == nil {
// 		errChan <- fmt.Errorf("orders collection is nil")
// 		return
// 	}

// 	pipeline := []bson.M{
// 		{"$match": bson.M{"created_at": bson.M{"$gte": start, "$lte": end}}},
// 		{"$group": bson.M{
// 			"_id":               "$employee_id",
// 			"orders_processed":  bson.M{"$sum": 1},
// 			"revenue_generated": bson.M{"$sum": "$total_price"},
// 		}},
// 	}

// 	cursor, err := ordersCol.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		errChan <- fmt.Errorf("failed to fetch order data: %v", err)
// 		return
// 	}
// 	defer cursor.Close(ctx)

// 	var results []bson.M
// 	if err = cursor.All(ctx, &results); err != nil {
// 		errChan <- fmt.Errorf("failed to decode order data: %v", err)
// 		return
// 	}

// 	var employeeRankings []EmployeeProductivity
// 	var totalRevenue float64

// 	for _, res := range results {
// 		orders, _ := res["orders_processed"].(int32)
// 		revenue, _ := res["revenue_generated"].(float64)
// 		employeeRankings = append(employeeRankings, EmployeeProductivity{
// 			EmployeeID:       fmt.Sprintf("%v", res["_id"]),
// 			OrdersProcessed:  int(orders),
// 			RevenueGenerated: revenue,
// 			EfficiencyScore:  revenue / float64(orders), // Dummy calculation
// 			AccuracyRate:     0.95,                      // Dummy value
// 			AvgOrderTime:     5.0,                       // Dummy value (in minutes)
// 		})
// 		totalRevenue += revenue
// 	}

// 	report.Productivity.EmployeeRankings = employeeRankings
// 	report.Sales.RevenueTrends = RevenueAnalysis{
// 		TotalRevenue: totalRevenue,
// 		RevenueByDayPart: map[string]float64{
// 			"Morning":   totalRevenue * 0.3,
// 			"Afternoon": totalRevenue * 0.5,
// 			"Evening":   totalRevenue * 0.2,
// 		},
// 		GrowthRate:       0.1,                // Dummy value
// 		RecurringRevenue: totalRevenue * 0.2, // Dummy value
// 	}
// }

// // fetchInventoryData simulates fetching product performance data.
// func fetchInventoryData(ctx context.Context, db map[string]*mongo.Collection, start, end time.Time, report *ExtendedReport, errChan chan error, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	topItems := []ProductRanking{
// 		{ProductID: "P001", ProductName: "Burger", TotalSales: 1500.0, TotalUnits: 300},
// 		{ProductID: "P002", ProductName: "Pizza", TotalSales: 1200.0, TotalUnits: 200},
// 	}

// 	worstItems := []ProductRanking{
// 		{ProductID: "P010", ProductName: "Salad", TotalSales: 300.0, TotalUnits: 50},
// 	}

// 	report.Sales.ProductPerformance = ProductMetrics{
// 		TopSellingItems:   topItems,
// 		WorstSellingItems: worstItems,
// 		MarginAnalysis:    map[string]float64{"Burger": 0.3, "Pizza": 0.25},
// 		WasteAnalysis:     map[string]float64{"Burger": 0.05, "Pizza": 0.03},
// 	}
// }

// // fetchFraudData simulates detection and aggregation of fraud-related activities.
// func fetchFraudData(ctx context.Context, db map[string]*mongo.Collection, start, end time.Time, report *ExtendedReport, errChan chan error, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	suspicious := []SuspiciousActivity{
// 		{EmployeeID: "E999", ActivityType: "Multiple Orders", Occurrences: 5, LastOccurred: time.Now()},
// 		{EmployeeID: "E888", ActivityType: "Duplicate Check-ins", Occurrences: 3, LastOccurred: time.Now()},
// 	}
// 	anomaly := AnomalySummary{
// 		TotalFraudCases: 2,
// 		PotentialLoss:   500.0,
// 		MostCommonType:  "Multiple Orders",
// 	}

// 	report.FraudIndicators = FraudAnalysis{
// 		SuspiciousActivities: suspicious,
// 		AnomalySummary:       anomaly,
// 	}
// }

// // fetchHistoricalData simulates retrieving historical data for predictions and seasonal trends.
// func fetchHistoricalData(ctx context.Context, db map[string]*mongo.Collection, start, end time.Time, report *ExtendedReport, errChan chan error, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	forecast := ForecastPrediction{
// 		ExpectedOrders:  200,
// 		ExpectedRevenue: 40000.0,
// 		ConfidenceLevel: 0.8,
// 	}
// 	staffing := []StaffingRecommendation{
// 		{ShiftName: "Morning", RecommendedStaff: 5, ExpectedDemand: 100},
// 		{ShiftName: "Evening", RecommendedStaff: 3, ExpectedDemand: 50},
// 	}

// 	report.PredictedMetrics = PredictiveMetrics{
// 		NextWeekForecast:        forecast,
// 		StaffingRecommendations: staffing,
// 	}

// 	report.SeasonalTrends = []SeasonalTrend{
// 		{Period: "Winter", TrendIndicator: "High", ImpactMetrics: map[string]float64{"revenue": 0.15}},
// 		{Period: "Summer", TrendIndicator: "Low", ImpactMetrics: map[string]float64{"revenue": -0.05}},
// 	}
// }

// // calculateDerivedMetrics computes additional metrics based on the fetched data.
// func calculateDerivedMetrics(report *ExtendedReport) {
// 	for i, emp := range report.Productivity.EmployeeRankings {
// 		if emp.OrdersProcessed > 0 {
// 			report.Productivity.EmployeeRankings[i].EfficiencyScore = emp.RevenueGenerated / float64(emp.OrdersProcessed)
// 		} else {
// 			report.Productivity.EmployeeRankings[i].EfficiencyScore = 0
// 		}
// 	}
// }

// // generatePredictions adjusts forecast metrics based on historical trends.
// func generatePredictions(report *ExtendedReport) {
// 	report.PredictedMetrics.NextWeekForecast.ExpectedOrders += 10
// 	report.PredictedMetrics.NextWeekForecast.ExpectedRevenue *= 1.05
// }

// // applySortingRankings sorts rankings such as employee productivity in descending order.
// func applySortingRankings(report *ExtendedReport) {
// 	sort.SliceStable(report.Productivity.EmployeeRankings, func(i, j int) bool {
// 		return report.Productivity.EmployeeRankings[i].RevenueGenerated > report.Productivity.EmployeeRankings[j].RevenueGenerated
// 	})
// }

// // SaveReportCSV saves the extended report as a CSV file in the folder structure Records/<period>.
// // The CSV file is named with the period and a timestamp.
// func SaveReportCSV(report *ExtendedReport, period string) error {
// 	baseDir := "Records"
// 	periodDir := filepath.Join(baseDir, strings.ToLower(period))
// 	if err := os.MkdirAll(periodDir, os.ModePerm); err != nil {
// 		return fmt.Errorf("failed to create directory %s: %v", periodDir, err)
// 	}
// 	filename := filepath.Join(periodDir, fmt.Sprintf("extended_report_%s_%s.csv", period, time.Now().Format("20060102_150405")))
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file %s: %v", filename, err)
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Write header row
// 	header := []string{"Section", "Metric", "Value"}
// 	if err := writer.Write(header); err != nil {
// 		return fmt.Errorf("failed to write header: %v", err)
// 	}

// 	// Example: Write Attendance Summary metrics
// 	att := report.Attendance.Summary
// 	rows := [][]string{
// 		{"Attendance", "TotalWorkHours", fmt.Sprintf("%.2f", att.TotalWorkHours)},
// 		{"Attendance", "AvgDailyAttendance", fmt.Sprintf("%.2f", att.AvgDailyAttendance)},
// 		{"Attendance", "LateArrivals", fmt.Sprintf("%d", att.LateArrivals)},
// 		{"Attendance", "EarlyDepartures", fmt.Sprintf("%d", att.EarlyDepartures)},
// 		{"Attendance", "AbsenteeismRate", fmt.Sprintf("%.2f", att.AbsenteeismRate)},
// 	}
// 	for _, row := range rows {
// 		if err := writer.Write(row); err != nil {
// 			return fmt.Errorf("failed to write row: %v", err)
// 		}
// 	}

// 	log.Printf("CSV report saved: %s", filename)
// 	return nil
// }
