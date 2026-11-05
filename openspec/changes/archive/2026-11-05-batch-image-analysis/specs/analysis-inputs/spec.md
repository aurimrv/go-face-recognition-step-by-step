## ADDED Requirements

### Requirement: Default analysis target when no arguments are given
When the program is invoked with no path arguments, the system SHALL analyze `images/unknown.jpg` and write the annotated image to `images/result.jpg`.

#### Scenario: No-argument run
- **WHEN** the user runs the program with no path arguments
- **THEN** the system analyzes `images/unknown.jpg` and saves the output as `images/result.jpg`

### Requirement: Explicit image file arguments
The system SHALL accept one or more filesystem paths to image files as arguments and analyze each JPEG file, writing `result_<stem>.jpg` in the same directory as that source file, where `<stem>` is the source filename without its extension.

#### Scenario: Single explicit JPEG
- **WHEN** the user runs the program with argument `images/a.jpg`
- **THEN** the system analyzes that file and saves `images/result_a.jpg`

#### Scenario: Multiple explicit JPEGs
- **WHEN** the user runs the program with arguments `images/a.jpg` and `images/b.jpg`
- **THEN** the system analyzes each file and saves `images/result_a.jpg` and `images/result_b.jpg`

#### Scenario: Explicit path uses stem naming not the default result name
- **WHEN** the user runs the program with argument `images/unknown.jpg`
- **THEN** the system saves `images/result_unknown.jpg` and MUST NOT replace the no-arg-only name `images/result.jpg` for that invocation

### Requirement: Directory arguments include subfolders
The system SHALL accept directory path arguments and recursively discover JPEG files under each directory (including subfolders). For each discovered JPEG, the system SHALL write `result_<stem>.jpg` beside that source file.

#### Scenario: Album directory
- **WHEN** the user runs the program with a directory argument containing JPEG files in nested subfolders
- **THEN** each discovered JPEG is analyzed and a `result_<stem>.jpg` is written next to that source file

#### Scenario: Mixed file and directory arguments
- **WHEN** the user runs the program with both file and directory path arguments
- **THEN** the system analyzes the explicit files and all JPEGs discovered under the directories

### Requirement: JPEG-only inputs
The system SHALL treat only files with JPEG extensions (case-insensitive `.jpg` and `.jpeg`) as analysis inputs when discovering files under directories. Non-JPEG files MUST NOT be selected by directory discovery.

#### Scenario: Directory contains non-JPEG media
- **WHEN** a walked directory contains both `.jpg` and `.avif` files
- **THEN** the system analyzes the `.jpg` files and does not select the `.avif` files as inputs

### Requirement: Skip previously generated result images
The system SHALL NOT select as analysis inputs any file whose base name is `result.jpg` or whose base name matches the pattern `result_*.jpg` produced by this tool.

#### Scenario: Second run over the same album
- **WHEN** a directory already contains `result_photo.jpg` next to `photo.jpg` and the user analyzes that directory again
- **THEN** the system analyzes `photo.jpg` and MUST NOT analyze `result_photo.jpg`

### Requirement: Reuse recognizer across jobs
The system SHALL initialize the face recognizer and known-person samples once per process run and reuse them for every analysis job in that run.

#### Scenario: Multiple jobs in one invocation
- **WHEN** the resolved job list contains more than one image
- **THEN** each image is classified using the same initialized recognizer and sample set without reloading persons or models between jobs

### Requirement: Batch continues after per-image failure
When the invocation includes explicit path arguments (one or more), the system SHALL continue processing remaining jobs after a per-image failure (including decode failure or inability to find faces), record that failure, and exit with a non-zero status if any job failed.

#### Scenario: One bad file among several
- **WHEN** the user passes multiple image paths and one image cannot be decoded or yields no faces
- **THEN** the system reports the failure for that image, continues with the remaining images, and exits non-zero if any image failed

#### Scenario: Successful jobs still write outputs
- **WHEN** some jobs succeed and others fail in the same explicit-argument run
- **THEN** successful jobs still write their `result_<stem>.jpg` outputs
