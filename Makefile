include ./vendor/actiontech.cloud/universe/ucommon/v4/build/Makefile.variables

## Makefile Content
# 1.Parameter Definition And Check
# 2.Code Check
# 3.Stripped Dependence
# 4.Golang Binary Compile
# 5.RPM Build/Upload
# 6.K8s Docker Images Build
# 7.Frontend(Optional)


################################## 1.Parameter Definition And Check ##########################################
## Dynamic Parameter
# The ftp host and user/pass used to upload rpm package, can be overwrite by: `make DOCKER_IMAGE=image:tag`
DOCKER_IMAGE  ?= docker-registry:5000/actiontech/universe-compiler-go1.14.1-centos6
DOTNET_DOCKER_IMAGE ?= docker-registry:5000/actiontech/universe-compiler-dotnetcore2.1
TEST_DOCKER_IMAGE  ?= docker-registry:5000/actiontech/universe-compiler-go1.14.1-ubuntu-with-docker
K8S_DOCKER_IMAGE_BUILD ?= docker-registry:5000/actiontech/universe-compiler-go1.14.1-ubuntu-with-docker

## Static Parameter, should not be overwrite
CAP = CAP_CHOWN,CAP_SYS_RESOURCE,CAP_SETUID,CAP_SETGID+eip
RUN_ON_DMP = true
PROJECT_NAME = sqle
SUB_PROJECT_NAME = sqle_sqlserver
override VERSION=4.20.11.0
GOBIN = ${shell pwd}/bin
K8S_DOCKER_IMAGE_GENERATED = docker-registry:5000/actiontech/k8s/$(PROJECT_NAME):v$(VERSION)
default: install
DOTNET_TARGET = centos.7-x64
PARSER_PATH   = ${shell pwd}/vendor/github.com/pingcap/parser
SQLE_LDFLAGS   = ${LDFLAGS}" -X 'main.runOnDmpStr=${RUN_ON_DMP}'"

######################################## 2.Code Check ####################################################
## Static Code Analysis
vet: swagger
	GOOS=$(GOOS) GOARCH=amd64 go vet $$(GOOS=${GOOS} GOARCH=${GOARCH} go list ./...)
	GOOS=$(GOOS) GOARCH=amd64 go vet ./vendor/actiontech.cloud/...

## Unit Test
test: swagger parser
	GOOS=$(GOOS) GOARCH=amd64 go test -v ./$(PROJECT_NAME)


docker_test: pull_image
	CTN_NAME="universe_docker_test_$$RANDOM" && \
    $(DOCKER) run -d --entrypoint /sbin/init --add-host docker-registry:${DOCKER_REGISTRY}  --privileged --name $${CTN_NAME} -v $(shell pwd):/universe/sqle --rm -w /universe/sqle $(DOCKER_IMAGE) && \
    $(DOCKER) exec $${CTN_NAME} make test vet ; \
    $(DOCKER) stop $${CTN_NAME}


#################################### 3.Stripped Dependence ##############################################
# All stripped dependence should upload to ftp before doing rpm packing.
# FTP Directory structure shoule be:
# aarch64 RHEL7 -- ${RELEASE_FTPD_HOST}/deploy-stripped/linux_aarch64/el7/
# aarch64 RHEL8 -- ${RELEASE_FTPD_HOST}/deploy-stripped/linux_aarch64/el8/
# arm64 RHEL7 -- ${RELEASE_FTPD_HOST}/deploy-stripped/linux_amd64/el7/
# arm64 RHEL8 -- ${RELEASE_FTPD_HOST}/deploy-stripped/linux_amd64/el8/

# sqle had no stripped dependence
################################### 4.Golang Binary Compile #############################################
# For go mod
sync_vendor:
	GOPROXY=$(GOPROXY) GONOSUMDB=$(GONOSUMDB) go mod vendor

# Generic
docker_clean:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(DOCKER_IMAGE) -c "cd /universe && make clean ${MAKEFLAGS}"
docker_install:
	$(DOCKER) run -v $(shell pwd):/universe --rm $(DOCKER_IMAGE) -c "cd /universe && make install $(MAKEFLAGS)"
