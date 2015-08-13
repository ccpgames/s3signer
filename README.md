# s3signer

A simple web service to return S3 signed URLs.


## Usage

```bash
$ curl http://localhost:8080/mybucket/myfile/
https://s3-us-west-1.amazonaws.com/mybucket/myfile?AWSAccessKeyId=RANDOMLETTERS&Expires=1234567890&Signature=MORERANDOMLETTERS
```

The signed URLs returned are valid for `5` seconds.


## Configuration

Configuration is via environment variables. The [the goamz.aws.EnvAuth](https://godoc.org/github.com/mitchellh/goamz/aws#EnvAuth)
variables need to be provided, and s3signer uses one additional variable `AWS_REGION` to
specify the region by name. Something like `us-west-1`, ([as defined here](https://godoc.org/github.com/mitchellh/goamz/aws#pkg-variables)).


## Buckets, files

All buckets that the credentials have access to are cached on startup. If a new
bucket is created, the service will need to be restarted. Files inside of buckets
are searched for on each request.

Nested files in buckets are not supported at this time.
