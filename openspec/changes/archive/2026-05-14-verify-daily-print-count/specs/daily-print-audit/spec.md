## ADDED Requirements

### Requirement: Per-field negative delta detection
The audit tool SHALL detect when any individual metric field (print, copy, scan, fax) has a negative daily delta.

#### Scenario: Negative print delta detected
- **WHEN** the audit runs and a printer's `print_count` delta (current - previous) is negative
- **THEN** the audit reports it as a negative delta anomaly with field type "print"

#### Scenario: Negative copy delta detected
- **WHEN** the audit runs and a printer's `copy_count` delta (current - previous) is negative
- **THEN** the audit reports it as a negative delta anomaly with field type "copy"

#### Scenario: Negative scan delta detected
- **WHEN** the audit runs and a printer's `scan_total` delta (current - previous) is negative
- **THEN** the audit reports it as a negative delta anomaly with field type "scan"

#### Scenario: Negative fax delta detected
- **WHEN** the audit runs and a printer's `fax_count` delta (current - previous) is negative
- **THEN** the audit reports it as a negative delta anomaly with field type "fax"

### Requirement: Daily total consistency check
The audit tool SHALL verify that `daily_total` equals the sum of individual daily deltas (daily_print + daily_copy + daily_scan + daily_fax).

#### Scenario: daily_total matches sum of deltas
- **WHEN** `daily_total` equals `daily_print + daily_copy + daily_scan + daily_fax`
- **THEN** no anomaly is reported for this record

#### Scenario: daily_total does not match sum of deltas
- **WHEN** `daily_total` does NOT equal `daily_print + daily_copy + daily_scan + daily_fax`
- **THEN** the audit reports a consistency mismatch with both the expected and actual values

### Requirement: Per-field daily average anomaly
The audit tool SHALL check each individual metric field's daily average, not just total.

#### Scenario: Per-field daily average exceeds threshold
- **WHEN** any field's daily delta divided by gap days exceeds the configurable threshold (default 500)
- **THEN** the audit reports a per-field daily average anomaly with the field name and calculated average

#### Scenario: Threshold is configurable via URL parameter
- **WHEN** the page URL contains `?daily_avg_threshold=XXX`
- **THEN** the audit uses XXX as the daily average threshold instead of the default 500

### Requirement: Ground truth comparison mode
The audit tool SHALL support a mode that recalculates daily deltas directly from `printer_metrics` and compares them against `daily_stats`.

#### Scenario: Ground truth mode is enabled
- **WHEN** the page URL contains `?mode=ground_truth`
- **THEN** the audit fetches all `printer_metrics` records and recalculates deltas client-side

#### Scenario: Ground truth mismatch is reported
- **WHEN** a recalculated delta differs from the corresponding `daily_stats` row
- **THEN** the audit reports a ground truth mismatch with both values side by side

#### Scenario: Ground truth mode shows summary
- **WHEN** the ground truth audit completes
- **THEN** a summary shows: total records compared, matching records, mismatching records

### Requirement: Duplicate record detection
The audit tool SHALL detect if there are multiple `printer_metrics` records for the same printer on the same day.

#### Scenario: Duplicate records exist
- **WHEN** two or more `printer_metrics` records share the same `printer_id` and `recorded_at`
- **THEN** the audit reports a duplicate record anomaly with the printer ID, date, and count of duplicates

#### Scenario: No duplicate records
- **WHEN** each printer has at most one record per day
- **THEN** no duplicate anomaly is reported

### Requirement: CSV export of anomalies
The audit tool SHALL allow exporting the anomaly list to CSV.

#### Scenario: User activates export
- **WHEN** the user clicks the "Export CSV" button
- **THEN** a CSV file containing all current anomaly data is downloaded

#### Scenario: CSV includes all relevant columns
- **WHEN** the CSV is generated
- **THEN** it includes: unit, location, IP, model, date, anomaly type, current value, previous value, gap days, daily average, and risk level
