## Purpose

Define how the dashboard presents unit navigation and keeps right-side operational content scoped to the selected unit or all units.

## Requirements

### Requirement: Unit navigation menu
The dashboard SHALL provide a unit navigation menu outside the main dashboard content area, with an `All` option and one option for each available printer unit.

#### Scenario: Menu is populated after data load
- **WHEN** the dashboard successfully loads printer data with one or more unit values
- **THEN** the unit navigation menu displays `All` and each distinct unit exactly once

#### Scenario: All is selected by default
- **WHEN** the dashboard first renders after successful data load
- **THEN** `All` is the active unit selection

### Requirement: Scoped dashboard rendering
The dashboard SHALL render the right-side content according to the active unit selection. The all-units scope SHALL show organization-wide summaries and aggregate charts without the bottom printer table. A selected-unit scope SHALL show summaries and charts for that unit, display the bottom printer table for that unit, and show `依單位列印狀況` as per-printer print statistics inside the selected unit.

#### Scenario: All units selected
- **WHEN** the user selects `All`
- **THEN** summary cards, daily trend chart, yearly statistics chart, unit aggregation chart, supply warnings, and incidents are calculated from all loaded printers
- **AND** the `年度總印量統計` section appears before the `依單位列印狀況` section
- **AND** the `依單位列印狀況` section is visible as an all-units unit aggregation chart
- **AND** the bottom printer table is hidden

#### Scenario: Single unit selected
- **WHEN** the user selects a specific unit
- **THEN** summary cards, daily trend chart, yearly statistics chart, supply warnings, incidents, and printer table are calculated only from printers belonging to that unit
- **AND** the bottom printer table is visible
- **AND** the `依單位列印狀況` chart shows per-printer print statistics for printers inside that selected unit

#### Scenario: Switching unit scope
- **WHEN** the user changes from one unit selection to another
- **THEN** the active menu state and all right-side dashboard sections update to match the new selection without reloading the page
- **AND** the bottom printer table visibility updates to match whether the active scope is `All` or a specific unit
- **AND** daily trend, yearly statistics, and `依單位列印狀況` chart data update from cached loaded data for the active scope

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

### Requirement: Selected-unit printer monthly chart
When a specific unit is selected, the dashboard SHALL display each printer in that unit and its month-by-month print volume for the current year in the `依單位列印狀況` section.

#### Scenario: Selected unit has printers with current-year history
- **WHEN** the user selects a unit with printers that have current-year daily history rows
- **THEN** the `依單位列印狀況` chart shows monthly print-volume values grouped by printer for that year
- **AND** no printers outside the selected unit are included

#### Scenario: Selected unit has no current-year history
- **WHEN** the user selects a unit whose printers have no current-year daily history rows
- **THEN** the selected-unit chart renders an empty or zero-data state without falling back to all-unit data

### Requirement: Printer list name field
The bottom printer list SHALL use table labels and row values that match the displayed operational data. The location field SHALL be labeled `地點`, the warranty column SHALL NOT be displayed, the color column SHALL show color capability, and the IP column SHALL show the printer IP address.

#### Scenario: Location header renders
- **WHEN** the bottom printer list renders
- **THEN** the table header formerly labeled `印表機` is labeled `地點`

#### Scenario: Warranty column is removed
- **WHEN** the bottom printer list renders
- **THEN** no `保固` header is displayed
- **AND** no warranty-days row cell is displayed in the table body

#### Scenario: Color column displays color capability
- **WHEN** a printer row renders
- **THEN** the `彩色` column displays `✔` for color-capable printers or `✘` for non-color printers
- **AND** it does not display the printer IP address

#### Scenario: IP column displays address
- **WHEN** a printer row renders
- **THEN** the `IP` column displays the printer IP address
- **AND** it does not display `✔` or `✘`

### Requirement: Homepage yearly statistics chart
The dashboard SHALL include a yearly statistics chart that summarizes print volume by year for the active scope and appears before the unit or printer print-status chart.

#### Scenario: All scope yearly chart renders
- **WHEN** the dashboard is in the `All` scope and loaded printer trend data contains one or more years
- **THEN** the yearly statistics chart displays annual total print volume for all loaded printers
- **AND** the yearly statistics chart appears before the `依單位列印狀況` unit aggregation chart

#### Scenario: Selected unit yearly chart renders
- **WHEN** the dashboard is in a selected-unit scope and that unit's loaded printer trend data contains one or more years
- **THEN** the yearly statistics chart displays annual total print volume for printers in that selected unit only
- **AND** it does not include printers from other units

#### Scenario: Yearly chart uses cached trend data
- **WHEN** the yearly statistics chart is rendered
- **THEN** it aggregates from already loaded printer trend rows for the active scope
- **AND** it does not issue additional Supabase queries

#### Scenario: Yearly chart refreshes after data reload
- **WHEN** the user manually refreshes dashboard data and the fetch succeeds
- **THEN** the yearly statistics chart is recalculated from the refreshed loaded data for the active scope

#### Scenario: No yearly data
- **WHEN** the active scope has no loaded trend rows containing yearly totals
- **THEN** the yearly chart area renders an empty state instead of stale or unrelated data
