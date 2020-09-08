#!/bin/make -f

NDDSHOME=/home/johan/rti_connext_dds-6.0.0
RTILIBDIR=$(NDDSHOME)/lib/x64Linux3gcc4.8.2

all: ./examplepub ./examplesub ./verificationpub ./verificationsub

idl: example/src/example_constants.go example/src/example_constants.go verification/src/one_constants.go verification/src/two_constants.go

#
# GoDdsGen
#
goddsgen: main/goddsgen.go generate/*.go parse/*.go
	go build main/goddsgen.go

#
# Example
#
example/src:
	mkdir -p example/src

example/src/example.xml: example/src example/idl/example.idl
	rm -rf example/src/example.xml
	$(NDDSHOME)/bin/rtiddsgen -d example/src -I exampleidl -convertToXml example/idl/example.idl

example/src/example.c: example/src example/idl/example.idl
	rm -rf example/src/example.h example/src/example.c example/src/exampleSupport.* example/src/examplePlugin.*
	$(NDDSHOME)/bin/rtiddsgen -d example/src -I example/idl -create typefiles  -language c example/idl/example.idl

example/src/example_constants.go: goddsgen example/src/example.c example/src/example.xml
	rm -rf example/src/*.go
	./goddsgen example/src/example.xml $(NDDSHOME) $(RTILIBDIR) example/src example

./examplepub: example/src/example_constants.go example/src/example.c example/pub/examplepub.go
	go build example/pub/examplepub.go

./examplesub: example/src/example_constants.go example/src/example.c example/sub/examplesub.go
	go build example/sub/examplesub.go

#
# Verification
#
verification/src:
	mkdir -p verification/src

verification/src/one.xml: verification/src verification/idlOne/one.idl
	rm -rf verification/src/one.xml
	$(NDDSHOME)/bin/rtiddsgen -d verification/src -I verification/idlOne -convertToXml verification/idlOne/one.idl

verification/src/two.xml: verification/src verification/idlTwo/two.idl
	rm -rf verification/src/two.xml
	$(NDDSHOME)/bin/rtiddsgen -d verification/src -I verification/idlOne -I verification/idlTwo -convertToXml verification/idlTwo/two.idl

verification/src/one.c: verification/src verification/idlOne/one.idl
	rm -rf verification/src/one.h verification/src/one.c verification/src/oneSupport.* verification/src/onePlugin.*
	$(NDDSHOME)/bin/rtiddsgen -d verification/src -I verification/idlOne -create typefiles  -language c verification/idlOne/one.idl

verification/src/two.c: verification/src verification/idlTwo/two.idl
	rm -rf verification/src/two.h verification/src/two.c verification/src/twoSupport.* verification/src/twoPlugin.*
	$(NDDSHOME)/bin/rtiddsgen -d verification/src -I verification/idlOne -I verification/idlTwo -create typefiles  -language c verification/idlTwo/two.idl

verification/src/one.go: goddsgen verification/src/one.c verification/src/one.xml
	rm -rf verification/src/one*.go
	./goddsgen verification/src/one.xml $(NDDSHOME) $(RTILIBDIR) verification/src eb

verification/src/two.go: goddsgen verification/src/two.c verification/src/two.xml
	rm -rf verification/src/two*.go
	./goddsgen verification/src/two.xml $(NDDSHOME) $(RTILIBDIR) verification/src eb

./verificationpub: verification/src/one.go verification/src/one.c verification/src/two.go verification/src/two.c verification/pub/verificationpub.go
	go build verification/pub/verificationpub.go

./verificationsub: verification/src/one.go verification/src/one.c verification/src/two.go verification/src/two.c verification/sub/verificationsub.go
	go build verification/sub/verificationsub.go

clean:
	rm -rf example/src verification/src
	rm -f ./examplepub ./examplesub ./verificationpub ./verificationsub ./goddsgen
