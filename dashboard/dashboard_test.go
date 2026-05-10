package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/prism/maturity"
)

func TestGenerateDashboard(t *testing.T) {
	// Load the operations maturity model
	specFile := filepath.Join("..", "maturity-models", "operations.json")
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
	specFile := filepath.Join("..", "maturity-models", "security.json")
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
	specFile := filepath.Join("..", "maturity-models", "operations.json")
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
