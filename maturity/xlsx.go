package maturity

import (
	"fmt"
	"sort"
	"strings"

	"github.com/plexusone/omniframe"
	"github.com/xuri/excelize/v2"
)

// XLSXGenerator generates Excel reports from maturity specs.
type XLSXGenerator struct {
	spec *Spec
	file *excelize.File
}

// NewXLSXGenerator creates a new XLSX generator.
func NewXLSXGenerator(spec *Spec) *XLSXGenerator {
	return &XLSXGenerator{
		spec: spec,
		file: excelize.NewFile(),
	}
}

// Generate creates the XLSX file with all sheets.
func (g *XLSXGenerator) Generate() error {
	// Create sheets
	if err := g.createRequirementsSheet(); err != nil {
		return fmt.Errorf("failed to create requirements sheet: %w", err)
	}

	if err := g.createSLOsSheet(); err != nil {
		return fmt.Errorf("failed to create SLOs sheet: %w", err)
	}

	if err := g.createThresholdMatrixSheet(); err != nil {
		return fmt.Errorf("failed to create threshold matrix sheet: %w", err)
	}

	if err := g.createFrameworkMappingsSheet(); err != nil {
		return fmt.Errorf("failed to create framework mappings sheet: %w", err)
	}

	if err := g.createProgressSheet(); err != nil {
		return fmt.Errorf("failed to create progress sheet: %w", err)
	}

	if err := g.createLevelDefinitionsSheet(); err != nil {
		return fmt.Errorf("failed to create level definitions sheet: %w", err)
	}

	// Delete the default "Sheet1"
	_ = g.file.DeleteSheet("Sheet1")

	return nil
}

// SaveAs saves the XLSX file to the specified path.
func (g *XLSXGenerator) SaveAs(filename string) error {
	return g.file.SaveAs(filename)
}

// createRequirementsSheet creates the Requirements (Enablers) sheet.
func (g *XLSXGenerator) createRequirementsSheet() error {
	sheetName := "Requirements"
	index, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}
	g.file.SetActiveSheet(index)

	// Headers
	headers := []string{
		"ID", "Domain", "Level", "Name", "Description", "Type",
		"Layer", "Team", "Effort", "Status", "Enables", "Depends On",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "L1", headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]
		for _, level := range domain.Levels {
			for _, e := range level.Enablers {
				g.setCellValue(sheetName, fmt.Sprintf("A%d", row), e.ID)
				g.setCellValue(sheetName, fmt.Sprintf("B%d", row), domainName)
				g.setCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("M%d", level.Level))
				g.setCellValue(sheetName, fmt.Sprintf("D%d", row), e.Name)
				g.setCellValue(sheetName, fmt.Sprintf("E%d", row), e.Description)
				g.setCellValue(sheetName, fmt.Sprintf("F%d", row), e.Type)
				g.setCellValue(sheetName, fmt.Sprintf("G%d", row), e.Layer)
				g.setCellValue(sheetName, fmt.Sprintf("H%d", row), e.Team)
				g.setCellValue(sheetName, fmt.Sprintf("I%d", row), e.Effort)
				g.setCellValue(sheetName, fmt.Sprintf("J%d", row), "-") // Status tracked via PRISM Maturity State
				g.setCellValue(sheetName, fmt.Sprintf("K%d", row), strings.Join(e.CriteriaIDs, ", "))
				g.setCellValue(sheetName, fmt.Sprintf("L%d", row), strings.Join(e.DependsOn, ", "))

				row++
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 12, "C": 8, "D": 35, "E": 50, "F": 15,
		"G": 12, "H": 20, "I": 12, "J": 15, "K": 30, "L": 25,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	// Auto filter
	g.setAutoFilter(sheetName, "A1:L1")

	return nil
}

