package dashboard

import (
	"testing"

	capstack "github.com/grokify/prism-capability"
	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/maturity"
)

func TestAggregate(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		method   AggregationMethod
		expected float64
	}{
		{
			name:     "empty values returns 1.0",
			values:   []float64{},
			method:   AggregationMin,
			expected: 1.0,
		},
		{
			name:     "single value min",
			values:   []float64{3.0},
			method:   AggregationMin,
			expected: 3.0,
		},
		{
			name:     "single value avg",
			values:   []float64{3.0},
			method:   AggregationAvg,
			expected: 3.0,
		},
		{
			name:     "multiple values min",
			values:   []float64{2.0, 4.0, 3.0},
			method:   AggregationMin,
			expected: 2.0,
		},
		{
			name:     "multiple values avg",
			values:   []float64{2.0, 4.0, 3.0},
			method:   AggregationAvg,
			expected: 3.0,
		},
		{
			name:     "fractional values min",
			values:   []float64{2.5, 3.5, 4.5},
			method:   AggregationMin,
			expected: 2.5,
		},
		{
			name:     "fractional values avg",
			values:   []float64{2.0, 3.0, 4.0},
			method:   AggregationAvg,
			expected: 3.0,
		},
		{
			name:     "unknown method defaults to min",
			values:   []float64{2.0, 4.0, 3.0},
			method:   AggregationMethod("unknown"),
			expected: 2.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Aggregate(tc.values, tc.method)
			if got != tc.expected {
				t.Errorf("Aggregate(%v, %s) = %v, want %v", tc.values, tc.method, got, tc.expected)
			}
		})
	}
}

func TestBuildSLIToCapabilityIndex(t *testing.T) {
	cs := &capstack.CapabilityStack{
		Capabilities: []capstack.Capability{
			{
				ID:      "cap-sast",
				Name:    "SAST",
				LayerID: "code",
				PRISMRef: &capstack.PRISMRef{
					SLIIDs: []string{"sli-sast-coverage", "sli-sast-findings"},
				},
			},
			{
				ID:      "cap-dast",
				Name:    "DAST",
				LayerID: "runtime",
				PRISMRef: &capstack.PRISMRef{
					SLIIDs: []string{"sli-dast-coverage"},
				},
			},
			{
				ID:      "cap-sbom",
				Name:    "SBOM",
				LayerID: "code",
				PRISMRef: &capstack.PRISMRef{
					SLIIDs: []string{"sli-sbom-coverage", "sli-sast-coverage"}, // shared SLI
				},
			},
			{
				ID:      "cap-no-prism",
				Name:    "No PRISM",
				LayerID: "code",
				// No PRISMRef
			},
		},
	}

	agg := NewMaturityAggregator(nil, cs, nil, AggregationMin)
	index := agg.sliCapIndex

	// Check sli-sast-coverage maps to both cap-sast and cap-sbom
	capIDs := index["sli-sast-coverage"]
	if len(capIDs) != 2 {
		t.Errorf("expected 2 capabilities for sli-sast-coverage, got %d", len(capIDs))
	}

	// Check sli-dast-coverage maps to cap-dast only
	capIDs = index["sli-dast-coverage"]
	if len(capIDs) != 1 || capIDs[0] != "cap-dast" {
		t.Errorf("expected [cap-dast] for sli-dast-coverage, got %v", capIDs)
	}

	// Check non-existent SLI
	capIDs = index["sli-nonexistent"]
	if len(capIDs) != 0 {
		t.Errorf("expected empty for sli-nonexistent, got %v", capIDs)
	}
}

func TestBuildSLIToCapabilityIndex_NilCapStack(t *testing.T) {
	agg := NewMaturityAggregator(nil, nil, nil, AggregationMin)
	index := agg.sliCapIndex

	if index == nil {
		t.Error("index should not be nil even with nil capStack")
	}
	if len(index) != 0 {
		t.Errorf("expected empty index, got %d entries", len(index))
	}
}

