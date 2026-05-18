package prism

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{"valid security", DomainSecurity, false},
		{"valid operations", DomainOperations, false},
		{"empty", "", true},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDomain(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDomain(%q) error = %v, wantErr %v", tt.domain, err, tt.wantErr)
			}
		})
	}
}

func TestValidateStage(t *testing.T) {
	tests := []struct {
		name    string
		stage   string
		wantErr bool
	}{
		{"valid design", StageDesign, false},
		{"valid build", StageBuild, false},
		{"valid test", StageTest, false},
		{"valid runtime", StageRuntime, false},
		{"valid response", StageResponse, false},
		{"empty", "", true},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStage(tt.stage)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStage(%q) error = %v, wantErr %v", tt.stage, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{"valid prevention", CategoryPrevention, false},
		{"valid detection", CategoryDetection, false},
		{"valid response", CategoryResponse, false},
		{"valid reliability", CategoryReliability, false},
		{"valid efficiency", CategoryEfficiency, false},
		{"valid quality", CategoryQuality, false},
		{"empty", "", true},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory(%q) error = %v, wantErr %v", tt.category, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMaturityLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		wantErr bool
	}{
		{"level 1", MaturityLevel1, false},
		{"level 2", MaturityLevel2, false},
		{"level 3", MaturityLevel3, false},
		{"level 4", MaturityLevel4, false},
		{"level 5", MaturityLevel5, false},
		{"level 0", 0, true},
		{"level 6", 6, true},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMaturityLevel(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMaturityLevel(%d) error = %v, wantErr %v", tt.level, err, tt.wantErr)
			}
		})
	}
}

func TestMetricCalculateStatus(t *testing.T) {
	tests := []struct {
		name           string
		metric         Metric
		expectedStatus string
	}{
		{
			name: "higher better - green",
			metric: Metric{
				Current:        95,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Green: 90, Yellow: 70, Red: 50},
			},
			expectedStatus: StatusGreen,
		},
		{
			name: "higher better - yellow",
			metric: Metric{
				Current:        75,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Green: 90, Yellow: 70, Red: 50},
			},
			expectedStatus: StatusYellow,
		},
		{
			name: "higher better - red",
			metric: Metric{
				Current:        40,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Green: 90, Yellow: 70, Red: 50},
			},
			expectedStatus: StatusRed,
		},
		{
			name: "lower better - green",
			metric: Metric{
				Current:        5,
				TrendDirection: TrendLowerBetter,
				Thresholds:     &Thresholds{Green: 10, Yellow: 30, Red: 50},
			},
			expectedStatus: StatusGreen,
		},
		{
			name: "lower better - yellow",
			metric: Metric{
				Current:        20,
				TrendDirection: TrendLowerBetter,
				Thresholds:     &Thresholds{Green: 10, Yellow: 30, Red: 50},
			},
			expectedStatus: StatusYellow,
		},
		{
			name: "lower better - red",
			metric: Metric{
				Current:        60,
				TrendDirection: TrendLowerBetter,
				Thresholds:     &Thresholds{Green: 10, Yellow: 30, Red: 50},
			},
			expectedStatus: StatusRed,
		},
		{
			name: "no thresholds",
			metric: Metric{
				Current:        50,
				TrendDirection: TrendHigherBetter,
				Thresholds:     nil,
			},
			expectedStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.metric.CalculateStatus()
			if status != tt.expectedStatus {
				t.Errorf("CalculateStatus() = %q, want %q", status, tt.expectedStatus)
			}
		})
	}
}

