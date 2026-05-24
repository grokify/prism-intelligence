// Package render provides maturity-aware rendering for capability stacks.
package render

import (
	"fmt"

	caprender "github.com/grokify/prism-capability/render"
	"github.com/grokify/prism-maturity/dashboard"
)

// maturityColors maps maturity levels to badge colors.
var maturityColors = map[int]string{
	1: "#ef4444", // red - M1 Initial
	2: "#f59e0b", // amber - M2 Developing
	3: "#eab308", // yellow - M3 Defined
	4: "#22c55e", // green - M4 Managed
	5: "#3b82f6", // blue - M5 Optimizing
}

// maturityTextColors maps maturity levels to text colors.
var maturityTextColors = map[int]string{
	1: "#ffffff",
	2: "#000000",
	3: "#000000",
	4: "#ffffff",
	5: "#ffffff",
}

// BuildMaturityOverlay creates an OverlayProvider from maturity aggregation data.
// This allows capability stack renderers to display maturity levels as badges.
func BuildMaturityOverlay(agg *dashboard.MaturityAggregator) caprender.OverlayProvider {
	if agg == nil {
		return nil
	}

	overlays := make(caprender.OverlayProvider)

	for _, capMat := range agg.GetCapabilityMaturities() {
		level := int(capMat.AggregateLevel)
		if level < 1 {
			level = 1
		}
		if level > 5 {
			level = 5
		}

		badgeColor := maturityColors[level]
		if badgeColor == "" {
			badgeColor = "#6366f1" // indigo default
		}
		badgeTextColor := maturityTextColors[level]
		if badgeTextColor == "" {
			badgeTextColor = "#ffffff"
		}

		// Format badge text
		badgeText := formatMaturityBadge(capMat.AggregateLevel)

		// Build tooltip
		tooltip := fmt.Sprintf("Maturity: %s", badgeText)
		if len(capMat.SLIIDs) > 0 {
			tooltip += fmt.Sprintf(" (%d SLIs)", len(capMat.SLIIDs))
		}

		overlays[capMat.CapabilityID] = caprender.CapabilityOverlay{
			BadgeText:      badgeText,
			BadgeColor:     badgeColor,
			BadgeTextColor: badgeTextColor,
			TooltipExtra:   tooltip,
		}
	}

	return overlays
}

// formatMaturityBadge formats the maturity level for display.
func formatMaturityBadge(level float64) string {
	if level == float64(int(level)) {
		return fmt.Sprintf("M%d", int(level))
	}
	return fmt.Sprintf("M%.1f", level)
}

// BuildLayerOverlay creates an OverlayProvider for layer-level maturity.
// This is useful when rendering layer summaries rather than individual capabilities.
func BuildLayerOverlay(agg *dashboard.MaturityAggregator) caprender.OverlayProvider {
	if agg == nil {
		return nil
	}

	overlays := make(caprender.OverlayProvider)

	for _, layerMat := range agg.GetLayerMaturities() {
		level := int(layerMat.AggregateLevel)
		if level < 1 {
			level = 1
		}
		if level > 5 {
			level = 5
		}

		badgeColor := maturityColors[level]
		if badgeColor == "" {
			badgeColor = "#6366f1"
		}
		badgeTextColor := maturityTextColors[level]
		if badgeTextColor == "" {
			badgeTextColor = "#ffffff"
		}

		badgeText := formatMaturityBadge(layerMat.AggregateLevel)
		tooltip := fmt.Sprintf("Layer Maturity: %s (%d capabilities)",
			badgeText, len(layerMat.Capabilities))

		overlays[layerMat.LayerID] = caprender.CapabilityOverlay{
			BadgeText:      badgeText,
			BadgeColor:     badgeColor,
			BadgeTextColor: badgeTextColor,
			TooltipExtra:   tooltip,
		}
	}

	return overlays
}
