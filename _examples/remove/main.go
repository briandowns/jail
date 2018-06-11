package main

import (
	"fmt"
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
	var jid int
	var err error
	go func() {
		jid, err = jail.Jail(o)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		time.Sleep(30 * time.Second)
	}()
	fmt.Printf("removing JID: %d\n", jid)
	if err := jail.Remove(jid); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
