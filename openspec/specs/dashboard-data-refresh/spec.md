## Purpose

Define when the dashboard loads Supabase data, when it reuses cached page-session data, and how the user can manually refresh the loaded data.

## Requirements

### Requirement: Initial dashboard data load
The dashboard SHALL load data from Supabase when an authenticated user enters the dashboard for the first time in the current page session.

#### Scenario: First authenticated entry
- **WHEN** an authenticated session is detected and dashboard data has not been loaded in the current page session
- **THEN** the dashboard fetches printer, supply, metric, and incident data from Supabase

### Requirement: Cached render on focus or repeated auth notification
The dashboard SHALL NOT reload Supabase data solely because the page regains focus or receives a repeated auth notification for the same authenticated user.

#### Scenario: Page focus returns after data is loaded
- **WHEN** dashboard data has already been loaded and the page regains focus
- **THEN** the dashboard keeps using cached data without issuing new Supabase data fetches

#### Scenario: Same user auth callback repeats
- **WHEN** dashboard data has already been loaded for the current authenticated user and an auth callback reports the same user again
- **THEN** the dashboard keeps using cached data without issuing new Supabase data fetches

### Requirement: Manual refresh reload
The dashboard SHALL provide a user-triggered refresh action that reloads data from Supabase and re-renders the current dashboard scope.

#### Scenario: User clicks refresh
- **WHEN** the user activates the dashboard refresh action
- **THEN** the dashboard fetches fresh printer, supply, metric, and incident data from Supabase
- **AND** the dashboard re-renders the currently selected unit scope after the fetch succeeds

#### Scenario: Last update reflects fetch time
- **WHEN** a Supabase data fetch succeeds
- **THEN** the dashboard updates the last-update display to the successful fetch time
