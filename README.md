![](icon.svg)

# id

[![GoDoc](https://godoc.org/github.com/ljpx/id?status.svg)](https://godoc.org/github.com/ljpx/id)

Package `id` defines a globally unique 128-bit identifier.

## Usage Example

Calling `id.New()` produces a new, unique `id.ID`:

```go
userID := id.New()
fmt.Printf("%v", userID)
```

```bash
15bffdea7ae1138c6c57ad5fb01ed08a
```

## Structure

An `id.ID` is actually defined as `[16]byte` under the hood.  Generating a new
`id.ID` takes the current time, with nanosecond precision, and converts this to
a `uint64`.  This becomes the first 64 bits of the identifier.  The last 64 bits
are read from a CSPRNG.

## Interfaces

The `id.ID` type satisfies a number of interfaces:

- `fmt.Stringer` - Returns the identifier in string format, e.g. `"15bffdea7ae1138c6c57ad5fb01ed08a"`

- `json.Marshaler` - Marshals to a JSON string, like above.

- `json.Unmarshaler` - Marshals from a JSON string.

- `driver.Value` - Converts to a string to be stored in SQL databases.

- `sql.Scanner` - Parses strings from SQL queries.