// createSLOsSheet creates the SLOs (Criteria) sheet with framework columns.
func (g *XLSXGenerator) createSLOsSheet() error {
	sheetName := "SLOs"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Collect all unique frameworks across all criteria (sorted alphabetically)
	frameworks := g.collectAllFrameworks()

	// Headers - base columns plus framework columns
	headers := []string{
		"ID", "Domain", "Level", "Name", "Metric", "Type", "Operator",
		"Target", "Unit", "Current", "Met", "Layer", "Category", "Required",
	}
	// Add framework columns
	for _, fw := range frameworks {
		headers = append(headers, fw)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"548235"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	endHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setCellStyle(sheetName, "A1", endHeaderCell, headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				g.setCellValue(sheetName, fmt.Sprintf("A%d", row), c.ID)
				g.setCellValue(sheetName, fmt.Sprintf("B%d", row), domainName)
				g.setCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("M%d", level.Level))
				g.setCellValue(sheetName, fmt.Sprintf("D%d", row), c.Name)
				g.setCellValue(sheetName, fmt.Sprintf("E%d", row), c.GetMetricName(g.spec))

				// Type column - Quantitative or Qualitative (resolve from SLI)
				isQual := c.IsQualitativeWithSpec(g.spec)
				criterionType := "Quantitative"
				if isQual {
					criterionType = "Qualitative"
				}
				g.setCellValue(sheetName, fmt.Sprintf("F%d", row), criterionType)

				g.setCellValue(sheetName, fmt.Sprintf("G%d", row), OperatorSymbol(c.Operator))

				// Target - different display for qualitative
				if isQual {
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), "Tracked")
				} else {
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), c.Target)
				}

				g.setCellValue(sheetName, fmt.Sprintf("I%d", row), c.GetUnit(g.spec))

				// Current value/status tracked via PRISM Maturity State
				g.setCellValue(sheetName, fmt.Sprintf("J%d", row), "-")
				g.setCellValue(sheetName, fmt.Sprintf("K%d", row), "-")

				g.setCellValue(sheetName, fmt.Sprintf("L%d", row), c.GetLayer(g.spec))
				g.setCellValue(sheetName, fmt.Sprintf("M%d", row), c.GetCategory(g.spec))

				required := "Yes"
				if !c.Required && c.Weight > 0 {
					required = "No"
				}
				g.setCellValue(sheetName, fmt.Sprintf("N%d", row), required)

				// Framework columns - show control reference if mapped (resolve from SLI)
				frameworkRefs := make(map[string]string)
				for _, fm := range c.GetFrameworkMappings(g.spec) {
					frameworkRefs[fm.Framework] = fm.Reference
				}
				for i, fw := range frameworks {
					col, _ := excelize.CoordinatesToCellName(15+i, row) // Column O onwards
					if ref, ok := frameworkRefs[fw]; ok {
						g.setCellValue(sheetName, col, ref)
					} else {
						g.setCellValue(sheetName, col, "-")
					}
				}

				row++
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 12, "C": 8, "D": 30, "E": 35, "F": 12,
		"G": 10, "H": 12, "I": 10, "J": 12, "K": 8, "L": 12, "M": 12, "N": 10,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}
	// Set framework column widths
	for i := range frameworks {
		col, _ := excelize.ColumnNumberToName(15 + i)
		g.setColWidth(sheetName, col, col, 15)
	}

	// Auto filter
	endFilterCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setAutoFilter(sheetName, "A1:"+endFilterCell)

	return nil
}

