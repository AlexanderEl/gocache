# GoCache - lighting fast caching service

A very simple and lightweight caching service that is capable of handling high loads.

Initially created for educational/fun purposes, but it has proven itself to be highly performant for use in heavy workloads.

## Functionality
- GET
- SET
- HAS
- DELETE

## Features
- **Concurrent Safety** - all operations are ACID compliant
- **Smart expiration** - cleans up expired entries

## Supported Data Formats
Currently only supports strings for the keys and values

# Performance Benchmark
### **1.276M** operations per second


### Data Values:
```
Total:  175,693,037
Gets:    30,002,476
Sets:    60,000,000
Has:     29,997,524
Deletes: 55,693,037

Duration: 137.606s
```

The setting of data was controlled while Get, Set, Delete were randomly chosen and executed.

This did require over 14GB of memory to hold the 1.5M elements.

With some fine tuning of the benchmark tests, I am positive this can be pushed further by optimizing the work per worker and reducing overhead introduced in the test.

### Hardware used for Performance
- OS: **Ubuntu 24.04.2 LTS**
- CPU: **i9-14900HX**
- Ram: **32GB**

## License
This package is licensed under the MIT License - see the LICENSE file for details.

## Reference
Project inspired by a video from [anthdm](https://github.com/anthdm)
