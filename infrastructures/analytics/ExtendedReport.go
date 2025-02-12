package analytics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ──────────────────────────────────────────
// DATA STRUCTURES
// ──────────────────────────────────────────

// Order-related structures
type FoodOrder struct {
	FoodId   primitive.ObjectID `json:"food_id" bson:"food_id"`
	Quantity float64            `json:"quantity" bson:"quantity"`
}

type DrinkOrder struct {
	DrinkId  primitive.ObjectID `json:"drink_id" bson:"drink_id"`
	Quantity float64            `json:"quantity" bson:"quantity"`
}

type Order struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	EmployeeId primitive.ObjectID `json:"employee_id" bson:"employee_id"`
	Foods      []FoodOrder        `json:"foods" bson:"foods"`
	Drinks     []DrinkOrder       `json:"drinks" bson:"drinks"`

	TotalPrice float64   `json:"total_price" bson:"total_price"`
	Status     string    `json:"status" bson:"status"` // e.g., "completed", "pending"
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

// Attendance-related structure
type Attendance struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeID primitive.ObjectID `json:"employee_id" bson:"employee_id"`
	Time       time.Time          `json:"time" bson:"time"`
	Type       string             `json:"type" bson:"type"` // expected values: "in", "out"
}

// Daily order metrics (orders per day)
type DailyOrderMetrics struct {
	Date         string  `json:"date"`
	OrdersCount  int     `json:"orders_count"`
	TotalRevenue float64 `json:"total_revenue"`
}

// ──────────────────────────────────────────
// REPORT STRUCTURES
// ──────────────────────────────────────────

type DetailedReport struct {
	OrderMetrics       OrderMetrics       `json:"order_metrics"`
	ItemMetrics        ItemMetrics        `json:"item_metrics"`
	EmployeeMetrics    []EmployeeMetrics  `json:"employee_metrics"`
	AttendanceMetrics  AttendanceMetrics  `json:"attendance_metrics"`
	OperationalMetrics OperationalMetrics `json:"operational_metrics"`
}

type OrderMetrics struct {
	TotalOrders         int              `json:"total_orders"`
	TotalRevenue        float64          `json:"total_revenue"`
	AverageOrderValue   float64          `json:"average_order_value"`
	MedianOrderValue    float64          `json:"median_order_value"`
	MinOrderValue       float64          `json:"min_order_value"`
	MaxOrderValue       float64          `json:"max_order_value"`
	OrderCompletionRate float64          `json:"order_completion_rate"` // percentage of completed orders
	OrdersByStatus      map[string]int   `json:"orders_by_status"`
	PeakOrderHours      []HourOrderCount `json:"peak_order_hours"`
	// Orders per employee, per day etc. can be derived from other sections.
}

type HourOrderCount struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

type ItemMetrics struct {
	// For both foods and drinks, we calculate total quantity sold.
	BestSellingFoods  []ItemRanking `json:"best_selling_foods"`
	BestSellingDrinks []ItemRanking `json:"best_selling_drinks"`
}

type ItemRanking struct {
	ItemId        string  `json:"item_id"`
	TotalQuantity float64 `json:"total_quantity"`
	// Additional fields like orders count could be added here.
}

type EmployeeMetrics struct {
	EmployeeID            string  `json:"employee_id"`
	OrdersProcessed       int     `json:"orders_processed"`
	TotalSales            float64 `json:"total_sales"`
	AverageOrderValue     float64 `json:"average_order_value"`
	AverageProcessingTime string  `json:"average_processing_time"` // e.g., "3m45s"
	ProcessingTimeStdDev  string  `json:"processing_time_std_dev"` // e.g., "1.23 minutes"
	TotalWorkDuration     string  `json:"total_work_duration"`     // total work time from paired checkin/out
}

type AttendanceMetrics struct {
	TotalCheckins   int                           `json:"total_checkins"`
	TotalCheckouts  int                           `json:"total_checkouts"`
	EmployeeRecords map[string]EmployeeAttendance `json:"employee_attendance"`
	WorkDurations   map[string]time.Duration      `json:"-"` // internal use
}

type EmployeeAttendance struct {
	EmployeeID   string `json:"employee_id"`
	Checkins     int    `json:"checkins"`
	Checkouts    int    `json:"checkouts"`
	LateCheckins int    `json:"late_checkins"`
}

type OperationalMetrics struct {
	EmployeeOrderToAttendanceRatio []EmployeeOrderAttendanceRatio `json:"employee_order_attendance_ratio"`
	IdleEmployees                  []string                       `json:"idle_employees"`
}

