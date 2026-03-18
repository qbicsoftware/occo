# OpenCode Helper CLI

## Overview

This document is the single traceable product source for the OpenCode helper CLI. It captures the initial product direction, requirements, first feature set, and the traceability structure that later user stories must follow.

Status: Draft for V1  
Owner: TBD  
Scope baseline: V1

---

## Decision Log

### DEC-001 - CLI is the sole official distribution channel for V1

The helper CLI is the only supported distribution mechanism for official OpenCode setup assets in V1.

Included in the CLI release bundle:
- Official config presets
- Official inter-agent schema contracts
- Version metadata needed for validation and provenance

Implications:
- Users do not need to clone this repository to use the official setup assets
- The CLI release is the canonical end-user artifact
- `self-update` updates the CLI and bundled assets together
- Live remote schema fetching is out of scope for V1

---

## PRD

### PRD-001 - Product Goal

Provide a small helper CLI that bootstraps and maintains a local OpenCode project setup using official bundled config presets and official bundled inter-agent schemas, with safe validation and self-update behavior tied to CLI releases.

### PRD-002 - Problem Statement

Current OpenCode project setup depends on repository-local assets and setup knowledge that are not yet exposed through a single user-facing distribution and setup workflow. Users may not want to clone the repository locally just to obtain supported presets and schema files.

### PRD-003 - Target Users

- Maintainers who publish and evolve the official OpenCode setup
- Developers who want to initialize a project with a supported OpenCode setup
- Contributors who need a repeatable, validated local setup flow

### PRD-004 - User Value

The helper CLI should let a user:
- install one official tool
- inspect available presets
- apply a supported local configuration
- install required schemas
- validate setup health
- update to the latest supported release

### PRD-005 - Success Criteria

V1 is successful when:
- A user can install the CLI without cloning this repository
- A user can initialize a project from bundled presets and schemas
- A user can validate whether the project setup is healthy or drifted
- A user can update the CLI and bundled assets via GitHub release upgrade
- Installed project assets can be traced back to the CLI release version

---

## Scope

### In Scope

- Bundling official config presets with the CLI release
- Bundling official inter-agent schemas with the CLI release
- Project-local setup from bundled assets
- Safe install/apply behavior for project files
- Validation and diagnostics
- Self-update from the latest GitHub release
- Version/provenance reporting for bundled assets

### Out of Scope

- Direct end-user fetching of raw assets from GitHub releases
- Live remote schema syncing during normal setup
- Automatic migration of arbitrary user-customized configs
- Windows support unless explicitly added later
- Full remote orchestration or hosted service behavior

---

## Requirements

### Functional Requirements

#### REQ-F-001 - Bundled Presets

The CLI shall include the official OpenCode config presets as bundled release assets.

Depends on:
- DEC-001

#### REQ-F-002 - Bundled Schemas

The CLI shall include the official inter-agent handoff and result schemas as bundled release assets.

Depends on:
- DEC-001

#### REQ-F-003 - Preset Discovery

The CLI shall let users list the bundled config presets available in the installed CLI release.

Depends on:
- REQ-F-001

#### REQ-F-004 - Preset Application

The CLI shall let users materialize a selected bundled preset into a target project as `opencode.json` or equivalent configured output.

Depends on:
- REQ-F-001
- REQ-F-003

#### REQ-F-005 - Schema Installation

The CLI shall let users install the bundled schemas into the target project or supported scope.

Depends on:
- REQ-F-002

#### REQ-F-006 - Setup Validation

The CLI shall validate that the target project contains the required OpenCode setup files and that installed assets are compatible with the bundled release contents.

Depends on:
- REQ-F-001
- REQ-F-002

#### REQ-F-007 - Diagnostics

The CLI shall provide diagnostics for missing files, drift, invalid setup state, and likely operator mistakes.

Depends on:
- REQ-F-006

#### REQ-F-008 - Self Update

The CLI shall update itself from the latest GitHub release of the project.

Depends on:
- DEC-001

#### REQ-F-009 - Version Reporting

The CLI shall report its own version and the version or identity of the bundled asset set.

Depends on:
- REQ-F-001
- REQ-F-002
- REQ-F-008

#### REQ-F-010 - Safe File Handling

The CLI shall avoid overwriting existing user files by default and shall require explicit opt-in for destructive replacement behavior.

Depends on:
- REQ-F-004
- REQ-F-005

### Non-Functional Requirements

#### REQ-NF-001 - Minimal Footprint

The CLI should remain small and operationally simple, with minimal runtime dependencies.

#### REQ-NF-002 - Offline-Friendly Operation

Normal project setup after CLI installation should not require network access.

Depends on:
- DEC-001

#### REQ-NF-003 - Deterministic Release Bundle

Each CLI release shall pin a specific compatible set of bundled presets and schemas.

Depends on:
- DEC-001

#### REQ-NF-004 - Clear Automation Semantics

The CLI should provide stable exit codes and machine-friendly command behavior for validation and diagnostics flows.

