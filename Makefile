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
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signup functions/signup/main.go
	mkdir -p bin
	build-lambda-zip -o bin/signup.zip signup
	rm signup

signin:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signin functions/signin/main.go
	mkdir -p bin
	build-lambda-zip -o bin/signin.zip signin
	rm signin

refresh:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o refresh functions/refresh/main.go
	mkdir -p bin
	build-lambda-zip -o bin/refresh.zip refresh
	rm refresh

profile:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o profile functions/profile/main.go
	mkdir -p bin
	build-lambda-zip -o bin/profile.zip profile
	rm profile

general:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o general functions/general/main.go
	mkdir -p bin
	build-lambda-zip -o bin/general.zip general
	rm general

deviceList:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o devicelist functions/deviceList/main.go
	mkdir -p bin
	build-lambda-zip -o bin/devicelist.zip devicelist
	rm devicelist

deviceUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceUpdate functions/deviceUpdate/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceUpdate.zip deviceUpdate
	rm deviceUpdate

deviceUpdateMeta:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceUpdateMeta functions/deviceUpdateMeta/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceUpdateMeta.zip deviceUpdateMeta
	rm deviceUpdateMeta

cognitoRegister:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoRegister functions/cognitoRegister/main.go
	mkdir -p bin
	build-lambda-zip -o bin/cognitoRegister.zip cognitoRegister
	rm cognitoRegister

cognitoProfileList:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoProfileList functions/cognitoProfileList/main.go
	mkdir -p bin
	build-lambda-zip -o bin/cognitoProfileList.zip cognitoProfileList
	rm cognitoProfileList

cognitoProfileUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoProfileUpdate functions/cognitoProfileUpdate/main.go
	mkdir -p bin
	build-lambda-zip -o bin/cognitoProfileUpdate.zip cognitoProfileUpdate
	rm cognitoProfileUpdate

deviceAdd:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceAdd functions/deviceAdd/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceAdd.zip deviceAdd
	rm deviceAdd

deviceGet:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceGet functions/deviceGet/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceGet.zip deviceGet
	rm deviceGet

map:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o map functions/map/main.go
	mkdir -p bin
	build-lambda-zip -o bin/map.zip map
	rm map

zoneUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o zoneUpdate functions/zoneUpdate/main.go
	mkdir -p bin
	build-lambda-zip -o bin/zoneUpdate.zip zoneUpdate
	rm zoneUpdate