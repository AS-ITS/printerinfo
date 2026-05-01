## MODIFIED Requirements

### Requirement: Scoped dashboard rendering
The dashboard SHALL render the right-side content according to the active unit selection. The all-units scope SHALL show organization-wide summaries and aggregate charts without the bottom printer table. A selected-unit scope SHALL show summaries for that unit, replace the unit aggregation chart with a current-year monthly print-volume chart by printer, and display the bottom printer table for that unit.

#### Scenario: All units selected
- **WHEN** the user selects `All`
- **THEN** summary cards, trend chart, yearly statistics chart, unit aggregation chart, supply warnings, and incidents are calculated from all loaded printers
- **AND** the `年度總印量統計` section appears before the `依單位列印狀況` section
- **AND** the `依單位列印狀況` section is visible as an all-units unit aggregation chart
- **AND** the bottom printer table is hidden

#### Scenario: Single unit selected
- **WHEN** the user selects a specific unit
- **THEN** summary cards, trend chart, supply warnings, incidents, and printer table are calculated only from printers belonging to that unit
- **AND** the bottom printer table is visible
- **AND** the `依單位列印狀況` unit aggregation view is replaced by a current-year monthly print-volume chart for each printer in the selected unit

#### Scenario: Switching unit scope
- **WHEN** the user changes from one unit selection to another
- **THEN** the active menu state and all right-side dashboard sections update to match the new selection without reloading the page
- **AND** the bottom printer table visibility updates to match whether the active scope is `All` or a specific unit

### Requirement: Homepage yearly statistics chart
The all-units homepage dashboard SHALL include a yearly statistics chart that summarizes print volume by year across all loaded printers and appears before the unit aggregation chart.

#### Scenario: All scope yearly chart renders
- **WHEN** the dashboard is in the `All` scope and loaded printer trend data contains one or more years
- **THEN** the homepage displays a yearly statistics chart with annual total print volume for each year
- **AND** the yearly statistics chart appears before the `依單位列印狀況` unit aggregation chart

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
