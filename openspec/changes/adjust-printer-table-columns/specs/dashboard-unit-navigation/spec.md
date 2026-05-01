## MODIFIED Requirements

### Requirement: Printer list name field
The bottom printer list SHALL use table labels and row values that match the displayed operational data. The location field SHALL be labeled `地點`, the warranty column SHALL NOT be displayed, the color column SHALL show color capability, and the IP column SHALL show the printer IP address.

#### Scenario: Location header renders
- **WHEN** the bottom printer list renders
- **THEN** the table header formerly labeled `印表機` is labeled `地點`

#### Scenario: Warranty column is removed
- **WHEN** the bottom printer list renders
- **THEN** no `保固` header is displayed
- **AND** no warranty-days row cell is displayed in the table body

#### Scenario: Color column displays color capability
- **WHEN** a printer row renders
- **THEN** the `彩色` column displays `✔` for color-capable printers or `✘` for non-color printers
- **AND** it does not display the printer IP address

#### Scenario: IP column displays address
- **WHEN** a printer row renders
- **THEN** the `IP` column displays the printer IP address
- **AND** it does not display `✔` or `✘`
