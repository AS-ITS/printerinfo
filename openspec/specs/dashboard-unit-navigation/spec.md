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
The dashboard SHALL render the right-side content according to the active unit selection. The all-units scope SHALL show organization-wide summaries and the unit aggregation chart. A selected-unit scope SHALL show summaries for that unit and replace the unit aggregation chart with a current-year monthly print-volume chart by printer.

#### Scenario: All units selected
- **WHEN** the user selects `All`
- **THEN** summary cards, trend chart, unit aggregation chart, supply warnings, incidents, and printer table are calculated from all loaded printers
- **AND** the `依單位列印狀況` section is visible as an all-units unit aggregation chart

#### Scenario: Single unit selected
- **WHEN** the user selects a specific unit
- **THEN** summary cards, trend chart, supply warnings, incidents, and printer table are calculated only from printers belonging to that unit
- **AND** the `依單位列印狀況` unit aggregation view is replaced by a current-year monthly print-volume chart for each printer in the selected unit

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

### Requirement: Selected-unit printer monthly chart
When a specific unit is selected, the dashboard SHALL display each printer in that unit and its month-by-month print volume for the current year.

#### Scenario: Selected unit has printers with current-year history
- **WHEN** the user selects a unit with printers that have current-year daily history rows
- **THEN** the selected-unit chart shows monthly print-volume values grouped by printer for that year

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
The all-units homepage dashboard SHALL include a yearly statistics chart that summarizes print volume by year across all loaded printers.

#### Scenario: All scope yearly chart renders
- **WHEN** the dashboard is in the `All` scope and loaded printer trend data contains one or more years
- **THEN** the homepage displays a yearly statistics chart with annual total print volume for each year

#### Scenario: Yearly chart uses cached trend data
- **WHEN** the yearly statistics chart is rendered
- **THEN** it aggregates from already loaded printer trend rows
- **AND** it does not issue additional Supabase queries

#### Scenario: Yearly chart refreshes after data reload
- **WHEN** the user manually refreshes dashboard data and the fetch succeeds
- **THEN** the yearly statistics chart is recalculated from the refreshed loaded data

#### Scenario: No yearly data
- **WHEN** the dashboard is in the `All` scope and no loaded trend rows contain yearly totals
- **THEN** the yearly chart area renders an empty state instead of stale or unrelated data
