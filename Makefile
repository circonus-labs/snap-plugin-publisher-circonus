# lock to the version required by snap plugin library
# govendor fetch google.golang.org/grpc@0032a855ba5c8a3c8e0d71c2deef354b70af1584

default:
	$(MAKE) deps
	$(MAKE) all
deps:
	bash -c "govendor sync"
test:
	bash -c "./scripts/test.sh $(TEST)"
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