#### REQ-NF-005 - Cross-Platform Baseline

V1 should support macOS and Linux.

#### REQ-NF-006 - Traceable Provenance

Installed assets should be attributable to the CLI release that provided them.

Depends on:
- REQ-F-009

---

## First Features

### FEAT-001 - Project Initialization

Description:
- Bootstrap a target project using bundled official assets

Likely command shape:
- `opencode-helper init`

Satisfies:
- REQ-F-001
- REQ-F-002
- REQ-F-004
- REQ-F-005
- REQ-F-010

### FEAT-002 - Preset Listing

Description:
- List bundled config presets included in the current CLI release

Likely command shape:
- `opencode-helper preset list`

Satisfies:
- REQ-F-003
- REQ-F-009

### FEAT-003 - Preset Selection and Apply

Description:
- Apply a bundled preset to a target project

Likely command shape:
- `opencode-helper preset use <name>`

Satisfies:
- REQ-F-004
- REQ-F-010

### FEAT-004 - Schema Install

Description:
- Install the bundled handoff/result schemas into the project or supported scope

Likely command shape:
- `opencode-helper schema install`

Satisfies:
- REQ-F-005
- REQ-F-010

### FEAT-005 - Validation

Description:
- Validate that the local project setup matches supported expectations

Likely command shape:
- `opencode-helper validate`

Satisfies:
- REQ-F-006
- REQ-F-009
- REQ-NF-004
- REQ-NF-006

### FEAT-006 - Doctor

Description:
- Diagnose drift, missing files, incompatible state, and likely remediation paths

Likely command shape:
- `opencode-helper doctor`

Satisfies:
- REQ-F-007
- REQ-NF-004

### FEAT-007 - Version

Description:
- Report CLI version and bundled asset identity

Likely command shape:
- `opencode-helper version`

Satisfies:
- REQ-F-009
- REQ-NF-006

### FEAT-008 - Self Update

Description:
- Update the CLI from the latest GitHub release, including its bundled asset set

Likely command shape:
- `opencode-helper self-update`

Satisfies:
- REQ-F-008
- REQ-F-009
- REQ-NF-003

---

## Traceability Matrix

| ID | Type | Links To |
|---|---|---|
| DEC-001 | Decision | PRD-001, REQ-F-001, REQ-F-002, REQ-F-008, REQ-NF-002, REQ-NF-003 |
| PRD-001 | PRD | REQ-F-001 to REQ-F-010, REQ-NF-001 to REQ-NF-006 |
| REQ-F-001 | Functional Requirement | FEAT-001, FEAT-002, FEAT-003 |
| REQ-F-002 | Functional Requirement | FEAT-001, FEAT-004 |
| REQ-F-003 | Functional Requirement | FEAT-002 |
| REQ-F-004 | Functional Requirement | FEAT-001, FEAT-003 |
| REQ-F-005 | Functional Requirement | FEAT-001, FEAT-004 |
| REQ-F-006 | Functional Requirement | FEAT-005 |
| REQ-F-007 | Functional Requirement | FEAT-006 |
| REQ-F-008 | Functional Requirement | FEAT-008 |
| REQ-F-009 | Functional Requirement | FEAT-002, FEAT-005, FEAT-007, FEAT-008 |
| REQ-F-010 | Functional Requirement | FEAT-001, FEAT-003, FEAT-004 |
| REQ-NF-001 | Non-Functional Requirement | FEAT-001 to FEAT-008 |
| REQ-NF-002 | Non-Functional Requirement | FEAT-001 to FEAT-007 |
| REQ-NF-003 | Non-Functional Requirement | FEAT-008 |
| REQ-NF-004 | Non-Functional Requirement | FEAT-005, FEAT-006 |
| REQ-NF-005 | Non-Functional Requirement | FEAT-001 to FEAT-008 |
| REQ-NF-006 | Non-Functional Requirement | FEAT-005, FEAT-007 |

---

## User Story Backlog Placeholder

User stories will be added later and must:
- have stable IDs using the format `US-###`
- reference at least one feature ID
- reference at least one requirement ID
- include acceptance criteria
- identify whether the story is user-facing, maintainer-facing, or release-engineering-facing

Template:

### US-001 - Title

Type:
- User-facing | Maintainer-facing | Release-engineering-facing

Related features:
- FEAT-xxx

Related requirements:
- REQ-F-xxx
- REQ-NF-xxx

Story:
- As a ...
- I want ...
- So that ...

Acceptance criteria:
- Given ...
- When ...
- Then ...

---

## Open Questions for Post-V1

- Should the CLI support exporting presets under names other than `opencode.json`?
- Should schema install support both project-local and user-global scopes in V1, or should V1 stay project-local only?
- Should `self-update` support pinned release channels or version constraints?
- Should the CLI embed asset checksums and manifest metadata directly for stronger provenance reporting?
- Should custom user overlays on top of official presets be supported in V2?
