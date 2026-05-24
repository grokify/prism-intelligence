package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	capstack "github.com/grokify/prism-capability"
	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/maturity"
)

func TestGenerateDashboard(t *testing.T) {
	// Load the operations maturity model
	specFile := filepath.Join("..", "examples", "operations", "model.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	if dashboard.ID == "" {
		t.Error("Dashboard ID is empty")
	}

	if dashboard.Title == "" {
		t.Error("Dashboard title is empty")
	}

	if len(dashboard.Widgets) == 0 {
		t.Error("Dashboard has no widgets")
	}

	if len(dashboard.DataSources) == 0 {
		t.Error("Dashboard has no data sources")
	}

	// Verify we have different widget types
	widgetTypes := make(map[string]int)
	for _, w := range dashboard.Widgets {
		widgetTypes[w.Type]++
	}

	if widgetTypes["metric"] == 0 {
		t.Error("Dashboard has no metric widgets")
	}

	if widgetTypes["chart"] == 0 {
		t.Error("Dashboard has no chart widgets")
	}

	// Export to JSON for inspection
	jsonBytes, err := dashboard.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal dashboard: %v", err)
	}

	// Verify it's valid JSON
	var check map[string]any
	if err := json.Unmarshal(jsonBytes, &check); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Save to temp file for inspection
	tmpFile := filepath.Join(os.TempDir(), "prism-dashboard-test.json")
	if err := os.WriteFile(tmpFile, jsonBytes, 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated dashboard: %s (%d bytes)", tmpFile, len(jsonBytes))
	}
}

func TestGenerateSecurityDashboard(t *testing.T) {
	// Load the security maturity model
	specFile := filepath.Join("..", "examples", "security", "model.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	// Check SLI tables are generated
	hasSLITable := false
	for _, w := range dashboard.Widgets {
		if w.Type == "table" {
			hasSLITable = true
			break
		}
	}

	if !hasSLITable {
		t.Error("Dashboard has no SLI tables")
	}

	// Export to JSON
	jsonBytes, err := dashboard.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal dashboard: %v", err)
	}

	tmpFile := filepath.Join(os.TempDir(), "prism-security-dashboard.json")
	if err := os.WriteFile(tmpFile, jsonBytes, 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated security dashboard: %s (%d bytes)", tmpFile, len(jsonBytes))
	}
}

func TestEmptySpec(t *testing.T) {
	gen := NewGenerator(&maturity.Spec{
		Domains: map[string]*maturity.DomainModel{},
	})

	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate empty dashboard: %v", err)
	}

	if dashboard == nil {
		t.Error("Dashboard is nil")
	}
}

func TestNilSpec(t *testing.T) {
	gen := NewGenerator(nil)
	_, err := gen.Generate()
	if err == nil {
		t.Error("Expected error for nil spec")
	}
}

func TestGenerateHTML(t *testing.T) {
	specFile := filepath.Join("..", "examples", "operations", "model.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	html, err := dashboard.ToHTML(DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	if len(html) == 0 {
		t.Error("HTML output is empty")
	}

	// Check for required elements
	if !contains(html, "<!DOCTYPE html>") {
		t.Error("HTML missing DOCTYPE")
	}
	if !contains(html, "echarts") {
		t.Error("HTML missing ECharts")
	}
	if !contains(html, "dashboard") {
		t.Error("HTML missing dashboard data")
	}

	// Save for inspection
	tmpFile := filepath.Join(os.TempDir(), "prism-dashboard.html")
	if err := os.WriteFile(tmpFile, []byte(html), 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated HTML: %s (%d bytes)", tmpFile, len(html))
	}
}

