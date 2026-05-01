## 1. Legend Navigation

- [x] 1.1 Add an all-units chart legend click handler in `renderUnitAggregationChart`.
- [x] 1.2 Resolve the clicked legend entry to the corresponding unit name.
- [x] 1.3 Call the existing `filterByUnit(unit)` path when a valid unit legend entry is clicked.
- [x] 1.4 Preserve selected-unit monthly printer chart legend behavior without unit navigation.

## 2. State And Data Behavior

- [x] 2.1 Verify legend navigation updates the active left-side unit menu state.
- [x] 2.2 Verify legend navigation updates all scoped dashboard sections through cached data.
- [x] 2.3 Verify legend navigation does not issue new Supabase fetches.

## 3. Verification

- [x] 3.1 Add or run a JavaScript behavior check for all-units legend navigation.
- [x] 3.2 Verify selected-unit chart legend clicks do not navigate units.
- [x] 3.3 Run JavaScript syntax validation and existing Go tests.
