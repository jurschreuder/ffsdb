Fixed Float Slice Database
==========================

Super simple Golang DB to be used when iterating over []float64 slices of fixed length, when they do not fit in RAM anymore.

# Why?

I use this for machine learning, when some times you have a lot of float slices you want to iterate over.

## install
```
go get github.com/jurschreuder/ffsdb
```

## create a new database 
create at path test.db\
save []float64 slices of length 256\
overwrite old database
```go
fdb, err := NewFfsdb("test.db", 256, true) // (filepath, []float64 length, overwrite old file)
if err != nil {
    panic(err)
}
defer fdb.Close()
```

## add new entry to the database
where foo is a []float64
```go
foo := make([]float64, 256)
err := fdb.Add(foo)
```

## read all entries in the datebase from the beginning
where foo is a []float64
```go
fdb.Rewind()
ok := true
for ok {
    foo, ok = fdb.ReadNext()
}
```

## read an entry at a specific id
```go
id = int64(100)
vals, err := fdb.ReadId(id)
```

## update an entry at a specific id
```go
id = int64(100)
foo = make([]float64, 256)
err := fdb.Update(id, foo)
```

## performance
for []float64 with length 256\
on 2018 mac book
```
ffsdb_test.go:27: added 100000 in time: 1.62976014s
ffsdb_test.go:41: read 100000 in time: 580.782646ms
```