// createThresholdMatrixSheet creates a pivot-style sheet showing SLIs with thresholds for each maturity level.
// This provides a human-readable view of how thresholds progress across M1-M5.
func (g *XLSXGenerator) createThresholdMatrixSheet() error {
	sheetName := "Threshold Matrix"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{"Category", "Tags", "Frameworks", "SLI Name", "Unit", "M1", "M2", "M3", "M4", "M5"}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2E75B6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "J1", headerStyle)

	// Build a map of SLI ID -> level -> threshold string
	// Structure: sliID -> { level: ">=99.9%", ... }
	type sliInfo struct {
		id         string
		name       string
		category   string
		unit       string
		tags       []string       // sorted, deduplicated tags
		frameworks []string       // sorted list of framework names
		thresholds map[int]string // level -> formatted threshold
	}

	sliMap := make(map[string]*sliInfo)
	var sliOrder []string // preserve order of first appearance

	domainNames := g.sortedDomainNames()
	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]
		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				if c.SLIID == "" {
					continue
				}

				// Get or create SLI info
				info, exists := sliMap[c.SLIID]
				if !exists {
					sli := c.GetSLI(g.spec)
					name := c.SLIID
					category := ""
					unit := ""
					var tags []string
					var frameworks []string
					if sli != nil {
						name = sli.Name
						category = sli.Category
						unit = sli.Unit
						// Get normalized (sorted, deduplicated) tags
						tags = sli.GetNormalizedTags()
						// Collect framework names from SLI mappings
						fwSet := make(map[string]bool)
						for _, fm := range sli.FrameworkMappings {
							fwSet[fm.Framework] = true
						}
						for fw := range fwSet {
							frameworks = append(frameworks, fw)
						}
						sort.Strings(frameworks)
					}
					info = &sliInfo{
						id:         c.SLIID,
						name:       name,
						category:   category,
						unit:       unit,
						tags:       tags,
						frameworks: frameworks,
						thresholds: make(map[int]string),
					}
					sliMap[c.SLIID] = info
					sliOrder = append(sliOrder, c.SLIID)
				}

				// Format threshold
				threshold := g.formatThreshold(c, info.unit)
				info.thresholds[level.Level] = threshold
			}
		}
	}

	// Sort by category (NIST CSF order), then by SLI order within category
	catWeights := CategorySortWeight()

	// Build SLI order map from spec categories
	sliOrderMap := make(map[string]int) // sliID -> order within category
	for _, cat := range g.spec.Categories {
		for idx, sliID := range cat.SLIOrder {
			sliOrderMap[sliID] = idx
		}
	}

	sort.Slice(sliOrder, func(i, j int) bool {
		a, b := sliMap[sliOrder[i]], sliMap[sliOrder[j]]

		// First sort by category weight (NIST CSF order)
		weightA := catWeights[a.category]
		weightB := catWeights[b.category]
		if weightA == 0 {
			weightA = 100 // Unknown categories sort last
		}
		if weightB == 0 {
			weightB = 100
		}
		if weightA != weightB {
			return weightA < weightB
		}

		// Within same category, sort by SLI order if defined
		orderA, hasOrderA := sliOrderMap[a.id]
		orderB, hasOrderB := sliOrderMap[b.id]
		if hasOrderA && hasOrderB {
			return orderA < orderB
		}
		if hasOrderA {
			return true // SLIs with order come first
		}
		if hasOrderB {
			return false
		}

		// Fall back to alphabetical by name
		return a.name < b.name
	})

	// Data rows
	row := 2
	for _, sliID := range sliOrder {
		info := sliMap[sliID]

		g.setCellValue(sheetName, fmt.Sprintf("A%d", row), info.category)
		g.setCellValue(sheetName, fmt.Sprintf("B%d", row), strings.Join(info.tags, ", "))
		g.setCellValue(sheetName, fmt.Sprintf("C%d", row), strings.Join(info.frameworks, ", "))
		g.setCellValue(sheetName, fmt.Sprintf("D%d", row), info.name)
		g.setCellValue(sheetName, fmt.Sprintf("E%d", row), info.unit)

		// M1-M5 thresholds (columns F-J)
		for level := 1; level <= 5; level++ {
			col, _ := excelize.CoordinatesToCellName(5+level, row)
			if threshold, ok := info.thresholds[level]; ok {
				g.setCellValue(sheetName, col, threshold)
			} else {
				g.setCellValue(sheetName, col, "-")
			}
		}

		row++
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 15, "B": 25, "C": 25, "D": 35, "E": 10, "F": 15, "G": 15, "H": 15, "I": 15, "J": 15,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	// Style threshold columns with center alignment
	if row > 2 {
		centerStyle, _ := g.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center"},
		})
		endCell := fmt.Sprintf("J%d", row-1)
		g.setCellStyle(sheetName, "F2", endCell, centerStyle)
	}

	// Auto filter
	g.setAutoFilter(sheetName, "A1:J1")

	return nil
}

