# Values and basic types

## Theory

Go gives common values concrete types such as `string`, `int`, `float64`, and
`bool`. Function signatures make those choices visible. Formatting verbs also
reflect types: `%s` formats a string, `%d` an integer, `%.2f` a floating-point
number, and `%t` a boolean.

```go
label := fmt.Sprintf("%s: %.1f", "temperature", 21.5)
```

Arithmetic operands should have compatible types. A Celsius-to-Fahrenheit
conversion can stay in `float64` so fractional temperatures are preserved.

## Task

Implement `ProfileSummary` so it produces the exact tested sentence and uses
`statusLabel` to map the boolean state to `active` or `inactive`. Implement
`CelsiusToFahrenheit` without rounding the result.

## Run

```bash
make 01_core_01_values_and_types
```

## Further practice

Add a Kelvin conversion and decide how your API should represent values below
absolute zero.

