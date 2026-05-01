## Why

The bottom printer table currently contains a warranty column that is no longer needed, and the color-capability indicator is rendered under the IP column while IP appears under the color column. The column labels and values should match their operational meaning so users can scan the table without confusion.

## What Changes

- Remove the warranty column from the bottom printer table header and rows.
- Swap the rendered values for the `彩色` and `IP` columns so `彩色` shows `✔` or `✘`, and `IP` shows the printer IP address.
- Rename the table header `印表機` to `地點`.
- Preserve the existing printer row actions, responsive hidden-column behavior, and remaining data columns.

## Capabilities

### New Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: Refines the printer table display requirements for column labels, removed warranty column, and correct color/IP column values.

## Impact

- Affected code: `index.html`
- Affected UI: bottom printer table header and row markup
- No data model, Supabase query, Go script, authentication, or deployment workflow changes expected
