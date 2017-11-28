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

ANDROID_HOME=$(HOME)/android-sdk-linux
export ANDROID_HOME
downloadGoMobile:
	go get golang.org/x/mobile/cmd/gomobile
	sudo apt-get install -qq libstdc++6:i386 lib32z1 expect
	cd ~ ;curl -L https://gist.githubusercontent.com/xiaokangwang/4a0f19476d86213ef6544aa45b3d2808/raw/c23f4e3ca83d3a97b05a53924d7634ff4f80e434/ubuntu-cli-install-android-sdk.sh | sudo bash -
	ls ~/android-sdk-linux/
	gomobile init -ndk ~/android-ndk-r15c;gomobile bind -v  -tags json github.com/xiaokangwang/AndroidLibV2ray

buildVGO:
	git clone https://github.com/xiaokangwang/V2RayGO.git
	ln libv2ray.aar V2RayGO/libv2ray/libv2ray.aar
	cd V2RayGO;echo "sdk.dir=$(ANDROID_HOME)" > local.properties
	cd V2RayGO; ./gradlew
	sudo apt install zipalign
	cd V2RayGO/app/build/outputs/apk/release; $(ANDROID_HOME)/build-tools/27.0.1/zipalign -v -p 4 app-release-unsigned.apk app-release-unsigned-aligned.apk
	cd V2RayGO/app/build/outputs/apk/release; keytool -genkey -v -keystore temp-release-key.jks -keyalg RSA -keysize 2048 -validity 10000 -alias tempkey
	cd V2RayGO/app/build/outputs/apk/release; $(ANDROID_HOME)/build-tools/27.0.1/apksigner sign --ks temp-release-key.jks --out app-release.apk app-release-unsigned-aligned.apk

BuildMobile:
	@echo Stub

all: asset pb shippedBinary fetchDep
	@echo DONE
