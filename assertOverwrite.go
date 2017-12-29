package libv2ray

var overridedAssets map[string](string)

func SetAssetsOverride(path string, data string) {
	if overridedAssets == nil {
		overridedAssets = make(map[string](string))
	}
	overridedAssets[path] = data
}

func ClearAssetsOverride(path string) {
	if overridedAssets == nil {
		overridedAssets = make(map[string](string))
	}
	delete(overridedAssets, path)
}