func TestCalculateCapabilityMaturity_NoSLIs(t *testing.T) {
	cs := &capstack.CapabilityStack{
		Capabilities: []capstack.Capability{
			{
				ID:      "cap-no-sli",
				Name:    "No SLIs",
				LayerID: "code",
			},
		},
	}

	agg := NewMaturityAggregator(nil, cs, nil, AggregationMin)
	cap := cs.Capabilities[0]
	result := agg.CalculateCapabilityMaturity(cap)

	if result.AggregateLevel != 1.0 {
		t.Errorf("expected 1.0 for capability with no SLIs, got %v", result.AggregateLevel)
	}
	if len(result.SLIIDs) != 0 {
		t.Errorf("expected empty SLIIDs, got %v", result.SLIIDs)
	}
}

func TestCalculateCapabilityMaturity_WithSLIs(t *testing.T) {
	spec := &maturity.Spec{
		Domains: map[string]*maturity.DomainModel{
			"security": {
				Name: "Security",
				Levels: []maturity.Level{
					{Level: 1, Criteria: []maturity.Criterion{{SLIID: "sli-sast", Operator: "gte", Target: 0}}},
					{Level: 2, Criteria: []maturity.Criterion{{SLIID: "sli-sast", Operator: "gte", Target: 50}}},
					{Level: 3, Criteria: []maturity.Criterion{{SLIID: "sli-sast", Operator: "gte", Target: 80}}},
				},
			},
		},
	}

	cs := &capstack.CapabilityStack{
		Capabilities: []capstack.Capability{
			{
				ID:      "cap-sast",
				Name:    "SAST",
				LayerID: "code",
				PRISMRef: &capstack.PRISMRef{
					SLIIDs: []string{"sli-sast"},
				},
			},
		},
	}

	stateDoc := &prism.PRISMDocument{
		SLIState: prism.SLIStateMap{
			"sli-sast": &prism.SLIState{
				Windows: map[string]*prism.WindowState{
					"30d": {Value: 75}, // Should achieve M2 (>=50) but not M3 (>=80)
				},
			},
		},
	}

	agg := NewMaturityAggregator(spec, cs, stateDoc, AggregationMin)
	cap := cs.Capabilities[0]
	result := agg.CalculateCapabilityMaturity(cap)

	if result.AggregateLevel != 2.0 {
		t.Errorf("expected 2.0 for 75%% coverage (M2 threshold), got %v", result.AggregateLevel)
	}
	if len(result.SLIIDs) != 1 {
		t.Errorf("expected 1 SLI, got %v", len(result.SLIIDs))
	}
}

func TestCalculateLayerMaturity_NoCapabilities(t *testing.T) {
	cs := &capstack.CapabilityStack{
		Layers: []capstack.Layer{
			{ID: "empty-layer", Name: "Empty Layer", Order: 1},
		},
		Capabilities: []capstack.Capability{}, // No capabilities
	}

	agg := NewMaturityAggregator(nil, cs, nil, AggregationMin)
	layer := cs.Layers[0]
	result := agg.CalculateLayerMaturity(layer)

	if result.AggregateLevel != 1.0 {
		t.Errorf("expected 1.0 for layer with no capabilities, got %v", result.AggregateLevel)
	}
	if len(result.Capabilities) != 0 {
		t.Errorf("expected empty capabilities, got %v", len(result.Capabilities))
	}
}

