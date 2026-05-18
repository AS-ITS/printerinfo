## 1. Per-field negative delta detection

- [ ] 1.1 Add per-field (print/copy/scan/fax) negative delta check in runAudit()
- [ ] 1.2 Add "negative_delta" badge type and summary card for negative delta count
- [ ] 1.3 Add per-field negative delta entries to allAnomalies with risk "danger"

## 2. Daily total consistency check

- [ ] 2.1 Add daily_total consistency check (daily_total vs sum of daily_print/copy/scan/fax) in runAudit()
- [ ] 2.2 Add "consistency_mismatch" badge type and summary card
- [ ] 2.3 Include expected vs actual values in the anomaly table row

## 3. Per-field daily average anomaly

- [ ] 3.1 Add per-field daily average calculation (each field delta / gapDays) in runAudit()
- [ ] 3.2 Add per-field daily average anomaly detection with configurable threshold
- [ ] 3.3 Add "daily_avg_field" badge type and summary card
- [ ] 3.4 Add ?daily_avg_threshold=XXX URL param support (default 500)

## 4. Ground truth comparison mode

- [ ] 4.1 Add ?mode=ground_truth URL param detection and conditional logic
- [ ] 4.2 Fetch all printer_metrics records directly in ground_truth mode
- [ ] 4.3 Implement client-side delta recalculation matching daily_stats logic
- [ ] 4.4 Compare recalculated deltas against daily_stats rows and report mismatches
- [ ] 4.5 Add ground truth summary: total compared, matching, mismatching counts

## 5. Duplicate record detection

- [ ] 5.1 Add duplicate detection logic (same printer_id + recorded_at count > 1)
- [ ] 5.2 Add "duplicate" badge type and summary card
- [ ] 5.3 Include printer ID, date, and duplicate count in anomaly table

## 6. CSV export

- [ ] 6.1 Add "Export CSV" button to the UI
- [ ] 6.2 Implement CSV generation with Blob URL download
- [ ] 6.3 Include all anomaly columns: unit, location, IP, model, date, anomaly type, values, gap days, daily avg, risk

## 7. Summary UI updates

- [ ] 7.1 Update renderSummary() to show new anomaly categories
- [ ] 7.2 Update audit table columns to include anomaly type column
- [ ] 7.3 Ensure RWD layout still works with new columns
