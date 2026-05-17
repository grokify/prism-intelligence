# Maturity Model Workflow

This guide walks you through the complete maturity model workflow: defining a model, tracking state, and generating dashboards.

## Overview

PRISM uses three document types for maturity modeling:

| Document | Purpose | Question Answered |
|----------|---------|-------------------|
| **Model** | Definitions | What does good look like at each level? |
| **State** | Measurement | Where are we now? |
| **Dashboard** | Visualization | How do we present progress? |

```
Model (SLIs, thresholds) → State (measurements) → Dashboard (visualization)
```

## Step 1: Define Your Maturity Model

A maturity model defines Service Level Indicators (SLIs) and thresholds for each maturity level (M1-M5).

Create `model.json`:

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism-maturity-model.schema.json",
  "metadata": {
    "name": "Operations Maturity Model",
    "description": "Defines what good looks like for operations"
  },
  "slis": {
    "sli-availability": {
      "id": "sli-availability",
      "name": "Service Availability",
      "metricName": "availability_pct",
      "unit": "%",
      "type": "quantitative",
      "category": "reliability",
      "tags": ["runtime-defense"]
    },
    "sli-error-rate": {
      "id": "sli-error-rate",
      "name": "Error Rate",
      "metricName": "error_rate_pct",
      "unit": "%",
      "type": "quantitative",
      "category": "reliability"
    },
    "sli-mttr": {
      "id": "sli-mttr",
      "name": "Mean Time to Recovery",
      "metricName": "mttr_minutes",
      "unit": "min",
      "type": "quantitative",
      "category": "response"
    }
  },
  "domains": {
    "operations": {
      "name": "Operations",
      "description": "Operational reliability and performance",
      "levels": [
        {
          "level": 1,
          "name": "Reactive",
          "description": "Ad-hoc operations, firefighting mode",
          "criteria": [
            {"id": "avail-m1", "sliId": "sli-availability", "operator": ">=", "target": 95}
          ]
        },
        {
          "level": 2,
          "name": "Basic",
          "description": "Basic monitoring in place",
          "criteria": [
            {"id": "avail-m2", "sliId": "sli-availability", "operator": ">=", "target": 99},
            {"id": "error-m2", "sliId": "sli-error-rate", "operator": "<=", "target": 1}
          ]
        },
        {
          "level": 3,
          "name": "Defined",
          "description": "SLOs defined and tracked",
          "criteria": [
            {"id": "avail-m3", "sliId": "sli-availability", "operator": ">=", "target": 99.5},
            {"id": "error-m3", "sliId": "sli-error-rate", "operator": "<=", "target": 0.5},
            {"id": "mttr-m3", "sliId": "sli-mttr", "operator": "<=", "target": 60}
          ]
        },
        {
          "level": 4,
          "name": "Managed",
          "description": "Data-driven, proactive management",
          "criteria": [
            {"id": "avail-m4", "sliId": "sli-availability", "operator": ">=", "target": 99.9},
            {"id": "error-m4", "sliId": "sli-error-rate", "operator": "<=", "target": 0.1},
            {"id": "mttr-m4", "sliId": "sli-mttr", "operator": "<=", "target": 30}
          ]
        },
        {
          "level": 5,
          "name": "Optimizing",
          "description": "Continuous improvement, industry-leading",
          "criteria": [
            {"id": "avail-m5", "sliId": "sli-availability", "operator": ">=", "target": 99.99},
            {"id": "error-m5", "sliId": "sli-error-rate", "operator": "<=", "target": 0.01},
            {"id": "mttr-m5", "sliId": "sli-mttr", "operator": "<=", "target": 15}
          ]
        }
      ]
    }
  }
}
```

### Validate Your Model

```bash
prism maturity model validate model.json
```

### Lint for Dashboard Issues

Check for issues that affect dashboard display:

```bash
prism maturity model lint model.json
```

## Step 2: Track Current State

A state document captures current measurements for your SLIs.

Create `state.json`:

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism-maturity-state.schema.json",
  "metadata": {
    "name": "Operations State Q2 2026",
    "maturityModelRef": "./model.json",
    "assessedAt": "2026-05-15",
    "assessedBy": "SRE Team"
  },
  "sloWindows": ["7d", "30d"],
  "sliState": {
    "sli-availability": {
      "sliId": "sli-availability",
      "qualitativeState": "alerting",
      "windows": {
        "30d": { "value": 99.5, "timestamp": "2026-05-15T00:00:00Z" }
      }
    },
    "sli-error-rate": {
      "sliId": "sli-error-rate",
      "qualitativeState": "measured",
      "windows": {
        "30d": { "value": 0.3, "timestamp": "2026-05-15T00:00:00Z" }
      }
    },
    "sli-mttr": {
      "sliId": "sli-mttr",
      "qualitativeState": "measured",
      "windows": {
        "30d": { "value": 45, "timestamp": "2026-05-15T00:00:00Z" }
      }
    }
  },
  "maturityState": {
    "operations": {
      "domainId": "operations",
      "current": {
        "level": 3,
        "achievedAt": "2026-03-15",
        "note": "All M3 criteria met"
      },
      "target": {
        "level": 5,
        "targetDate": "2027-06-30"
      }
    }
  }
}
```

