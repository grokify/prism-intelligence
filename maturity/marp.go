package maturity

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// MarpGenerator generates Marp presentations from maturity specs.
type MarpGenerator struct {
	spec *Spec
}

// NewMarpGenerator creates a new Marp generator.
func NewMarpGenerator(spec *Spec) *MarpGenerator {
	return &MarpGenerator{spec: spec}
}

// Generate creates the Marp presentation content.
func (g *MarpGenerator) Generate() (string, error) {
	var sb strings.Builder

	// Write frontmatter
	sb.WriteString(marpFrontmatter(g.spec.Metadata.Name))

	// Executive Overview
	sb.WriteString(g.execOverviewSlides())

	// Domain sections
	domainOrder := []string{"security", "operational-excellence", "quality", "product", "ai"}
	for _, domainKey := range domainOrder {
		domain, ok := g.spec.Domains[domainKey]
		if !ok {
			continue
		}
		sb.WriteString(g.domainSlides(domainKey, domain))
	}

	// Summary and appendix
	sb.WriteString(g.summarySlides())
	sb.WriteString(g.appendixSlides())

	return sb.String(), nil
}

// SaveAs saves the presentation to a file.
func (g *MarpGenerator) SaveAs(filename string) error {
	content, err := g.Generate()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(content), 0600)
}

func marpFrontmatter(title string) string {
	return fmt.Sprintf(`---
marp: true
theme: default
paginate: true
header: '%s'
footer: 'Confidential'
---

<!-- _class: lead -->

# %s

## A Unified Framework for B2B SaaS Health Metrics

**Security | Operational Excellence | Quality | Product | AI**

---

`, title, title)
}

func (g *MarpGenerator) execOverviewSlides() string {
	var sb strings.Builder

	sb.WriteString(`<!-- _class: lead -->

# Executive Overview

---

# The Challenge

## Fragmented Metrics Across the Organization

- **Security** tracks vulnerabilities, compliance, threat detection
- **Engineering** tracks DORA metrics, SLOs, incidents
- **Quality** tracks coverage, defects, test results
- **Product** tracks adoption, activation, churn
- **AI** tracks adoption, productivity, governance

**Result:** No unified view of organizational health

---

# Maturity Model

## 5 Levels of Capability

| Level | Name | Description |
|-------|------|-------------|
| **M1** | Reactive | Ad-hoc, firefighting, heroics |
| **M2** | Basic | Documented, some repeatability |
| **M3** | Defined | Standardized, consistent execution |
| **M4** | Managed | Data-driven, measured, controlled |
| **M5** | Optimizing | Continuous improvement, automated |

---

# Current State Summary

## Where We Are Today

`)

	// Build current state table from domains (state tracking done via PRISM Maturity State documents)
	sb.WriteString("| Domain | Levels Defined |\n")
	sb.WriteString("|--------|----------------|\n")

	domainOrder := []string{"security", "operational-excellence", "quality", "product", "ai"}
	for _, domainKey := range domainOrder {
		domain, ok := g.spec.Domains[domainKey]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("| **%s** | M1-M%d |\n",
			domain.Name, len(domain.Levels)))
	}

	sb.WriteString("\n---\n\n")
	return sb.String()
}

func (g *MarpGenerator) domainSlides(domainKey string, domain *DomainModel) string {
	var sb strings.Builder

	// Domain title slide
	sb.WriteString(fmt.Sprintf(`<!-- _class: lead -->

# %s Domain

**%s**

---

`, domain.Name, domain.Description))

	// SLI-based criteria table (using new SLI structure)
	if len(g.spec.SLIs) > 0 {
		// Collect SLIs for this domain
		domainSLIs := g.collectDomainSLIs(domainKey, domain)
		if len(domainSLIs) > 0 {
			sb.WriteString(fmt.Sprintf("# %s SLIs\n\n## Service Level Indicators\n\n", domain.Name))
			sb.WriteString("| SLI | Type | Unit | Category |\n")
			sb.WriteString("|-----|------|------|----------|\n")

			for _, sli := range domainSLIs {
				sb.WriteString(fmt.Sprintf("| **%s** | %s | %s | %s |\n",
					sli.Name,
					formatSLIType(sli.SLIType),
					formatUnit(sli.Unit),
					formatCategory(sli.Category),
				))
			}
			sb.WriteString("\n---\n\n")
		}
	}

	// Level definitions
	sb.WriteString(fmt.Sprintf("# %s Maturity Levels\n\n## What Each Level Means\n\n", domain.Name))
	sb.WriteString("| Level | Name | Description |\n")
	sb.WriteString("|-------|------|-------------|\n")

	for _, level := range domain.Levels {
		sb.WriteString(fmt.Sprintf("| **M%d** | %s | %s |\n",
			level.Level, level.Name, truncate(level.Description, 50)))
	}
	sb.WriteString("\n---\n\n")

	// Enablers/Roadmap
	sb.WriteString(fmt.Sprintf("# %s Roadmap\n\n## Key Initiatives\n\n", domain.Name))
	sb.WriteString("| Project | Level | Type | Layer |\n")
	sb.WriteString("|---------|-------|------|-------|\n")

	enablers := g.collectEnablersFromModel(domain)
	for i, e := range enablers {
		if i >= 5 {
			break
		}
		sb.WriteString(fmt.Sprintf("| %s | M%d | %s | %s |\n",
			e.Name, e.Level, formatEnablerType(e.Type), formatLayer(e.Layer)))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

type enablerWithLevel struct {
	Enabler
	Level int
}

// collectEnablersFromModel collects enablers from the maturity model definition.
// State tracking (status) should be done via PRISM Maturity State documents.
func (g *MarpGenerator) collectEnablersFromModel(domain *DomainModel) []enablerWithLevel {
	var enablers []enablerWithLevel

	for _, level := range domain.Levels {
		for _, e := range level.Enablers {
			enablers = append(enablers, enablerWithLevel{
				Enabler: Enabler{
					ID:          e.ID,
					Name:        e.Name,
					Description: e.Description,
					Type:        e.Type,
					Layer:       e.Layer,
				},
				Level: level.Level,
			})
		}
	}

	// Sort by level, then by name
	sort.Slice(enablers, func(i, j int) bool {
		if enablers[i].Level != enablers[j].Level {
			return enablers[i].Level < enablers[j].Level
		}
		return enablers[i].Name < enablers[j].Name
	})

	return enablers
}

// collectDomainSLIs collects SLIs referenced by criteria in a domain.
func (g *MarpGenerator) collectDomainSLIs(_ string, domain *DomainModel) []*SLI {
	seen := make(map[string]bool)
	var slis []*SLI

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			if criterion.SLIID != "" && !seen[criterion.SLIID] {
				if sli, ok := g.spec.SLIs[criterion.SLIID]; ok && sli != nil {
					slis = append(slis, sli)
					seen[criterion.SLIID] = true
				}
			}
		}
	}

	// Sort by name
	sort.Slice(slis, func(i, j int) bool {
		return slis[i].Name < slis[j].Name
	})

	return slis
}

