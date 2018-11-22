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

func TestSeek(t *testing.T) {
	// create a new database at path 'test.db'
	fdb, err := NewFfsdb("test.db", 2, true)
	if err != nil {
		t.Fatal(err)
	}

	// example float64 slice to save in the database
	foo := make([]float64, 2)

	start := time.Now()
	for i := 0; i < 10; i++ {
		// add new entry to the database
		foo[0] = float64(i)
		err := fdb.Add(foo)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("added 10 in time:", time.Since(start))

	// test if we rewind if we still write at the end of the file
	fdb.Rewind()

	start = time.Now()
	for i := 10; i < 20; i++ {
		// add new entry to the database
		foo[0] = float64(i)
		err := fdb.Add(foo)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("added 10 in time:", time.Since(start))

	fdb.Rewind()
	start = time.Now()
	i := 0

	// read all entries from the database
	ok := true
	for ok {
		foo, ok = fdb.ReadNext()
		if ok {
			t.Log(i, "-", foo)
			i++
		}
	}
	t.Log("read", i, "in time:", time.Since(start))

	// read at a specific id
	bar, err := fdb.ReadId(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id 10:", bar)

	// update an id
	foo = []float64{100., 100.}
	err = fdb.Update(10, foo)
	if err != nil {
		t.Fatal(err)
	}
	bar, err = fdb.ReadId(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id 10 after update:", bar)

}
