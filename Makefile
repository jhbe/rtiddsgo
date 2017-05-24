#!/bin/make -f

# NDDS_QOS_PROFILE = example/my_qos.xml

all: example

# ===========================================================================
#  Build the parser tool.
# ===========================================================================

parse/parser.go: $(filter-out parse/parser.go, $(wildcard parse/*.go)) parse/parser.y  Makefile
    # Generate the parser.go file from the parser.y goyacc file.
	go generate parse/parse.go
	rm parse/y.output

goddsgen: main/goddsgen.go $(wildcard parse/*.go) $(wildcard generate/*.go) parse/parser.go Makefile
	go build main/goddsgen.go

# ===========================================================================
#  Build the example pre-requisites.
# ===========================================================================

example/mymessage.c: example/mymessage.idl Makefile
	rm -rf example/*.h example/*.c
	$(NDDSHOME)/bin/rtiddsgen -create typefiles -d example -I example -language C example/mymessage.idl

example/mymodule_mymessage.go example/mymodule_myerror.go: goddsgen example/mymessage.idl Makefile
	rm -f  example/mymodule_*.go
	./goddsgen example/mymessage.idl $(NDDSHOME) ../example mymessage example example
	go fmt ./example

# ===========================================================================
#  Build the example/publisher example.
# ===========================================================================

./publisher: example/mymodule_mymessage.go example/mymodule_myerror.go example/mymessage.c main/publisher.go Makefile
	go build main/publisher.go
	
# ===========================================================================
#  Build the example/subscriber example.
# ===========================================================================

./subscriber: example/mymodule_mymessage.go example/mymodule_myerror.go example/mymessage.c main/subscriber.go Makefile
	go build main/subscriber.go

# ===========================================================================

example: ./publisher ./subscriber
	
# ===========================================================================

clean:
	rm -rf test/mymessage.go parseidl/parser.go parseidl/y.output test_go/*
	rm -rf example/mymodule_mymessage.go example/*.h example/*.c
	rm -rf ./publisher ./subscriber ./goddsgen
