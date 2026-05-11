package maturity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMarp(t *testing.T) {
	specFile := "../examples/organization/model.json"
	outputFile := filepath.Join(os.TempDir(), "organization-maturity-generated.md")

	if err := GenerateMarp(specFile, outputFile); err != nil {
		t.Fatalf("Failed to generate Marp: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify basic structure
	contentStr := string(content)

	// Should have Marp frontmatter
	if !strings.Contains(contentStr, "marp: true") {
		t.Error("Missing marp frontmatter")
	}

	// Should have all domains
	domains := []string{"Security", "Operational Excellence", "Quality", "Product", "AI"}
	for _, domain := range domains {
		if !strings.Contains(contentStr, domain+" Domain") {
			t.Errorf("Missing domain section: %s", domain)
		}
	}

	// Should NOT have PRISM references
	if strings.Contains(contentStr, "PRISM") {
		t.Error("Found PRISM reference - should be removed")
	}

	info, _ := os.Stat(outputFile)
	t.Logf("Generated Marp: %s (%d bytes)", outputFile, info.Size())
	t.Logf("Slide count: %d", strings.Count(contentStr, "\n---\n")+1)

	// Copy to docs/presentations for review (constant path, safe)
	docsOutput := filepath.Clean("docs/presentations/organization-maturity-generated.md")
	if err := os.WriteFile(docsOutput, content, 0600); err != nil { //nolint:gosec // constant path
		t.Logf("Could not copy to docs: %v", err)
	} else {
		t.Logf("Copied to: %s", docsOutput)
	}
}
