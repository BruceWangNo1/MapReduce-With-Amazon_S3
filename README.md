## MapReduce implementation of distributed systems in Golang with Amazon S3
This project is inspired by the MIT course [6.824: Distributed Systems](https://pdos.csail.mit.edu/6.824/index.html) and most of the code is from the course lab assignment.

I intend to add features so that this could be deployed on a master and a certain number of slaves with a shared storage utilizing Amazon S3(Simple Storage Service).

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

```
/src/sophie
go test -run TestListBuckets
go test -run TestCreateBuckets
go test -run TestUploadFileStream
go test -run TestDownloadFile
```

Start three terminals and run the following three commands respectively (you should try to run the last two simultaneously becuase they finish really quickly):

```
/src/main
go run wc_s3.go master distributed pg
go run wc_s3.go worker localhost:7777 localhost:7778
go run wc_s3.go worker localhost:7777 localhost:7779
```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## License
The MIT License