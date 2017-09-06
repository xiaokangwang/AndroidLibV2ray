package libv2ray

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type StabilityAssist struct {
	Engaged      bool
	Forground    bool
	intr         chan int
	sharedClient *http.Client
	SASS         StabilityAssistSupportSet
	armed        bool
}

func (sa *StabilityAssist) monitor() {
	select {
	case <-sa.intr:
	default:
	}
	log.Println("Monitor Engaged")
	first := true
	waitsec := 15
	for sa.Engaged {
		if first {
			waitsec = 30
			first = false
		} else {
			if !sa.Forground {
				waitsec = (waitsec * 3) / 2
				if waitsec > 500 {
					waitsec = 500
				}
			} else {
				waitsec = 15
			}
		}
		currentWait := time.Tick(time.Duration(waitsec) * time.Second)
		select {
		case <-currentWait:
			sa.probe()
		case s := <-sa.intr:
			if s != 0 {
				sa.probe()
			}
			continue
		}
	}
}

func (sa *StabilityAssist) ProbeNowOpi() {
	select {
	case sa.intr <- 1:
	default:
	}
}

func (sa *StabilityAssist) probe() {
	if sa.sharedClient == nil {
		tr := &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 600 * time.Second,
		}
		sa.sharedClient = &http.Client{Transport: tr, Timeout: 20 * time.Second}
	}
	log.Println("Probing")
	res, err := sa.sharedClient.Get("https://static.kkdev.org/CONNTEST")
	if err != nil {
		log.Println("Probing fail1", err)
		sa.probeFail()
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Probing fail2", err)
		sa.probeFail()
		return
	}
	res.Body.Close()
	sa.armed = true
	log.Println("Armed")
}

func (sa *StabilityAssist) probeFail() {
	if sa.armed {
		sa.Engaged = false
		sa.SASS.OnProbeFailed()
	}
}

func (sa *StabilityAssist) Start() {
	sa.Engaged = true
	sa.armed = false
	select {
	case sa.intr <- 0:
	default:
	}
	go sa.monitor()
}

func (sa *StabilityAssist) Stop() {
	sa.Engaged = false
	select {
	case sa.intr <- 0:
	default:
	}
}

func GetStabilityAssist() *StabilityAssist {
	log.SetPrefix("SAS")
	return &StabilityAssist{intr: make(chan int, 1), Forground: true}
}

type StabilityAssistSupportSet interface {
	OnProbeFailed()
}
