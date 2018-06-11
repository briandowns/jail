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
	var j int32
	go func(jid *int32) {
		jailID, err := jail.Jail(o)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		*jid = jailID
		time.Sleep(30 * time.Second)
	}(&j)
	time.Sleep(5 * time.Second)
	fmt.Printf("removing JID: %d\n", j)
	if err := jail.Remove(j); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
