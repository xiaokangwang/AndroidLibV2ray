pb:
		@echo "pb Start"
		cd configure && make pb
asset:
	mkdir assets
	cd assets;wget https://raw.githubusercontent.com/v2ray/v2ray-core/e60de73c704d46d91633035e6b06184f7186a4e0/tools/release/config/geosite.dat
	cd assets;wget https://github.com/v2ray/v2ray-core/blob/e60de73c704d46d91633035e6b06184f7186a4e0/tools/release/config/geoip.dat?raw=true

shippedBinary:
	cd shippedBinary; $(MAKE) shippedBinary

fetchDep:
	-go get -u github.com/xiaokangwang/V2RayConfigureFileUtil

all: asset pb shippedBinary fetchDep
	@echo DONE
