## ADDED Requirements

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
