## ADDED Requirements

### Requirement: Unit navigation menu
The dashboard SHALL provide a unit navigation menu outside the main dashboard content area, with an `All` option and one option for each available printer unit.

#### Scenario: Menu is populated after data load
- **WHEN** the dashboard successfully loads printer data with one or more unit values
- **THEN** the unit navigation menu displays `All` and each distinct unit exactly once

#### Scenario: All is selected by default
- **WHEN** the dashboard first renders after successful data load
- **THEN** `All` is the active unit selection

### Requirement: Scoped dashboard rendering
The dashboard SHALL render the right-side content according to the active unit selection.

#### Scenario: All units selected
- **WHEN** the user selects `All`
- **THEN** summary cards, trend chart, unit chart, supply warnings, incidents, and printer table are calculated from all loaded printers

#### Scenario: Single unit selected
- **WHEN** the user selects a specific unit
- **THEN** summary cards, trend chart, unit chart, supply warnings, incidents, and printer table are calculated only from printers belonging to that unit

#### Scenario: Switching unit scope
- **WHEN** the user changes from one unit selection to another
- **THEN** the active menu state and all right-side dashboard sections update to match the new selection without reloading the page

### Requirement: Responsive unit navigation layout
The dashboard SHALL keep the unit navigation usable on desktop and mobile viewport sizes without overlapping the dashboard content.

#### Scenario: Desktop layout
- **WHEN** the dashboard is viewed on a desktop-width viewport
- **THEN** the unit navigation appears on the left and the scoped dashboard content appears on the right

#### Scenario: Mobile layout
- **WHEN** the dashboard is viewed on a narrow mobile-width viewport
- **THEN** the unit navigation remains accessible and does not overlap or hide the scoped dashboard content

### Requirement: Stable data source behavior
The dashboard SHALL use the existing loaded Supabase data for unit filtering and SHALL NOT require an additional Supabase query when changing the selected unit.

#### Scenario: Selecting a unit after initial load
- **WHEN** the user selects a different unit after the dashboard has loaded
- **THEN** the dashboard updates from cached loaded data without issuing new data fetches
