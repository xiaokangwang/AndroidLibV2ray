Platdep=shippedBinary.386 shippedBinary.amd64 shippedBinary.arm64 shippedBinary.arm

shippedBinaryDep: $GOPATH/bin/go-bindata
		go get -u github.com/jteeuwen/go-bindata/...

shippedBinary.%:
	$GOPATH/bin/go-bindata -nometadata -nomemcopy -pkg shippedBinarys -o ./binary_$*.go -tags $* ArchIndep/ ArchDep/$*/

shippedBinary: $(Platdep)
	@echo "Done"

clean:
	-rm binary*
