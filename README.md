## Prerequisites
> Install [GoLang](https://golang.org/) v1.10+  
> Install [dep](#dep-installation) - go dependency management tool
> Build [lambda-build-zip](#lambda-build-zip) - used to create zip archives

### Dep installation
Install [dep](https://github.com/golang/dep)
```sh
    export DEP_RELEASE_TAG=v0.5.0 
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

### lambda-build-zip
Install [lambda-build-zip](https://github.com/aws/aws-lambda-go)

#### lambda-build-zip (possible) issues
After installation lambda-build-zip should be accessible from command promopt. 
If that's not the case, you might want to check if lambda-build-zip is in `$PATH` (`echo $PATH`).
As a last resort, copy `lambda-build-zip` from `$GOPATH/bin` to `/usr/local/bin`.

#### AWS CLI installation
AWS SAM has issues with AWS CLI versions >1.16.145~
```sh
virtualenv -p python3.7 cli_env
source cli_env/bin/activate
pip install awscli==1.16.141
pip install --user aws-sam-cli
```

### Makefile
`make vendor`        - get dependencies with `dep`
`make LAMBDA_NAME` - build LAMBDA_NAME lambda into the `/bin` dir. Example:

```sh
    make register
```

## Run local
How to debug on your local machine?

Go to project root, and type:
```sh
    ENV=local go run ./functions/[SPECIFIC_FUNCTION]/*.go
```
`SPECIFIC_FUNCTION` is one of the packages in `functions` directory.

Example:
```sh
    ENV=local go run ./functions/signup/*.go 
```

## Environment (.env) file
Place `.env` file into the project root dir.

Required environment variables:
```sh
    COGNITO_REGION=us-east-1
    COGNITO_USER_POOL_ID=us-east-xxxxxxxxxxx
    COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxx
```

If you are running app in local environment, you need:
```sh
    AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxxx
    AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Project postman documentation
[![Go to Postman API Documentation](https://run.pstmn.io/button.svg)](https://documenter.getpostman.com/view/5704418/SzRuYBqS?version=latest)

## Built With

- [go](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.

## Install aws cli