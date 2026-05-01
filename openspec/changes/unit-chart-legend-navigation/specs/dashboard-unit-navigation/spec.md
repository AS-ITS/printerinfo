## ADDED Requirements

### Requirement: Unit chart legend navigation
When the dashboard is showing the all-units `依單位列印狀況` chart, each unit legend entry SHALL navigate to that unit's scoped dashboard view.

#### Scenario: Click unit legend entry
- **WHEN** the dashboard is in the `All` scope and the user clicks a legend entry for a unit in the `依單位列印狀況` chart
- **THEN** the dashboard selects that unit
- **AND** the result matches clicking the same unit in the left-side unit menu

#### Scenario: Legend navigation reuses cached data
- **WHEN** the user navigates to a unit by clicking the all-units chart legend
- **THEN** the dashboard updates from cached loaded data without issuing new Supabase data fetches

#### Scenario: Selected-unit chart legend does not navigate units
- **WHEN** the dashboard is already scoped to a unit and shows the per-printer monthly chart
- **THEN** clicking chart legend entries does not navigate to a different unit
