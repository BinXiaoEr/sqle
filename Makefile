################################## Parameter Definition And Check ##########################################
override GIT_VERSION    		= $(shell git rev-parse --abbrev-ref HEAD)${CUSTOM} $(shell git rev-parse HEAD)
override GIT_COMMIT     		= $(shell git rev-parse HEAD)
override PROJECT_NAME 			= sqle
override LDFLAGS 				= -ldflags "-X 'main.version=${GIT_VERSION}'"
override DOCKER         		= $(shell which docker)
override GOOS           		= linux
override OS_VERSION 			= el7
override GO_BUILD_FLAGS 		= -mod=vendor
override RPM_USER_GROUP_NAME 	= actiontech
override RPM_USER_NAME 			= actiontech-universe

GOARCH         		= amd64
RPMBUILD_TARGET		= x86_64

ifeq ($(GOARCH), arm64)
    RPMBUILD_TARGET = aarch64
endif

ifeq ($(IS_PRODUCTION_RELEASE),true)
# When performing a publishing operation, two cases:
# 1. if there is tag on current commit, means that
# 	 we release new version on current branch just now.
#    Set rpm name with tag name(v1.2109.0 -> 1.2109.0).
#
# 2. if there is no tag on current commit, means that
#    current branch is on process.
#    Set rpm name with current branch name(release-1.2109.x-ee or release-1.2109.x -> 1.2109.x).
    PROJECT_VERSION = $(shell if [ "$$(git tag --points-at HEAD | tail -n1)" ]; then git tag --points-at HEAD | tail -n1 | sed 's/v\(.*\)/\1/'; else git rev-parse --abbrev-ref HEAD | sed 's/release-\(.*\)/\1/' | tr '-' '\n' | head -n1; fi)
else
#    When performing daily packaging, set rpm name with current branch name(release-1.2109.x-ee or release-1.2109.x -> 1.2109.x).
    PROJECT_VERSION = $(shell git rev-parse --abbrev-ref HEAD | sed 's/release-\(.*\)/\1/' | tr '-' '\n' | head -n1)
endif

EDITION ?= ce
GO_BUILD_TAGS = dummyhead
ifeq ($(EDITION),ee)
    GO_BUILD_TAGS :=$(GO_BUILD_TAGS),enterprise
else ifeq ($(EDITION),trial)
    GO_BUILD_TAGS :=$(GO_BUILD_TAGS),trial
endif
RELEASE = qa
ifeq ($(RELEASE),rel)
    GO_BUILD_TAGS :=$(GO_BUILD_TAGS),release
endif

override RPM_NAME = $(PROJECT_NAME)-$(EDITION)-$(PROJECT_VERSION).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm
TARGET_USER ?=
ifdef TARGET_USER
	override RPM_NAME := $(PROJECT_NAME)-$(EDITION)-$(TARGET_USER)-$(PROJECT_VERSION).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm
endif

## The docker registry to pull compiler image, can be overwrite by: `make DOCKER_REGISTRY=10.0.0.1`
DOCKER_REGISTRY ?= 10.186.18.20

## Dynamic Parameter
GOLANGCI_LINT_IMAGE ?=golangci/golangci-lint:v1.45.2
SCSPELL_IMAGE ?=gerrywastaken/scspell
GO_COMPILER_IMAGE ?= golang:1.19.6
RPM_BUILD_IMAGE ?= rpmbuild/centos7

## Static Parameter, should not be overwrite
GOBIN = ${shell pwd}/bin
PARSER_PATH   = ${shell pwd}/vendor/github.com/pingcap/parser
LOCALE_PATH   = ${shell pwd}/sqle/locale
PLUGIN_LOCALE_PATH   = ${shell pwd}/sqle/driver/mysql/plocale

## Arm Build
ARM_CGO_BUILD_FLAG =
ifeq ($(EDITION)_$(GOARCH),ee_arm64)
    ARM_CGO_BUILD_FLAG = CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc
endif

default: install

######################################## i18n ##########################################################
install_i18n_tool:
	GOBIN=$(GOBIN) go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest

extract_i18n:
	cd ${LOCALE_PATH} && $(GOBIN)/goi18n extract -sourceLanguage zh
	cd ${PLUGIN_LOCALE_PATH} && $(GOBIN)/goi18n extract -sourceLanguage zh

