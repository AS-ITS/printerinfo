## MODIFIED Requirements

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

### Requirement: Selected-unit printer monthly chart
When a specific unit is selected, the dashboard SHALL display each printer in that unit and its month-by-month print volume for the current year in the `依單位列印狀況` section.

#### Scenario: Selected unit has printers with current-year history
- **WHEN** the user selects a unit with printers that have current-year daily history rows
- **THEN** the `依單位列印狀況` chart shows monthly print-volume values grouped by printer for that year
- **AND** no printers outside the selected unit are included

#### Scenario: Selected unit has no current-year history
- **WHEN** the user selects a unit whose printers have no current-year daily history rows
- **THEN** the selected-unit chart renders an empty or zero-data state without falling back to all-unit data

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
