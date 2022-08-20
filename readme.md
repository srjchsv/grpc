### educational - gRPC client-server

- to run server use `make server`.
Logs are in the stdout.
### For now two rpc services made
1. CreateLaptop - creates a random laptop and saves it in servers in memory storage.
2. SearchLaptop - filters and searches for laptops with multiple conditions.
- `make client` runs client,executes CreateLaptop rpc- ten random laptops will be created and saved in the servers in memory storage. And then SearchLaptop executes- search request will be made, with some random input. Result shown in the stdout.

- run tests with `make test` 

Im learning gRPC from the TECH SCHOOL educational series. And be pushing code, while going throught the course.
TECH SCHOOL complete gRPC course link:
`https://dev.to/techschoolguru/series/7311`