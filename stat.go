package libv2ray

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type InterfaceInfo struct {
	RxByte, RxPacket, TxByte, TxPacket int
}

func (sc *StatControler) CollectInterfaceInfo() error {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return err
	}
	d, _ := ioutil.ReadAll(f)
	s := strings.Split(strings.TrimSpace(string(d)), "\n")
	for _, scc := range s {
		subsc := strings.Split(strings.TrimSpace(scc), " ")
		//log.Println("CollectInterfaceInfo Examing ", scc)

		if subsc[0] == sc.InterfaceTarget {
			subscclean := make([]string, 0, 5)
			for _, csubsc := range subsc {
				if csubsc != "" {
					subscclean = append(subscclean, csubsc)
				}
			}
			//spew.Dump(subscclean)
			var infoCandidcate InterfaceInfo
			var convError error
			infoCandidcate.RxByte, convError = strconv.Atoi(subscclean[1])
			if convError != nil {
				return convError
			}
			infoCandidcate.RxPacket, convError = strconv.Atoi(subscclean[2])
			if convError != nil {
				return convError
			}
			infoCandidcate.TxByte, convError = strconv.Atoi(subscclean[9])
			if convError != nil {
				return convError
			}
			infoCandidcate.TxPacket, convError = strconv.Atoi(subscclean[10])
			if convError != nil {
				return convError
			}
			sc.CollectedInterfaceInfo = &infoCandidcate
		}
	}
	return nil
}

func (v *V2RayPoint) GetStatControler() *StatControler {
	return &StatControler{InterfaceTarget: "tun0:"}
}

type StatControler struct {
	InterfaceTarget        string
	CollectedInterfaceInfo *InterfaceInfo
}
