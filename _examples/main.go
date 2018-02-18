package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/briandowns/jail"
)

func main() {
	jo := &jail.JailOpts{
		Path:     "/path/to/jail",
		Name:     "jailname",
		Hostname: "hostname",
		Chdir:    true,
	}

	jid, err := jail.Jail(jo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("JID: %d\n", jid)

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