func TestMetricProgressToTarget(t *testing.T) {
	tests := []struct {
		name     string
		baseline float64
		current  float64
		target   float64
		want     float64
	}{
		{"zero progress", 0, 0, 100, 0},
		{"full progress", 0, 100, 100, 1.0},
		{"half progress", 0, 50, 100, 0.5},
		{"over target", 0, 150, 100, 1.0},
		{"negative progress", 100, 50, 0, 0.5},
		{"same baseline and target", 50, 50, 50, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metric{
				Baseline: tt.baseline,
				Current:  tt.current,
				Target:   tt.target,
			}
			got := m.ProgressToTarget()
			if got != tt.want {
				t.Errorf("ProgressToTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricValidation(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metric
		wantErrs int
	}{
		{
			name: "valid metric",
			metric: Metric{
				Name:       "Test Metric",
				Domain:     DomainSecurity,
				Stage:      StageRuntime,
				Category:   CategoryDetection,
				MetricType: MetricTypeRate,
			},
			wantErrs: 0,
		},
		{
			name: "missing required fields",
			metric: Metric{
				Name: "",
			},
			wantErrs: 5, // name, domain, stage, category, metricType
		},
		{
			name: "invalid domain",
			metric: Metric{
				Name:       "Test",
				Domain:     "invalid",
				Stage:      StageRuntime,
				Category:   CategoryDetection,
				MetricType: MetricTypeRate,
			},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.metric.Validate()
			if len(errs) != tt.wantErrs {
				t.Errorf("Validate() returned %d errors, want %d: %v", len(errs), tt.wantErrs, errs)
			}
		})
	}
}

func TestPRISMDocumentValidation(t *testing.T) {
	tests := []struct {
		name     string
		doc      PRISMDocument
		wantErrs bool
	}{
		{
			name: "valid document",
			doc: PRISMDocument{
				Metrics: []Metric{
					{
						ID:         "m1",
						Name:       "Test Metric",
						Domain:     DomainSecurity,
						Stage:      StageRuntime,
						Category:   CategoryDetection,
						MetricType: MetricTypeRate,
					},
				},
			},
			wantErrs: false,
		},
		{
			name: "empty metrics",
			doc: PRISMDocument{
				Metrics: []Metric{},
			},
			wantErrs: true,
		},
		{
			name: "duplicate metric IDs",
			doc: PRISMDocument{
				Metrics: []Metric{
					{ID: "m1", Name: "Test 1", Domain: DomainSecurity, Stage: StageRuntime, Category: CategoryDetection, MetricType: MetricTypeRate},
					{ID: "m1", Name: "Test 2", Domain: DomainSecurity, Stage: StageRuntime, Category: CategoryDetection, MetricType: MetricTypeRate},
				},
			},
			wantErrs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.doc.Validate()
			if errs.HasErrors() != tt.wantErrs {
				t.Errorf("Validate() hasErrors = %v, want %v: %v", errs.HasErrors(), tt.wantErrs, errs)
			}
		})
	}
}

func TestMaturityModel(t *testing.T) {
	model := NewMaturityModelWithCells()

	// Test default levels
	if len(model.Levels) != 5 {
		t.Errorf("expected 5 maturity levels, got %d", len(model.Levels))
	}

	// Test cells created for all domain/stage combinations
	expectedCells := len(AllDomains()) * len(AllStages())
	if len(model.Cells) != expectedCells {
		t.Errorf("expected %d cells, got %d", expectedCells, len(model.Cells))
	}

	// Test GetCell
	cell := model.GetCell(DomainSecurity, StageRuntime)
	if cell == nil {
		t.Error("GetCell returned nil for valid domain/stage")
	}
	if cell.Domain != DomainSecurity || cell.Stage != StageRuntime {
		t.Errorf("GetCell returned wrong cell: %v", cell)
	}

	// Test SetCellLevel
	err := model.SetCellLevel(DomainSecurity, StageRuntime, MaturityLevel4)
	if err != nil {
		t.Errorf("SetCellLevel failed: %v", err)
	}
	cell = model.GetCell(DomainSecurity, StageRuntime)
	if cell.CurrentLevel != MaturityLevel4 {
		t.Errorf("SetCellLevel didn't update level: got %d, want %d", cell.CurrentLevel, MaturityLevel4)
	}

	// Test CalculateMaturityScore
	score := cell.CalculateMaturityScore()
	expectedScore := 0.8 // Level 4 / 5
	if score != expectedScore {
		t.Errorf("CalculateMaturityScore() = %v, want %v", score, expectedScore)
	}
}

func TestCustomerAwarenessData(t *testing.T) {
	data := NewCustomerAwarenessData("2024-Q1")

	// Set counts
	_ = data.SetCount(AwarenessUnaware, 10)
	_ = data.SetCount(AwarenessAwareNotActing, 20)
	_ = data.SetCount(AwarenessAwareRemediating, 30)
	_ = data.SetCount(AwarenessAwareRemediated, 40)

	// Test TotalCount
	total := data.TotalCount()
	if total != 100 {
		t.Errorf("TotalCount() = %d, want 100", total)
	}

	// Test UnawareRate
	unawareRate := data.UnawareRate()
	if unawareRate != 0.1 {
		t.Errorf("UnawareRate() = %v, want 0.1", unawareRate)
	}

	// Test ProactiveDetectionRate
	detectionRate := data.ProactiveDetectionRate()
	if detectionRate != 0.9 {
		t.Errorf("ProactiveDetectionRate() = %v, want 0.9", detectionRate)
	}

	// Test ProactiveResolutionRate
	resolutionRate := data.ProactiveResolutionRate()
	if resolutionRate != 0.4 {
		t.Errorf("ProactiveResolutionRate() = %v, want 0.4", resolutionRate)
	}

	// Test AwarenessScore
	// Formula: (unaware * 0.0) + (notActing * 0.25) + (remediating * 0.5) + (remediated * 1.0)
	// = (0.1 * 0.0) + (0.2 * 0.25) + (0.3 * 0.5) + (0.4 * 1.0) = 0 + 0.05 + 0.15 + 0.4 = 0.6
	awarenessScore := data.AwarenessScore()
	if awarenessScore < 0.59 || awarenessScore > 0.61 {
		t.Errorf("AwarenessScore() = %v, want ~0.6", awarenessScore)
	}

	// Test Summary
	summary := data.Summary()
	if summary.TotalCustomers != 100 {
		t.Errorf("Summary.TotalCustomers = %d, want 100", summary.TotalCustomers)
	}
}

func TestPRISMScoreCalculation(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{
				ID:         "m1",
				Name:       "Security Coverage",
				Domain:     DomainSecurity,
				Stage:      StageRuntime,
				Category:   CategoryDetection,
				MetricType: MetricTypeCoverage,
				Baseline:   0,
				Current:    80,
				Target:     100,
			},
			{
				ID:         "m2",
				Name:       "Availability",
				Domain:     DomainOperations,
				Stage:      StageRuntime,
				Category:   CategoryReliability,
				MetricType: MetricTypeRate,
				Baseline:   99,
				Current:    99.9,
				Target:     99.99,
			},
		},
		Maturity: NewMaturityModelWithCells(),
	}

	// Set some maturity levels
	_ = doc.Maturity.SetCellLevel(DomainSecurity, StageRuntime, MaturityLevel4)
	_ = doc.Maturity.SetCellLevel(DomainOperations, StageRuntime, MaturityLevel3)

	// Calculate score
	config := DefaultScoreConfig()
	score := doc.CalculatePRISMScore(config, nil)

	// Verify score structure
	if score.Overall < 0 || score.Overall > 1 {
		t.Errorf("Overall score out of range: %v", score.Overall)
	}
	if score.SecurityScore < 0 || score.SecurityScore > 1 {
		t.Errorf("SecurityScore out of range: %v", score.SecurityScore)
	}
	if score.OperationsScore < 0 || score.OperationsScore > 1 {
		t.Errorf("OperationsScore out of range: %v", score.OperationsScore)
	}
	if score.Interpretation == "" {
		t.Error("Interpretation should not be empty")
	}

	// Test with awareness data
	awareness := NewCustomerAwarenessData("2024-Q1")
	_ = awareness.SetCount(AwarenessUnaware, 20)
	_ = awareness.SetCount(AwarenessAwareRemediating, 30)
	_ = awareness.SetCount(AwarenessAwareRemediated, 50)

	scoreWithAwareness := doc.CalculatePRISMScore(config, awareness)
	if scoreWithAwareness.AwarenessScore == 1.0 {
		t.Error("AwarenessScore should reflect awareness data")
	}
}

