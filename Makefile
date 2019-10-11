VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || true)


ifndef ECFLOW_BUILD_DIR
$(error ECFLOW_BUILD_DIR is not set)
endif

ifndef ECFLOW_SOURCE_DIR
$(error ECFLOW_SOURCE_DIR is not set)
endif

ifndef BOOST_LIB_DIR
$(error BOOST_LIB_DIR is not set)
endif

CGO_CXXFLAGS +=  -I$(pwd)/ecflow_client_standalone
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/ACore/src
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/ANattr/src
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/ANode/src
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/Base/src
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/Base/src/cts
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/Base/src/stc
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/CSim/src
CGO_CXXFLAGS +=  -I${ECFLOW_SOURCE_DIR}/Client/src
CGO_CXXFLAGS +=  -Wno-deprecated-declarations
export CGO_CXXFLAGS

CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/Client -llibclient
CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/Base -lbase
CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/CSim -llibsimu
CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/ANode -lnode
CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/ANattr -lnodeattr
CGO_LDFLAGS += -L${ECFLOW_BUILD_DIR}/ACore -lcore
CGO_LDFLAGS += -L${BOOST_LIB_DIR} -lboost_system-mt -lboost_filesystem-mt -lboost_date_time-mt -lboost_program_options-mt -lboost_serialization-mt -lboost_thread-mt -lboost_regex-mt
export CGO_LDFLAGS

export BIN_PATH := $(shell pwd)/bin

.PHONY: cmd

cmd:
	go build \
		-ldflags "-X \"github.com/perillaroc/ecflow-watchman.Version=${VERSION}\" \
    	-X \"github.com/perillaroc/ecflow-watchman.BuildTime=${BUILD_TIME}\" \
    	-X \"github.com/perillaroc/ecflow-watchman.GitCommit=${GIT_COMMIT}\" " \
		-o bin/ecflow_watchman \
		cmd/main.go