// formatThreshold formats a criterion's threshold for display.
func (g *XLSXGenerator) formatThreshold(c Criterion, unit string) string {
	// Qualitative criteria
	if c.Operator == "exists" {
		return "Tracked"
	}

	// Format numeric threshold with operator and unit
	symbol := OperatorSymbol(c.Operator)
	if unit != "" {
		return fmt.Sprintf("%s%v%s", symbol, c.Target, unit)
	}
	return fmt.Sprintf("%s%v", symbol, c.Target)
}

// collectAllFrameworks returns all unique frameworks across all criteria, sorted alphabetically.
// Resolves framework mappings from both inline criterion mappings and referenced SLIs.
func (g *XLSXGenerator) collectAllFrameworks() []string {
	frameworkSet := make(map[string]bool)

	for _, domain := range g.spec.Domains {
		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Use GetFrameworkMappings to resolve from SLI if needed
				for _, fm := range c.GetFrameworkMappings(g.spec) {
					frameworkSet[fm.Framework] = true
				}
			}
		}
	}

	var frameworks []string
	for fw := range frameworkSet {
		frameworks = append(frameworks, fw)
	}
	sort.Strings(frameworks)
	return frameworks
}

// createFrameworkMappingsSheet creates a detailed Framework Mappings sheet (Option 4).
func (g *XLSXGenerator) createFrameworkMappingsSheet() error {
	sheetName := "Framework Mappings"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{
		"SLO ID", "SLO Name", "Domain", "Level", "Framework", "Reference",
		"Control Name", "Baseline", "Version", "Status",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"7030A0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "J1", headerStyle)

	// Data rows - one row per SLO-framework mapping
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Get framework mappings (resolve from SLI if needed)
				frameworkMappings := c.GetFrameworkMappings(g.spec)
				if len(frameworkMappings) == 0 {
					continue
				}

				// Create a row for each framework mapping
				for _, fm := range frameworkMappings {
					g.setCellValue(sheetName, fmt.Sprintf("A%d", row), c.ID)
					g.setCellValue(sheetName, fmt.Sprintf("B%d", row), c.Name)
					g.setCellValue(sheetName, fmt.Sprintf("C%d", row), domainName)
					g.setCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("M%d", level.Level))
					g.setCellValue(sheetName, fmt.Sprintf("E%d", row), fm.Framework)
					g.setCellValue(sheetName, fmt.Sprintf("F%d", row), fm.Reference)
					g.setCellValue(sheetName, fmt.Sprintf("G%d", row), fm.Name)
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), fm.Baseline)
					g.setCellValue(sheetName, fmt.Sprintf("I%d", row), fm.Version)
					g.setCellValue(sheetName, fmt.Sprintf("J%d", row), "-") // Status tracked via PRISM Maturity State

					row++
				}
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 30, "C": 12, "D": 8, "E": 15, "F": 15,
		"G": 35, "H": 12, "I": 10, "J": 10,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	// Auto filter
	g.setAutoFilter(sheetName, "A1:J1")

	return nil
}