func TestMaturityBullet(t *testing.T) {
	// Test NewMaturityBullet
	bullet := NewMaturityBullet("Availability", "99.5% uptime", 3.5, 5)

	if bullet.Title != "Availability" {
		t.Errorf("Expected title 'Availability', got '%s'", bullet.Title)
	}
	if bullet.Subtitle != "99.5% uptime" {
		t.Errorf("Expected subtitle '99.5%% uptime', got '%s'", bullet.Subtitle)
	}
	if len(bullet.Ranges) != 3 {
		t.Errorf("Expected 3 ranges, got %d", len(bullet.Ranges))
	}
	if len(bullet.Measures) != 1 || bullet.Measures[0] != 3.5 {
		t.Errorf("Expected measures [3.5], got %v", bullet.Measures)
	}
	if len(bullet.Markers) != 1 || bullet.Markers[0] != 5 {
		t.Errorf("Expected markers [5], got %v", bullet.Markers)
	}

	// Test MaturityLevel
	cases := []struct {
		value    float64
		expected string
	}{
		{5.0, "M5"},
		{4.5, "M4"},
		{4.0, "M4"},
		{3.5, "M3"},
		{2.0, "M2"},
		{1.0, "M1"},
		{0.5, "M0"},
	}

	for _, tc := range cases {
		got := MaturityLevel(tc.value)
		if got != tc.expected {
			t.Errorf("MaturityLevel(%v) = %s, want %s", tc.value, got, tc.expected)
		}
	}

	// Test MaturityStatus
	statusCases := []struct {
		value    float64
		expected string
	}{
		{5.0, "green"},
		{4.5, "yellow"},
		{3.0, "red"},
	}

	for _, tc := range statusCases {
		got := MaturityStatus(tc.value)
		if got != tc.expected {
			t.Errorf("MaturityStatus(%v) = %s, want %s", tc.value, got, tc.expected)
		}
	}
}

func TestMaturityBulletCSS(t *testing.T) {
	css := GetMaturityBulletCSS()

	// Check for required CSS classes
	requiredClasses := []string{
		".bullet",
		".range.s0",
		".range.s1",
		".range.s2",
		".measure.s0",
		"#fee2e2", // red
		"#fef3c7", // yellow
		"#dcfce7", // green
	}

	for _, class := range requiredClasses {
		if !contains(css, class) {
			t.Errorf("CSS missing '%s'", class)
		}
	}
}

func TestGenerateBullets(t *testing.T) {
	specFile := filepath.Join("..", "examples", "operations", "model.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	bulletData := gen.GenerateMaturityBullets()

	if len(bulletData.Bullets) == 0 {
		t.Error("Expected bullets to be generated")
	}

	// Check JSON serialization
	jsonBytes, err := bulletData.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal bullet data: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Error("JSON output is empty")
	}
}

