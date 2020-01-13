#go_apps = bin/register

#bin/% : functions/%.go
#		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

#build: $(go_apps) | vendor

PACKAGED_TEMPLATE = packaged.yaml
S3_BUCKET = artifacts.cityo.io
STACK_NAME = CityOS
TEMPLATE = template.yaml

ifdef OS
    package_lambda = build-lambda-zip -o
    FixPath = $(subst /,\,$1)
else
    ifeq ($(shell uname), Darwin)
      package_lambda = zip -p
      FixPath = $1
	endif
	ifeq ($(shell uname), Linux)
	  package_lambda = zip -p
      FixPath = $1
   endif
endif

.PHONY : all
all : register signup signup signin refresh profile general deviceList deviceListMinimal deviceUpdate \
	deviceUpdateMeta cognitoRegister cognitoProfileList cognitoProfileUpdate deviceAdd deviceGet \
	deviceDel map zoneUpdate validateEmail seeder schemaGet chartLiveDevice chartCache chartHour \
	chartSave chartHourDevice chartHourAll chartCacheDay chartDay chartDayDevice chartDayAll chartHasData \
	chartCacheSix chartSix chartSixDevice chartSixAll chartLiveAll forgotPasswordStart forgotPasswordEnd notifications \
	cityList cityDel cityAdd verifyAndRedirect customizeMessage

vendor: Gopkg.toml
		dep ensure -v

register:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o register functions/register/main.go
	mkdir -p bin
	$(package_lambda) bin/register.zip register
	rm register

signup:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signup functions/signup/main.go
	mkdir -p bin
	$(package_lambda) bin/signup.zip signup
	rm signup

signin:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o signin functions/signin/main.go
	mkdir -p bin
	$(package_lambda) bin/signin.zip signin
	rm signin

refresh:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o refresh functions/refresh/main.go
	mkdir -p bin
	$(package_lambda) bin/refresh.zip refresh
	rm refresh

profile:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o profile functions/profile/main.go
	mkdir -p bin
	$(package_lambda) bin/profile.zip profile
	rm profile

general:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o general functions/general/main.go
	mkdir -p bin
	$(package_lambda) bin/general.zip general
	rm general

deviceList:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceList functions/deviceList/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceList.zip deviceList
	rm deviceList

deviceListMinimal:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceListMinimal functions/deviceListMinimal/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceListMinimal.zip deviceListMinimal
	rm deviceListMinimal

deviceUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceUpdate functions/deviceUpdate/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceUpdate.zip deviceUpdate
	rm deviceUpdate

deviceUpdateMeta:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceUpdateMeta functions/deviceUpdateMeta/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceUpdateMeta.zip deviceUpdateMeta
	rm deviceUpdateMeta

cognitoRegister:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoRegister functions/cognitoRegister/main.go
	mkdir -p bin
	$(package_lambda) bin/cognitoRegister.zip cognitoRegister
	rm cognitoRegister

cognitoProfileList:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoProfileList functions/cognitoProfileList/main.go
	mkdir -p bin
	$(package_lambda) bin/cognitoProfileList.zip cognitoProfileList
	rm cognitoProfileList

cognitoProfileUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cognitoProfileUpdate functions/cognitoProfileUpdate/main.go
	mkdir -p bin
	$(package_lambda) bin/cognitoProfileUpdate.zip cognitoProfileUpdate
	rm cognitoProfileUpdate

deviceAdd:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceAdd functions/deviceAdd/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceAdd.zip deviceAdd
	rm deviceAdd

deviceGet:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceGet functions/deviceGet/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceGet.zip deviceGet
	rm deviceGet

deviceDel:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o deviceDel functions/deviceDel/main.go
	mkdir -p bin
	$(package_lambda) bin/deviceDel.zip deviceDel
	rm deviceDel

map:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o map functions/map/main.go
	mkdir -p bin
	$(package_lambda) bin/map.zip map
	rm map

zoneUpdate:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o zoneUpdate functions/zoneUpdate/main.go
	mkdir -p bin
	$(package_lambda) bin/zoneUpdate.zip zoneUpdate
	rm zoneUpdate

validateEmail:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o validateEmail functions/validateEmail/main.go
	mkdir -p bin
	$(package_lambda) bin/validateEmail.zip validateEmail
	rm validateEmail

seeder:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o seeder functions/seeder/main.go
	mkdir -p bin
	$(package_lambda) bin/seeder.zip seeder
	rm seeder

schemaGet:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o schemaGet functions/schemaGet/main.go
	mkdir -p bin
	$(package_lambda) bin/schemaGet.zip schemaGet
	rm schemaGet

chartLiveDevice:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartLiveDevice functions/chartLiveDevice/main.go
	mkdir -p bin
	$(package_lambda) bin/chartLiveDevice.zip chartLiveDevice
	rm chartLiveDevice

chartLiveAll:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartLiveAll functions/chartLiveAll/main.go
	mkdir -p bin
	$(package_lambda) bin/chartLiveAll.zip chartLiveAll
	rm chartLiveAll

