# PRISM Examples

This directory contains example documents demonstrating PRISM's two main document types.

## Document Types

### 1. Maturity Models (`maturity-models/`)

**Schema:** `maturity.schema.json`

Maturity models define the **criteria for each maturity level (M1-M5)**. They are reference documents that describe what "good" looks like at each level.

**Structure:**

```
maturity-models/
├── operations/
│   ├── model.json      # Maturity level definitions
│   └── dashboard.html  # Generated visualization
├── security/
│   └── model.json
└── organization/
    └── model.json
```

**Contents:**

- `slis` - Service Level Indicators with thresholds for each level
- `domains` - Domain areas (e.g., Reliability, Incident Management)
- `levels` - M1 through M5 with criteria and thresholds
- `assessments` - Current state evaluation (optional)

**Example SLI with level thresholds:**

```json
{
  "slis": {
    "sli-availability": {
      "name": "Service Availability",
      "unit": "%",
      "sliType": "availability"
    }
  },
  "domains": {
    "reliability": {
      "levels": [
        {"level": 1, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 95}]},
        {"level": 2, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99}]},
        {"level": 3, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.5}]},
        {"level": 4, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.9}]},
        {"level": 5, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.99}]}
      ]
    }
  }
}
```

### 2. PRISM Documents (`prism-documents/`)

**Schema:** `prism.schema.json`

PRISM documents track the **current state** of your system against metrics, goals, and roadmaps.

**Contents:**

- `metrics` - Current metric values with SLOs and thresholds
- `goals` - Strategic objectives with target levels
- `phases` - Time-based planning periods
- `initiatives` - Projects that drive improvement
- `layers` - System architecture layers (requirements, code, infra, runtime)
- `roadmap` - Multi-phase planning

## How They Relate

```
┌─────────────────────────────────────────────────────────────────┐
│                      Maturity Model                             │
│  Defines: "What does M3 availability look like? (≥99.5%)"      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ references
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PRISM Document                             │
│  Tracks: "Current availability is 99.7% → We're at M3"         │
│  Plans: "Goal: Reach M4 (99.9%) by Q2"                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ generates
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Dashboard                                │
│  Visualizes: Progress, gaps, trends                            │
└─────────────────────────────────────────────────────────────────┘
```

## Observability Methodologies

Maturity models use standard observability methodologies to classify SLIs:

| Methodology | SLI Types | Focus |
|-------------|-----------|-------|
| **RED** | Rate, Errors, Duration | User experience (request-driven) |
| **USE** | Utilization, Saturation, Errors | Infrastructure health |
| **Golden Signals** | Latency, Traffic, Errors, Saturation | SRE overview |

## Example Domains

| Domain | Model | Description |
|--------|-------|-------------|
| Operations | `operations/model.json` | Reliability, deployment, monitoring |
| Security | `security/model.json` | Prevention, detection, response |
| Organization | `organization/model.json` | Team structure, processes |

## Generating Dashboards

To generate a dashboard from a maturity model:

```bash
# Using Go test (temporary)
go test ./dashboard -run TestGenerateHTML

# Future CLI (planned)
prism dashboard examples/maturity-models/operations/model.json \
  -o examples/maturity-models/operations/dashboard.html
```

## File Naming Conventions

| File | Purpose |
|------|---------|
| `model.json` | Maturity model definition |
| `dashboard.html` | Generated visualization (gitignored) |
| `assessment.json` | Current state evaluation (optional) |
