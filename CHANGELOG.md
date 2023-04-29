# Changelog

## [Unreleased]

## [v3.0.0] - 2023-04-29

### Changed
- Config struct `PostLoad` hooks are now automatically called by `Config.Load`. [#17](https://github.com/go-nacelle/config/pull/17)

### Added
- Added `LoadError` to be returned when loading fails. [#17](https://github.com/go-nacelle/config/pull/17)
- Added `PostLoadError` to be returned when a config struct's `PostLoad` hook fails. [#17](https://github.com/go-nacelle/config/pull/17)
- Added `SerializeError` to be returned when value serialization into a config struct field fails. [#17](https://github.com/go-nacelle/config/pull/17)

### Removed
- `Config.PostLoad` was removed from the public API. [#17](https://github.com/go-nacelle/config/pull/17)

## [v2.0.1] - 2022-10-10

### Added

- Added `WithConfig` and `FromContext`. [#15](https://github.com/go-nacelle/config/pull/15)

## [v2.0.0] - 2021-05-31

### Added

- Added `Describe` method to `Config` interface. [#8](https://github.com/go-nacelle/config/pull/8)
- Added `WithLogger` and `WithMaskedKeys` to replace `NewLoggingConfig`. [#11](https://github.com/go-nacelle/config/pull/11)

### Removed

- Removed mocks package. [#9](https://github.com/go-nacelle/config/pull/9)
- Removed `MustLoad` from `Config` interface. [#10](https://github.com/go-nacelle/config/pull/10)
- Removed `NewLoggingConfig`. [#11](https://github.com/go-nacelle/config/pull/11)

### Changed

- Split `Load` method in the `Config` interface into `Load` and `PostLoad` methods. [#7](https://github.com/go-nacelle/config/pull/7)
- The `Config` interface is now a struct with the same name and set of methods. [#12](https://github.com/go-nacelle/config/pull/12)

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

### Changed

- Migrated from [efritz/zubrin](https://github.com/efritz/zubrin).

[Unreleased]: https://github.com/go-nacelle/config/compare/v2.0.1...HEAD
[v1.0.0]: https://github.com/go-nacelle/config/releases/tag/v1.0.0
[v1.1.0]: https://github.com/go-nacelle/config/compare/v1.0.0...v1.1.0
[v1.2.0]: https://github.com/go-nacelle/config/compare/v1.1.0...v1.2.0
[v1.2.1]: https://github.com/go-nacelle/config/compare/v1.2.0...v1.2.1
[v2.0.0]: https://github.com/go-nacelle/config/compare/v1.2.1...v2.0.0
[v2.0.1]: https://github.com/go-nacelle/config/compare/v2.0.0...v2.0.1
