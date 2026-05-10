package prism

import (
	"slices"
	"testing"
)

func TestAllSLITypes(t *testing.T) {
	types := AllSLITypes()
	if len(types) != 8 {
		t.Errorf("AllSLITypes() returned %d types, want 8", len(types))
	}

	expected := []string{
		SLITypeAvailability,
		SLITypeLatency,
		SLITypeErrorRate,
		SLITypeThroughput,
		SLITypeSaturation,
		SLITypeUtilization,
		SLITypeQuality,
		SLITypeFreshness,
	}

	for _, e := range expected {
		if !slices.Contains(types, e) {
			t.Errorf("AllSLITypes() missing %q", e)
		}
	}
}

func TestAllMethodologies(t *testing.T) {
	methodologies := AllMethodologies()
	if len(methodologies) != 3 {
		t.Errorf("AllMethodologies() returned %d methodologies, want 3", len(methodologies))
	}

	expected := []string{
		MethodologyGoldenSignals,
		MethodologyRED,
		MethodologyUSE,
	}

	for _, e := range expected {
		if !slices.Contains(methodologies, e) {
			t.Errorf("AllMethodologies() missing %q", e)
		}
	}
}

func TestGoldenSignalsSLITypes(t *testing.T) {
	types := GoldenSignalsSLITypes()
	if len(types) != 4 {
		t.Errorf("GoldenSignalsSLITypes() returned %d types, want 4", len(types))
	}

	expected := []string{SLITypeLatency, SLITypeThroughput, SLITypeErrorRate, SLITypeSaturation}
	for _, e := range expected {
		if !slices.Contains(types, e) {
			t.Errorf("GoldenSignalsSLITypes() missing %q", e)
		}
	}
}

func TestREDSLITypes(t *testing.T) {
	types := REDSLITypes()
	if len(types) != 3 {
		t.Errorf("REDSLITypes() returned %d types, want 3", len(types))
	}

	expected := []string{SLITypeThroughput, SLITypeErrorRate, SLITypeLatency}
	for _, e := range expected {
		if !slices.Contains(types, e) {
			t.Errorf("REDSLITypes() missing %q", e)
		}
	}
}

func TestUSESLITypes(t *testing.T) {
	types := USESLITypes()
	if len(types) != 3 {
		t.Errorf("USESLITypes() returned %d types, want 3", len(types))
	}

	expected := []string{SLITypeUtilization, SLITypeSaturation, SLITypeErrorRate}
	for _, e := range expected {
		if !slices.Contains(types, e) {
			t.Errorf("USESLITypes() missing %q", e)
		}
	}
}

func TestSLITypesForMethodology(t *testing.T) {
	tests := []struct {
		methodology string
		wantLen     int
	}{
		{MethodologyGoldenSignals, 4},
		{MethodologyRED, 3},
		{MethodologyUSE, 3},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.methodology, func(t *testing.T) {
			types := SLITypesForMethodology(tt.methodology)
			gotLen := len(types)
			if gotLen != tt.wantLen {
				t.Errorf("SLITypesForMethodology(%q) returned %d types, want %d", tt.methodology, gotLen, tt.wantLen)
			}
		})
	}
}

func TestMethodologiesForSLIType(t *testing.T) {
	tests := []struct {
		sliType  string
		contains []string
	}{
		{SLITypeLatency, []string{MethodologyGoldenSignals, MethodologyRED}},
		{SLITypeThroughput, []string{MethodologyGoldenSignals, MethodologyRED}},
		{SLITypeErrorRate, []string{MethodologyGoldenSignals, MethodologyRED, MethodologyUSE}},
		{SLITypeSaturation, []string{MethodologyGoldenSignals, MethodologyUSE}},
		{SLITypeUtilization, []string{MethodologyUSE}},
		{SLITypeAvailability, nil},
		{SLITypeQuality, nil},
		{SLITypeFreshness, nil},
	}

	for _, tt := range tests {
		t.Run(tt.sliType, func(t *testing.T) {
			methodologies := MethodologiesForSLIType(tt.sliType)
			for _, expected := range tt.contains {
				if !slices.Contains(methodologies, expected) {
					t.Errorf("MethodologiesForSLIType(%q) missing %q", tt.sliType, expected)
				}
			}
			if tt.contains == nil && len(methodologies) != 0 {
				t.Errorf("MethodologiesForSLIType(%q) returned %v, want empty", tt.sliType, methodologies)
			}
		})
	}
}

