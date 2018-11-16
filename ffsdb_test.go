package ffsdb

import (
	"testing"
	"time"
)

func TestFfsdb(t *testing.T) {
	// create a new database at path 'test.db'
	fdb, err := NewFfsdb("test.db", 256, true)
	if err != nil {
		t.Fatal(err)
	}

	// example float64 slice to save in the database
	foo := make([]float64, 256)

	start := time.Now()
	for i := 0; i < 100e3; i++ {

		// add new entry to the database
		err := fdb.Add(foo)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("added 100000 in time:", time.Since(start))

	// read the datebase from the beginning
	fdb.Rewind()

	start = time.Now()
	i := 0

	// read all entries from the database
	ok := true
	for ok {
		foo, ok = fdb.ReadNext()
		i++
	}
	t.Log("read", i-1, "in time:", time.Since(start))
}
