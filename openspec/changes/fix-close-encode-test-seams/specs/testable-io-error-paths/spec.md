## ADDED Requirements

### Requirement: Drawer I/O failure seams
The `entity` drawer implementation SHALL route JPEG file close and JPEG encode operations used by `loadImage` and `SaveImage` through package-level seams that default to the standard library `Close` and `jpeg.Encode` behavior, so tests can force close-after-decode and encode-failure branches without changing recognition behavior.

#### Scenario: Default production behavior unchanged
- **WHEN** the application loads and saves images without tests replacing seams
- **THEN** drawer I/O uses real file close and JPEG encode

#### Scenario: Close-after-decode failure testable
- **WHEN** a test replaces the drawer close seam to return an error after a successful decode
- **THEN** `loadImage` surfaces that close error

#### Scenario: Encode failure testable
- **WHEN** a test replaces the JPEG encode seam to return an error
- **THEN** `SaveImage` returns that encode error

### Requirement: Load-persons directory close seam
The `LoadPersonsUseCase` implementation SHALL close the opened persons directory through a package-level seam that defaults to `(*os.File).Close`, so tests can force the close-error logging branch.

#### Scenario: Default close uses real Close
- **WHEN** persons are loaded outside tests that replace the seam
- **THEN** the persons directory file is closed with the real `Close` method

#### Scenario: Close-error branch testable
- **WHEN** a test replaces the persons-dir close seam to return an error after a successful load
- **THEN** Execute still returns the loaded persons (close failure remains non-fatal) and the close-error path is exercised

### Requirement: Coverage verification for sealed branches
After seams and tests are in place, the project SHALL be able to reach 100% statement coverage for the previously blocked drawer close/encode and load-persons close branches under `go test ./... -coverprofile=… -coverpkg=./...`.

#### Scenario: Blocked statements covered
- **WHEN** the seam-based tests run as part of the module coverprofile measurement
- **THEN** the close-after-decode, encode-error, and persons-dir close-error statements are covered