start_trans_i18n:
	cd ${LOCALE_PATH} && touch translate.en.toml && $(GOBIN)/goi18n merge -sourceLanguage=zh active.*.toml
	cd ${PLUGIN_LOCALE_PATH} && touch translate.en.toml && $(GOBIN)/goi18n merge -sourceLanguage=zh active.*.toml

end_trans_i18n:
	cd ${LOCALE_PATH} && $(GOBIN)/goi18n merge active.en.toml translate.en.toml && rm -rf translate.en.toml
	cd ${PLUGIN_LOCALE_PATH} && $(GOBIN)/goi18n merge active.en.toml translate.en.toml && rm -rf translate.en.toml

######################################## Code Check ####################################################
## Static Code Analysis
vet: swagger
	GOOS=$(GOOS) GOARCH=amd64 go vet $$(GOOS=${GOOS} GOARCH=${GOARCH} go list ./...)

## Unit Test
test:
	cd $(PROJECT_NAME) && GOOS=$(GOOS) GOARCH=amd64 go test -v ./... -count 1

clean:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go clean

install: install_sqled install_scannerd

install_sqled: swagger
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) ${LDFLAGS} -tags $(GO_BUILD_TAGS) -o $(GOBIN)/sqled ./$(PROJECT_NAME)/cmd/sqled

install_scannerd:
	$(ARM_CGO_BUILD_FLAG) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) ${LDFLAGS} -tags $(GO_BUILD_TAGS) -o $(GOBIN)/scannerd ./$(PROJECT_NAME)/cmd/scannerd

dlv_install:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -gcflags "all=-N -l" $(GO_BUILD_FLAGS) ${LDFLAGS} -tags $(GO_BUILD_TAGS) -o $(GOBIN)/sqled ./$(PROJECT_NAME)/cmd/sqled
swagger:
	GOARCH=amd64 go build -o ${shell pwd}/bin/swag ${shell pwd}/build/swag/main.go
	rm -rf ${shell pwd}/sqle/docs
	${shell pwd}/bin/swag init -g ./$(PROJECT_NAME)/api/app.go -o ${shell pwd}/sqle/docs

parser:
	cd build/goyacc && GOOS=${GOOS} GOARCH=amd64 GOBIN=$(GOBIN) go install
	$(GOBIN)/goyacc -o /dev/null ${PARSER_PATH}/parser.y
	$(GOBIN)/goyacc -o ${PARSER_PATH}/parser.go ${PARSER_PATH}/parser.y 2>&1 | egrep "(shift|reduce)/reduce" | awk '{print} END {if (NR > 0) {print "Find conflict in parser.y. Please check y.output for more information."; exit 1;}}'
	rm -f y.output

	@if [ $(ARCH) = $(LINUX) ]; \
	then \
		sed -i -e 's|//line.*||' -e 's/yyEofCode/yyEOFCode/' ${PARSER_PATH}/parser.go; \
	elif [ $(ARCH) = $(MAC) ]; \
	then \
		/usr/bin/sed -i "" 's|//line.*||' ${PARSER_PATH}/parser.go; \
		/usr/bin/sed -i "" 's/yyEofCode/yyEOFCode/' ${PARSER_PATH}/parser.go; \
	fi

	@awk 'BEGIN{print "// Code generated by goyacc DO NOT EDIT."} {print $0}' ${PARSER_PATH}/parser.go > tmp_parser.go && mv tmp_parser.go ${PARSER_PATH}/parser.go;

upload:
	curl --ftp-create-dirs -T $(shell pwd)/$(RPM_NAME) ftp://$(RELEASE_FTPD_HOST)/actiontech-$(PROJECT_NAME)/$(EDITION)/$(RELEASE)/$(PROJECT_VERSION)/$(RPM_NAME)
	curl --ftp-create-dirs -T $(shell pwd)/$(RPM_NAME).md5 ftp://$(RELEASE_FTPD_HOST)/actiontech-$(PROJECT_NAME)/$(EDITION)/$(RELEASE)/$(PROJECT_VERSION)/$(RPM_NAME).md5

###################################### docker #####################################################
docker_lint:
	$(DOCKER) run -v $(shell pwd):/universe -w /universe --rm $(GOLANGCI_LINT_IMAGE) golangci-lint run -c ./.golangci.yml --timeout=20m