// createProgressSheet creates the Progress Summary sheet.
// Note: Actual progress tracking should be done via PRISM Maturity State documents.
// This sheet shows the maturity model structure without current state data.
func (g *XLSXGenerator) createProgressSheet() error {
	sheetName := "Progress"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{
		"Domain", "Levels Defined", "Total Criteria", "Total Enablers",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"7030A0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "D1", headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		g.setCellValue(sheetName, fmt.Sprintf("A%d", row), domain.Name)
		g.setCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("M1-M%d", len(domain.Levels)))

		// Count total criteria and enablers
		totalCriteria := 0
		totalEnablers := 0
		for _, level := range domain.Levels {
			totalCriteria += len(level.Criteria)
			totalEnablers += len(level.Enablers)
		}
		g.setCellValue(sheetName, fmt.Sprintf("C%d", row), totalCriteria)
		g.setCellValue(sheetName, fmt.Sprintf("D%d", row), totalEnablers)

		row++
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 20, "B": 15, "C": 15, "D": 15,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	return nil
}

// createLevelDefinitionsSheet creates the Level Definitions sheet.
func (g *XLSXGenerator) createLevelDefinitionsSheet() error {
	sheetName := "Level Definitions"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Build headers dynamically from domains
	domainNames := g.sortedDomainNames()
	headers := []string{"Level", "Name"}
	for _, d := range domainNames {
		headers = append(headers, g.spec.Domains[d].Name)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"305496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	endCol, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setCellStyle(sheetName, "A1", endCol, headerStyle)

	// Data rows for levels M1-M5
	levelNames := DefaultLevelNames()
	for level := 1; level <= 5; level++ {
		row := level + 1
		g.setCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("M%d", level))
		g.setCellValue(sheetName, fmt.Sprintf("B%d", row), levelNames[level])

		// Description for each domain
		for col, domainName := range domainNames {
			domain := g.spec.Domains[domainName]
			levelDef, found := domain.GetLevel(level)
			desc := ""
			if found {
				desc = levelDef.Description
			}
			cell, _ := excelize.CoordinatesToCellName(col+3, row)
			g.setCellValue(sheetName, cell, desc)
		}
	}

	// Set column widths
	g.setColWidth(sheetName, "A", "A", 8)
	g.setColWidth(sheetName, "B", "B", 12)
	for i := range domainNames {
		col, _ := excelize.ColumnNumberToName(i + 3)
		g.setColWidth(sheetName, col, col, 50)
	}

	// Enable text wrap for description columns
	wrapStyle, _ := g.file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	endCol, _ = excelize.CoordinatesToCellName(len(headers), 6)
	g.setCellStyle(sheetName, "C2", endCol, wrapStyle)

	return nil
}

// Helper methods

// setCellValue wraps SetCellValue and ignores errors for simplicity.
// Cell value setting errors in excelize are rare and non-fatal.
func (g *XLSXGenerator) setCellValue(sheet, cell string, value interface{}) {
	_ = g.file.SetCellValue(sheet, cell, value)
}

// setCellStyle wraps SetCellStyle and ignores errors.
func (g *XLSXGenerator) setCellStyle(sheet, startCell, endCell string, styleID int) {
	_ = g.file.SetCellStyle(sheet, startCell, endCell, styleID)
}

// setColWidth wraps SetColWidth and ignores errors.
func (g *XLSXGenerator) setColWidth(sheet, startCol, endCol string, width float64) {
	_ = g.file.SetColWidth(sheet, startCol, endCol, width)
}

// setAutoFilter wraps AutoFilter and ignores errors.
func (g *XLSXGenerator) setAutoFilter(sheet, rangeRef string) {
	_ = g.file.AutoFilter(sheet, rangeRef, nil)
}

func (g *XLSXGenerator) sortedDomainNames() []string {
	var names []string
	for name := range g.spec.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GenerateXLSX is a convenience function to generate XLSX from a spec file.
func GenerateXLSX(specFile, outputFile string) error {
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		return err
	}

	gen := NewXLSXGenerator(spec)
	if err := gen.Generate(); err != nil {
		return err
	}

	return gen.SaveAs(outputFile)
}

