#go_apps = bin/register

#bin/% : functions/%.go
#		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

#build: $(go_apps) | vendor

vendor: Gopkg.toml
		dep ensure -v

register:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o register functions/register/main.go
	mkdir -p bin
	build-lambda-zip -o bin/register.zip register
	rm register

signup:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signup functions/signup/*.go
	mkdir -p bin
	build-lambda-zip -o bin/signup.zip signup
	rm signup

signin:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signin functions/signin/*.go
	mkdir -p bin
	build-lambda-zip -o bin/signin.zip signin
	rm signin

profile:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o profile functions/profile/*.go
	mkdir -p bin
	build-lambda-zip -o bin/profile.zip profile
	rm profile

general:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o general functions/general/main.go
	mkdir -p bin
	build-lambda-zip -o bin/general.zip general
	rm general