func TestSLIIsGoldenSignal(t *testing.T) {
	tests := []struct {
		name string
		sli  *SLI
		want bool
	}{
		{"nil SLI", nil, false},
		{"empty SLI type", &SLI{}, false},
		{"latency", &SLI{SLIType: SLITypeLatency}, true},
		{"throughput", &SLI{SLIType: SLITypeThroughput}, true},
		{"error_rate", &SLI{SLIType: SLITypeErrorRate}, true},
		{"saturation", &SLI{SLIType: SLITypeSaturation}, true},
		{"availability", &SLI{SLIType: SLITypeAvailability}, false},
		{"utilization", &SLI{SLIType: SLITypeUtilization}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sli.IsGoldenSignal()
			if got != tt.want {
				t.Errorf("SLI.IsGoldenSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSLIIsRED(t *testing.T) {
	tests := []struct {
		name string
		sli  *SLI
		want bool
	}{
		{"nil SLI", nil, false},
		{"empty SLI type", &SLI{}, false},
		{"throughput", &SLI{SLIType: SLITypeThroughput}, true},
		{"error_rate", &SLI{SLIType: SLITypeErrorRate}, true},
		{"latency", &SLI{SLIType: SLITypeLatency}, true},
		{"saturation", &SLI{SLIType: SLITypeSaturation}, false},
		{"utilization", &SLI{SLIType: SLITypeUtilization}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sli.IsRED()
			if got != tt.want {
				t.Errorf("SLI.IsRED() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSLIIsUSE(t *testing.T) {
	tests := []struct {
		name string
		sli  *SLI
		want bool
	}{
		{"nil SLI", nil, false},
		{"empty SLI type", &SLI{}, false},
		{"utilization", &SLI{SLIType: SLITypeUtilization}, true},
		{"saturation", &SLI{SLIType: SLITypeSaturation}, true},
		{"error_rate", &SLI{SLIType: SLITypeErrorRate}, true},
		{"latency", &SLI{SLIType: SLITypeLatency}, false},
		{"throughput", &SLI{SLIType: SLITypeThroughput}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sli.IsUSE()
			if got != tt.want {
				t.Errorf("SLI.IsUSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSLIMethodologies(t *testing.T) {
	tests := []struct {
		name     string
		sli      *SLI
		contains []string
	}{
		{"nil SLI", nil, nil},
		{"empty SLI type", &SLI{}, nil},
		{"error_rate - all three", &SLI{SLIType: SLITypeErrorRate}, []string{MethodologyGoldenSignals, MethodologyRED, MethodologyUSE}},
		{"latency - two", &SLI{SLIType: SLITypeLatency}, []string{MethodologyGoldenSignals, MethodologyRED}},
		{"utilization - one", &SLI{SLIType: SLITypeUtilization}, []string{MethodologyUSE}},
		{"availability - none", &SLI{SLIType: SLITypeAvailability}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sli.Methodologies()
			for _, expected := range tt.contains {
				if !slices.Contains(got, expected) {
					t.Errorf("SLI.Methodologies() = %v, missing %q", got, expected)
				}
			}
			if tt.contains == nil && len(got) != 0 {
				t.Errorf("SLI.Methodologies() = %v, want empty", got)
			}
		})
	}
}

func TestValidateSLIType(t *testing.T) {
	tests := []struct {
		name    string
		sliType string
		wantErr bool
	}{
		{"empty - optional", "", false},
		{"valid availability", SLITypeAvailability, false},
		{"valid latency", SLITypeLatency, false},
		{"valid error_rate", SLITypeErrorRate, false},
		{"valid throughput", SLITypeThroughput, false},
		{"valid saturation", SLITypeSaturation, false},
		{"valid utilization", SLITypeUtilization, false},
		{"valid quality", SLITypeQuality, false},
		{"valid freshness", SLITypeFreshness, false},
		{"invalid", "invalid_type", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSLIType(tt.sliType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSLIType(%q) error = %v, wantErr %v", tt.sliType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMethodology(t *testing.T) {
	tests := []struct {
		name        string
		methodology string
		wantErr     bool
	}{
		{"empty - optional", "", false},
		{"valid GOLDEN_SIGNALS", MethodologyGoldenSignals, false},
		{"valid RED", MethodologyRED, false},
		{"valid USE", MethodologyUSE, false},
		{"invalid", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMethodology(tt.methodology)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMethodology(%q) error = %v, wantErr %v", tt.methodology, err, tt.wantErr)
			}
		})
	}
}

func TestMetricValidateSLIType(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metric
		wantErrs int
	}{
		{
			name: "valid metric with SLI type",
			metric: Metric{
				Name:       "Test",
				Domain:     DomainOperations,
				Stage:      StageRuntime,
				Category:   CategoryReliability,
				MetricType: MetricTypeLatency,
				SLI:        &SLI{SLIType: SLITypeLatency},
			},
			wantErrs: 0,
		},
		{
			name: "invalid SLI type",
			metric: Metric{
				Name:       "Test",
				Domain:     DomainOperations,
				Stage:      StageRuntime,
				Category:   CategoryReliability,
				MetricType: MetricTypeLatency,
				SLI:        &SLI{SLIType: "invalid"},
			},
			wantErrs: 1,
		},
		{
			name: "no SLI - valid",
			metric: Metric{
				Name:       "Test",
				Domain:     DomainOperations,
				Stage:      StageRuntime,
				Category:   CategoryReliability,
				MetricType: MetricTypeLatency,
			},
			wantErrs: 0,
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

func TestAnalyzeSLICoverage(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", Name: "M1", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", Name: "M2", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m3", Name: "M3", SLI: &SLI{SLIType: SLITypeErrorRate}},
			{ID: "m4", Name: "M4"}, // No SLI
		},
	}

	coverage := doc.AnalyzeSLICoverage()

	if coverage.TotalMetrics != 4 {
		t.Errorf("TotalMetrics = %d, want 4", coverage.TotalMetrics)
	}
	if coverage.MetricsWithSLI != 3 {
		t.Errorf("MetricsWithSLI = %d, want 3", coverage.MetricsWithSLI)
	}
	if coverage.ByType[SLITypeLatency] != 2 {
		t.Errorf("ByType[latency] = %d, want 2", coverage.ByType[SLITypeLatency])
	}
	if coverage.ByType[SLITypeErrorRate] != 1 {
		t.Errorf("ByType[error_rate] = %d, want 1", coverage.ByType[SLITypeErrorRate])
	}
	if coverage.CoverageRatio != 0.75 {
		t.Errorf("CoverageRatio = %f, want 0.75", coverage.CoverageRatio)
	}
	// Should have 6 missing types (all except latency and error_rate)
	if len(coverage.MissingTypes) != 6 {
		t.Errorf("MissingTypes count = %d, want 6", len(coverage.MissingTypes))
	}
}

func TestAnalyzeSLICoverageByLayer(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", Layer: LayerRuntime, SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", Layer: LayerRuntime, SLI: &SLI{SLIType: SLITypeErrorRate}},
			{ID: "m3", Layer: LayerCode},
		},
	}

	coverage := doc.AnalyzeSLICoverageByLayer()

	runtimeCoverage := coverage[LayerRuntime]
	if runtimeCoverage.TotalMetrics != 2 {
		t.Errorf("Runtime TotalMetrics = %d, want 2", runtimeCoverage.TotalMetrics)
	}
	if runtimeCoverage.MetricsWithSLI != 2 {
		t.Errorf("Runtime MetricsWithSLI = %d, want 2", runtimeCoverage.MetricsWithSLI)
	}

	codeCoverage := coverage[LayerCode]
	if codeCoverage.TotalMetrics != 1 {
		t.Errorf("Code TotalMetrics = %d, want 1", codeCoverage.TotalMetrics)
	}
	if codeCoverage.MetricsWithSLI != 0 {
		t.Errorf("Code MetricsWithSLI = %d, want 0", codeCoverage.MetricsWithSLI)
	}
}

func TestAnalyzeMethodologyCoverage(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", SLI: &SLI{SLIType: SLITypeThroughput}},
			{ID: "m3", SLI: &SLI{SLIType: SLITypeErrorRate}},
			// Missing saturation for complete Golden Signals
		},
	}

	// Test Golden Signals coverage
	gsCoverage := doc.AnalyzeMethodologyCoverage(MethodologyGoldenSignals)
	if gsCoverage == nil {
		t.Fatal("AnalyzeMethodologyCoverage returned nil for Golden Signals")
	}
	if len(gsCoverage.RequiredTypes) != 4 {
		t.Errorf("Golden Signals RequiredTypes = %d, want 4", len(gsCoverage.RequiredTypes))
	}
	if len(gsCoverage.CoveredTypes) != 3 {
		t.Errorf("Golden Signals CoveredTypes = %d, want 3", len(gsCoverage.CoveredTypes))
	}
	if len(gsCoverage.MissingTypes) != 1 {
		t.Errorf("Golden Signals MissingTypes = %d, want 1", len(gsCoverage.MissingTypes))
	}
	if gsCoverage.MissingTypes[0] != SLITypeSaturation {
		t.Errorf("Golden Signals missing type = %q, want %q", gsCoverage.MissingTypes[0], SLITypeSaturation)
	}
	if gsCoverage.IsComplete {
		t.Error("Golden Signals IsComplete = true, want false")
	}

	// Test RED coverage (should be complete)
	redCoverage := doc.AnalyzeMethodologyCoverage(MethodologyRED)
	if !redCoverage.IsComplete {
		t.Error("RED IsComplete = false, want true")
	}
	if redCoverage.CoverageRatio != 1.0 {
		t.Errorf("RED CoverageRatio = %f, want 1.0", redCoverage.CoverageRatio)
	}

	// Test unknown methodology
	unknownCoverage := doc.AnalyzeMethodologyCoverage("UNKNOWN")
	if unknownCoverage != nil {
		t.Error("AnalyzeMethodologyCoverage should return nil for unknown methodology")
	}
}

func TestAnalyzeAllMethodologyCoverage(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", SLI: &SLI{SLIType: SLITypeThroughput}},
			{ID: "m3", SLI: &SLI{SLIType: SLITypeErrorRate}},
			{ID: "m4", SLI: &SLI{SLIType: SLITypeSaturation}},
			{ID: "m5", SLI: &SLI{SLIType: SLITypeUtilization}},
		},
	}

	allCoverage := doc.AnalyzeAllMethodologyCoverage()

	if len(allCoverage) != 3 {
		t.Errorf("AnalyzeAllMethodologyCoverage returned %d entries, want 3", len(allCoverage))
	}

	for _, methodology := range AllMethodologies() {
		if _, exists := allCoverage[methodology]; !exists {
			t.Errorf("AnalyzeAllMethodologyCoverage missing %q", methodology)
		}
	}

	// All three methodologies should be complete with these metrics
	if !allCoverage[MethodologyGoldenSignals].IsComplete {
		t.Error("Golden Signals should be complete")
	}
	if !allCoverage[MethodologyRED].IsComplete {
		t.Error("RED should be complete")
	}
	if !allCoverage[MethodologyUSE].IsComplete {
		t.Error("USE should be complete")
	}
}

func TestGetMetricsBySLIType(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", Name: "M1", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", Name: "M2", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m3", Name: "M3", SLI: &SLI{SLIType: SLITypeErrorRate}},
			{ID: "m4", Name: "M4"},
		},
	}

	latencyMetrics := doc.GetMetricsBySLIType(SLITypeLatency)
	if len(latencyMetrics) != 2 {
		t.Errorf("GetMetricsBySLIType(latency) returned %d metrics, want 2", len(latencyMetrics))
	}

	errorMetrics := doc.GetMetricsBySLIType(SLITypeErrorRate)
	if len(errorMetrics) != 1 {
		t.Errorf("GetMetricsBySLIType(error_rate) returned %d metrics, want 1", len(errorMetrics))
	}

	availMetrics := doc.GetMetricsBySLIType(SLITypeAvailability)
	if len(availMetrics) != 0 {
		t.Errorf("GetMetricsBySLIType(availability) returned %d metrics, want 0", len(availMetrics))
	}
}

func TestGetMetricsByMethodology(t *testing.T) {
	doc := &PRISMDocument{
		Metrics: []Metric{
			{ID: "m1", SLI: &SLI{SLIType: SLITypeLatency}},
			{ID: "m2", SLI: &SLI{SLIType: SLITypeThroughput}},
			{ID: "m3", SLI: &SLI{SLIType: SLITypeErrorRate}},
			{ID: "m4", SLI: &SLI{SLIType: SLITypeSaturation}},
			{ID: "m5", SLI: &SLI{SLIType: SLITypeUtilization}},
			{ID: "m6", SLI: &SLI{SLIType: SLITypeAvailability}},
		},
	}

	gsMetrics := doc.GetMetricsByMethodology(MethodologyGoldenSignals)
	// Golden Signals includes: latency, throughput, error_rate, saturation
	if len(gsMetrics) != 4 {
		t.Errorf("GetMetricsByMethodology(GOLDEN_SIGNALS) returned %d metrics, want 4", len(gsMetrics))
	}

	redMetrics := doc.GetMetricsByMethodology(MethodologyRED)
	// RED includes: throughput, error_rate, latency
	if len(redMetrics) != 3 {
		t.Errorf("GetMetricsByMethodology(RED) returned %d metrics, want 3", len(redMetrics))
	}

	useMetrics := doc.GetMetricsByMethodology(MethodologyUSE)
	// USE includes: utilization, saturation, error_rate
	if len(useMetrics) != 3 {
		t.Errorf("GetMetricsByMethodology(USE) returned %d metrics, want 3", len(useMetrics))
	}

	unknownMetrics := doc.GetMetricsByMethodology("UNKNOWN")
	if unknownMetrics != nil {
		t.Errorf("GetMetricsByMethodology(UNKNOWN) returned %v, want nil", unknownMetrics)
	}
}
