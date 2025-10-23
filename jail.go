/*-
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2026 Brian J. Downs
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

// build +FreeBSD

// Package jail provides the ability to lock a process
// or Goroutine into a FreeBSD jail.
package jail

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

const EtcdConfigFile = "/etc/jail.conf"

const (
	sysJail       = 338
	sysJailAttach = 436
	sysJailGet    = 506
	sysJailSet    = 507
	sysJailRemove = 508
)

const (
	// CreateFlag Create a new jail. If a jid or name parameters exists, they
	// must not refer to an existing jail.
	CreateFlag = uintptr(0x01)

	// UpdateFlag Modify an existing jail. One of the jid or name parameters must
	// exist, and must refer to an existing jail. If both JAIL_CREATE and JAIL_UPDATE
	// are set, a jail will be created if it does not yet exist, and modified if
	// it does exist.
	UpdateFlag = uintptr(0x02)

	// AttachFlag In addition to creating or modifying the jail, attach the current
	// process to it, as with the jail_attach() system call.
	AttachFlag = uintptr(0x04)

	// DyingFlag Allow setting a jail that is in the process of being removed.
	DyingFlag = uintptr(0x08)

	// SetMaskFlag ...
	SetMaskFlag = uintptr(0x0f)

	// GetMaskFlag ...
	GetMaskFlag = uintptr(0x08)
)

// jailAPIVersion is the current jail API version.
const jailAPIVersion uint32 = 2

// MaxChildJails is the maximum number of jails
// for the system.
const MaxChildJails int64 = 999999

// jail contains the data that will be passed into
// the jail(2) syscall.
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

// Opts holds the options to be passed in to
// create the new jail.
type Opts struct {
	Version  uint32
	Path     string
	Name     string
	Hostname string
	IP4      string
	Chdir    bool
}

const (
	JailRawValue    = 0x01
	JailBool        = 0x02
	JailParamNoBool = 0x04
	JailParamSys    = 0x80
)

// JailParam
type JailParam struct {
	Name       string
	Value      interface{}
	ValueLen   int
	ElemLen    int
	CtlType    int
	StructType int
	Flags      int
}

// typedef uint32_t in_addr_t
type inAddrT uint32

// inAddr
type inAddr struct {
	sAddr inAddrT
}

// validate makes sure the required fields are present.
func (o *Opts) validate() error {
	if o.Path == "" {
		return errors.New("missing path")
	}
	if o.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

// Jail takes the given parameters, validates, and creates a new jail.
func Jail(o *Opts) (int32, error) {
	if err := o.validate(); err != nil {
		return 0, err
	}

	jn, err := unix.BytePtrFromString(o.Name)
	if err != nil {
		return 0, err
	}

	jp, err := unix.BytePtrFromString(o.Path)
	if err != nil {
		return 0, err
	}

	hn, err := unix.BytePtrFromString(o.Name)
	if err != nil {
		return 0, err
	}

	j := &jail{
		Version:  o.Version,
		Path:     uintptr(unsafe.Pointer(jp)),
		Hostname: uintptr(unsafe.Pointer(hn)),
		Name:     uintptr(unsafe.Pointer(jn)),
	}
	if o.IP4 != "" {
		uint32ip := ip2int(net.ParseIP(o.IP4))
		ia := &inAddr{
			sAddr: inAddrT(uint32ip),
		}
		j.IP4s = 1
		j.IP6s = uint32(0)
		j.IP4 = uintptr(unsafe.Pointer(ia))
	}

	r1, _, e1 := unix.Syscall(sysJail, uintptr(unsafe.Pointer(j)), 0, 0)
	if e1 != 0 {
		switch int(e1) {
		case ErrJailPermDenied:
			return 0, fmt.Errorf("unprivileged user: %d", e1)
		case ErrJailFaultOutsideOfAllocatedSpace:
			return 0, fmt.Errorf("fault outside of allocation space: %d", e1)
		case ErrJailInvalidVersion:
			return 0, fmt.Errorf("invalid version: %d", e1)
		case ErrjailNoFreeJIDFound:
			return 0, fmt.Errorf("no free JID found: %d", e1)
		case ErrJailNoSuchFileDirectory:
			return 0, fmt.Errorf("No such file or directory: %s\n", o.Path)
		}
		return 0, fmt.Errorf("%d", e1)
	}

	if o.Chdir {
		if err := os.Chdir("/"); err != nil {
			return 0, err
		}
	}

	return int32(r1), nil
}

// Clone creates a new version of the previously created jail.
func (j *jail) Clone() (int, error) {
	nj := &jail{
		Version:  j.Version,
		Path:     j.Path,
		Name:     j.Name,
		Hostname: j.Hostname,
	}

	r1, _, e1 := unix.Syscall(sysJail, uintptr(unsafe.Pointer(nj)), 0, 0)
	if e1 != 0 {
		return 0, fmt.Errorf("%d", e1)
	}

	return int(r1), nil
}

// ID returns the JID of the corresponding jail.
func ID(name string) (int32, error) {
	params := NewParams()
	params.Add("name", name)

	if err := Get(params, 0); err != nil {
		return -1, err
	}

	return params["jid"].(int32), nil
}

// Name returns the name of the corresponding jail.
func Name(id int32) (string, error) {
	params := NewParams()
	params.Add("jid", id)

	if err := Get(params, 0); err != nil {
		return "", err
	}

	return params["name"].(string), nil
}

// validParams contains a list of the valid parameters that
// are allowed to be used when calling the Set or Get functions.
// TODO(briandowns) add more as they are identified.
var validParams = []string{
	"jid",
	"name",
	"dying",
	"persist",
	"nopersist",
}

// isValidParam verifies that the given parameter is valid.
func isValidParam(param string) bool {
	for _, p := range validParams {
		if param == p {
			return true
		}
	}

	return false
}

// Params contains the individual settings passed in to either get
// or set a jail.
type Params map[string]interface{}

// NewParams creates a new value of type Params by
// initializing the underlying map.
func NewParams() Params {
	return make(map[string]interface{})
}

// Add adds the given key and value to the params map.
func (p Params) Add(k string, v interface{}) error {
	if p == nil {
		return errors.New("cannot assign values to nil map")
	}

	if !isValidParam(k) {
		return errors.New("invalid parameter: " + k)
	}

	if _, ok := p[k]; !ok {
		p[k] = v
		return nil
	}

	return fmt.Errorf("key of %q already set with value of %v", k, p[k])
}

// Validate is used to make sure that the params assigned
// are indeed correct and usable. This has been exposed for
// a caller to do validation as well as the package interally.
func (p Params) Validate() error {
	return nil
}

// buildIovec takes the containing map value and builds
// out a slice of syscall.Iovec.
func (p Params) buildIovec() ([]unix.Iovec, error) {
	iovec := make([]unix.Iovec, len(p))
	var itr int

	for k, v := range p {
		ib, err := unix.BytePtrFromString(k)
		if err != nil {
			return nil, err
		}

		rv := reflect.ValueOf(v)
		var size uint64

		switch rv.Kind() {
		case reflect.String:
			s := string(rv.String())
			size = uint64(unsafe.Sizeof(s))
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Bool:
			size = uint64(unsafe.Sizeof(rv.UnsafeAddr()))
		default:
			return nil, errors.New("invalid value passed in for key: " + k)
		}

		iovec[itr] = unix.Iovec{
			Base: ib,
			Len:  size,
		}

		itr++
	}

	return iovec, nil
}

// Set creates	a new jail, or modifies	an existing
// one, and optionally locks the current process in it.
func Set(params Params, flags uintptr) error {
	iov, err := params.buildIovec()
	if err != nil {
		return err
	}

	return getSet(sysJailSet, iov[0], flags)
}

// Get retrieves a matching jail based on the provided params.
func Get(params Params, flags uintptr) error {
	iov, err := params.buildIovec()
	if err != nil {
		return err
	}

	return getSet(sysJailGet, iov[0], flags)
}

// getSet performas the given syscall with the params and flags provided.
func getSet(call int, iov unix.Iovec, flags uintptr) error {
	_, _, e1 := unix.Syscall(uintptr(call), uintptr(unsafe.Pointer(&iov)), 0, flags)
	if e1 != 0 {
		switch call {
		case sysJailGet:
			switch int(e1) {
			case ErrJailGetFaultOutsideOfAllocatedSpace:
				return fmt.Errorf("fault outside of allocated space: %d", e1)
			case enoent:
				return fmt.Errorf("jail referred to either does not exist or is inaccessible: %d", e1)
			case einval:
				return fmt.Errorf("invalid param provided: %d", e1)
			}
		case sysJailSet:
			switch int(e1) {
			case eperm:
				return fmt.Errorf("not allowed or restricted: %d", e1)
			case ErrJailSetFaultOutsideOfAllocatedSpace:
				return fmt.Errorf("fault outside of allocated space: %d", e1)
			case ErrJailSetParamNotExist, ErrJailSetParamWrongSize:
				return fmt.Errorf("invalid param provided: %d", e1)
			case ErrJailSetUpdateFlagNotSet:
				return fmt.Errorf("set update flag not set: %d", e1)
			case ErrJailSetNameTooLong:
				return fmt.Errorf("set name too long: %d", e1)
			case ErrJailSetNoIDsLeft:
				return fmt.Errorf("no JID's left: %d", e1)
			}
		}
	}

	return nil
}

// Attach receives a jail ID and attempts to attach the current
// process to that jail.
func Attach(jailID int32) error {
	return attachRemove(sysJailAttach, jailID)
}

// Remove receives a jail ID and attempts to remove the associated jail.
func Remove(jailID int32) error {
	return attachRemove(sysJailRemove, jailID)
}

// attachRemove
func attachRemove(call, jailID int32) error {
	jid := uintptr(unsafe.Pointer(&jailID))
	_, _, e1 := unix.Syscall(uintptr(call), jid, 0, 0)
	if e1 != 0 {
		switch int(e1) {
		case ErrJailAttachUnprivilegedUser:
			return fmt.Errorf("unprivileged user")
		case ErrjailAttachJIDNotExist:
			return fmt.Errorf("JID does not exist")
		}
	}

	return nil
}

// ip2int converts the given IP address to an uint32.
func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.LittleEndian.Uint32(ip[12:16])
	}

	return binary.LittleEndian.Uint32(ip)
}

// uint32ip converts an uint32 representation of a string into an IP.
func uint32ip(nn uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)

	return ip.String()
}
