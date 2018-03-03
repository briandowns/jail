package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/jail"
)

func main() {
	jo := &jail.JailOpts{
		Path:     "/zroot/jails/build",
		Name:     "jailname",
		Hostname: "hostname",
		IP4:      "192.168.0.200",
		Chdir:    true,
	}

	jid, err := jail.Jail(jo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("JID: %d\n", jid)
	time.Sleep(30 * time.Second)
	fmt.Println("/ director listing in jail")
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
