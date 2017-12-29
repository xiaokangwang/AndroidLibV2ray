package libv2ray

var overridedAssets map[string](string)

/*SetAssetsOverride will define file path @[data] to
provide the asset named @[path].

SetAssetsOverride("geoip.dat","/data/data/appdir/dat/geoip-override.dat")
will ask libv2ray to serve v2ray "geoip.dat" with
file "/data/data/appdir/dat/geoip-override.dat" in file system
*/
func SetAssetsOverride(path string, data string) {
	if overridedAssets == nil {
		overridedAssets = make(map[string](string))
	}
	overridedAssets[path] = data
}

/*ClearAssetsOverride can disable an override
at given @[path].

ClearAssetsOverride("geoip.dat")
will ask libv2ray to serve vanilla asset file.
*/
func ClearAssetsOverride(path string) {
	if overridedAssets == nil {
		overridedAssets = make(map[string](string))
	}
	delete(overridedAssets, path)
}
