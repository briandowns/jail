// build +FreeBSD
package jail

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	sysJail       = 338
	sysJailAttach = 436
	sysJailGet    = 506
	sysJailSet    = 507
	sysJailRemove = 508
)

const jailAPIVersion = 2

// IP4
type IP4 struct{}

// IP6
type IP6 struct{}

type Jail struct {
	Version  uint32
	Path     uintptr
	Name     uintptr
	Hostname uintptr
	/*IP4s     uint32
	IP6s     uint32
	IP4r     interface{}
	IP6      interface{}*/
}

// JailOpts hlds the options to be passed in to
// create the new jail
type JailOpts struct {
	Version  int
	Path     string
	Name     string
	Hostname string
	Chdir    bool
}

// validate makes sure the required fields are present
func (j *JailOpts) validate() error {
	if j.Version != jailAPIVersion {
		return errors.New("invalid version")
	}
	if j.Path == "" {
		return errors.New("missing path")
	}
	if j.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

// New takes the given parameters, validates, and creates a new jail
func New(jo *JailOpts) (int, error) {
	if err := jo.validate(); err != nil {
		return 0, err
	}
	jn, err := syscall.BytePtrFromString(jo.Name)
	if err != nil {
		return 0, err
	}
	jp, err := syscall.BytePtrFromString(jo.Path)
	if err != nil {
		return 0, err
	}
	hn, err := syscall.BytePtrFromString(jo.Name)
	if err != nil {
		return 0, err
	}
	jail := &Jail{
		Version:  uint32(jo.Version),
		Path:     uintptr(unsafe.Pointer(jp)),
		Name:     uintptr(unsafe.Pointer(jn)),
		Hostname: uintptr(unsafe.Pointer(hn)),
	}
	r1, _, e1 := syscall.Syscall(sysJail, uintptr(unsafe.Pointer(jail)), 0, 0)
	if e1 != 0 {
		return 0, fmt.Errorf("%d", e1)
	}
	if jo.Chdir {
		if err := os.Chdir(jo.Path); err != nil {
			return 0, err
		}
	}
	return int(r1), nil
}

// Clone creates a new version of the previously created jail
func (j *Jail) Clone() (*Jail, int, error) {
	nj := &Jail{
		Version:  j.Version,
		Path:     j.Path,
		Name:     j.Name,
		Hostname: j.Hostname,
	}
	r1, _, e1 := syscall.Syscall(sysJail, uintptr(unsafe.Pointer(nj)), 0, 0)
	if e1 != 0 {
		return nil, 0, fmt.Errorf("%d", e1)
	}
	return nj, int(r1), nil
}

// IOVEC
type IOVEC struct {
	IOVBase interface{}
	IOVLen  int64
}

// NamedPars
type NamedPars map[string]interface{}

// validate makes sure the given map is usable
func (NamedPars) validate() error {
	return nil
}

// JailGet
func JailGet() error {
	_, _, e1 := syscall.Syscall(sysJailGet, 0, 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// JailSet
func JailSet() error {
	_, _, e1 := syscall.Syscall(sysJailSet, 0, 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// JailAttach
func JailAttach(jid int) error {
	_, _, e1 := syscall.Syscall(sysJailAttach, uintptr(jid), 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// JailRemove
func JailRemove(jid int) error {
	_, _, e1 := syscall.Syscall(sysJailRemove, uintptr(jid), 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}
