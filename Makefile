pb:
	  go get -u github.com/golang/protobuf/protoc-gen-go
		@echo "pb Start"
		cd configure && make pb
asset:
	mkdir assets
	cd assets;curl https://raw.githubusercontent.com/v2ray/v2ray-core/e60de73c704d46d91633035e6b06184f7186a4e0/tools/release/config/geosite.dat > geosite.dat
	cd assets;curl https://github.com/v2ray/v2ray-core/blob/e60de73c704d46d91633035e6b06184f7186a4e0/tools/release/config/geoip.dat?raw=true > geoip.dat

shippedBinary:
	cd shippedBinarys; $(MAKE) shippedBinary

fetchDep:
	-go get -u github.com/xiaokangwang/V2RayConfigureFileUtil
	-cd $(GOPATH)/src/github.com/xiaokangwang/V2RayConfigureFileUtil;$(MAKE) all
	go get -u github.com/xiaokangwang/V2RayConfigureFileUtil
	-go get -u github.com/xiaokangwang/AndroidLibV2ray
	-cd $(GOPATH)/src/github.com/xiaokangwang/libV2RayAuxiliaryURL; $(MAKE) all
	go get -u github.com/xiaokangwang/AndroidLibV2ray

downloadGoMobile:
	go get golang.org/x/mobile/cmd/gomobile
	sudo apt-get install -qq libstdc++6:i386 lib32z1 expect
	mkdir ~/andrsdk; cd ~/andrsdk ;curl -L https://gist.githubusercontent.com/xiaokangwang/4a0f19476d86213ef6544aa45b3d2808/raw/ebab0c3fe856e6163e49c0b0fb8bc8b98fb8ab4a/ubuntu-cli-install-android-sdk.sh | sudo bash -
all: asset pb shippedBinary fetchDep
	@echo DONE
