#go_apps = bin/register

#bin/% : functions/%.go
#		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

#build: $(go_apps) | vendor

vendor: Gopkg.toml
		dep ensure

register:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o register functions/register/main.go
	mkdir -p bin
	build-lambda-zip -o bin/register.zip register
	rm register

general:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o general functions/general/main.go
	mkdir -p bin
	build-lambda-zip -o bin/general.zip general
	rm general