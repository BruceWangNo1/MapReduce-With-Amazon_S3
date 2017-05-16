Sophie - Implementation of A Distributed MapReduce library
============================

This project is inspired by the MIT course [6.824: Distributed Systems](https://pdos.csail.mit.edu/6.824/index.html) and a large portion of the code based on lab assignment.

This library could be used to deploy MapReduce tasks on a cluster with the help of AWS including Amazon EC2 (one instance as a master), Amazon S3 (shared storage), Amazon ECS (one cluster of workers), and Amazon ECR (for storing our worker docker image). Before you begin your trial, you should have a good understanding of all these things above.
This library is fault-tolerant to worker failures and can scale to accommodate increasing needs with little overhead cost.


## Project Logo: Sophie
Courtesy of my enthusiastic roommate Ligu.
![Sophie_logo Designed by Ligu](src/panel/public/static/sophie_logo_1_800px.png)


## Environment
The master should be deployed on an Amazon EC2 instance.
The workers should be deployed on an Amazon ECS cluster of docker instances from a worker docker image uploaded.
We show you how to use this library and run this project by the classic word count example.


## Usage Example
1. Set up Amazon Credentials in your terminal and sophie/s3.go.
2. Master: Run command like this in your EC2 instance.
	```
	go run src/main/wc_s3.go master 192.168.1.1:7777 random 3 3
	```
3. Worker: Use ./Dockerfile to build the worker docker image and pushed to the ECS repository. Please change the corresponding parts to accommodate your specification. Launch a cluste of docker instances and specify your task definition and run your task.
4. Now your task is running.                                                     


## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

If you have any questions about this project, feel free to contact me.

## License
The MIT License