#go_apps = bin/register

#bin/% : functions/%.go
#		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

#build: $(go_apps) | vendor

.PHONY : all
all : register signup signup signin refresh profile general deviceList deviceListMinimal deviceUpdate \
	deviceUpdateMeta cognitoRegister cognitoProfileList cognitoProfileUpdate deviceAdd deviceGet \
	deviceDel map zoneUpdate validateEmail seeder schemaGet chartLiveDevice chartCache

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

deviceListMinimal:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceListMinimal functions/deviceListMinimal/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceListMinimal.zip deviceListMinimal
	rm deviceListMinimal

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

deviceDel:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceDel functions/deviceDel/main.go
	mkdir -p bin
	build-lambda-zip -o bin/deviceDel.zip deviceDel
	rm deviceDel

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

validateEmail:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o validateEmail functions/validateEmail/main.go
	mkdir -p bin
	build-lambda-zip -o bin/validateEmail.zip validateEmail
	rm validateEmail

seeder:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o seeder functions/seeder/main.go
	mkdir -p bin
	build-lambda-zip -o bin/seeder.zip seeder
	rm seeder

schemaGet:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o schemaGet functions/schemaGet/main.go
	mkdir -p bin
	build-lambda-zip -o bin/schemaGet.zip schemaGet
	rm schemaGet

chartLiveDevice:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartLiveDevice functions/chartLiveDevice/main.go
	mkdir -p bin
	build-lambda-zip -o bin/chartLiveDevice.zip chartLiveDevice
	rm chartLiveDevice

chartCache:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartCache functions/chartCache/main.go
	mkdir -p bin
	build-lambda-zip -o bin/chartCache.zip chartCache
	rm chartCache

test:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test functions/test/main.go
	mkdir -p bin
	build-lambda-zip -o bin/test.zip test
	rm test