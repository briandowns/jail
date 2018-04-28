// build +FreeBSD

// Package jail provides the ability to lock a process
// or Goroutine into a FreeBSD jail
package jail

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
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

const (
	// CreateFlag Create a new jail. If a jid or name parameters exists, they must
	// not refer to an existing jail.
	CreateFlag = uintptr(0x01)

	// UpdateFlag Modify an existing jail. One of the jid or name parameters must
	// exist, and must refer to an existing jail. If both JAIL_CREATE and JAIL_UPDATE
	// are set, a jail will be created if it does not yet exist, and modified if it does exist.
	UpdateFlag = uintptr(0x02)

	// AttachFlag In addition to creating or modifying the jail, attach the current process
	// to it, as with the jail_attach() system call.
	AttachFlag = uintptr(0x04)

	// DyingFlag Allow setting a jail that is in the process of being removed.
	DyingFlag = uintptr(0x08)

	// SetMaskFlag ...
	SetMaskFlag = uintptr(0x0f)

	// GetMaskFlag ...
	GetMaskFlag = uintptr(0x08)
)

const jailAPIVersion uint32 = 2

type jail struct {
	Version  uint32
	Path     uintptr
	Name     uintptr
	Hostname uintptr
	IP4s     uint32
	IP6s     uint32
	IP4      uintptr
	IP6      uintptr
}

// Opts hlds the options to be passed in to
// create the new jail
type Opts struct {
	Version  int
	Path     string
	Name     string
	Hostname string
	IP4      string
	Chdir    bool
}

// typedef uint32_t in_addr_t
type inAddrT uint32

// inAddr
type inAddr struct {
	sAddr inAddrT
}

// validate makes sure the required fields are present
func (o *Opts) validate() error {
	if o.Path == "" {
		return errors.New("missing path")
	}
	if o.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

// Jail takes the given parameters, validates, and creates a new jail
func Jail(o *Opts) (int, error) {
	if err := o.validate(); err != nil {
		return 0, err
	}
	jn, err := syscall.BytePtrFromString(o.Name)
	if err != nil {
		return 0, err
	}
	jp, err := syscall.BytePtrFromString(o.Path)
	if err != nil {
		return 0, err
	}
	hn, err := syscall.BytePtrFromString(o.Name)
	if err != nil {
		return 0, err
	}
	/*var uint32ip uint32
	if jo.IP4 != "" {
		uip, err := ipToUint32(jo.IP4)
		if err != nil {
			return 0, err
		}
		uint32ip = uip
	}*/
	//t := tester()
	//t := uint32(3232235720)
	/*x := ((t>>24)&0xff)    |
	             ((t<<8)&0xff0000) |
	             ((t>>8)&0xff00)   |
	             ((t<<24)&0xff000000)

		iat := inAddrT(x)*/
	oip := net.ParseIP("192.168.0.200")
	nip := oip.To4()
	//x1 := binary.LittleEndian.Uint32(nip)
	//x1 := uint32(uint(nip[0]) | uint(nip[1])<<8 | uint(nip[2])<<16 | uint(nip[3])<<32)
	x1 := uint32(uint(nip[3]) | uint(nip[2])<<8 | uint(nip[1])<<16 | uint(nip[0])<<32)
	ia := &inAddr{sAddr: x1}

	fmt.Printf("%+v\n", ia)

	j := &jail{
		Version:  uint32(0),
		Path:     uintptr(unsafe.Pointer(jp)),
		Hostname: uintptr(unsafe.Pointer(hn)),
		Name:     uintptr(unsafe.Pointer(jn)),
		IP4s:     uint32(1),
		IP6s:     uint32(0),
		IP4:      uintptr(unsafe.Pointer(ia)),
	}
	r1, _, e1 := syscall.Syscall(sysJail, uintptr(unsafe.Pointer(j)), 0, 0)
	if e1 != 0 {
		return 0, fmt.Errorf("%d", e1)
	}
	if o.Chdir {
		if err := os.Chdir("/"); err != nil {
			return 0, err
		}
	}
	return int(r1), nil
}

// Clone creates a new version of the previously created jail
func (j *jail) Clone() (int, error) {
	nj := &jail{
		Version:  j.Version,
		Path:     j.Path,
		Name:     j.Name,
		Hostname: j.Hostname,
	}
	r1, _, e1 := syscall.Syscall(sysJail, uintptr(unsafe.Pointer(nj)), 0, 0)
	if e1 != 0 {
		return 0, fmt.Errorf("%d", e1)
	}
	return int(r1), nil
}

// Params contains the individual settings passed in to either get
// or set a jail
type Params map[string]interface{}

// NewParams creates a new value of type Params by
// initializing the underluing map
func NewParams() Params {
	return make(map[string]interface{})
}

// Add adds the given key and value to the params map
func (p Params) Add(k string, v interface{}) error {
	if p == nil {
		return errors.New("cannot assign values to nil map")
	}
	if _, ok := p[k]; !ok {
		p[k] = v
		return nil
	}
	return fmt.Errorf("key of %s already set with value of %v", k, p[k])
}

// Set
func Set(params Params, flags uintptr) error {
	iovec := uintptr(unsafe.Pointer(params))
	_, _, e1 := syscall.Syscall(sysJailSet, iovec, 0, flags)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// Get
func Get(params Params, flags uintptr) error {
	iovec := uintptr(unsafe.Pointer(params))
	_, _, e1 := syscall.Syscall(sysJailGet, iovec, 0, flags)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// Attach receives a jail ID and attempts to attach the current
// process to that jail
func Attach(jailID int) error {
	jid := uintptr(unsafe.Pointer(&jailID))
	_, _, e1 := syscall.Syscall(sysJailAttach, jid, 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

// Remove receives a jail ID and attempts to remove the associated jail
func Remove(jailID int) error {
	jid := uintptr(unsafe.Pointer(&jailID))
	_, _, e1 := syscall.Syscall(sysJailRemove, jid, 0, 0)
	if e1 != 0 {
		return fmt.Errorf("%d", e1)
	}
	return nil
}

func tester() uint32 {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(3232235720))
	return binary.BigEndian.Uint32(buf)
}

// ipToUint32 converts a string representation of an IP address into
// an uint32
//func ipToUint32(sip string) (uint32, error) {
/*var ip net.IP
	ip, _, err := net.ParseCIDR(sip)
	if err != nil {
		return 0, err
	}
        fmt.Println(ip.String())*/
// convert string to net.IP
/*if len(ip) == 16 {
	return binary.BigEndian.Uint32(ip[12:16]), nil
}*/
//ip := net.ParseIP(sip)
//fmt.Println(len(ip))
//fmt.Println(ip.String())
//return binary.LittleEndian.Uint32(ip), nil
//buf := make([]byte, 8)
//return binary.BigEndian.PutUint32(buf, v), nil
//return binary.BigEndian.Uint32(ip[12:16]), nil
//}

// uint32ip converts an uint32 representation of a string into an IP
func uint32ip(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip.String()
}
