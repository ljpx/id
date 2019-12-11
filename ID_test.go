package id

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ljpx/test"
	_ "github.com/mattn/go-sqlite3"
)

func TestNewIDIsUnique(t *testing.T) {
	// Arrange and Act.
	id1 := New()
	id2 := New()

	// Assert.
	test.That(t, id1).IsNotEqualTo(Empty)
	test.That(t, id2).IsNotEqualTo(Empty)
	test.That(t, id1).IsNotEqualTo(id2)
}

func TestEmptyIDIsAllZeros(t *testing.T) {
	// Arrange and Act and Assert.
	for i := 0; i < ByteSize; i++ {
		test.That(t, Empty[i]).IsEqualTo(byte(0))
	}
}

func TestIsValid(t *testing.T) {
	// Arrange.
	id := New()

	// Act and Assert.
	test.That(t, id.IsValid()).IsTrue()
	test.That(t, Empty.IsValid()).IsFalse()
}

func TestStringAndParseSymmetric(t *testing.T) {
	// Arrange.
	id1 := New()

	// Act.
	hexID := id1.String()
	id2, err := Parse(hexID)

	// Assert.
	test.That(t, err).IsNil()
	test.That(t, id2).IsEqualTo(id1)
}

func TestParseInvalidHex(t *testing.T) {
	// Arrange.
	hexID := strings.Repeat("F", 31) + "S"

	// Act.
	id, err := Parse(hexID)

	// Assert.
	test.That(t, id).IsEqualTo(Empty)
	test.That(t, err).IsNotNil()
}

func TestParseWrongByteSize(t *testing.T) {
	// Arrange.
	hexID := strings.Repeat("F", 34)

	// Act.
	id, err := Parse(hexID)

	// Assert.
	test.That(t, id).IsEqualTo(Empty)
	test.That(t, err).IsNotNil()
}

func TestMarshalAndUnmarshalSymmetric(t *testing.T) {
	// Arrange.
	x1 := &typeWithID{SomeID: New()}

	// Act.
	rawJSON, err := json.Marshal(x1)
	test.That(t, err).IsNil()

	x2 := &typeWithID{}
	err = json.Unmarshal(rawJSON, x2)
	test.That(t, err).IsNil()

	// Assert.
	test.That(t, x2.SomeID).IsEqualTo(x1.SomeID)
}

func TestUnmarshalInvalidHex(t *testing.T) {
	// Arrange.
	rawInvalidID := []byte(`"FFFFFH`)

	// Act.
	newID := &ID{}
	err := newID.UnmarshalJSON(rawInvalidID)

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestDatabaseUsage(t *testing.T) {
	// Arrange.
	rand.Seed(time.Now().UnixNano())
	tempFile := fmt.Sprintf("id_tests_%v.db", rand.Int63())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%v", tempFile))
	test.That(t, err).IsNil()

	defer os.Remove(tempFile)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE dummy_table (
			i  INTEGER PRIMARY KEY,
			id TEXT
		);
	`)
	test.That(t, err).IsNil()

	id1 := New()

	// Act.
	_, err = db.Exec("INSERT INTO dummy_table (i, id) VALUES(0, ?);", id1)
	test.That(t, err).IsNil()

	_, err = db.Exec("INSERT INTO dummy_table (i, id) VALUES(1, NULL);")
	test.That(t, err).IsNil()

	_, err = db.Exec("INSERT INTO dummy_table (i, id) VALUES(2, ?);", "Hello, World!")
	test.That(t, err).IsNil()

	id2 := New()
	id3 := New()
	id4 := New()
	id5 := New()

	err = db.QueryRow("SELECT id FROM dummy_table WHERE i = 0;").Scan(&id2)
	test.That(t, err).IsNil()

	err = db.QueryRow("SELECT id FROM dummy_table WHERE i = 1;").Scan(&id3)
	test.That(t, err).IsNil()

	err = db.QueryRow("SELECT id FROM dummy_table WHERE i = 2;").Scan(&id4)
	test.That(t, err).IsNotNil()

	err = db.QueryRow("SELECT i FROM dummy_table WHERE i = 2;").Scan(&id5)
	test.That(t, err).IsNotNil()

	// Assert.
	test.That(t, id2).IsEqualTo(id1)
	test.That(t, id3).IsEqualTo(Empty)
}

// -----------------------------------------------------------------------------

type typeWithID struct {
	SomeID ID
}
