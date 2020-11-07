# Changelog

## [Unreleased]

## [v1.2.1] - 2020-09-30

### Removed

- Removed dependency on [aphistic/sweet](https://github.com/aphistic/sweet) by rewriting tests to use [testify](https://github.com/stretchr/testify). [#5](https://github.com/go-nacelle/config/pull/5)

## [v1.2.0] - 2020-04-01

### Added

- Added `FlagSourcer` that reads configuration values from the command line. [#3](https://github.com/go-nacelle/config/pull/3)
- Added `Init` method to `Config` and `Sourcer`. [#4](https://github.com/go-nacelle/config/pull/4)

## [v1.1.0] - 2019-09-05

### Added

- Added options to supply a filesystem adapter to glob, file, and directory sourcers. [#2](https://github.com/go-nacelle/config/pull/2)

## [v1.0.0] - 2019-06-17

[Unreleased]: https://github.com/go-nacelle/config/compare/v1.2.1...HEAD
[v1.0.0]: https://github.com/go-nacelle/config/releases/tag/v1.0.0
[v1.1.0]: https://github.com/go-nacelle/config/compare/v1.0.0...v1.1.0
[v1.2.0]: https://github.com/go-nacelle/config/compare/v1.1.0...v1.2.0
[v1.2.1]: https://github.com/go-nacelle/config/compare/v1.2.0...v1.2.1
