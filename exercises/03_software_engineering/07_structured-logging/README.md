# Structured Logging

## Theory

Structured logs produce machine-parseable output with key-value pairs. Log levels (Debug, Info, Warn, Error) let you filter output at runtime. A context-carrying logger propagates request-scoped fields to all child loggers via `With(key, value)`.

## Task

Implement a Logger with levels, fields, and With() for child loggers. Log output goes to an io.Writer in the format: `[LEVEL] key=val key2=val2 msg`

## Run

```bash
make 03_software_engineering_07_structured_logging
```

## Further practice

Add leveled output writers (stderr for errors), rotating file output.
