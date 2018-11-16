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
```

## add new entry to the database
where foo is a []float64
```go
err := fdb.Add(foo)
```

## read all entries in the datebase from the beginning
wher foo is a []float64
```go
fdb.Rewind()
ok := true
for ok {
    foo, ok = fdb.ReadNext()
}
```

### ToDo
make ReadNext() faster