func TestInterpretScore(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{0.95, "Elite"},
		{0.9, "Elite"},
		{0.85, "Strong"},
		{0.75, "Strong"},
		{0.6, "Medium"},
		{0.5, "Medium"},
		{0.3, "Weak"},
		{0.25, "Weak"},
		{0.1, "Critical"},
		{0.0, "Critical"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := InterpretScore(tt.score)
			if got != tt.want {
				t.Errorf("InterpretScore(%v) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

func TestMaturityLevelName(t *testing.T) {
	tests := []struct {
		level int
		want  string
	}{
		{1, "Reactive"},
		{2, "Basic"},
		{3, "Defined"},
		{4, "Managed"},
		{5, "Optimizing"},
		{0, "Unknown"},
		{6, "Unknown"},
	}

	for _, tt := range tests {
		got := MaturityLevelName(tt.level)
		if got != tt.want {
			t.Errorf("MaturityLevelName(%d) = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestGetMetricsByDomain(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{Name: "M1", Domain: DomainSecurity},
			{Name: "M2", Domain: DomainOperations},
			{Name: "M3", Domain: DomainSecurity},
		},
	}

	securityMetrics := doc.GetMetricsByDomain(DomainSecurity)
	if len(securityMetrics) != 2 {
		t.Errorf("GetMetricsByDomain(security) returned %d metrics, want 2", len(securityMetrics))
	}

	opsMetrics := doc.GetMetricsByDomain(DomainOperations)
	if len(opsMetrics) != 1 {
		t.Errorf("GetMetricsByDomain(operations) returned %d metrics, want 1", len(opsMetrics))
	}
}

func TestGetMetricByID(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", Name: "Metric 1"},
			{ID: "m2", Name: "Metric 2"},
		},
	}

	m := doc.GetMetricByID("m1")
	if m == nil {
		t.Error("GetMetricByID returned nil for existing ID")
	}
	if m.Name != "Metric 1" {
		t.Errorf("GetMetricByID returned wrong metric: %v", m.Name)
	}

	m = doc.GetMetricByID("nonexistent")
	if m != nil {
		t.Error("GetMetricByID should return nil for nonexistent ID")
	}
}

func TestDataPointTimestampRoundTrip(t *testing.T) {
	// Test various timestamp formats and timezone handling
	tests := []struct {
		name      string
		timestamp time.Time
	}{
		{
			name:      "UTC",
			timestamp: time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:      "with nanoseconds",
			timestamp: time.Date(2024, 3, 15, 10, 30, 0, 123456789, time.UTC),
		},
		{
			name:      "epoch",
			timestamp: time.Unix(0, 0).UTC(),
		},
		{
			name:      "recent",
			timestamp: time.Now().UTC().Truncate(time.Second),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := DataPoint{
				Timestamp: tt.timestamp,
				Value:     42.5,
				Note:      "test note",
			}

			// Marshal to JSON
			data, err := json.Marshal(original)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// Unmarshal back
			var roundTripped DataPoint
			if err := json.Unmarshal(data, &roundTripped); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			// Compare (truncate to second precision since JSON doesn't preserve nanos precisely)
			originalTrunc := original.Timestamp.Truncate(time.Second)
			roundTrippedTrunc := roundTripped.Timestamp.Truncate(time.Second)
			if !originalTrunc.Equal(roundTrippedTrunc) {
				t.Errorf("timestamp mismatch: original=%v, roundTripped=%v", original.Timestamp, roundTripped.Timestamp)
			}

			if roundTripped.Value != original.Value {
				t.Errorf("value mismatch: original=%v, roundTripped=%v", original.Value, roundTripped.Value)
			}

			if roundTripped.Note != original.Note {
				t.Errorf("note mismatch: original=%q, roundTripped=%q", original.Note, roundTripped.Note)
			}
		})
	}
}

func TestDataPointTimestampJSONFormat(t *testing.T) {
	// Verify the JSON format is RFC3339
	dp := DataPoint{
		Timestamp: time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC),
		Value:     100,
	}

	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	// Should contain RFC3339 format timestamp
	if !strings.Contains(jsonStr, "2024-03-15T10:30:45Z") {
		t.Errorf("expected RFC3339 format, got: %s", jsonStr)
	}
}

func TestMeetsSLO(t *testing.T) {
	tests := []struct {
		name   string
		metric Metric
		want   bool
	}{
		{
			name:   "no SLO defined",
			metric: Metric{Current: 50},
			want:   true,
		},
		{
			name: "no operator defined",
			metric: Metric{
				Current: 50,
				SLO:     &SLO{Target: ">=90%"},
			},
			want: true,
		},
		{
			name: "GTE meets",
			metric: Metric{
				Current: 99.9,
				SLO:     &SLO{Operator: SLOOperatorGTE, Value: 99.9},
			},
			want: true,
		},
		{
			name: "GTE fails",
			metric: Metric{
				Current: 99.8,
				SLO:     &SLO{Operator: SLOOperatorGTE, Value: 99.9},
			},
			want: false,
		},
		{
			name: "LTE meets",
			metric: Metric{
				Current: 100,
				SLO:     &SLO{Operator: SLOOperatorLTE, Value: 100},
			},
			want: true,
		},
		{
			name: "LTE fails",
			metric: Metric{
				Current: 150,
				SLO:     &SLO{Operator: SLOOperatorLTE, Value: 100},
			},
			want: false,
		},
		{
			name: "GT meets",
			metric: Metric{
				Current: 100,
				SLO:     &SLO{Operator: SLOOperatorGT, Value: 99},
			},
			want: true,
		},
		{
			name: "GT fails at boundary",
			metric: Metric{
				Current: 99,
				SLO:     &SLO{Operator: SLOOperatorGT, Value: 99},
			},
			want: false,
		},
		{
			name: "LT meets",
			metric: Metric{
				Current: 5,
				SLO:     &SLO{Operator: SLOOperatorLT, Value: 10},
			},
			want: true,
		},
		{
			name: "EQ meets",
			metric: Metric{
				Current: 100,
				SLO:     &SLO{Operator: SLOOperatorEQ, Value: 100},
			},
			want: true,
		},
		{
			name: "EQ fails",
			metric: Metric{
				Current: 99,
				SLO:     &SLO{Operator: SLOOperatorEQ, Value: 100},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.metric.MeetsSLO()
			if got != tt.want {
				t.Errorf("MeetsSLO() = %v, want %v", got, tt.want)
			}
		})
	}
}