docker_scspell:
	$(DOCKER) run -v $(shell pwd):/universe -w /universe --rm $(SCSPELL_IMAGE) sh -c "python scspell.py sqle"


docker_test:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "cd /universe && make test ${MAKEFLAGS}"

docker_check: docker_lint docker_scspell docker_test

docker_clean:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "cd /universe && make clean ${MAKEFLAGS}"

# todo 升级golang版本后，git获取版本号失败，临时添加"git config --global --add safe.directory /universe"解决
docker_install:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "git config --global --add safe.directory /universe && cd /universe && make install $(MAKEFLAGS)"

docker_install_sqled:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "cd /universe && make install_sqled $(MAKEFLAGS)"

docker_install_scannerd:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "cd /universe && make install_scannerd $(MAKEFLAGS)"


docker_rpm: docker_install
	$(DOCKER) run -v $(shell pwd):/universe/sqle --user root --rm $(RPM_BUILD_IMAGE) sh -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; \
	(tar zcf ${PROJECT_NAME}.tar.gz /universe --transform 's/universe/${PROJECT_NAME}-$(GIT_COMMIT)/' >/tmp/build.log 2>&1) && \
	(rpmbuild --define 'group_name $(RPM_USER_GROUP_NAME)' --define 'user_name $(RPM_USER_NAME)' \
	--define 'commit $(GIT_COMMIT)' --define 'os_version $(OS_VERSION)' \
	--target $(RPMBUILD_TARGET)  -bb --with qa /universe/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && \
	(cat ~/rpmbuild/RPMS/$(RPMBUILD_TARGET)/${PROJECT_NAME}-$(GIT_COMMIT)-qa.$(OS_VERSION).$(RPMBUILD_TARGET).rpm) || (cat /tmp/build.log && exit 1)" > $(RPM_NAME) && \
	md5sum $(RPM_NAME) > $(RPM_NAME).md5

override SQLE_DOCKER_IMAGE ?= actiontech/$(PROJECT_NAME)-$(EDITION):$(PROJECT_VERSION)

docker_image:
	cp $(shell pwd)/$(RPM_NAME) $(shell pwd)/sqle.rpm
	$(DOCKER) build  -t $(SQLE_DOCKER_IMAGE) -f ./docker-images/sqle/Dockerfile .

docker_start:
	cd ./docker-images/sqle && SQLE_IMAGE=$(SQLE_DOCKER_IMAGE) docker-compose up -d

docker_stop:
	cd ./docker-images/sqle && docker-compose down

docker_rpm_with_dms: docker_install 
	$(DOCKER) run -v  $(dir $(CURDIR))dms:/universe/dms -v $(shell pwd):/universe/sqle --user root --rm $(RPM_BUILD_IMAGE) sh -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; \
	(tar zcf ${PROJECT_NAME}.tar.gz /universe/sqle /universe/dms --transform 's/universe/${PROJECT_NAME}-$(GIT_COMMIT)/' > /tmp/build.log 2>&1) && \
	(rpmbuild --define 'group_name $(RPM_USER_GROUP_NAME)' --define 'user_name $(RPM_USER_NAME)' \
	--define 'commit $(GIT_COMMIT)' --define 'os_version $(OS_VERSION)' \
	--define 'edition $(EDITION)' \
	--target $(RPMBUILD_TARGET)  -bb --with qa /universe/sqle/build/sqled_with_dms.spec >> /tmp/build.log 2>&1) && \
	(cat ~/rpmbuild/RPMS/$(RPMBUILD_TARGET)/${PROJECT_NAME}-$(GIT_COMMIT)-qa.$(OS_VERSION).$(RPMBUILD_TARGET).rpm) || (cat /tmp/build.log && exit 1)" > $(RPM_NAME) && \
	md5sum $(RPM_NAME) > $(RPM_NAME).md5


###################################### ui #####################################################
fill_ui_dir:
	# fill ui dir, it is used by rpm build.
	mkdir -p ./ui/static

.PHONY: help
help:
	$(warning ---------------------------------------------------------------------------------)
	$(warning Supported Variables And Values:)
	$(warning ---------------------------------------------------------------------------------)
	$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
# 需要下载modvendor，下载命令:go get -u github.com/goware/modvendor
go_mod_vendor:
	go mod vendor
	modvendor -copy="**/*.c **/*.h" -v	
