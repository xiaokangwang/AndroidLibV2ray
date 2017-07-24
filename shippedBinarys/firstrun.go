package shippedBinarys

import (
	"os"
	"strconv"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
)

type FirstRun struct {
	status *CoreI.Status
}

func (v *FirstRun) checkIfRcExist() error {
	datadir := v.status.GetDataDir()

	if _, err := os.Stat(datadir + strconv.Itoa(CoreI.CheckVersion())); !os.IsNotExist(err) {
		return nil
	}
	IndepDir, err := AssetDir("ArchIndep")
	if err != nil {
		return err
	}
	for _, fn := range IndepDir {
		RestoreAsset(datadir, fn)
	}
	DepDir, err := AssetDir("ArchDep")
	if err != nil {
		return err
	}
	for _, fn := range DepDir {
		DepDir2, err := AssetDir("ArchDep/" + fn)
		if err != nil {
			return err
		}
		for _, FND := range DepDir2 {
			RestoreAsset(datadir, "ArchDep/"+fn+"/"+FND)
		}
	}
	s, _ := os.Create(datadir + strconv.Itoa(CoreI.CheckVersion()))
	s.Close()

	return nil
}
