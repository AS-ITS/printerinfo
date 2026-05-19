## MODIFIED Requirements

### Requirement: Scoped dashboard rendering
The dashboard SHALL render the right-side content according to the active unit selection. The all-units scope SHALL show organization-wide summaries and aggregate charts without the bottom printer table. A selected-unit scope SHALL show summaries and charts for that unit, display the bottom printer table for that unit, and show `依單位列印狀況` as per-printer print statistics inside the selected unit.

#### Scenario: Daily trend chart renders all available data
- **WHEN** the dashboard loads and `printer_metrics` contains rows after 4/15
- **THEN** the daily trend chart MUST display all available daily data points including dates after 4/15
- **AND** the daily_total value MUST be correctly calculated using the delta between consecutive non-zero readings

#### Scenario: Daily trend chart handles data gaps
- **WHEN** there are missing dates in `printer_metrics` (no readings for some days)
- **THEN** the daily trend chart MUST NOT display false zero values for missing dates
- **AND** the daily trend chart MUST show only dates where at least one printer has a valid metric reading

#### Scenario: Daily trend chart with custom date range
- **WHEN** the user selects a custom date range that includes dates after 4/15
- **THEN** the daily trend chart MUST correctly aggregate daily_total for all dates within the selected range
- **AND** the daily_total calculation MUST NOT be affected by data collection gaps
