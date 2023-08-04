# goapplambda


Example of how to deploy go-app application using s3, lambda, and cloudfront.

## setup steps

* Create lambda function with a function url enabled. (goapplambda in this example)
* Create a S3 bucket to hold the web resources. (mlctrez-goapplambda in this example)
* Create a cloudfront distribution with the default origin pointing to the lambda function url.
* Add an origin to the cloudfront distribution for the s3 bucket.
* Add a behavior matching /web/* to be routed to the s3 bucket origin.
* Create a cloudfront function applied to the /web/* behavior to set the wasm size header. (goapplambda-wasmsize in this example) 
* Optionally, add a custom CNAME and SSL certificate to be used instead of the default xxx.cloudfront.net name.

## demo
[example deployment](https://goapplambda.mlctrez.com/)


[![Go Report Card](https://goreportcard.com/badge/github.com/mlctrez/goapplambda)](https://goreportcard.com/report/github.com/mlctrez/goapplambda)

created by [tigwen](https://github.com/mlctrez/tigwen)
