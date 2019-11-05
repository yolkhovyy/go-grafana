# backend-code-exercise

A timeseries is a series of data points in ascending time order,
where "point" is a combination of a numeric value and unix timestamp.

Whereas our users send us "messy" timeseries data (duplicate timestamps, time between points not consistent, etc),
our query responses should be "clean": no duplicate timestamps, even spacing between timestamps, etc.

So we need a function that converts messy timeseries data to clean timeseries data. It should also only return data in the range queried for.
It looks like:

```
func Fix(in []Point, from, to, interval uint32) []Point
```

* The input and output timeseries data is implemented as a slice of points, where each point is an instance of the included Point struct.
* The interval is a number of seconds such as 10, 60, etc.
* The from/to are unix timestamps.

The input data has the following properties:
1) The points are always in increasing timestamp order.
2) The spacing between the timestamps is not always consistent.
3) There might be gaps (missing points).
4) It may include points that are outside the time range we want to select.

The returned data must have these properties:
1) Points must be equally spaced by the provided interval (in seconds).
2) Points must have timestamps that are multiples of the interval.
   Meaning for interval=60, timestamps would be 1000000020, 1000000080, 1000000140, etc.
   If the input timestamp is not a multiple of the interval, it should be adjusted to the next one that is.
3) If multiple points correspond to the same timestamp, the first one wins.
4) Points must have a timestamp >= from and < end.
5) Points without values should use the value `math.NaN()`.

See the included code for scaffolding code and unit tests.