### Validate State Against Model

```bash
prism maturity state validate state.json --model model.json
```

## Step 3: Generate Outputs

### HTML Dashboard

Generate an interactive dashboard with bullet charts and progress tracking:

```bash
prism maturity model dashboard model.json --state state.json -f html -o dashboard.html
```

The dashboard includes:

- **Overall maturity** bullet chart showing current vs target level
- **SLI bullet charts** grouped by category (NIST CSF order)
- **Progress to target** showing completion status per SLI
- **SLI definitions** table with thresholds for each maturity level

### Excel Spreadsheet

Generate an XLSX file with multiple sheets:

```bash
prism maturity model xlsx model.json -o model.xlsx
```

Sheets include:

- **Domains** - Domain overview with maturity levels
- **SLOs** - All SLI/SLO definitions
- **Threshold Matrix** - Pivot view of M1-M5 thresholds per SLI
- **Framework Mappings** - SLI to framework mappings (NIST CSF, etc.)

### Markdown Report

Generate a markdown report for documentation:

```bash
prism maturity model report model.json -o report.md
```

## NIST CSF Category Ordering

PRISM automatically sorts categories by NIST CSF 2.0 canonical order:

1. **Govern** - Governance and oversight
2. **Identify** - Asset and risk identification
3. **Protect** - Safeguards and controls
4. **Detect** - Monitoring and detection
5. **Respond** - Incident response
6. **Recover** - Recovery and restoration

Operations-focused categories (reliability, efficiency, quality) sort after NIST CSF categories.

## Using Tags for Classification

Tags enable multi-dimensional classification orthogonal to category:

```json
{
  "slis": {
    "vuln-scan": {
      "name": "Vulnerability Scan Coverage",
      "category": "detect",
      "tags": ["shift-left", "vulnerability-management"]
    }
  }
}
```

Recommended tags:

| Tag | Description |
|-----|-------------|
| `ai` | AI/ML-specific security |
| `shift-left` | Design/build-time controls |
| `supply-chain` | Software supply chain security |
| `runtime-defense` | Production-time protection |
| `vulnerability-management` | Vulnerability handling |
| `incident-response` | Incident handling |

## Complete Example

See the `examples/operations/` directory for a complete working example:

```bash
# View the model
cat examples/operations/model.json

# View the state
cat examples/operations/state-q2-2026.json

# Generate dashboard
prism maturity model dashboard examples/operations/model.json \
  --state examples/operations/state-q2-2026.json \
  -f html -o dashboard.html

# Generate Excel
prism maturity model xlsx examples/operations/model.json -o model.xlsx
```

## Next Steps

- [SLI/SLO Schema](../schema/slos.md) - Detailed SLI configuration options
- [CLI Reference](../cli/maturity-model.md) - Full CLI documentation
- [v0.7.0 Release Notes](../releases/v0.7.0.md) - Latest features
