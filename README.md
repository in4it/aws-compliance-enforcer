# AWS Compliance

## Deploy

* Build and copy the binary to s3
```
make # or download the latest binary from the github release page
aws s3api create-bucket --bucket lambda-scripts-eu-west-1-accountid --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1
aws s3 cp main.zip s3://compliance-enforcer-lambda-eu-west-1-accountid/compliance-enforcer
```

* Apply the cloudformation template
