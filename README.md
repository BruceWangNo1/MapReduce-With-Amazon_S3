Sophie - MapReduce implementation of distributed systems in Golang with Amazon S3
============================

This project is inspired by the MIT course [6.824: Distributed Systems](https://pdos.csail.mit.edu/6.824/index.html) and most of the code is from the course lab assignment.

I intend to add features so that this could be deployed on a master and a certain number of slaves with a shared storage utilizing Amazon S3 (Simple Storage Service).

## Project Logo: Sophie
Courtesy of my enthusiastic roommate Ligu.
![Sophie_logo Designed by Ligu](src/panel/public/static/sophie_logo_1_800px.png)

## Environment
At first I developed and deployed this project solely on my laptop. At Usage 3, I spinned up an AWS S3 service as a shared storage for the program and thus S3 was involved. From Usage 4, I started to sync my code to AWS EC2 service to gain better response time with S3.

## Usage

1. Original Lab Test
	
	```
	/src/mapreduce
	go test -run SequentialMany
	go test -run Sequential
	go test -run TestBasic
	go test -run Failure
	```
2. Original Work Count Test

	```
	/src/main
	bash ./test-mr.sh
	```
3. Sophie AWS S3 Service Test

	```
	/src/sophie
	go test -run TestListBuckets
	go test -run TestCreateBuckets
	go test -run TestUploadFileStream
	go test -run TestDownloadFile
	```
	
4. Start three terminals and run the following three commands respectively 	(you should try to run the last two simultaneously becuase they finish 	really quickly. Notice that you need to comment out the panel service in 	wc_s3.go.)

	```
	/src/main
	go run wc_s3.go master localhost:7777 pg
	go run wc_s3.go worker localhost:7777 localhost:7778
	go run wc_s3.go worker localhost:7777 localhost:7779
	```
5. Web UI in progress. Go to localhost:8000 to view web UI

	```
	$GOPATH
	./test.sh
	```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

If you have any questions about this project, feel free to contact me.

## License
The MIT License