chartHourDevice:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartHourDevice functions/chartHourDevice/main.go
	mkdir -p bin
	$(package_lambda) bin/chartHourDevice.zip chartHourDevice
	rm chartHourDevice

chartHourAll:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartHourAll functions/chartHourAll/main.go
	mkdir -p bin
	$(package_lambda) bin/chartHourAll.zip chartHourAll
	rm chartHourAll

chartCache:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartCache functions/chartCache/main.go
	mkdir -p bin
	$(package_lambda) bin/chartCache.zip chartCache
	rm chartCache

chartCacheDay:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartCacheDay functions/chartCacheDay/main.go
	mkdir -p bin
	$(package_lambda) bin/chartCacheDay.zip chartCacheDay
	rm chartCacheDay

chartCacheSix:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartCacheSix functions/chartCacheSix/main.go
	mkdir -p bin
	$(package_lambda) bin/chartCacheSix.zip chartCacheSix
	rm chartCacheSix

chartSix:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartSix functions/chartSix/main.go
	mkdir -p bin
	$(package_lambda) bin/chartSix.zip chartSix
	rm chartSix

chartHour:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartHour functions/chartHour/main.go
	mkdir -p bin
	$(package_lambda) bin/chartHour.zip chartHour
	rm chartHour

chartDay:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartDay functions/chartDay/main.go
	mkdir -p bin
	$(package_lambda) bin/chartDay.zip chartDay
	rm chartDay

chartDayDevice:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartDayDevice functions/chartDayDevice/main.go
	mkdir -p bin
	$(package_lambda) bin/chartDayDevice.zip chartDayDevice
	rm chartDayDevice

chartHasData:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartHasData functions/chartHasData/main.go
	mkdir -p bin
	$(package_lambda) bin/chartHasData.zip chartHasData
	rm chartHasData

chartDayAll:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartDayAll functions/chartDayAll/main.go
	mkdir -p bin
	$(package_lambda) bin/chartDayAll.zip chartDayAll
	rm chartDayAll

chartSixDevice:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartSixDevice functions/chartSixDevice/main.go
	mkdir -p bin
	$(package_lambda) bin/chartSixDevice.zip chartSixDevice
	rm chartSixDevice

chartSixAll:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartSixAll functions/chartSixAll/main.go
	mkdir -p bin
	$(package_lambda) bin/chartSixAll.zip chartSixAll
	rm chartSixAll

chartSave:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chartSave functions/chartSave/main.go
	mkdir -p bin
	$(package_lambda) bin/chartSave.zip chartSave
	rm chartSave

test:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o test functions/test/main.go
	mkdir -p bin
	$(package_lambda) bin/test.zip test
	rm test

notifications:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o notifications functions/notifications/main.go
	mkdir -p bin
	$(package_lambda) bin/notifications.zip notifications
	rm notifications

forgotPasswordStart:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o forgotPasswordStart functions/forgotPasswordStart/main.go
	mkdir -p bin
	$(package_lambda) bin/forgotPasswordStart.zip forgotPasswordStart
	rm forgotPasswordStart

forgotPasswordEnd:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o forgotPasswordEnd functions/forgotPasswordEnd/main.go
	mkdir -p bin
	$(package_lambda) bin/forgotPasswordEnd.zip forgotPasswordEnd
	rm forgotPasswordEnd


verifyAndRedirect:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o verifyAndRedirect functions/verifyAndRedirect/main.go
	mkdir -p bin
	$(package_lambda) bin/verifyAndRedirect.zip verifyAndRedirect
	rm verifyAndRedirect

customizeMessage:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o customizeMessage functions/customizeMessage/main.go
	mkdir -p bin
	$(package_lambda) bin/customizeMessage.zip customizeMessage
	rm customizeMessage

cityList:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cityList functions/cityList/main.go
	mkdir -p bin
	$(package_lambda) bin/cityList.zip cityList
	rm cityList

cityDel:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cityDel functions/cityDel/main.go
	mkdir -p bin
	$(package_lambda) bin/cityDel.zip cityDel
	rm cityDel

cityAdd:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cityAdd functions/cityAdd/main.go
	mkdir -p bin
	$(package_lambda) bin/cityAdd.zip cityAdd
	rm cityAdd

# cityUpdate:
# 	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cityUpdate functions/cityUpdate/main.go
# 	mkdir -p bin
# 	$(package_lambda) bin/cityUpdate.zip cityUpdate
# 	rm cityUpdate	

.PHONY: deploy_swagger
deploy_swagger:
	aws s3 cp swagger.yaml s3://artifacts.cityo.io/CityOS/swagger.yaml

#.PHONY: package
package: all
	sam package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --s3-prefix $(STACK_NAME) --output-template-file $(PACKAGED_TEMPLATE)

.PHONY: deploy
deploy: deploy_swagger package
	sam deploy --stack-name $(STACK_NAME) --template-file $(PACKAGED_TEMPLATE) --capabilities CAPABILITY_NAMED_IAM
