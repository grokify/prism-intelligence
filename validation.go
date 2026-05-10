package prism

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ValidationError represents a validation error with context.
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("%s: %s (value: %q)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

// HasErrors returns true if there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// ValidateDomain validates a domain value.
func ValidateDomain(domain string) error {
	if domain == "" {
		return errors.New("domain is required")
	}
	if !slices.Contains(AllDomains(), domain) {
		return fmt.Errorf("invalid domain %q, must be one of: %s", domain, strings.Join(AllDomains(), ", "))
	}
	return nil
}

// ValidateStage validates a stage value.
func ValidateStage(stage string) error {
	if stage == "" {
		return errors.New("stage is required")
	}
	if !slices.Contains(AllStages(), stage) {
		return fmt.Errorf("invalid stage %q, must be one of: %s", stage, strings.Join(AllStages(), ", "))
	}
	return nil
}

// ValidateCategory validates a category value.
func ValidateCategory(category string) error {
	if category == "" {
		return errors.New("category is required")
	}
	if !slices.Contains(AllCategories(), category) {
		return fmt.Errorf("invalid category %q, must be one of: %s", category, strings.Join(AllCategories(), ", "))
	}
	return nil
}

// ValidateLayer validates a layer value.
func ValidateLayer(layer string) error {
	if layer == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllLayers(), layer) {
		return fmt.Errorf("invalid layer %q, must be one of: %s", layer, strings.Join(AllLayers(), ", "))
	}
	return nil
}

// ValidateQualityVertical validates an ISO 25010 quality vertical value.
func ValidateQualityVertical(vertical string) error {
	if vertical == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllQualityVerticals(), vertical) {
		return fmt.Errorf("invalid quality vertical %q, must be one of: %s", vertical, strings.Join(AllQualityVerticals(), ", "))
	}
	return nil
}

// ValidateTeamType validates a team type value.
func ValidateTeamType(teamType string) error {
	if teamType == "" {
		return errors.New("team type is required")
	}
	if !slices.Contains(AllTeamTypes(), teamType) {
		return fmt.Errorf("invalid team type %q, must be one of: %s", teamType, strings.Join(AllTeamTypes(), ", "))
	}
	return nil
}

// ValidateMaturityLevel validates a maturity level value.
func ValidateMaturityLevel(level int) error {
	if level < MaturityLevel1 || level > MaturityLevel5 {
		return fmt.Errorf("invalid maturity level %d, must be between %d and %d", level, MaturityLevel1, MaturityLevel5)
	}
	return nil
}

// ValidateAwarenessState validates an awareness state value.
func ValidateAwarenessState(state string) error {
	if state == "" {
		return errors.New("awareness state is required")
	}
	if !slices.Contains(AllAwarenessStates(), state) {
		return fmt.Errorf("invalid awareness state %q, must be one of: %s", state, strings.Join(AllAwarenessStates(), ", "))
	}
	return nil
}

// ValidateFramework validates a framework value.
func ValidateFramework(framework string) error {
	if framework == "" {
		return errors.New("framework is required")
	}
	if !slices.Contains(AllFrameworks(), framework) {
		return fmt.Errorf("invalid framework %q, must be one of: %s", framework, strings.Join(AllFrameworks(), ", "))
	}
	return nil
}

// ValidateMetricType validates a metric type value.
func ValidateMetricType(metricType string) error {
	if metricType == "" {
		return errors.New("metric type is required")
	}
	if !slices.Contains(AllMetricTypes(), metricType) {
		return fmt.Errorf("invalid metric type %q, must be one of: %s", metricType, strings.Join(AllMetricTypes(), ", "))
	}
	return nil
}

// ValidateTrendDirection validates a trend direction value.
func ValidateTrendDirection(trend string) error {
	if trend == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllTrendDirections(), trend) {
		return fmt.Errorf("invalid trend direction %q, must be one of: %s", trend, strings.Join(AllTrendDirections(), ", "))
	}
	return nil
}

// ValidateStatus validates a status value.
func ValidateStatus(status string) error {
	if status == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllStatuses(), status) {
		return fmt.Errorf("invalid status %q, must be one of: %s", status, strings.Join(AllStatuses(), ", "))
	}
	return nil
}

// ValidateWindow validates an SLO window value.
func ValidateWindow(window string) error {
	if window == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllWindows(), window) {
		return fmt.Errorf("invalid window %q, must be one of: %s", window, strings.Join(AllWindows(), ", "))
	}
	return nil
}