type EmployeeOrderAttendanceRatio struct {
	EmployeeID string  `json:"employee_id"`
	Ratio      float64 `json:"ratio"` // orders per checkin
}

// Additional insights for overall sales trends and cohort analysis.
type SalesTrend struct {
	MonthlyRevenue map[string]float64  `json:"monthly_revenue"`
	MonthlyOrders  map[string]int      `json:"monthly_orders"`
	DailyMetrics   []DailyOrderMetrics `json:"daily_metrics"`
	GrowthRate     float64             `json:"growth_rate"`
}

type CohortData struct {
	CohortLabel       string  `json:"cohort_label"`
	OrderCount        int     `json:"order_count"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

type EmployeeEfficiency struct {
	EmployeeID      string  `json:"employee_id"`
	EfficiencyRatio float64 `json:"efficiency_ratio"` // orders per checkin
}

// ExtendedReport is the final structure including all analytics.
type ExtendedReport struct {
	DetailedReport
	SalesTrend                SalesTrend           `json:"sales_trend"`
	PredictiveRevenue         float64              `json:"predictive_revenue"`
	CohortAnalysis            []CohortData         `json:"cohort_analysis"`
	EmployeeEfficiencyRanking []EmployeeEfficiency `json:"employee_efficiency_ranking"`
	EmployeeRevenueRanking    []EmployeeMetrics    `json:"employee_revenue_ranking"`
}

// ──────────────────────────────────────────
// HELPER FUNCTIONS
// ──────────────────────────────────────────

// standardDeviation computes the standard deviation of a slice of float64 values.
func standardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))
	sumSq := 0.0
	for _, v := range values {
		sumSq += (v - mean) * (v - mean)
	}
	return math.Sqrt(sumSq / float64(len(values)))
}

// median computes the median of a slice of float64 values.
func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	n := len(values)
	if n%2 == 1 {
		return values[n/2]
	}
	return (values[n/2-1] + values[n/2]) / 2.0
}

// calculateWorkDurations groups attendance records by employee,
// pairs "in" and "out" events (in order) and sums the durations.
func calculateWorkDurations(attendances []Attendance) map[string]time.Duration {
	workDurations := make(map[string]time.Duration)
	// Group records by employee
	empRecords := make(map[string][]Attendance)
	for _, rec := range attendances {
		empID := rec.EmployeeID.Hex()
		empRecords[empID] = append(empRecords[empID], rec)
	}
	// For each employee, sort records and pair "in" with the next "out"
	for empID, records := range empRecords {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Time.Before(records[j].Time)
		})
		var total time.Duration
		var lastIn time.Time
		inProgress := false
		for _, rec := range records {
			if rec.Type == "in" {
				// start a new shift if not already in progress
				if !inProgress {
					lastIn = rec.Time
					inProgress = true
				}
			} else if rec.Type == "out" && inProgress {
				total += rec.Time.Sub(lastIn)
				inProgress = false
			}
		}
		workDurations[empID] = total
	}
	return workDurations
}

// calculateDailyOrderMetrics groups orders by day.
func calculateDailyOrderMetrics(orders []Order) []DailyOrderMetrics {
	daily := make(map[string]DailyOrderMetrics)
	for _, order := range orders {
		day := order.CreatedAt.Format("2006-01-02")
		if m, ok := daily[day]; ok {
			m.OrdersCount++
			m.TotalRevenue += order.TotalPrice
			daily[day] = m
		} else {
			daily[day] = DailyOrderMetrics{
				Date:         day,
				OrdersCount:  1,
				TotalRevenue: order.TotalPrice,
			}
		}
	}
	var results []DailyOrderMetrics
	for _, m := range daily {
		results = append(results, m)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date < results[j].Date
	})
	return results
}

// ──────────────────────────────────────────
// CALCULATION FUNCTIONS
// ──────────────────────────────────────────

// calculateOrderMetrics computes overall order metrics.
func calculateOrderMetrics(orders []Order) OrderMetrics {
	totalOrders := len(orders)
	totalRevenue := 0.0
	sumOrderValue := 0.0
	var orderValues []float64
	ordersByStatus := make(map[string]int)
	peakHourMap := make(map[int]int)
	completedOrders := 0
	var minOrderValue, maxOrderValue float64

	for i, order := range orders {
		totalRevenue += order.TotalPrice
		sumOrderValue += order.TotalPrice
		orderValues = append(orderValues, order.TotalPrice)
		if i == 0 || order.TotalPrice < minOrderValue {
			minOrderValue = order.TotalPrice
		}
		if i == 0 || order.TotalPrice > maxOrderValue {
			maxOrderValue = order.TotalPrice
		}
		ordersByStatus[order.Status]++
		if order.Status == "completed" {
			completedOrders++
		}
		hour := order.CreatedAt.Hour()
		peakHourMap[hour]++
	}

	averageOrderValue := 0.0
	if totalOrders > 0 {
		averageOrderValue = sumOrderValue / float64(totalOrders)
	}
	medianOrderValue := median(orderValues)
	orderCompletionRate := 0.0
	if totalOrders > 0 {
		orderCompletionRate = float64(completedOrders) / float64(totalOrders) * 100
	}

	var peakOrderHours []HourOrderCount
	for hour, count := range peakHourMap {
		peakOrderHours = append(peakOrderHours, HourOrderCount{Hour: hour, Count: count})
	}
	sort.Slice(peakOrderHours, func(i, j int) bool {
		return peakOrderHours[i].Count > peakOrderHours[j].Count
	})

	return OrderMetrics{
		TotalOrders:         totalOrders,
		TotalRevenue:        totalRevenue,
		AverageOrderValue:   averageOrderValue,
		MedianOrderValue:    medianOrderValue,
		MinOrderValue:       minOrderValue,
		MaxOrderValue:       maxOrderValue,
		OrderCompletionRate: orderCompletionRate,
		OrdersByStatus:      ordersByStatus,
		PeakOrderHours:      peakOrderHours,
	}
}

// calculateItemMetrics aggregates sales quantities for food and drinks.
func calculateItemMetrics(orders []Order) ItemMetrics {
	foodSales := make(map[string]float64)
	drinkSales := make(map[string]float64)
	for _, order := range orders {
		for _, food := range order.Foods {
			key := food.FoodId.Hex()
			foodSales[key] += food.Quantity
		}
		for _, drink := range order.Drinks {
			key := drink.DrinkId.Hex()
			drinkSales[key] += drink.Quantity
		}
	}
	return ItemMetrics{
		BestSellingFoods:  rankItems(foodSales),
		BestSellingDrinks: rankItems(drinkSales),
	}
}

func rankItems(itemMap map[string]float64) []ItemRanking {
	var rankings []ItemRanking
	for id, qty := range itemMap {
		rankings = append(rankings, ItemRanking{ItemId: id, TotalQuantity: qty})
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].TotalQuantity > rankings[j].TotalQuantity
	})
	return rankings
}

// calculateEmployeeMetrics computes performance per employee from order data,
// and later enriches with total work duration.
func calculateEmployeeMetrics(orders []Order, workDurations map[string]time.Duration) []EmployeeMetrics {
	employeeData := make(map[string]*EmployeeMetrics)
	employeeTimestamps := make(map[string][]time.Time)

	for _, order := range orders {
		empID := order.EmployeeId.Hex()
		if _, exists := employeeData[empID]; !exists {
			employeeData[empID] = &EmployeeMetrics{EmployeeID: empID}
		}
		emp := employeeData[empID]
		emp.OrdersProcessed++
		emp.TotalSales += order.TotalPrice
		employeeTimestamps[empID] = append(employeeTimestamps[empID], order.CreatedAt)
	}

	for empID, emp := range employeeData {
		if emp.OrdersProcessed > 0 {
			emp.AverageOrderValue = emp.TotalSales / float64(emp.OrdersProcessed)
		}
		timestamps := employeeTimestamps[empID]
		if len(timestamps) > 1 {
			sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })
			var totalGap time.Duration
			var gaps []float64
			for i := 1; i < len(timestamps); i++ {
				gap := timestamps[i].Sub(timestamps[i-1])
				totalGap += gap
				gaps = append(gaps, gap.Minutes())
			}
			avgGap := totalGap / time.Duration(len(timestamps)-1)
			emp.AverageProcessingTime = avgGap.String()
			stdDev := standardDeviation(gaps)
			emp.ProcessingTimeStdDev = fmt.Sprintf("%.2f minutes", stdDev)
		} else {
			emp.AverageProcessingTime = "N/A"
			emp.ProcessingTimeStdDev = "N/A"
		}
		// Enrich with work duration if available
		if dur, ok := workDurations[empID]; ok {
			emp.TotalWorkDuration = dur.String()
		} else {
			emp.TotalWorkDuration = "0s"
		}
	}

	var employees []EmployeeMetrics
	for _, v := range employeeData {
		employees = append(employees, *v)
	}
	// Sort by total sales descending (for revenue ranking)
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].TotalSales > employees[j].TotalSales
	})
	return employees
}

// calculateAttendanceMetrics computes check-in/out counts, late arrivals,
// and returns a map of total work durations per employee.
func calculateAttendanceMetrics(attendances []Attendance) AttendanceMetrics {
	totalCheckins := 0
	totalCheckouts := 0
	employeeRecords := make(map[string]EmployeeAttendance)
	for _, record := range attendances {
		empID := record.EmployeeID.Hex()
		rec, exists := employeeRecords[empID]
		if !exists {
			rec = EmployeeAttendance{EmployeeID: empID}
		}
		if record.Type == "in" {
			totalCheckins++
			rec.Checkins++
			if record.Time.Hour() >= 9 {
				rec.LateCheckins++
			}
		} else if record.Type == "out" {
			totalCheckouts++
			rec.Checkouts++
		}
		employeeRecords[empID] = rec
	}
	return AttendanceMetrics{
		TotalCheckins:   totalCheckins,
		TotalCheckouts:  totalCheckouts,
		EmployeeRecords: employeeRecords,
	}
}

// calculateOperationalMetrics computes ratios of orders processed to attendance.
func calculateOperationalMetrics(empMetrics []EmployeeMetrics, attMetrics AttendanceMetrics) OperationalMetrics {
	employeeOrders := make(map[string]int)
	for _, emp := range empMetrics {
		employeeOrders[emp.EmployeeID] = emp.OrdersProcessed
	}
	employeeAttendance := make(map[string]int)
	for empID, record := range attMetrics.EmployeeRecords {
		employeeAttendance[empID] = record.Checkins
	}
	var ratios []EmployeeOrderAttendanceRatio
	var idleEmployees []string
	for empID, attCount := range employeeAttendance {
		orders := employeeOrders[empID]
		ratio := 0.0
		if attCount > 0 {
			ratio = float64(orders) / float64(attCount)
		}
		ratios = append(ratios, EmployeeOrderAttendanceRatio{EmployeeID: empID, Ratio: ratio})
		if attCount > 0 && orders == 0 {
			idleEmployees = append(idleEmployees, empID)
		}
	}
	sort.Slice(ratios, func(i, j int) bool {
		return ratios[i].Ratio > ratios[j].Ratio
	})
	return OperationalMetrics{
		EmployeeOrderToAttendanceRatio: ratios,
		IdleEmployees:                  idleEmployees,
	}
}

// calculateSalesTrend computes monthly revenue, order counts, and daily metrics.
func calculateSalesTrend(orders []Order) SalesTrend {
	monthlyRevenue := make(map[string]float64)
	monthlyOrders := make(map[string]int)
	for _, order := range orders {
		month := order.CreatedAt.Format("Jan 2006")
		monthlyRevenue[month] += order.TotalPrice
		monthlyOrders[month]++
	}
	dailyMetrics := calculateDailyOrderMetrics(orders)
	var months []string
	for m := range monthlyRevenue {
		months = append(months, m)
	}
	sort.Strings(months)
	growthRate := 0.0
	if len(months) >= 2 {
		first := monthlyRevenue[months[0]]
		last := monthlyRevenue[months[len(months)-1]]
		if first > 0 {
			growthRate = (last - first) / first * 100
		}
	}
	return SalesTrend{
		MonthlyRevenue: monthlyRevenue,
		MonthlyOrders:  monthlyOrders,
		DailyMetrics:   dailyMetrics,
		GrowthRate:     growthRate,
	}
}

// forecastRevenue provides a naive forecast for next month's revenue.
func forecastRevenue(trend SalesTrend) float64 {
	var months []string
	for m := range trend.MonthlyRevenue {
		months = append(months, m)
	}
	if len(months) == 0 {
		return 0
	}
	sort.Strings(months)
	lastMonthRevenue := trend.MonthlyRevenue[months[len(months)-1]]
	forecast := lastMonthRevenue + (trend.GrowthRate/100)*lastMonthRevenue/float64(len(months))
	return forecast
}

// calculateCohortAnalysis groups orders by month and computes cohort metrics.
func calculateCohortAnalysis(orders []Order) []CohortData {
	cohortMap := make(map[string][]Order)
	for _, order := range orders {
		cohort := order.CreatedAt.Format("Jan 2006")
		cohortMap[cohort] = append(cohortMap[cohort], order)
	}
	var cohorts []CohortData
	for cohort, ordersInCohort := range cohortMap {
		totalRevenue := 0.0
		for _, o := range ordersInCohort {
			totalRevenue += o.TotalPrice
		}
		avgOrder := 0.0
		if len(ordersInCohort) > 0 {
			avgOrder = totalRevenue / float64(len(ordersInCohort))
		}
		cohorts = append(cohorts, CohortData{
			CohortLabel:       cohort,
			OrderCount:        len(ordersInCohort),
			TotalRevenue:      totalRevenue,
			AverageOrderValue: avgOrder,
		})
	}
	sort.Slice(cohorts, func(i, j int) bool {
		t1, _ := time.Parse("Jan 2006", cohorts[i].CohortLabel)
		t2, _ := time.Parse("Jan 2006", cohorts[j].CohortLabel)
		return t1.Before(t2)
	})
	return cohorts
}

// calculateEmployeeEfficiency computes orders per checkin for each employee.
func calculateEmployeeEfficiency(empMetrics []EmployeeMetrics, attMetrics AttendanceMetrics) []EmployeeEfficiency {
	var rankings []EmployeeEfficiency
	for empID, record := range attMetrics.EmployeeRecords {
		orders := 0
		for _, emp := range empMetrics {
			if emp.EmployeeID == empID {
				orders = emp.OrdersProcessed
				break
			}
		}
		ratio := 0.0
		if record.Checkins > 0 {
			ratio = float64(orders) / float64(record.Checkins)
		}
		rankings = append(rankings, EmployeeEfficiency{
			EmployeeID:      empID,
			EfficiencyRatio: ratio,
		})
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].EfficiencyRatio > rankings[j].EfficiencyRatio
	})
	return rankings
}

// calculateEmployeeRevenueRanking sorts employee metrics by total sales.
func calculateEmployeeRevenueRanking(empMetrics []EmployeeMetrics) []EmployeeMetrics {
	// Create a copy
	ranking := make([]EmployeeMetrics, len(empMetrics))
	copy(ranking, empMetrics)
	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].TotalSales > ranking[j].TotalSales
	})
	return ranking
}

// ──────────────────────────────────────────
// REPORT GENERATION FUNCTION
// ──────────────────────────────────────────

// GenerateExtendedReport fetches orders and attendance records for the given period,
// calculates detailed KPIs and analytics including item analysis, attendance work duration,
// daily order metrics, and employee rankings, and returns an ExtendedReport.
func GenerateExtendedReport(db map[string]mongo.Collection, startDate, endDate time.Time) (*ExtendedReport, error) {
	ctx := context.TODO()

	// Fetch orders
	ordersColl := db["orders"]
	orderFilter := bson.M{"created_at": bson.M{"$gte": startDate, "$lte": endDate}}
	orderCursor, err := ordersColl.Find(ctx, orderFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	var orders []Order
	if err = orderCursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}
	orderCursor.Close(ctx)

	// Fetch attendance records
	attendanceColl := db["attendance"]
	attendanceFilter := bson.M{"time": bson.M{"$gte": startDate, "$lte": endDate}}
	attendanceCursor, err := attendanceColl.Find(ctx, attendanceFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %v", err)
	}
	var attendances []Attendance
	if err = attendanceCursor.All(ctx, &attendances); err != nil {
		return nil, fmt.Errorf("failed to decode attendance: %v", err)
	}
	attendanceCursor.Close(ctx)

	// Compute work durations from attendance records.
	workDurations := calculateWorkDurations(attendances)

	// Calculate core metrics.
	orderMetrics := calculateOrderMetrics(orders)
	itemMetrics := calculateItemMetrics(orders)
	attendanceMetrics := calculateAttendanceMetrics(attendances)
	operationalMetrics := calculateOperationalMetrics(
		calculateEmployeeMetrics(orders, workDurations), attendanceMetrics)
	// Calculate employee metrics (enriched with work duration).
	employeeMetrics := calculateEmployeeMetrics(orders, workDurations)
	detailedReport := DetailedReport{
		OrderMetrics:       orderMetrics,
		ItemMetrics:        itemMetrics,
		EmployeeMetrics:    employeeMetrics,
		AttendanceMetrics:  attendanceMetrics,
		OperationalMetrics: operationalMetrics,
	}

	// Additional analytics.
	salesTrend := calculateSalesTrend(orders)
	predictiveRevenue := forecastRevenue(salesTrend)
	cohortAnalysis := calculateCohortAnalysis(orders)
	employeeEfficiencyRanking := calculateEmployeeEfficiency(employeeMetrics, attendanceMetrics)
	employeeRevenueRanking := calculateEmployeeRevenueRanking(employeeMetrics)

	extReport := &ExtendedReport{
		DetailedReport:            detailedReport,
		SalesTrend:                salesTrend,
		PredictiveRevenue:         predictiveRevenue,
		CohortAnalysis:            cohortAnalysis,
		EmployeeEfficiencyRanking: employeeEfficiencyRanking,
		EmployeeRevenueRanking:    employeeRevenueRanking,
	}
	return extReport, nil
}
