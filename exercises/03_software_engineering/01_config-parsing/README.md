# Configuration parsing

## Theory

Real programs read configuration from multiple sources: compiled-in defaults,
configuration files, and environment variables. The standard overlay pattern
is: load defaults first, apply the config file on top, then apply environment
variables on top of that. This lets every layer override the one below.

Simple text parsing avoids external dependencies. Split the raw config text on
newlines, match each line for `Key=Value` pairs, and populate a configuration
struct. Validations—such as rejecting invalid port numbers—ensure the resulting
config is usable.

```go
cfg := Defaults()
cfg = overlayFile(cfg, fileContents)
cfg = overlayEnv(cfg, env)
```

## Task

Implement `LoadConfig` so that it builds a `Config` value by starting with
defaults, overlaying values from a simple `Key=Value` config text, and
overriding non-empty values from environment variables. Return an error when
the port is out of range.

No external formats (JSON, TOML, YAML) are required—use simple string
splitting on newlines and `=` characters.

## Run

```bash
make 03_software_engineering_01_config_parsing
```

## Further practice

Add support for a `--config` flag that reads from a file path, and implement
nested config keys using `.` notation (e.g. `database.host=localhost`).