// ValidateGoalStatus validates a goal status value.
func ValidateGoalStatus(status string) error {
	if status == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllGoalStatuses(), status) {
		return fmt.Errorf("invalid goal status %q, must be one of: %s", status, strings.Join(AllGoalStatuses(), ", "))
	}
	return nil
}

// ValidatePhaseStatus validates a phase status value.
func ValidatePhaseStatus(status string) error {
	if status == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllPhaseStatuses(), status) {
		return fmt.Errorf("invalid phase status %q, must be one of: %s", status, strings.Join(AllPhaseStatuses(), ", "))
	}
	return nil
}

// ValidateQuarter validates a quarter value.
func ValidateQuarter(quarter string) error {
	if quarter == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllQuarters(), quarter) {
		return fmt.Errorf("invalid quarter %q, must be one of: %s", quarter, strings.Join(AllQuarters(), ", "))
	}
	return nil
}

// ValidateInitiativeStatus validates an initiative status value.
func ValidateInitiativeStatus(status string) error {
	if status == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllInitiativeStatuses(), status) {
		return fmt.Errorf("invalid initiative status %q, must be one of: %s", status, strings.Join(AllInitiativeStatuses(), ", "))
	}
	return nil
}

// ValidateSLOOperator validates an SLO operator value.
func ValidateSLOOperator(operator string) error {
	if operator == "" {
		return nil // Optional field
	}
	validOps := []string{SLOOperatorGTE, SLOOperatorLTE, SLOOperatorEQ, SLOOperatorGT, SLOOperatorLT}
	if !slices.Contains(validOps, operator) {
		return fmt.Errorf("invalid SLO operator %q, must be one of: %s", operator, strings.Join(validOps, ", "))
	}
	return nil
}

// ValidateSLIType validates an SLI type value.
func ValidateSLIType(sliType string) error {
	if sliType == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllSLITypes(), sliType) {
		return fmt.Errorf("invalid SLI type %q, must be one of: %s", sliType, strings.Join(AllSLITypes(), ", "))
	}
	return nil
}

// ValidateMethodology validates an observability methodology value.
func ValidateMethodology(methodology string) error {
	if methodology == "" {
		return nil // Optional field
	}
	if !slices.Contains(AllMethodologies(), methodology) {
		return fmt.Errorf("invalid methodology %q, must be one of: %s", methodology, strings.Join(AllMethodologies(), ", "))
	}
	return nil
}

// Validate validates a Metric and returns validation errors.
func (m *Metric) Validate() ValidationErrors {
	var errs ValidationErrors

	if m.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	if err := ValidateDomain(m.Domain); err != nil {
		errs = append(errs, ValidationError{Field: "domain", Value: m.Domain, Message: err.Error()})
	}

	if err := ValidateStage(m.Stage); err != nil {
		errs = append(errs, ValidationError{Field: "stage", Value: m.Stage, Message: err.Error()})
	}

	if err := ValidateCategory(m.Category); err != nil {
		errs = append(errs, ValidationError{Field: "category", Value: m.Category, Message: err.Error()})
	}

	if err := ValidateLayer(m.Layer); err != nil {
		errs = append(errs, ValidationError{Field: "layer", Value: m.Layer, Message: err.Error()})
	}

	if err := ValidateQualityVertical(m.QualityVertical); err != nil {
		errs = append(errs, ValidationError{Field: "qualityVertical", Value: m.QualityVertical, Message: err.Error()})
	}

	if err := ValidateMetricType(m.MetricType); err != nil {
		errs = append(errs, ValidationError{Field: "metricType", Value: m.MetricType, Message: err.Error()})
	}

	if err := ValidateTrendDirection(m.TrendDirection); err != nil {
		errs = append(errs, ValidationError{Field: "trendDirection", Value: m.TrendDirection, Message: err.Error()})
	}

	if err := ValidateStatus(m.Status); err != nil {
		errs = append(errs, ValidationError{Field: "status", Value: m.Status, Message: err.Error()})
	}

	// Validate SLO window if present
	if m.SLO != nil {
		if err := ValidateWindow(m.SLO.Window); err != nil {
			errs = append(errs, ValidationError{Field: "slo.window", Value: m.SLO.Window, Message: err.Error()})
		}
	}

	// Validate SLI type if present
	if m.SLI != nil {
		if err := ValidateSLIType(m.SLI.SLIType); err != nil {
			errs = append(errs, ValidationError{Field: "sli.sliType", Value: m.SLI.SLIType, Message: err.Error()})
		}
	}

	// Validate framework mappings
	for i, fm := range m.FrameworkMappings {
		if err := ValidateFramework(fm.Framework); err != nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("frameworkMappings[%d].framework", i),
				Value:   fm.Framework,
				Message: err.Error(),
			})
		}
		if fm.Reference == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("frameworkMappings[%d].reference", i),
				Message: "is required",
			})
		}
	}

	return errs
}

