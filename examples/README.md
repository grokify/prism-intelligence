# PRISM Examples

This directory contains example documents demonstrating PRISM's two main document types.

## Document Types

### 1. PRISM Maturity Model (`maturity-models/`)

**Schema:** `prism-maturity-model.schema.json`

Maturity Models define the **criteria for each maturity level (M1-M5)**. They are reference documents that describe what "good" looks like at each level. They contain no current state.

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
- `qualitativeStates` - State progressions (e.g., "tracked" → "measured" → "alerting")

**Example SLI with level thresholds and qualitative states:**

```json
{
  "slis": {
    "sli-availability": {
      "name": "Service Availability",
      "unit": "%",
      "sliType": "availability",
      "measurementType": "hybrid",
      "qualitativeStates": [
        { "id": "none", "label": "Not tracked", "order": 0 },
        { "id": "tracked", "label": "Tracked", "order": 1 },
        { "id": "measured", "label": "Measured with SLO", "order": 2 },
        { "id": "alerting", "label": "SLO + Alerting", "order": 3 }
      ]
    }
  },
  "domains": {
    "reliability": {
      "levels": [
        {"level": 1, "criteria": [{"sliId": "sli-availability", "type": "qualitative", "target": "tracked"}]},
        {"level": 2, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 95}]},
        {"level": 3, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.5}]},
        {"level": 4, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.9}]},
        {"level": 5, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.99}]}
      ]
    }
  }
}
```

### 2. PRISM Maturity State (`maturity-state/`)

**Schema:** `prism-maturity-state.schema.json`

Maturity State documents track the **current state** of your system, including historical values and future targets. They reference a Maturity Model to evaluate achievement.

**Structure:**

```
maturity-state/
├── operations/
│   └── state-q2-2026.json  # Q2 2026 state
├── organization/
│   └── state-q1-2026.json  # Q1 2026 state (multi-domain)
└── security/
    └── state-q2-2026.json  # Q2 2026 state
```

**Contents:**

- `maturityModelRef` - Reference to the Maturity Model
- `sliState` - Current values with temporal windows (7d, 30d, 90d, quarterly, annual)
- `maturityState` - Current level, target level, and history per domain
- `enablerState` - Progress on enablers/initiatives
- `goals` - Strategic objectives with target levels
- `phases` - Time-based planning periods
- `initiatives` - Projects that drive improvement

**Example state document:**

```json
{
  "metadata": {
    "name": "Operations Q2 2026",
    "maturityModelRef": "../maturity-models/operations/model.json"
  },
  "sloWindows": ["7d", "30d", "90d", "quarterly"],
  "sliState": {
    "sli-availability": {
      "qualitativeState": "measured",
      "windows": {
        "7d":  { "value": 99.97, "timestamp": "2026-05-10T00:00:00Z" },
        "30d": { "value": 99.92, "timestamp": "2026-05-10T00:00:00Z" }
      },
      "targets": {
        "Q2_2026": { "value": 99.9 }
      }
    }
  },
  "maturityState": {
    "reliability": {
      "current": { "level": 3, "achievedAt": "2026-03-15" },
      "target": { "level": 4, "targetDate": "2026-06-30" }
    }
  }
}
```

## How They Relate

```
┌─────────────────────────────────────────────────────────────────┐
│                   PRISM Maturity Model                          │
│  Defines: "What does M3 availability look like? (≥99.5%)"       │
│  Contains: SLI definitions, level criteria, qualitative states  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ references
                              │
┌─────────────────────────────────────────────────────────────────┐
│                   PRISM Maturity State                          │
│  Tracks: "Current availability is 99.7% → We're at M3"          │
│  Plans: "Goal: Reach M4 (99.9%) by Q2"                          │
│  History: "30d values: 99.88 → 99.90 → 99.92"                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ generates
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Dashboard                                 │
│  Visualizes: Progress, gaps, trends                             │
└─────────────────────────────────────────────────────────────────┘
```

## Temporal Windows

PRISM Maturity State supports multiple SLO windows for tracking:

| Window | Description |
|--------|-------------|
| `7d` | Rolling 7-day window |
| `30d` | Rolling 30-day window |
| `90d` | Rolling 90-day window |
| `quarterly` | Calendar quarter |
| `annual` | Calendar year |

## Qualitative vs Quantitative

PRISM supports three measurement types:

| Type | Description | Example |
|------|-------------|---------|
| `quantitative` | Numeric values only | Availability: 99.9% |
| `qualitative` | State-based only | Monitoring: "tracked" |
| `hybrid` | Both numeric and state | Start with "tracked", progress to 99.9% |

## Observability Methodologies

Maturity Models use standard observability methodologies to classify SLIs:

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
| `model.json` | PRISM Maturity Model definition |
| `state-*.json` | PRISM Maturity State tracking |
| `dashboard.html` | Generated visualization (gitignored) |
