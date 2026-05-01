## 1. Table Header Changes

- [x] 1.1 Rename the bottom table header from `印表機` to `地點`.
- [x] 1.2 Remove the `保固` header from the bottom table.

## 2. Table Row Changes

- [x] 2.1 Remove the warranty-days table cell from each printer row.
- [x] 2.2 Move the color capability output (`✔` or `✘`) into the `彩色` column.
- [x] 2.3 Move the printer IP address into the `IP` column.
- [x] 2.4 Preserve remaining row actions and data columns.

## 3. Verification

- [x] 3.1 Verify generated table markup has no `保固` header or warranty-days row cell.
- [x] 3.2 Verify the `彩色` column renders only `✔` or `✘`.
- [x] 3.3 Verify the `IP` column renders the printer IP address.
- [x] 3.4 Run JavaScript syntax validation and existing Go tests.
