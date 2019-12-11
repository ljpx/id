package id

import (
	"crypto/rand"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ID represents a unique 128-bit identifier.
type ID [ByteSize]byte

var _ fmt.Stringer = ID{}
var _ json.Marshaler = ID{}
var _ json.Unmarshaler = &ID{}
var _ driver.Valuer = ID{}
var _ sql.Scanner = &ID{}

// ByteSize defines the size of the identifiers produced by this package, in
// bytes.
const ByteSize = 16

// Empty represents the empty ID (a sequence of 128 0 bits).
var Empty = ID{}

// New creates a new ID sourced from the current time and a CSPRNG.
func New() ID {
	newID := ID{}
	timestamp := time.Now().UnixNano()

	binary.BigEndian.PutUint64(newID[:ByteSize/2], uint64(timestamp))
	rand.Read(newID[ByteSize/2:])

	return newID
}

// String returns the string representation of the identifier.
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// Parse parses the provided hex string, returning an error if parsing fails.
func Parse(hexID string) (ID, error) {
	rawID, err := hex.DecodeString(hexID)
	if err != nil {
		return Empty, fmt.Errorf("parsing ID failed: %w", err)
	}

	if len(rawID) != ByteSize {
		format := "parsing ID failed: expected byte length to be %v bytes, but was %v bytes"
		return Empty, fmt.Errorf(format, ByteSize, len(rawID))
	}

	newID := ID{}
	copy(newID[:], rawID)

	return newID, nil
}

// MarshalJSON returns the JSON value of an ID.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, id.String())), nil
}

// UnmarshalJSON returns the parsed ID from a JSON value.
func (id *ID) UnmarshalJSON(raw []byte) error {
	hexID := strings.ReplaceAll(string(raw), `"`, "")

	newID, err := Parse(hexID)
	if err != nil {
		return err
	}

	*id = newID
	return nil
}

// Value is implemented to be database friendly.
func (id ID) Value() (driver.Value, error) {
	return id.String(), nil
}

// Scan is used to read out of a database.
func (id *ID) Scan(src interface{}) error {
	if src == nil {
		*id = Empty
		return nil
	}

	hexID, ok := src.(string)
	if !ok {
		return errors.New("expected to scan a string to ID, but was not a string")
	}

	newID, err := Parse(hexID)
	if err != nil {
		return err
	}

	*id = newID
	return nil
}

// IsValid checks if an ID is valid.
func (id ID) IsValid() bool {
	return id != Empty
}
