# PRISM Execution Integration

PRISM Intelligence integrates with [PRISM Execution](https://github.com/grokify/prism-execution) to export roadmaps, OKRs, and V2MOMs for execution tracking.

## Overview

PRISM Intelligence serves as the source of truth for requirements (maturity models, SLOs), while PRISM Execution handles execution tracking:

| PRISM Intelligence (Source of Truth) | PRISM Execution (Execution) |
|--------------------------------------|----------------------------|
| Goals with maturity models | Objectives |
| SLOs and MetricCriteria | Key Results |
| Phases with GoalTargets | Roadmap Phases |
| Initiatives with DeploymentStatus | Deliverables with RolloutStatus |

## Export Formats

### Roadmap Export

Best for B2B SaaS with customer deployment tracking.

```bash
prism export roadmap prism.json -o roadmap.json
```

Features:

- Phase-based timeline
- Deliverables with rollout status
- Deployment vs adoption tracking
- Phased rollout waves

### OKR Export

Best for goal-focused tracking.

```bash
prism export okr prism.json -o okr.json
```

Features:

- One objective per maturity level to achieve
- Key results from SLOs and MetricCriteria
- Phase targets for temporal planning

### V2MOM Export

Best for Salesforce-style planning.

```bash
prism export v2mom prism.json -o v2mom.json
```

Features:

- Methods from goals
- Measures from SLOs
- Projects from initiatives

### Combined Export

Export roadmap with embedded OKRs:

```bash
prism export roadmap prism.json --with-okrs -o full-roadmap.json
```

## Data Mapping

### Goals to Objectives

Each PRISM goal can produce multiple OKR objectives—one for each maturity level to achieve:

```
Goal: "High Reliability"
  CurrentLevel: 3
  TargetLevel: 5

→ Objective: "Achieve Managed Level for Reliability" (M4)
→ Objective: "Achieve Optimizing Level for Reliability" (M5)
```

### SLOs to Key Results

SLO requirements at each maturity level become key results:

```
Goal M4 Requirements:
  - availability >= 99.9%
  - mttr <= 30m

→ Key Result: "availability meets M4 requirements (>=99.9%)"
→ Key Result: "mttr meets M4 requirements (<=30m)"
```

### Initiatives to Deliverables

PRISM initiatives become roadmap deliverables with rollout tracking:

```
Initiative: "Observability Platform"
  Status: in_progress
  DevCompletionPercent: 90
  DeploymentStatus:
    TotalCustomers: 50
    DeployedCustomers: 45

→ Deliverable:
    Title: "Observability Platform"
    Status: in_progress
    Rollout:
      TotalCustomers: 50
      DeployedCustomers: 45
      Status: rolling_out
```

## RolloutStatus

The RolloutStatus type tracks B2B SaaS feature deployment and adoption:

| Field | Description |
|-------|-------------|
| `totalCustomers` | Total customers in rollout scope |
| `deployedCustomers` | Customers with feature deployed (available) |
| `adoptedCustomers` | Customers actively using the feature |
| `status` | Current rollout stage |
| `startDate` | Rollout start date |
| `targetDate` | Target completion date |
| `waves` | Phased rollout waves |

### Rollout Stages

| Stage | Description |
|-------|-------------|
| `not_started` | Rollout hasn't begun |
| `rolling_out` | Actively deploying to customers |
| `deployed` | 100% deployed, adoption ongoing |
| `adopted` | Target adoption achieved |
| `paused` | Rollout temporarily paused |
| `rolled_back` | Feature rolled back |

### Metrics

```
Deployment %  = deployedCustomers / totalCustomers × 100
Adoption %    = adoptedCustomers / totalCustomers × 100
Adoption Rate = adoptedCustomers / deployedCustomers × 100
```

## Workflow

### 1. Define Requirements in PRISM Intelligence

```json
{
  "goals": [{
    "id": "goal-reliability",
    "name": "High Reliability",
    "currentLevel": 3,
    "targetLevel": 5,
    "maturityModel": {
      "levels": [{
        "level": 4,
        "name": "Managed",
        "requiredSLOs": [{"metricId": "availability"}],
        "metricCriteria": [{"metricId": "availability", "operator": "gte", "value": 99.9}]
      }]
    }
  }],
  "phases": [{
    "id": "phase-q1-2026",
    "name": "Q1 2026",
    "goalTargets": [{"goalId": "goal-reliability", "enterLevel": 3, "exitLevel": 4}]
  }],
  "initiatives": [{
    "id": "init-monitoring",
    "name": "Observability Platform",
    "phaseId": "phase-q1-2026",
    "deploymentStatus": {
      "totalCustomers": 50,
      "deployedCustomers": 45
    }
  }]
}
```

### 2. Export to Roadmap

```bash
prism export roadmap prism.json --with-okrs -o roadmap.json
```

### 3. Use with PRISM Execution

```bash
# Validate exported document
splan validate roadmap.json

# Render to HTML
splan render roadmap.json -o roadmap.html
```

## Library Usage

For programmatic access, use the export package:

```go
import (
    "github.com/grokify/prism-intelligence"
    "github.com/grokify/prism-intelligence/export"
)

// Load PRISM document
doc, _ := prism.LoadFile("prism.json")

// Export to roadmap
rm := export.ConvertToRoadmap(doc)

// Export with OKRs
full := export.ConvertToRoadmapWithOKRs(doc)

// Export OKRs only
okrs := export.ConvertToStructuredOKR(doc)
```

## Related Documentation

- [Export CLI](../cli/export.md) - CLI command reference
- [Goals](../concepts/goals.md) - Goal and maturity model concepts
- [Phases](../concepts/phases.md) - Phase planning concepts