// GenerateSimpleXLSX generates a simple XLSX report using omniframe.
// This provides a basic export without conditional styling.
func GenerateSimpleXLSX(specFile, outputFile string) error {
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		return err
	}

	fs := omniframe.NewFrameSet(spec.Metadata.Name)

	// Create Requirements frame
	reqFrame, err := buildRequirementsFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build requirements frame: %w", err)
	}
	if err := fs.AddFrame(reqFrame); err != nil {
		return err
	}

	// Create SLOs frame
	sloFrame, err := buildSLOsFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build SLOs frame: %w", err)
	}
	if err := fs.AddFrame(sloFrame); err != nil {
		return err
	}

	// Create Progress frame
	progressFrame, err := buildProgressFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build progress frame: %w", err)
	}
	if err := fs.AddFrame(progressFrame); err != nil {
		return err
	}

	return fs.WriteXLSX(outputFile)
}

func buildRequirementsFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"ID", "Domain", "Level", "Name", "Description", "Type",
		"Layer", "Team", "Effort", "Enables", "Depends On",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]
		for _, level := range domain.Levels {
			for _, e := range level.Enablers {
				rows = append(rows, []any{
					e.ID,
					domainName,
					fmt.Sprintf("M%d", level.Level),
					e.Name,
					e.Description,
					e.Type,
					e.Layer,
					e.Team,
					e.Effort,
					strings.Join(e.CriteriaIDs, ", "),
					strings.Join(e.DependsOn, ", "),
				})
			}
		}
	}

	frame, err := omniframe.FromRows("Requirements", columns, rows)
	if err != nil {
		return nil, err
	}

	// Set column widths
	_ = frame.SetColumnWidth("ID", 25)
	_ = frame.SetColumnWidth("Name", 35)
	_ = frame.SetColumnWidth("Description", 50)

	return frame, nil
}

func buildSLOsFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"ID", "Domain", "Level", "Name", "Metric", "Type", "Operator",
		"Target", "Unit", "Layer", "Category", "Required",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Determine type (resolve from SLI if needed)
				isQual := c.IsQualitativeWithSpec(spec)
				criterionType := "Quantitative"
				if isQual {
					criterionType = "Qualitative"
				}

				// Determine target display
				var targetDisplay any
				if isQual {
					targetDisplay = "Tracked"
				} else {
					targetDisplay = c.Target
				}

				required := "Yes"
				if !c.Required && c.Weight > 0 {
					required = "No"
				}

				rows = append(rows, []any{
					c.ID,
					domainName,
					fmt.Sprintf("M%d", level.Level),
					c.Name,
					c.GetMetricName(spec),
					criterionType,
					OperatorSymbol(c.Operator),
					targetDisplay,
					c.GetUnit(spec),
					c.GetLayer(spec),
					c.GetCategory(spec),
					required,
				})
			}
		}
	}

	frame, err := omniframe.FromRows("SLOs", columns, rows)
	if err != nil {
		return nil, err
	}

	_ = frame.SetColumnWidth("ID", 25)
	_ = frame.SetColumnWidth("Name", 30)
	_ = frame.SetColumnWidth("Metric", 35)
	_ = frame.SetColumnWidth("Type", 12)

	return frame, nil
}

func buildProgressFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"Domain", "Levels Defined", "Total Criteria", "Total Enablers",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]

		// Count total criteria and enablers
		totalCriteria := 0
		totalEnablers := 0
		for _, level := range domain.Levels {
			totalCriteria += len(level.Criteria)
			totalEnablers += len(level.Enablers)
		}

		rows = append(rows, []any{
			domain.Name,
			fmt.Sprintf("M1-M%d", len(domain.Levels)),
			totalCriteria,
			totalEnablers,
		})
	}

	frame, err := omniframe.FromRows("Progress", columns, rows)
	if err != nil {
		return nil, err
	}

	_ = frame.SetColumnWidth("Domain", 20)

	return frame, nil
}

func sortedDomainNamesFromSpec(spec *Spec) []string {
	var names []string
	for name := range spec.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
