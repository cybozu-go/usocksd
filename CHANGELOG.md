# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - 2016-09-01
### Added
- usocksd now has its own SOCKS implementation under `socks` directory.
- Support for SOCKS4/4a thanks to the new implementation.
- usocksd now adopts [github.com/cybozu-go/cmd][cmd] framework.  
  As a result, it implements [the common spec][spec] including graceful restart.

### Changed
- The default configuration file path is now `/etc/usocksd.toml`.
- `log.file` config key is renamed to `log.filename`.

[cmd]: https://github.com/cybozu-go/cmd
[spec]: https://github.com/cybozu-go/cmd/blob/master/README.md#specifications
[Unreleased]: https://github.com/cybozu-go/usocksd/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/cybozu-go/usocksd/compare/v0.1...v1.0.0
