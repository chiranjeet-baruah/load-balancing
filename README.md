# Consistent Hashing in Go

This repository contains an implementation of consistent hashing in Go. Consistent hashing is a strategy used in distributed systems to evenly distribute data across an arbitrary number of servers or nodes. This implementation includes a simulation of a consistent hashing ring with nodes being added and removed dynamically.

## Files in this Repository

1. `consistent.go`: This file contains the core implementation of the consistent hashing ring. It includes functions for adding and removing nodes from the ring, as well as a function to get the node that should handle a given request.

2. `consistent_test.go`: This file contains a test that simulates the operation of the consistent hashing ring. It includes three goroutines that add nodes, remove nodes, and request nodes from the ring, respectively.

## How to Run

To clone this repository and run the tests, follow these steps:

1. Clone the repository
2. Navigate to the repository directory
3. Run the tests with the following command:

```bash
go test -v ./...
```

## References and Acknowledgements
1. https://github.com/sent-hil/consistenthash
2. https://github.com/stathat/consistent