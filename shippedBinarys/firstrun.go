package shippedBinarys

import (
	"fmt"
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
		fmt.Println(err)
		return nil
	}
	IndepDir, err := AssetDir("ArchIndep")
	fmt.Print(IndepDir)
	if err != nil {
		return err
	}
	for _, fn := range IndepDir {
		err := RestoreAsset(datadir, "ArchIndep/"+fn)
		//GrantPremission
		os.Chmod(datadir+"ArchIndep/"+fn, 0700)
		fmt.Println(os.Symlink(datadir+"ArchIndep/"+fn, datadir+fn))
		fmt.Print(err)
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
			os.Chmod(datadir+"ArchDep/"+fn+"/"+FND, 0700)
			os.Symlink(datadir+"ArchDep/"+fn+"/"+FND, datadir+FND)
		}
	}
	s, _ := os.Create(datadir + strconv.Itoa(CoreI.CheckVersion()))
	s.Close()

	return nil
}

func (v *FirstRun) SetCoreI(status *CoreI.Status) {
	v.status = status
}

func (v *FirstRun) CheckAndExport() error {
	return v.checkIfRcExist()
}
