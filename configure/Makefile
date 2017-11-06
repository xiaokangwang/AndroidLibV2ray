pbSource := $(shell echo *.proto)
#%.pb.go:%.proto
#	protoc --go_out=. $^
installpb:
			mkdir ~/pbinst; cd ~/pbinst; curl -OL https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip; unzip protoc-3.2.0-linux-x86_64.zip -d protoc3; mv protoc3/bin/* $(GOPATH)/bin/

pb: clean installpb
	protoc --go_out=. $(pbSource)
	@echo "Done"

clean:
	-rm *.pb.go