func (g *MarpGenerator) summarySlides() string {
	return `<!-- _class: lead -->

# Summary and Next Steps

---

# Cross-Domain Summary

## Maturity Progression Plan

| Domain | Current | Q2 Target | Q4 Target |
|--------|---------|-----------|-----------|
| **Security** | M2 | M3 | M4 |
| **Operational Excellence** | M3 | M4 | M4 |
| **Quality** | M2 | M3 | M4 |
| **Product** | M2 | M3 | M4 |
| **AI** | M2 | M3 | M4 |

---

# Next Steps

## Immediate Actions

1. **Approve** maturity model as the unified metrics framework
2. **Conduct** baseline maturity assessment per domain
3. **Assign** domain overlay owners
4. **Kickoff** Q1 initiatives

---

`
}

func (g *MarpGenerator) appendixSlides() string {
	return `<!-- _class: lead -->

# Appendix

---

# Framework Mappings

## Industry Alignment

| Framework | Mapping |
|-----------|---------|
| **DORA** | Operational Excellence metrics |
| **SRE** | Golden signals (latency, traffic, errors, saturation) |
| **NIST CSF** | Security domain stages |
| **MITRE ATT&CK** | Security detection metrics |
| **ISO 25010** | Quality characteristics |

---

# Glossary

| Term | Definition |
|------|------------|
| **SLO** | Service Level Objective - target for a metric |
| **SLI** | Service Level Indicator - the measurement |
| **DORA** | DevOps Research and Assessment metrics |
| **MTTR** | Mean Time to Recovery |
| **CFR** | Change Failure Rate |
| **AIOps** | AI for IT Operations |
| **CSPM** | Cloud Security Posture Management |
`
}

func formatSLIType(sliType string) string {
	switch sliType {
	case "availability":
		return "Availability"
	case "latency":
		return "Latency"
	case "error_rate":
		return "Error Rate"
	case "throughput":
		return "Throughput"
	case "saturation":
		return "Saturation"
	case "utilization":
		return "Utilization"
	case "quality":
		return "Quality"
	case "freshness":
		return "Freshness"
	default:
		if sliType == "" {
			return "-"
		}
		return sliType
	}
}

func formatUnit(unit string) string {
	if unit == "" {
		return "-"
	}
	return unit
}

func formatCategory(category string) string {
	switch category {
	case "prevention":
		return "Prevention"
	case "detection":
		return "Detection"
	case "response":
		return "Response"
	case "reliability":
		return "Reliability"
	case "efficiency":
		return "Efficiency"
	default:
		if category == "" {
			return "-"
		}
		return category
	}
}

func formatEnablerType(t string) string {
	switch t {
	case TypeImplementation:
		return "Implementation"
	case TypeProcess:
		return "Process"
	case TypeTraining:
		return "Training"
	case TypeTooling:
		return "Tooling"
	default:
		if t == "" {
			return "-"
		}
		return t
	}
}

func formatLayer(layer string) string {
	if layer == "" {
		return "-"
	}
	return layer
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GenerateMarp is a convenience function to generate Marp from a spec file.
func GenerateMarp(specFile, outputFile string) error {
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		return err
	}

	gen := NewMarpGenerator(spec)
	return gen.SaveAs(outputFile)
}
