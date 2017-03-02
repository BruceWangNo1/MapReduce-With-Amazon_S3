## MapReduce implementation of distributed systems in Golang with Amazon S3
This project is inspired by the MIT course [6.824: Distributed Systems](https://pdos.csail.mit.edu/6.824/index.html) and most of the code is from the course lab assignment.

## Usage
```
/src/mapreduce
go test -run SequentialMany
go test -run Sequential
go test -run TestBasic
go test -run Failure
```

```
/src/main
bash ./test-mr.sh
```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## License
The MIT License