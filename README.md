# go-ibc
IP Blocklist Checker (IBC) written in Go

### Prerequisites

* Go (>= 1.12)

### Installing and Updating

```
go get -u github.com/fancytools/go-ibc
```

## Using the tool

```
go-ibc 127.0.0.1
```
```
go-ibc ::1
```
## CPU Profile and Benchmark

The go test tool for benchmark e.g. parsing of input, especially ipv4 and ipv6 addresses.
```
go test -timeout 20s -bench=BenchmarkParsing -cpuprofile=cpu.txt
```

The tool pprof shows the profiling results on a local interactive web page.
```
go tool pprof --http 127.0.0.1:8080 cpu.txt
```

If you hadn't enough, you can use go compiler decisions to optimize the software extremely.
```
go build -gcflags="-m -m"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* The tool is currently under development.