func TestDashboardWithoutSLITypes(t *testing.T) {
	// Create a minimal spec without SLI types
	spec := &maturity.Spec{
		Metadata: &maturity.SpecMetadata{
			Name:        "Basic Model",
			Description: "A model without SLI types",
		},
		Domains: map[string]*maturity.DomainModel{
			"ops": {
				Name: "Operations",
				Levels: []maturity.Level{
					{
						Level: 1,
						Name:  "Reactive",
						Criteria: []maturity.Criterion{
							{ID: "c1", Name: "Basic monitoring"},
							{ID: "c2", Name: "Manual deployments"},
						},
					},
					{
						Level: 2,
						Name:  "Basic",
						Criteria: []maturity.Criterion{
							{ID: "c3", Name: "Automated alerts"},
						},
					},
				},
			},
		},
		// No SLIs defined - should fall back to flat list
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	// Should have widgets even without SLI types
	if len(dashboard.Widgets) == 0 {
		t.Error("Dashboard should have widgets")
	}

	// Check for bullet widgets (should be flat list, not methodology-grouped)
	hasBullet := false
	for _, w := range dashboard.Widgets {
		if w.Type == "bullet" {
			hasBullet = true
			// Should NOT have methodology in title (no RED/USE/Golden Signals)
			if contains(w.Title, "RED") || contains(w.Title, "USE") || contains(w.Title, "Golden") {
				t.Error("Dashboard without SLI types should not have methodology headers")
			}
		}
	}

	if !hasBullet {
		t.Error("Dashboard should have bullet widgets")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(len(s) >= len(substr) && (s == substr ||
			len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGenerateWithCapabilityStack(t *testing.T) {
	spec := &maturity.Spec{
		Metadata: &maturity.SpecMetadata{
			Name:        "Security Model",
			Description: "Test model with capability stack",
		},
		SLIs: map[string]*maturity.SLI{
			"sli-sast": {ID: "sli-sast", Name: "SAST Coverage", Category: "protect"},
			"sli-dast": {ID: "sli-dast", Name: "DAST Coverage", Category: "detect"},
		},
		Domains: map[string]*maturity.DomainModel{
			"security": {
				Name: "Security",
				Levels: []maturity.Level{
					{
						Level: 1,
						Name:  "Reactive",
						Criteria: []maturity.Criterion{
							{SLIID: "sli-sast", Operator: "gte", Target: 0},
						},
					},
					{
						Level: 2,
						Name:  "Basic",
						Criteria: []maturity.Criterion{
							{SLIID: "sli-sast", Operator: "gte", Target: 50},
						},
					},
				},
			},
		},
	}

	cs := &capstack.CapabilityStack{
		Metadata: capstack.Metadata{Name: "Security Stack"},
		Layers: []capstack.Layer{
			{ID: "code", Name: "Code", Order: 1},
			{ID: "runtime", Name: "Runtime", Order: 2},
		},
		Capabilities: []capstack.Capability{
			{
				ID:       "cap-sast",
				Name:     "SAST",
				LayerID:  "code",
				PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-sast"}},
			},
			{
				ID:       "cap-dast",
				Name:     "DAST",
				LayerID:  "runtime",
				PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-dast"}},
			},
		},
	}

	gen := NewGenerator(spec).WithCapabilityStack(cs)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	// Check for layer-specific widgets
	hasLayerOverview := false
	hasLayerMetric := false
	hasLayerBullet := false

	for _, w := range dashboard.Widgets {
		if w.ID == "layer-maturity-overview" {
			hasLayerOverview = true
		}
		if contains(w.ID, "layer-") && w.Type == "metric" {
			hasLayerMetric = true
		}
		if contains(w.ID, "bullet-layer-") {
			hasLayerBullet = true
		}
	}

	if !hasLayerOverview {
		t.Error("Dashboard with capStack should have layer maturity overview")
	}
	if !hasLayerMetric {
		t.Error("Dashboard with capStack should have layer metric cards")
	}
	if !hasLayerBullet {
		t.Error("Dashboard with capStack should have layer bullet charts")
	}
}

func TestGenerateWithoutCapabilityStack(t *testing.T) {
	spec := &maturity.Spec{
		Metadata: &maturity.SpecMetadata{
			Name:        "Basic Model",
			Description: "Test model without capability stack",
		},
		Domains: map[string]*maturity.DomainModel{
			"ops": {
				Name: "Operations",
				Levels: []maturity.Level{
					{Level: 1, Name: "Reactive"},
				},
			},
		},
	}

	// No capability stack
	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	// Should NOT have layer-specific widgets
	for _, w := range dashboard.Widgets {
		if w.ID == "layer-maturity-overview" {
			t.Error("Dashboard without capStack should not have layer maturity overview")
		}
		if contains(w.ID, "bullet-layer-") {
			t.Error("Dashboard without capStack should not have layer bullet charts")
		}
	}

	// Should still have domain widgets
	hasDomainMetric := false
	for _, w := range dashboard.Widgets {
		if contains(w.ID, "domain-") && w.Type == "metric" {
			hasDomainMetric = true
			break
		}
	}
	if !hasDomainMetric {
		t.Error("Dashboard should have domain metric cards")
	}
}

func TestAggregationMethods(t *testing.T) {
	spec := &maturity.Spec{
		Domains: map[string]*maturity.DomainModel{
			"security": {
				Name: "Security",
				Levels: []maturity.Level{
					{Level: 1, Criteria: []maturity.Criterion{{SLIID: "sli-a", Operator: "gte", Target: 0}}},
					{Level: 2, Criteria: []maturity.Criterion{{SLIID: "sli-a", Operator: "gte", Target: 50}}},
					{Level: 3, Criteria: []maturity.Criterion{{SLIID: "sli-b", Operator: "gte", Target: 80}}},
				},
			},
		},
	}

	cs := &capstack.CapabilityStack{
		Layers: []capstack.Layer{
			{ID: "code", Name: "Code", Order: 1},
		},
		Capabilities: []capstack.Capability{
			{ID: "cap-a", Name: "Cap A", LayerID: "code", PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-a"}}},
			{ID: "cap-b", Name: "Cap B", LayerID: "code", PRISMRef: &capstack.PRISMRef{SLIIDs: []string{"sli-b"}}},
		},
	}

	stateDoc := &prism.PRISMDocument{
		SLIState: prism.SLIStateMap{
			"sli-a": &prism.SLIState{Windows: map[string]*prism.WindowState{"30d": {Value: 75}}},  // M2
			"sli-b": &prism.SLIState{Windows: map[string]*prism.WindowState{"30d": {Value: 100}}}, // M3
		},
	}

	// Test MIN aggregation
	genMin := NewGenerator(spec).
		WithCapabilityStack(cs).
		WithStateDocument(stateDoc).
		WithAggregationMethod(AggregationMin)

	dashboardMin, err := genMin.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard with MIN: %v", err)
	}
	if dashboardMin == nil {
		t.Fatal("Dashboard is nil")
	}

	// Test AVG aggregation
	genAvg := NewGenerator(spec).
		WithCapabilityStack(cs).
		WithStateDocument(stateDoc).
		WithAggregationMethod(AggregationAvg)

	dashboardAvg, err := genAvg.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard with AVG: %v", err)
	}
	if dashboardAvg == nil {
		t.Fatal("Dashboard is nil")
	}

	// Both should have layer widgets
	hasLayerWidgetMin := false
	hasLayerWidgetAvg := false

	for _, w := range dashboardMin.Widgets {
		if w.ID == "layer-maturity-overview" {
			hasLayerWidgetMin = true
			break
		}
	}
	for _, w := range dashboardAvg.Widgets {
		if w.ID == "layer-maturity-overview" {
			hasLayerWidgetAvg = true
			break
		}
	}

	if !hasLayerWidgetMin {
		t.Error("MIN aggregation dashboard missing layer overview")
	}
	if !hasLayerWidgetAvg {
		t.Error("AVG aggregation dashboard missing layer overview")
	}
}

func TestGenerateHTMLWithCapabilityStack(t *testing.T) {
	spec := &maturity.Spec{
		Metadata: &maturity.SpecMetadata{Name: "Test Model"},
		Domains: map[string]*maturity.DomainModel{
			"ops": {Name: "Operations", Levels: []maturity.Level{{Level: 1}}},
		},
	}

	cs := &capstack.CapabilityStack{
		Layers:       []capstack.Layer{{ID: "code", Name: "Code", Order: 1}},
		Capabilities: []capstack.Capability{{ID: "cap-a", Name: "Cap A", LayerID: "code"}},
	}

	gen := NewGenerator(spec).WithCapabilityStack(cs)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	html, err := dashboard.ToHTML(DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	// Check for layer-related content
	if !contains(html, "layer-maturity-overview") {
		t.Error("HTML should contain layer-maturity-overview")
	}

	// Save for inspection
	tmpFile := filepath.Join(os.TempDir(), "prism-dashboard-capstack.html")
	if err := os.WriteFile(tmpFile, []byte(html), 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated HTML with capstack: %s (%d bytes)", tmpFile, len(html))
	}
}
