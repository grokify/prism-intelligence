package prism

import (
	"fmt"

	core "github.com/grokify/prism-core"
)

// TeamType constants imported from prism-core (Team Topologies).
const (
	TeamTypeStreamAligned = core.TeamTypeStreamAligned
	TeamTypePlatform      = core.TeamTypePlatform
	TeamTypeEnabling      = core.TeamTypeEnabling
	TeamTypeOverlay       = core.TeamTypeOverlay
)

// AllTeamTypes returns all valid team type values.
func AllTeamTypes() []string {
	return []string{
		TeamTypeStreamAligned,
		TeamTypePlatform,
		TeamTypeEnabling,
		TeamTypeOverlay,
	}
}

// ValidTeamType checks if a team type is valid.
func ValidTeamType(teamType string) bool {
	return core.ValidTeamType(teamType)
}

// Team represents a team in the organization following Team Topologies patterns.
type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"` // stream_aligned, platform, enabling, overlay

	// Domain accountability (for overlay/enabling teams)
	Domain string `json:"domain,omitempty"` // security, operations, quality

	// Layer accountability (which layers this team is responsible for)
	LayerAccountability []string `json:"layerAccountability,omitempty"` // code, infra, runtime

	// Service ownership (for stream-aligned teams)
	ServiceIDs []string `json:"serviceIds,omitempty"`

	// Contact information
	Owner string `json:"owner,omitempty"`
	Slack string `json:"slack,omitempty"`
	Email string `json:"email,omitempty"`
}

// Validate validates a Team and returns validation errors.
func (t *Team) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if t.ID == "" {
		errs = append(errs, ValidationError{Field: "id", Message: "is required"})
	}

	if t.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	if !ValidTeamType(t.Type) {
		errs = append(errs, ValidationError{Field: "type", Value: t.Type, Message: "invalid team type"})
	}

	// Validate domain if specified
	if t.Domain != "" {
		if !ValidDomain(t.Domain) {
			errs = append(errs, ValidationError{Field: "domain", Value: t.Domain, Message: "invalid domain"})
		}
	}

	// Validate layer accountability
	for i, layer := range t.LayerAccountability {
		if !ValidLayer(layer) {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("layerAccountability[%d]", i),
				Value:   layer,
				Message: "invalid layer",
			})
		}
	}

	// Validate service references if document is provided
	if doc != nil {
		for i, serviceID := range t.ServiceIDs {
			if doc.GetServiceByID(serviceID) == nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("serviceIds[%d]", i),
					Value:   serviceID,
					Message: "references non-existent service ID",
				})
			}
		}
	}

	return errs
}
