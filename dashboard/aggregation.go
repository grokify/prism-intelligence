// Package dashboard provides HTML dashboard generation for PRISM maturity models.
package dashboard

import (
	"sort"

	capstack "github.com/grokify/prism-capability"
	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/maturity"
)

// AggregationMethod defines how to aggregate maturity levels.
type AggregationMethod string

const (
	// AggregationMin uses the minimum value (most conservative).
	AggregationMin AggregationMethod = "min"
	// AggregationAvg uses the average value.
	AggregationAvg AggregationMethod = "avg"
)

// SLIToCapabilityIndex maps SLI IDs to capability IDs that reference them.
type SLIToCapabilityIndex map[string][]string

// CapabilityMaturity holds aggregated maturity data for a single capability.
type CapabilityMaturity struct {
	CapabilityID   string            `json:"capabilityId"`
	CapabilityName string            `json:"capabilityName"`
	LayerID        string            `json:"layerId"`
	SLIIDs         []string          `json:"sliIds"`
	SLILevels      []float64         `json:"sliLevels"`
	AggregateLevel float64           `json:"aggregateLevel"`
	Method         AggregationMethod `json:"method"`
}

// LayerMaturity holds aggregated maturity data for a capability layer.
type LayerMaturity struct {
	LayerID        string               `json:"layerId"`
	LayerName      string               `json:"layerName"`
	Order          int                  `json:"order"`
	Capabilities   []CapabilityMaturity `json:"capabilities"`
	AggregateLevel float64              `json:"aggregateLevel"`
	Method         AggregationMethod    `json:"method"`
}

// MaturityAggregator calculates aggregate maturity across capabilities and layers.
type MaturityAggregator struct {
	spec     *maturity.Spec
	capStack *capstack.CapabilityStack
	stateDoc *prism.PRISMDocument
	method   AggregationMethod

	// Cached index
	sliCapIndex SLIToCapabilityIndex
}

// NewMaturityAggregator creates a new maturity aggregator.
func NewMaturityAggregator(
	spec *maturity.Spec,
	cs *capstack.CapabilityStack,
	stateDoc *prism.PRISMDocument,
	method AggregationMethod,
) *MaturityAggregator {
	if method == "" {
		method = AggregationMin
	}
	agg := &MaturityAggregator{
		spec:     spec,
		capStack: cs,
		stateDoc: stateDoc,
		method:   method,
	}
	agg.sliCapIndex = agg.BuildSLIToCapabilityIndex()
	return agg
}

// BuildSLIToCapabilityIndex creates a reverse index from SLI IDs to capability IDs.
func (a *MaturityAggregator) BuildSLIToCapabilityIndex() SLIToCapabilityIndex {
	index := make(SLIToCapabilityIndex)
	if a.capStack == nil {
		return index
	}

	for _, cap := range a.capStack.AllCapabilities() {
		if cap.PRISMRef != nil && len(cap.PRISMRef.SLIIDs) > 0 {
			for _, sliID := range cap.PRISMRef.SLIIDs {
				index[sliID] = append(index[sliID], cap.ID)
			}
		}
	}
	return index
}

// getSLILevel returns the maturity level for an SLI.
// Uses the state document if available, otherwise returns 1.0 (M1).
func (a *MaturityAggregator) getSLILevel(sliID string) float64 {
	if a.stateDoc == nil || a.spec == nil {
		return 1.0
	}

	// Get SLI value from state document
	var sliValue float64
	var hasValue bool

	if a.stateDoc.SLIState != nil {
		if state, ok := a.stateDoc.SLIState[sliID]; ok && state != nil {
			// Try 30d window first, then any available window
			if state.Windows != nil {
				if ws, ok := state.Windows["30d"]; ok && ws != nil {
					sliValue = ws.Value
					hasValue = true
				} else {
					for _, ws := range state.Windows {
						if ws != nil {
							sliValue = ws.Value
							hasValue = true
							break
						}
					}
				}
			}
		}
	}

	if !hasValue {
		return 1.0
	}

	// Calculate maturity level based on criteria
	highestLevel := 1.0
	for _, domain := range a.spec.Domains {
		for _, level := range domain.Levels {
			for _, criterion := range level.Criteria {
				if criterion.SLIID != sliID {
					continue
				}
				if criterion.CheckMet(sliValue) {
					levelNum := float64(level.Level)
					if levelNum > highestLevel {
						highestLevel = levelNum
					}
				}
			}
		}
	}

	return highestLevel
}

