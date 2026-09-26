# Changelog

All notable changes to `azvm` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

### Changed

### Fixed

## [0.5.1] - 2026-09-26

### Added

- Added a looping progress animation while commands are pending

## [0.5.0] - 2026-09-26

### Added

- Added `azvm doctor <vm>` for deterministic VM diagnostics
- Added structured diagnostic checks and findings
- Added human-readable diagnosis output based on collected VM information

## [0.4.0] - 2026-09-26

### Added

- Added `azvm inspect <vm>` for displaying VM compute and network information

## [0.3.0] - 2026-09-24

### Added

- Added `azvm start <vm>`
- Added `azvm stop <vm>`
- Added `azvm restart <vm>`

## [0.2.2] - 2026-09-24

### Changed

- Refactored list command to use context for Azure VM listing

## [0.2.1] - 2026-09-24

### Changed

- Added reusable VM lookup by name
- Improved `status` command implementation
- Improved context handling across Azure operations
- Refined Azure client abstractions
- Improved test coverage and maintainability

## [0.2.0] - 2026-09-24

### Added

- Added `azvm status <vm>`

## [0.1.0] - 2026-09-23

### Added

- Added Azure CLI authentication with `azvm login`
- Added `azvm list`
- Added Azure Virtual Machine discovery
