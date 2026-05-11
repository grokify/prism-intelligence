package prism

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExampleFiles(t *testing.T) {
	examples := []struct {
		filename      string
		expectedCount int
		domain        string
	}{
		{"examples/prism-documents/operations-metrics.json", 8, DomainOperations},
	}

	for _, ex := range examples {
		t.Run(ex.filename, func(t *testing.T) {
			// Read file
			data, err := os.ReadFile(ex.filename)
			if err != nil {
				t.Fatalf("failed to read %s: %v", ex.filename, err)
			}

			// Unmarshal
			var doc PRISMDocument
			if err := json.Unmarshal(data, &doc); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", ex.filename, err)
			}

			// Validate
			errs := doc.Validate()
			if errs.HasErrors() {
				t.Errorf("validation failed for %s: %v", ex.filename, errs)
			}

			// Check metric count
			if len(doc.Metrics) != ex.expectedCount {
				t.Errorf("expected %d metrics, got %d", ex.expectedCount, len(doc.Metrics))
			}

			// Check all metrics are in expected domain
			for _, m := range doc.Metrics {
				if m.Domain != ex.domain {
					t.Errorf("metric %q has domain %q, expected %q", m.Name, m.Domain, ex.domain)
				}
			}

			// Calculate score
			score := doc.CalculatePRISMScore(nil, nil)
			if score.Overall < 0 || score.Overall > 1 {
				t.Errorf("score out of range: %v", score.Overall)
			}

			// Check score interpretation is valid
			validInterpretations := map[string]bool{
				"Elite": true, "Strong": true, "Medium": true, "Weak": true, "Critical": true,
			}
			if !validInterpretations[score.Interpretation] {
				t.Errorf("invalid interpretation: %q", score.Interpretation)
			}
		})
	}
}

func TestAllExampleFilesExist(t *testing.T) {
	expectedFiles := []string{
		"examples/prism-documents/operations-metrics.json",
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected example file %s does not exist", f)
		}
	}
}

func TestExampleFilesHaveSchema(t *testing.T) {
	pattern := "examples/prism-documents/*.json"
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob %s: %v", pattern, err)
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read %s: %v", f, err)
			}

			var doc PRISMDocument
			if err := json.Unmarshal(data, &doc); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", f, err)
			}

			if doc.Schema == "" {
				t.Errorf("example %s missing $schema field", f)
			}
		})
	}
}

func TestExampleMetricsHaveRequiredFields(t *testing.T) {
	pattern := "examples/prism-documents/*.json"
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob %s: %v", pattern, err)
	}

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("failed to read %s: %v", f, err)
			}

			var doc PRISMDocument
			if err := json.Unmarshal(data, &doc); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", f, err)
			}

			for i, m := range doc.Metrics {
				if m.ID == "" {
					t.Errorf("metric[%d] in %s missing ID", i, f)
				}
				if m.Name == "" {
					t.Errorf("metric[%d] in %s missing name", i, f)
				}
				if m.Domain == "" {
					t.Errorf("metric[%d] (%s) in %s missing domain", i, m.Name, f)
				}
				if m.Stage == "" {
					t.Errorf("metric[%d] (%s) in %s missing stage", i, m.Name, f)
				}
				if m.Category == "" {
					t.Errorf("metric[%d] (%s) in %s missing category", i, m.Name, f)
				}
				if m.MetricType == "" {
					t.Errorf("metric[%d] (%s) in %s missing metricType", i, m.Name, f)
				}
			}
		})
	}
}
