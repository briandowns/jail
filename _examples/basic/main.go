package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/jail"
)

func main() {
	o := &jail.Opts{
		Version:  uint32(2),
		Path:     "/zroot/jails/build",
		Name:     "jailname",
		Hostname: "hostname",
		IP4:      "192.168.0.200",
		Chdir:    true,
	}
	jid, err := jail.Jail(o)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("JID: %d - / directory listing in jail", jid)

	// here so a `jls` can be ran seperately to see that the jail is running
	time.Sleep(30 * time.Second)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