install: swagger parser
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(SQLE_LDFLAGS) $(GO_BUILD_FLAGS) -tags $(GO_BUILD_TAGS) -o $(GOBIN)/sqled ./$(PROJECT_NAME)

build_sqlserver:
	cd ./sqle/sqlserver/SqlserverProtoServer && dotnet publish -c Release -r ${DOTNET_TARGET}

swagger:
	GOARCH=amd64 go build -o ${shell pwd}/$(PROJECT_NAME)/swag ${shell pwd}/build/swag/main.go
	rm -rf ${shell pwd}/sqle/docs
	${shell pwd}/$(PROJECT_NAME)/swag init -g ./$(PROJECT_NAME)/api/app.go -o ${shell pwd}/sqle/docs

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

# Clean
clean:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go clean

###################################### 5.RPM Build #####################################################
# Compiler docker image
pull_image:
	$(DOCKER) pull $(DOCKER_IMAGE)

docker_rpm: docker_rpm/sqle docker_rpm/sqle_sqlserver

docker_rpm/sqle: pull_image docker_install
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; \
	(tar zcf ${PROJECT_NAME}.tar.gz /universe --transform 's/universe/${PROJECT_NAME}-${VERSION}_$(GIT_COMMIT)/' >/tmp/build.log 2>&1) && \
	(rpmbuild --define 'caps ${CAP}' --define 'group_name $(RPM_USER_GROUP_NAME)' --define 'user_name $(RPM_USER_NAME)' --define 'runOnDmp ${RUN_ON_DMP}' \
	--define 'commit $(GIT_COMMIT)' --define 'os_version $(OS_VERSION)' \
	--target $(RPMBUILD_TARGET)  -bb --with qa /universe/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && \
	(cat /root/rpmbuild/RPMS/$(RPMBUILD_TARGET)/${PROJECT_NAME}-${VERSION}_$(GIT_COMMIT)-qa.$(OS_VERSION).$(RPMBUILD_TARGET).rpm) || (cat /tmp/build.log && exit 1)" > $(PROJECT_NAME).$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm

docker_rpm/sqle_sqlserver: pull_image docker_install
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOTNET_DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; \
	(tar zcf ${SUB_PROJECT_NAME}.tar.gz /universe --transform 's/universe/${SUB_PROJECT_NAME}-${VERSION}_$(GIT_COMMIT)/' >/tmp/build.log 2>&1) && \
	(rpmbuild --define '_dotnet_target ${DOTNET_TARGET}' --define '_git_version ${GIT_VERSION}' --define 'group_name $(RPM_USER_GROUP_NAME)' --define 'user_name $(RPM_USER_NAME)' \
	--define 'commit $(GIT_COMMIT)' --define 'os_version $(OS_VERSION)' \
	--target $(RPMBUILD_TARGET) -bb --with qa /universe/sqle/build/sqled_sqlserver.spec >>/tmp/build.log 2>&1) && \
	(cat /root/rpmbuild/RPMS/$(RPMBUILD_TARGET)/${SUB_PROJECT_NAME}-${VERSION}_$(GIT_COMMIT)-qa.$(OS_VERSION).$(RPMBUILD_TARGET).rpm) || (cat /tmp/build.log && exit 1)" > ${SUB_PROJECT_NAME}.$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm

upload:
	curl -T $(shell pwd)/$(PROJECT_NAME).$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm \
	ftp://$(RELEASE_FTPD_HOST)/actiontech-$(PROJECT_NAME)/qa/$(VERSION)/$(PROJECT_NAME)-$(VERSION).$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm --ftp-create-dirs
	curl -T $(shell pwd)/$(SUB_PROJECT_NAME).$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm \
	ftp://$(RELEASE_FTPD_HOST)/actiontech-$(PROJECT_NAME)/qa/$(VERSION)/$(SUB_PROJECT_NAME)-$(VERSION).$(CUSTOMER).$(RELEASE).$(OS_VERSION).$(RPMBUILD_TARGET).rpm  --ftp-create-dirs
############################### 6.K8s Docker Images Build ##############################################

# sqle had no supprot K8s docker images build

.PHONY: help
help:
	$(warning ---------------------------------------------------------------------------------)
	$(warning Supported Variables And Values:)
	$(warning ---------------------------------------------------------------------------------)
	$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
