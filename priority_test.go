package prism

import "testing"

func TestImportanceWeight(t *testing.T) {
	tests := []struct {
		importance string
		want       int
	}{
		{ImportanceCritical, 4},
		{ImportanceHigh, 3},
		{ImportanceMedium, 2},
		{ImportanceLow, 1},
		{"unknown", 2}, // defaults to medium
		{"", 2},        // defaults to medium
	}

	for _, tt := range tests {
		t.Run(tt.importance, func(t *testing.T) {
			got := ImportanceWeight(tt.importance)
			if got != tt.want {
				t.Errorf("ImportanceWeight(%q) = %d, want %d", tt.importance, got, tt.want)
			}
		})
	}
}

func TestDynamicPriorityWeight(t *testing.T) {
	tests := []struct {
		priority string
		want     int
	}{
		{PriorityP0, 4},
		{PriorityP1, 3},
		{PriorityP2, 2},
		{PriorityP3, 1},
		{"unknown", 2}, // defaults to P2
		{"", 2},        // defaults to P2
	}

	for _, tt := range tests {
		t.Run(tt.priority, func(t *testing.T) {
			got := DynamicPriorityWeight(tt.priority)
			if got != tt.want {
				t.Errorf("DynamicPriorityWeight(%q) = %d, want %d", tt.priority, got, tt.want)
			}
		})
	}
}

func TestCalculatePriority(t *testing.T) {
	tests := []struct {
		name         string
		importance   string
		currentLevel int
		targetLevel  int
		want         string
	}{
		// At or above target
		{"at target", ImportanceCritical, 3, 3, PriorityP3},
		{"above target", ImportanceCritical, 4, 3, PriorityP3},

		// Critical importance
		{"critical gap 1", ImportanceCritical, 2, 3, PriorityP1}, // 4 * 1 = 4 -> P1
		{"critical gap 2", ImportanceCritical, 1, 3, PriorityP0}, // 4 * 2 = 8 -> P0
		{"critical gap 3", ImportanceCritical, 1, 4, PriorityP0}, // 4 * 3 = 12 -> P0

		// High importance
		{"high gap 1", ImportanceHigh, 2, 3, PriorityP2}, // 3 * 1 = 3 -> P2
		{"high gap 2", ImportanceHigh, 1, 3, PriorityP1}, // 3 * 2 = 6 -> P1
		{"high gap 3", ImportanceHigh, 1, 4, PriorityP0}, // 3 * 3 = 9 -> P0

		// Medium importance
		{"medium gap 1", ImportanceMedium, 2, 3, PriorityP2}, // 2 * 1 = 2 -> P2
		{"medium gap 2", ImportanceMedium, 1, 3, PriorityP1}, // 2 * 2 = 4 -> P1
		{"medium gap 3", ImportanceMedium, 1, 4, PriorityP1}, // 2 * 3 = 6 -> P1
		{"medium gap 4", ImportanceMedium, 1, 5, PriorityP0}, // 2 * 4 = 8 -> P0

		// Low importance
		{"low gap 1", ImportanceLow, 2, 3, PriorityP3}, // 1 * 1 = 1 -> P3
		{"low gap 2", ImportanceLow, 1, 3, PriorityP2}, // 1 * 2 = 2 -> P2
		{"low gap 3", ImportanceLow, 1, 4, PriorityP2}, // 1 * 3 = 3 -> P2
		{"low gap 4", ImportanceLow, 1, 5, PriorityP1}, // 1 * 4 = 4 -> P1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePriority(tt.importance, tt.currentLevel, tt.targetLevel)
			if got != tt.want {
				t.Errorf("CalculatePriority(%q, %d, %d) = %q, want %q",
					tt.importance, tt.currentLevel, tt.targetLevel, got, tt.want)
			}
		})
	}
}

func TestPriorityRationale(t *testing.T) {
	tests := []struct {
		name         string
		importance   string
		currentLevel int
		targetLevel  int
		wantContains string
	}{
		{"at target", ImportanceCritical, 3, 3, "At or above target"},
		{"P0", ImportanceCritical, 1, 3, "Immediate action required"},
		{"P1", ImportanceHigh, 1, 3, "High priority improvement"},
		{"P2", ImportanceMedium, 2, 3, "Scheduled improvement"},
		{"P3", ImportanceLow, 2, 3, "Low priority enhancement"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriorityRationale(tt.importance, tt.currentLevel, tt.targetLevel)
			if !containsString(got, tt.wantContains) {
				t.Errorf("PriorityRationale(%q, %d, %d) = %q, want to contain %q",
					tt.importance, tt.currentLevel, tt.targetLevel, got, tt.wantContains)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAllImportanceLevels(t *testing.T) {
	levels := AllImportanceLevels()
	if len(levels) != 4 {
		t.Errorf("AllImportanceLevels() returned %d levels, want 4", len(levels))
	}
	// Verify descending order
	expected := []string{ImportanceCritical, ImportanceHigh, ImportanceMedium, ImportanceLow}
	for i, level := range levels {
		if level != expected[i] {
			t.Errorf("AllImportanceLevels()[%d] = %q, want %q", i, level, expected[i])
		}
	}
}

func TestAllPriorityLevels(t *testing.T) {
	levels := AllPriorityLevels()
	if len(levels) != 4 {
		t.Errorf("AllPriorityLevels() returned %d levels, want 4", len(levels))
	}
	// Verify descending order
	expected := []string{PriorityP0, PriorityP1, PriorityP2, PriorityP3}
	for i, level := range levels {
		if level != expected[i] {
			t.Errorf("AllPriorityLevels()[%d] = %q, want %q", i, level, expected[i])
		}
	}
}
