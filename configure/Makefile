pbSource := $(shell echo *.proto)
#%.pb.go:%.proto
#	protoc --go_out=. $^

pb: clean
	sudo apt-get install tree
	tree $(GOPATH)
	$(GOPATH)/bin/protoc --go_out=. $(pbSource)
	@echo "Done"

clean:
	-rm *.pb.go
