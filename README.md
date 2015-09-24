Read / Write S3 buckets.  Authentication uses the same file $HOME/.aws/credentials that the AWS CLI tool uses.

Bucket contents are encrypted using AES256 server side encryption.

## Build

```
make
```

## Help

```
$ ./s3rw -h
Usage of ./s3rw:
  -bucket-name string
    	Bucket name
  -filename string
    	file to get or put
  -object-key string
    	Item key in bucket
  -op string
    	get or put
  -profile string
    	AWS profile from $HOME/.aws/credentials to use
  -region string
    	AWS region
  -version
    	Print version info and exit.
```