// CalculateCapabilityMaturity computes the aggregate maturity for a capability.
func (a *MaturityAggregator) CalculateCapabilityMaturity(cap capstack.Capability) CapabilityMaturity {
	result := CapabilityMaturity{
		CapabilityID:   cap.ID,
		CapabilityName: cap.Name,
		LayerID:        cap.LayerID,
		Method:         a.method,
		SLIIDs:         []string{},
		SLILevels:      []float64{},
	}

	if cap.PRISMRef == nil || len(cap.PRISMRef.SLIIDs) == 0 {
		result.AggregateLevel = 1.0
		return result
	}

	result.SLIIDs = cap.PRISMRef.SLIIDs

	// Collect SLI levels
	var levels []float64
	for _, sliID := range cap.PRISMRef.SLIIDs {
		level := a.getSLILevel(sliID)
		levels = append(levels, level)
	}
	result.SLILevels = levels

	// Aggregate
	result.AggregateLevel = Aggregate(levels, a.method)
	return result
}

// CalculateLayerMaturity computes the aggregate maturity for a layer.
func (a *MaturityAggregator) CalculateLayerMaturity(layer capstack.Layer) LayerMaturity {
	result := LayerMaturity{
		LayerID:      layer.ID,
		LayerName:    layer.Name,
		Order:        layer.Order,
		Method:       a.method,
		Capabilities: []CapabilityMaturity{},
	}

	if a.capStack == nil {
		result.AggregateLevel = 1.0
		return result
	}

	// Get capabilities for this layer
	caps := a.capStack.CapabilitiesByLayer(layer.ID)
	if len(caps) == 0 {
		result.AggregateLevel = 1.0
		return result
	}

	// Calculate maturity for each capability
	var capLevels []float64
	for _, cap := range caps {
		capMat := a.CalculateCapabilityMaturity(cap)
		result.Capabilities = append(result.Capabilities, capMat)
		capLevels = append(capLevels, capMat.AggregateLevel)
	}

	// Aggregate capability levels to layer level
	result.AggregateLevel = Aggregate(capLevels, a.method)
	return result
}

// GetLayerMaturities returns maturity data for all layers, sorted by Order.
func (a *MaturityAggregator) GetLayerMaturities() []LayerMaturity {
	if a.capStack == nil {
		return nil
	}

	var results []LayerMaturity
	for _, layer := range a.capStack.Layers {
		results = append(results, a.CalculateLayerMaturity(layer))
	}

	// Sort by Order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Order < results[j].Order
	})

	return results
}

// GetCapabilityMaturities returns maturity data for all capabilities.
func (a *MaturityAggregator) GetCapabilityMaturities() []CapabilityMaturity {
	if a.capStack == nil {
		return nil
	}

	var results []CapabilityMaturity
	for _, cap := range a.capStack.AllCapabilities() {
		results = append(results, a.CalculateCapabilityMaturity(cap))
	}

	return results
}

// Aggregate computes an aggregate value from a slice of values using the specified method.
func Aggregate(values []float64, method AggregationMethod) float64 {
	if len(values) == 0 {
		return 1.0 // Default to M1 if no values
	}

	switch method {
	case AggregationMin:
		minVal := values[0]
		for _, v := range values[1:] {
			if v < minVal {
				minVal = v
			}
		}
		return minVal

	case AggregationAvg:
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values))

	default:
		// Default to min for safety
		minVal := values[0]
		for _, v := range values[1:] {
			if v < minVal {
				minVal = v
			}
		}
		return minVal
	}
}