// Validate validates the entire PRISMDocument.
func (doc *PRISMDocument) Validate() ValidationErrors {
	var errs ValidationErrors

	if len(doc.Metrics) == 0 {
		errs = append(errs, ValidationError{Field: "metrics", Message: "at least one metric is required"})
	}

	// Validate each metric
	for i, m := range doc.Metrics {
		metricErrs := m.Validate()
		for _, e := range metricErrs {
			e.Field = fmt.Sprintf("metrics[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	// Validate maturity model if present
	if doc.Maturity != nil {
		maturityErrs := doc.Maturity.Validate()
		for _, e := range maturityErrs {
			e.Field = "maturity." + e.Field
			errs = append(errs, e)
		}
	}

	// Check for duplicate metric IDs
	seenIDs := make(map[string]int)
	for i, m := range doc.Metrics {
		if m.ID != "" {
			if prevIdx, exists := seenIDs[m.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("metrics[%d].id", i),
					Value:   m.ID,
					Message: fmt.Sprintf("duplicate ID, also used at metrics[%d]", prevIdx),
				})
			}
			seenIDs[m.ID] = i
		}
	}

	// Validate layers
	seenLayerIDs := make(map[string]int)
	for i, layer := range doc.Layers {
		layerErrs := layer.Validate()
		for _, e := range layerErrs {
			e.Field = fmt.Sprintf("layers[%d].%s", i, e.Field)
			errs = append(errs, e)
		}

		// Check for duplicate layer IDs
		if layer.ID != "" {
			if prevIdx, exists := seenLayerIDs[layer.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("layers[%d].id", i),
					Value:   layer.ID,
					Message: fmt.Sprintf("duplicate ID, also used at layers[%d]", prevIdx),
				})
			}
			seenLayerIDs[layer.ID] = i
		}
	}

	// Validate teams
	seenTeamIDs := make(map[string]int)
	for i, team := range doc.Teams {
		teamErrs := team.Validate(doc)
		for _, e := range teamErrs {
			e.Field = fmt.Sprintf("teams[%d].%s", i, e.Field)
			errs = append(errs, e)
		}

		// Check for duplicate team IDs
		if team.ID != "" {
			if prevIdx, exists := seenTeamIDs[team.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("teams[%d].id", i),
					Value:   team.ID,
					Message: fmt.Sprintf("duplicate ID, also used at teams[%d]", prevIdx),
				})
			}
			seenTeamIDs[team.ID] = i
		}
	}

	// Validate services
	seenServiceIDs := make(map[string]int)
	for i, service := range doc.Services {
		serviceErrs := service.Validate(doc)
		for _, e := range serviceErrs {
			e.Field = fmt.Sprintf("services[%d].%s", i, e.Field)
			errs = append(errs, e)
		}

		// Check for duplicate service IDs
		if service.ID != "" {
			if prevIdx, exists := seenServiceIDs[service.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("services[%d].id", i),
					Value:   service.ID,
					Message: fmt.Sprintf("duplicate ID, also used at services[%d]", prevIdx),
				})
			}
			seenServiceIDs[service.ID] = i
		}
	}

	// Validate metric service references
	for i, m := range doc.Metrics {
		if m.ServiceID != "" && doc.GetServiceByID(m.ServiceID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("metrics[%d].serviceId", i),
				Value:   m.ServiceID,
				Message: "references non-existent service ID",
			})
		}
	}

	// Validate OKR metric references
	for i, okr := range doc.OKRs {
		for j, metricID := range okr.MetricIDs {
			if doc.GetMetricByID(metricID) == nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("okrs[%d].metricIds[%d]", i, j),
					Value:   metricID,
					Message: "references non-existent metric ID",
				})
			}
		}
	}

	// Validate initiative metric references
	for i, init := range doc.Initiatives {
		for j, metricID := range init.MetricIDs {
			if doc.GetMetricByID(metricID) == nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("initiatives[%d].metricIds[%d]", i, j),
					Value:   metricID,
					Message: "references non-existent metric ID",
				})
			}
		}

		// Validate initiative status
		if err := ValidateInitiativeStatus(init.Status); err != nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("initiatives[%d].status", i),
				Value:   init.Status,
				Message: err.Error(),
			})
		}

		// Validate goal references
		for j, goalID := range init.GoalIDs {
			if doc.GetGoalByID(goalID) == nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("initiatives[%d].goalIds[%d]", i, j),
					Value:   goalID,
					Message: "references non-existent goal ID",
				})
			}
		}

		// Validate phase reference
		if init.PhaseID != "" && doc.GetPhaseByID(init.PhaseID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("initiatives[%d].phaseId", i),
				Value:   init.PhaseID,
				Message: "references non-existent phase ID",
			})
		}

		// Validate service reference
		if init.ServiceID != "" && doc.GetServiceByID(init.ServiceID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("initiatives[%d].serviceId", i),
				Value:   init.ServiceID,
				Message: "references non-existent service ID",
			})
		}
	}

	// Validate goals
	for i, goal := range doc.Goals {
		goalErrs := goal.Validate(doc)
		for _, e := range goalErrs {
			e.Field = fmt.Sprintf("goals[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	// Check for duplicate goal IDs
	seenGoalIDs := make(map[string]int)
	for i, g := range doc.Goals {
		if g.ID != "" {
			if prevIdx, exists := seenGoalIDs[g.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("goals[%d].id", i),
					Value:   g.ID,
					Message: fmt.Sprintf("duplicate ID, also used at goals[%d]", prevIdx),
				})
			}
			seenGoalIDs[g.ID] = i
		}
	}

	// Validate phases
	for i, phase := range doc.Phases {
		phaseErrs := phase.Validate(doc)
		for _, e := range phaseErrs {
			e.Field = fmt.Sprintf("phases[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	// Check for duplicate phase IDs
	seenPhaseIDs := make(map[string]int)
	for i, p := range doc.Phases {
		if p.ID != "" {
			if prevIdx, exists := seenPhaseIDs[p.ID]; exists {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("phases[%d].id", i),
					Value:   p.ID,
					Message: fmt.Sprintf("duplicate ID, also used at phases[%d]", prevIdx),
				})
			}
			seenPhaseIDs[p.ID] = i
		}
	}

	return errs
}

// Validate validates a Goal and returns validation errors.
func (g *Goal) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if g.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	if err := ValidateGoalStatus(g.Status); err != nil {
		errs = append(errs, ValidationError{Field: "status", Value: g.Status, Message: err.Error()})
	}

	if g.CurrentLevel != 0 {
		if err := ValidateMaturityLevel(g.CurrentLevel); err != nil {
			errs = append(errs, ValidationError{Field: "currentLevel", Value: fmt.Sprintf("%d", g.CurrentLevel), Message: err.Error()})
		}
	}

	if g.TargetLevel != 0 {
		if err := ValidateMaturityLevel(g.TargetLevel); err != nil {
			errs = append(errs, ValidationError{Field: "targetLevel", Value: fmt.Sprintf("%d", g.TargetLevel), Message: err.Error()})
		}
	}

	// Validate maturity model if present
	if g.MaturityModel != nil {
		mmErrs := g.MaturityModel.Validate(doc)
		for _, e := range mmErrs {
			e.Field = "maturityModel." + e.Field
			errs = append(errs, e)
		}
	}

	return errs
}