func TestCalculateLayerMaturity_WithCapabilities(t *testing.T) {
	spec := &maturity.Spec{
		Domains: map[string]*maturity.DomainModel{
			"security": {
				Name: "Security",
				Levels: []maturity.Level{
					{Level: 1, Criteria: []maturity.Criterion{{SLIID: "sli-a", Operator: "gte", Target: 0}}},
					{Level: 2, Criteria: []maturity.Criterion{{SLIID: "sli-a", Operator: "gte", Target: 50}}},
					{Level: 3, Criteria: []maturity.Criterion{{SLIID: "sli-a", Operator: "gte", Target: 80}}},
					{Level: 4, Criteria: []maturity.Criterion{{SLIID: "sli-b", Operator: "gte", Target: 90}}},
				},
			},
		},
	}

	cs := &capstack.CapabilityStack{
		Layers: []capstack.Layer{
			{ID: "code", Name: "Code", Order: 1},
		},
		Capabilities: []capstack.Capability{
			{
				ID:       "cap-a",
				Name:     "Cap A",
				LayerID:  "code",
				PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-a"}},
			},
			{
				ID:       "cap-b",
				Name:     "Cap B",
				LayerID:  "code",
				PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-b"}},
			},
		},
	}

	stateDoc := &prism.PRISMDocument{
		SLIState: prism.SLIStateMap{
			"sli-a": &prism.SLIState{Windows: map[string]*prism.WindowState{"30d": {Value: 85}}},  // M3
			"sli-b": &prism.SLIState{Windows: map[string]*prism.WindowState{"30d": {Value: 100}}}, // M4
		},
	}

	// Test with MIN aggregation
	aggMin := NewMaturityAggregator(spec, cs, stateDoc, AggregationMin)
	layer := cs.Layers[0]
	resultMin := aggMin.CalculateLayerMaturity(layer)

	if resultMin.AggregateLevel != 3.0 {
		t.Errorf("expected min=3.0 (min of M3 and M4), got %v", resultMin.AggregateLevel)
	}

	// Test with AVG aggregation
	aggAvg := NewMaturityAggregator(spec, cs, stateDoc, AggregationAvg)
	resultAvg := aggAvg.CalculateLayerMaturity(layer)

	if resultAvg.AggregateLevel != 3.5 {
		t.Errorf("expected avg=3.5 ((3+4)/2), got %v", resultAvg.AggregateLevel)
	}
}

func TestGetLayerMaturities_SortedByOrder(t *testing.T) {
	cs := &capstack.CapabilityStack{
		Layers: []capstack.Layer{
			{ID: "runtime", Name: "Runtime", Order: 3},
			{ID: "code", Name: "Code", Order: 1},
			{ID: "infra", Name: "Infrastructure", Order: 2},
		},
		Capabilities: []capstack.Capability{},
	}

	agg := NewMaturityAggregator(nil, cs, nil, AggregationMin)
	results := agg.GetLayerMaturities()

	if len(results) != 3 {
		t.Fatalf("expected 3 layers, got %d", len(results))
	}

	// Should be sorted by Order
	if results[0].LayerID != "code" {
		t.Errorf("expected first layer to be 'code' (order 1), got %s", results[0].LayerID)
	}
	if results[1].LayerID != "infra" {
		t.Errorf("expected second layer to be 'infra' (order 2), got %s", results[1].LayerID)
	}
	if results[2].LayerID != "runtime" {
		t.Errorf("expected third layer to be 'runtime' (order 3), got %s", results[2].LayerID)
	}
}

func TestGetLayerMaturities_NilCapStack(t *testing.T) {
	agg := NewMaturityAggregator(nil, nil, nil, AggregationMin)
	results := agg.GetLayerMaturities()

	if results != nil {
		t.Errorf("expected nil for nil capStack, got %v", results)
	}
}

func TestGetCapabilityMaturities(t *testing.T) {
	cs := &capstack.CapabilityStack{
		Capabilities: []capstack.Capability{
			{ID: "cap-a", Name: "Cap A", LayerID: "code"},
			{ID: "cap-b", Name: "Cap B", LayerID: "code"},
		},
		Foundational: []capstack.Capability{
			{ID: "cap-foundation", Name: "Foundation", LayerID: "cross-cutting"},
		},
	}

	agg := NewMaturityAggregator(nil, cs, nil, AggregationMin)
	results := agg.GetCapabilityMaturities()

	// Should include both regular and foundational capabilities
	if len(results) != 3 {
		t.Errorf("expected 3 capabilities (2 regular + 1 foundational), got %d", len(results))
	}
}

func TestAggregationMethodDefault(t *testing.T) {
	agg := NewMaturityAggregator(nil, nil, nil, "")

	if agg.method != AggregationMin {
		t.Errorf("expected default method to be 'min', got %s", agg.method)
	}
}
