# Perkakas Influx Library
This library helps you to send point to influx.
It supports sending one point at a time or you can create
batch points and send it at once.

# How To Send Point
```go
config := ClientConfig{
    Addr:               "http://localhost:8086",
    Database:           "myDB",
    Timeout:            5 * time.Second,
}

// Please note that tags is indexed.
// So please consider place a value to tags
// if you want to filter it on your query
// to achieve performance.
tags := Tags{
    "id": "123456",
    "price": "10000",
}

// Fields is not indexed, so you can consider it
// as additional data
fields := Fields{
    "name": "Mechanical Keyboard",
    "manufacturer": "Logitech",
}

c, err := NewClient(config)
if err != nil {
    t.Log(err)
    t.FailNow()
}

// You can say that the first parameter is the table name in SQL.
// The precision argument specifies the format/precision of any 
// returned timestamps.
c.WritePoints("products", tags, fields, "s")
```

# How To Send Batch Points
```go
config := ClientConfig{
    Addr:               "http://localhost:8086",
    Database:           "myDB",
    Timeout:            5 * time.Second,
}

tags := Tags{
    "id": "123456",
    "price": "10000",
}

fields := Fields{
    "name": "Mechanical Keyboard",
    "manufacturer": "Logitech",
}

c, err := NewClient(config)
if err != nil {
    t.Log(err)
    t.FailNow()
}

b, err := c.NewBatchPointsWriter("s")
if err != nil {
    t.Log(err)
    t.FailNow()
}

b.AddPoints("products", tags, fields) // you can add more points later

b.Write() // finally write the points
```

## Important Notes
Please close the client when your application exits