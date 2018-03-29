pb:
	  go get -u github.com/golang/protobuf/protoc-gen-go
		@echo "pb Start"
		cd configure && make pb
asset:
	mkdir assets
	cd assets;curl https://cdn.rawgit.com/v2ray/v2ray-core/e60de73c704d46d91633035e6b06184f7186a4e0/tools/release/config/geosite.dat > geosite.dat
	cd assets;curl https://cdn.rawgit.com/v2ray/v2ray-core/1777540e3d9eb7429c1ba72a93d8ef6c426bda13/release/config/geoip.dat > geoip.dat

shippedBinary:
	cd shippedBinarys; $(MAKE) shippedBinary

fetchDep:
	-go get -u github.com/xiaokangwang/V2RayConfigureFileUtil
	-cd $(GOPATH)/src/github.com/xiaokangwang/V2RayConfigureFileUtil;$(MAKE) all
	go get -u github.com/xiaokangwang/V2RayConfigureFileUtil
	-go get -u github.com/xiaokangwang/AndroidLibV2ray
	-cd $(GOPATH)/src/github.com/xiaokangwang/libV2RayAuxiliaryURL; $(MAKE) all
	-go get -u github.com/xiaokangwang/AndroidLibV2ray
	-cd $(GOPATH)/src/github.com/xiaokangwang/waVingOcean/configure; $(MAKE) pb
	go get -u github.com/xiaokangwang/AndroidLibV2ray

ANDROID_HOME=$(HOME)/android-sdk-linux
export ANDROID_HOME
PATH:=$(PATH):$(GOPATH)/bin
export PATH
downloadGoMobile:
	go get golang.org/x/mobile/cmd/...
	sudo apt-get install -qq libstdc++6:i386 lib32z1 expect
	cd ~ ;curl -L https://gist.githubusercontent.com/xiaokangwang/4a0f19476d86213ef6544aa45b3d2808/raw/ff5eb88663065d7159d6272f7b2eea0bd8b7425a/ubuntu-cli-install-android-sdk.sh | sudo bash - > /dev/null
	ls ~
	ls ~/android-sdk-linux/
	gomobile init -ndk ~/android-ndk-r15c;gomobile bind -v  -tags json github.com/xiaokangwang/AndroidLibV2ray

buildVGO:
	git clone https://github.com/xiaokangwang/V2RayGO.git
	ln libv2ray.aar V2RayGO/libv2ray/libv2ray.aar
	cd V2RayGO;echo "sdk.dir=$(ANDROID_HOME)" > local.properties
	cd V2RayGO; ./gradlew assembleRelease --stacktrace
	cd V2RayGO/app/build/outputs/apk/release; $(ANDROID_HOME)/build-tools/27.0.1/zipalign -v -p 4 app-release-unsigned.apk app-release-unsigned-aligned.apk
	cd V2RayGO/app/build/outputs/apk/release; keytool -genkey -v -keystore temp-release-key.jks -keyalg RSA -keysize 2048 -validity 365 -alias tempkey -dname "CN=vvv.kkdev.org, OU=VV, O=VVV, L=VVVVV, S=VV, C=VV" -storepass password -keypass password -noprompt
	cd V2RayGO/app/build/outputs/apk/release; $(ANDROID_HOME)/build-tools/27.0.1/apksigner sign --ks temp-release-key.jks  --ks-pass pass:password --key-pass pass:password --out app-release.apk app-release-unsigned-aligned.apk

BuildMobile:
	@echo Stub

all: asset pb shippedBinary fetchDep
	@echo DONE