// Validate validates a GoalMaturityModel and returns validation errors.
func (gmm *GoalMaturityModel) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if len(gmm.Levels) == 0 {
		errs = append(errs, ValidationError{Field: "levels", Message: "at least one level is required"})
	}

	seenLevels := make(map[int]bool)
	for i, level := range gmm.Levels {
		// Validate level number
		if err := ValidateMaturityLevel(level.Level); err != nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("levels[%d].level", i),
				Value:   fmt.Sprintf("%d", level.Level),
				Message: err.Error(),
			})
		}

		// Check for duplicate levels
		if seenLevels[level.Level] {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("levels[%d].level", i),
				Value:   fmt.Sprintf("%d", level.Level),
				Message: "duplicate level number",
			})
		}
		seenLevels[level.Level] = true

		// Validate level definition
		levelErrs := level.Validate(doc)
		for _, e := range levelErrs {
			e.Field = fmt.Sprintf("levels[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	return errs
}

// Validate validates a GoalMaturityLevel and returns validation errors.
func (gml *GoalMaturityLevel) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if gml.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	// Validate SLO requirements
	for i, sloReq := range gml.RequiredSLOs {
		if sloReq.MetricID == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("requiredSLOs[%d].metricId", i),
				Message: "is required",
			})
		} else if doc != nil && doc.GetMetricByID(sloReq.MetricID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("requiredSLOs[%d].metricId", i),
				Value:   sloReq.MetricID,
				Message: "references non-existent metric ID",
			})
		}
	}

	// Validate metric criteria
	for i, criterion := range gml.MetricCriteria {
		criterionErrs := criterion.Validate(doc)
		for _, e := range criterionErrs {
			e.Field = fmt.Sprintf("metricCriteria[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	return errs
}

// Validate validates a MetricCriterion and returns validation errors.
func (mc *MetricCriterion) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if mc.MetricID == "" {
		errs = append(errs, ValidationError{Field: "metricId", Message: "is required"})
	} else if doc != nil && doc.GetMetricByID(mc.MetricID) == nil {
		errs = append(errs, ValidationError{
			Field:   "metricId",
			Value:   mc.MetricID,
			Message: "references non-existent metric ID",
		})
	}

	if err := ValidateSLOOperator(mc.Operator); err != nil {
		errs = append(errs, ValidationError{Field: "operator", Value: mc.Operator, Message: err.Error()})
	} else if mc.Operator == "" {
		errs = append(errs, ValidationError{Field: "operator", Message: "is required"})
	}

	return errs
}

// Validate validates a Phase and returns validation errors.
func (p *Phase) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if p.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	if p.StartDate == "" {
		errs = append(errs, ValidationError{Field: "startDate", Message: "is required"})
	}

	if p.EndDate == "" {
		errs = append(errs, ValidationError{Field: "endDate", Message: "is required"})
	}

	if err := ValidatePhaseStatus(p.Status); err != nil {
		errs = append(errs, ValidationError{Field: "status", Value: p.Status, Message: err.Error()})
	}

	if err := ValidateQuarter(p.Quarter); err != nil {
		errs = append(errs, ValidationError{Field: "quarter", Value: p.Quarter, Message: err.Error()})
	}

	// Validate goal targets
	for i, target := range p.GoalTargets {
		if target.GoalID == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("goalTargets[%d].goalId", i),
				Message: "is required",
			})
		} else if doc != nil && doc.GetGoalByID(target.GoalID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("goalTargets[%d].goalId", i),
				Value:   target.GoalID,
				Message: "references non-existent goal ID",
			})
		}

		if err := ValidateMaturityLevel(target.EnterLevel); err != nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("goalTargets[%d].enterLevel", i),
				Value:   fmt.Sprintf("%d", target.EnterLevel),
				Message: err.Error(),
			})
		}

		if err := ValidateMaturityLevel(target.ExitLevel); err != nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("goalTargets[%d].exitLevel", i),
				Value:   fmt.Sprintf("%d", target.ExitLevel),
				Message: err.Error(),
			})
		}
	}

	// Validate swimlanes
	for i, sw := range p.Swimlanes {
		swErrs := sw.Validate(doc)
		for _, e := range swErrs {
			e.Field = fmt.Sprintf("swimlanes[%d].%s", i, e.Field)
			errs = append(errs, e)
		}
	}

	return errs
}

// Validate validates a Swimlane and returns validation errors.
func (sw *Swimlane) Validate(doc *PRISMDocument) ValidationErrors {
	var errs ValidationErrors

	if sw.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "is required"})
	}

	if sw.Domain != "" {
		if err := ValidateDomain(sw.Domain); err != nil {
			errs = append(errs, ValidationError{Field: "domain", Value: sw.Domain, Message: err.Error()})
		}
	}

	if sw.Stage != "" {
		if err := ValidateStage(sw.Stage); err != nil {
			errs = append(errs, ValidationError{Field: "stage", Value: sw.Stage, Message: err.Error()})
		}
	}

	// Validate initiative references
	for i, initID := range sw.InitiativeIDs {
		if doc != nil && doc.GetInitiativeByID(initID) == nil {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("initiativeIds[%d]", i),
				Value:   initID,
				Message: "references non-existent initiative ID",
			})
		}
	}

	return errs
}
