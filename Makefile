include ./vendor/actiontech.cloud/universe/ucommon/v3/build/Makefile.variables
GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOLIST=$(GOCMD) list
GOCLEAN=$(GOCMD) clean

GIT_VERSION   = $(shell git rev-parse --abbrev-ref HEAD) $(shell git rev-parse HEAD)
RPM_BUILD_BIN = $(shell type -p rpmbuild 2>/dev/null)
COMPILE_FLAG  =
DOCKER        = $(shell which docker)
DOCKER_IMAGE  = docker-registry:5000/actiontech/universe-compiler-go1.14.1-centos6
DOTNET_DOCKER_IMAGE = docker-registry:5000/actiontech/universe-compiler-dotnetcore2.1
DOTNET_TARGET = centos.7-x64
SQLE_LDFLAGS = -ldflags "-X 'main.version=\"${GIT_VERSION}\"' -X 'main.caps=${CAP}' -X 'main.defaultUser=${USER_NAME}' -X 'main.runOnDmpStr=${RUN_ON_DMP}'"

PROJECT_NAME = sqle
SUB_PROJECT_NAME = sqle_sqlserver
VERSION       = 9.9.9.9
CAP = CAP_CHOWN,CAP_SYS_RESOURCE,CAP_SETUID,CAP_SETGID+eip

PARSER_PATH   = ${shell pwd}/vendor/github.com/pingcap/parser
MAIN_MODULE   = ${shell pwd}
GOBIN         = ${MAIN_MODULE}/bin

.PHONY: build docs

default: build

pull_image:
    $(DOCKER) pull ${DOCKER_IMAGE}

install: swagger parser vet
	GOBIN=${GOBIN} GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${GOBIN}/sqled -mod=vendor ${SQLE_LDFLAGS} ${MAIN_MODULE}/${PROJECT_NAME}

build_sqlserver:
	cd ./sqle/sqlserver/SqlserverProtoServer && dotnet publish -c Release -r ${DOTNET_TARGET}

vet: swagger
	$(GOVET) $$($(GOLIST) ./... | grep -v vendor/)

test: swagger parser
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)

docker_rpm: pull_image
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${PROJECT_NAME}.tar.gz /universe --transform 's/universe/${PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild --define 'caps ${CAP}' --define 'runOnDmp true' -bb --with qa /universe/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${PROJECT_NAME}.x86_64.rpm
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOTNET_DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${SUB_PROJECT_NAME}.tar.gz /universe --transform 's/universe/${SUB_PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild --define '_dotnet_target ${DOTNET_TARGET}' --define '_git_version ${GIT_VERSION}' -bb --with qa /universe/sqle/build/sqled_sqlserver.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${SUB_PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${SUB_PROJECT_NAME}.x86_64.rpm

docker_rpm_without_dmp: pull_image
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${PROJECT_NAME}.tar.gz /universe --transform 's/universe/${PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild --define 'caps ${CAP}' --define 'runOnDmp false' -bb --with qa /universe/sqle/build/sqled.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${PROJECT_NAME}.x86_64.rpm
	$(DOCKER) run -v $(shell pwd):/universe/sqle --rm $(DOTNET_DOCKER_IMAGE) -c "(mkdir -p /root/rpmbuild/SOURCES >/dev/null 2>&1);cd /root/rpmbuild/SOURCES; (tar zcf ${SUB_PROJECT_NAME}.tar.gz /universe --transform 's/universe/${SUB_PROJECT_NAME}-${VERSION}/' >/tmp/build.log 2>&1) && (rpmbuild --define '_dotnet_target ${DOTNET_TARGET}' --define '_git_version ${GIT_VERSION}' -bb --with qa /universe/sqle/build/sqled_sqlserver.spec >>/tmp/build.log 2>&1) && (cat /root/rpmbuild/RPMS/x86_64/${SUB_PROJECT_NAME}-${VERSION}-qa.x86_64.rpm) || (cat /tmp/build.log && exit 1)" > ${SUB_PROJECT_NAME}.x86_64.rpm

docker_test: pull_image
	CTN_NAME="universe_docker_test_$$RANDOM" && \
    $(DOCKER) run -d --entrypoint /sbin/init --add-host docker-registry:${DOCKER_REGISTRY}  --privileged --name $${CTN_NAME} -v $(shell pwd):/universe/sqle --rm -w /universe/sqle $(DOCKER_IMAGE) && \
    $(DOCKER) exec $${CTN_NAME} make test ; \
    $(DOCKER) stop $${CTN_NAME}

upload:
	curl -T $(shell pwd)/${PROJECT_NAME}.x86_64.rpm -u admin:ftpadmin ftp://${RELEASE_FTPD_HOST}/actiontech-${PROJECT_NAME}/qa/${VERSION}/${PROJECT_NAME}-${VERSION}-qa.x86_64.rpm
	curl -T $(shell pwd)/${SUB_PROJECT_NAME}.x86_64.rpm -u admin:ftpadmin ftp://${RELEASE_FTPD_HOST}/actiontech-${PROJECT_NAME}/qa/${VERSION}/${SUB_PROJECT_NAME}-${VERSION}-qa.x86_64.rpm

parser:
	cd build/goyacc && GOOS=${GOOS} GOARCH=${GOARCH} GOBIN=$(GOBIN) go install
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

swagger:
	$(GOBUILD) -o $(MAIN_MODULE)/$(PROJECT_NAME)/swag $(MAIN_MODULE)/build/swag/main.go
	rm -rf $(MAIN_MODULE)/sqle/docs
	$(MAIN_MODULE)/$(PROJECT_NAME)/swag init -g ./$(PROJECT_NAME)/api/app.go -o $(MAIN_MODULE)/sqle/docs