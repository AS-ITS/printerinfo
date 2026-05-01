## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Selected-unit printer monthly chart
When a specific unit is selected, the dashboard SHALL display each printer in that unit and its month-by-month print volume for the current year.

#### Scenario: Selected unit has printers with current-year history
- **WHEN** the user selects a unit with printers that have current-year daily history rows
- **THEN** the selected-unit chart shows monthly print-volume values grouped by printer for that year

#### Scenario: Selected unit has no current-year history
- **WHEN** the user selects a unit whose printers have no current-year daily history rows
- **THEN** the selected-unit chart renders an empty or zero-data state without falling back to all-unit data

### Requirement: Printer list name field
The bottom printer list SHALL display the unit name in the prominent printer-name field.

#### Scenario: Printer row renders name field
- **WHEN** the bottom printer list renders a printer row
- **THEN** the prominent name text in the printer-name field shows that printer's unit name
- **AND** it does not use the printer/location label as the prominent name text
