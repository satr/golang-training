# Config: Environment Merge

## Theory

Configuration should follow a precedence chain: hardcoded defaults, then file values, then environment variables. Each layer overrides the one below. The `os` package provides `os.Environ` and `os.LookupEnv` for environment access, and `os/user` for the home directory.

## Task

Implement a Config struct that merges defaults, optional file paths, and environment variable overrides. Environment variables use the pattern `APP_<KEY>`.

## Run

```bash
make 03_software_engineering_12_config_env_merge
```

## Further practice

Add YAML/JSON file parsing, secrets provider integration.
