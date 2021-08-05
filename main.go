package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"gopkg.in/yaml.v2"
)

var isStart bool

type keyinfo struct {
	Mi     int64  `yaml:"mi"`
	Action string `yaml:"action"`
}

var actmap = map[string]interface{}{
	"downup": robotgo.KeyTap,
	"down":   robotgo.KeyToggle,
}

func main() {
	f := flag.String("f", "", "config_*.yaml, input the * part")
	flag.Parse()

	keyconf := config(*f)

	done := make(chan bool, len(keyconf))

	bind(done, keyconf)
}

func bind(done chan bool, keyconf map[string]keyinfo) {
	exit(done, keyconf)
	start(done, keyconf)
	pause(done, keyconf)
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func exit(done chan bool, keyconf map[string]keyinfo) {
	robotgo.EventHook(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		stop(done, len(keyconf))
		robotgo.EventEnd()
		os.Exit(0)
	})
}

func start(done chan bool, keyconf map[string]keyinfo) {
	robotgo.EventHook(hook.KeyDown, []string{"5"}, func(e hook.Event) {
		fmt.Println("start")
		if !isStart {
			isStart = true
			runKey(keyconf, done)
		}

	})
}

func pause(done chan bool, keyconf map[string]keyinfo) {
	robotgo.EventHook(hook.KeyDown, []string{"`"}, func(e hook.Event) {
		fmt.Println("pause")
		if isStart {
			isStart = false
			stop(done, len(keyconf))
		}
	})
}

func config(flag string) map[string]keyinfo {
	fName := "config_" + flag + ".yaml"
	fmt.Println(fName)

	bt, err := ioutil.ReadFile(fName)
	if err != nil {
		panic(err.Error())
	}
	keyconf := make(map[string]keyinfo)
	err = yaml.Unmarshal(bt, &keyconf)
	if err != nil {
		panic(err.Error())
	}

	return keyconf
}

func dynamicCall(act string, key string) {
	switch act {
	case "down":
		actmap[act].(func(string, ...string) string)(key, "down")
	case "downup":
		actmap[act].(func(string, ...interface{}) string)(key)
	}
}

func runKey(keyconf map[string]keyinfo, done <-chan bool) {

	for key, info := range keyconf {

		go func(key string, info keyinfo) {
			do := time.Duration(info.Mi) * time.Millisecond
			time.Sleep(time.Millisecond * 200)
			dynamicCall(info.Action, key)
			for {
				select {
				case <-done:
					robotgo.KeyToggle(key, "up")
					return
				case <-time.After(do):
					dynamicCall(info.Action, key)
				}
			}
		}(key, info)
	}
}

func stop(done chan<- bool, len int) {
	for i := 0; i < len; i++ {
		done <- true
	}
}
