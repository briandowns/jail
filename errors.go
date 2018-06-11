package jail

// errors in this file have been copied from https://github.com/freebsd/freebsd/blob/master/sys/sys/errno.h

const (
	eperm        = 1
	enoent       = 2
	efault       = 14
	eexist       = 17
	einval       = 22
	eagain       = 35
	enametoolong = 63
)

// The jail() system call will fail with one of the below errors
const (
	// ErrJailPermDenied [EPERM] This process is not allowed to create a jail,
	// either because it is not the super-user, or because it would exceed the
	// jail's children.max limit.
	ErrJailPermDenied = eperm

	// ErrJailFaultOutsideOfAllocatedSpace [EFAULT] jail points to an address
	// outside the allocated address space of the process.
	ErrJailFaultOutsideOfAllocatedSpace = efault

	// ErrJailInvalidVersion [EINVAL] The version number of the argument is not
	// correct.
	ErrJailInvalidVersion = einval

	// ErrjailNoFreeJIDFound [EAGAIN] No free JID could be found.
	ErrjailNoFreeJIDFound = eagain

        // ErrJailNoSuchFileDirectory [ENOENT] No such file or directory.  A component of a specified pathname
        // did not exist, or the pathname was   an empty string.
        ErrJailNoSuchFileDirectory = enoent
)

// The jail_set() system call will fail with one of the below errors
const (
	// ErrJailSetPermDenied [EPERM] This process is not allowed to create a jail,
	// either because it is not the super-user, or because it would exceed the
	// jail's children.max limit.
	ErrJailSetPermDenied = eperm

	// ErrJailSetPermRestricted [EPERM] A jail parameter was set to a less restrictive
	// value then the current environment.
	ErrJailSetPermRestricted = eperm

	// ErrJailSetFaultOutsideOfAllocatedSpace [EFAULT] Iov, or one of the addresses
	// contained within it, points to an address outside the allocated address space
	// of the process.
	ErrJailSetFaultOutsideOfAllocatedSpace = efault

	// ErrJailSetParamNotExist [ENOENT] The jail referred to by a jid or name parameter
	// does not exist, and the JAIL_CREATE flag is not set.
	ErrJailSetParamNotExist = enoent

	// ErrJailSetNotAccessibleProcInDiffJail [ENOENT] The jail referred to by a jid
	// is not accessible by the process, because the process is	in a different jail.
	ErrJailSetNotAccessibleProcInDiffJail = enoent

	// ErrJailSetUpdateFlagNotSet [EEXIST] The jail referred to by a jid or name
	// parameter exists, and the JAIL_UPDATE flag is not set.
	ErrJailSetUpdateFlagNotSet = eexist

	// Einval [EINVAL] A supplied parameter is the wrong size.
	ErrJailSetParamWrongSize = einval

	// ErrJailSetParamOutOfRange [EINVAL] A supplied parameter is out of range.
	ErrJailSetParamOutOfRange = einval

	// ErrJailSetStringNotNullTerminated [EINVAL] A supplied string parameter is
	// not null-terminated.
	ErrJailSetStringNotNullTerminated = einval

	// ErrJailSetUnknownParam [EINVAL] A supplied parameter name does not match
	// any known parameters.
	ErrJailSetUnknownParam = einval

	// ErrJailSetCreateOrUpdateNotSet [EINVAL] One of the JAIL_CREATE or JAIL_UPDATE
	// flags is not set.
	ErrJailSetCreateOrUpdateNotSet = einval

	// ErrJailSetNameTooLong [ENAMETOOLONG] A supplied string parameter is longer
	// than allowed.
	ErrJailSetNameTooLong = enametoolong

	// ErrJailSetNoIDsLeft [EAGAIN] There are no jail IDs left.
	ErrJailSetNoIDsLeft = eagain
)

// The jail_get() system call will fail with one of the below errors
const (
	// ErrJailGetFaultOutsideOfAllocatedSpace [EFAULT] Iov, or	one of the addresses
	// contained within it, points to an address outside the allocated address space
	// of the process.
	ErrJailGetFaultOutsideOfAllocatedSpace = efault

	// ErrJailGetNotExist [ENOENT] The jail referred to by jid or name parameter
	// does not exist.
	ErrJailGetNotExist = enoent

	// ErrJailGetNotAccessibleProcInDiffJail [ENOENT] The jail referred to by a
	// jid is not accessible by the process, because the process is in a different
	// jail.
	ErrJailGetNotAccessibleProcInDiffJail = enoent

	// [ENOENT] The lastjid parameter is greater than the highest current jail ID.
	ErrJailGetParamHigherThanCurJID = enoent

	// ErrJailGetParamWrongSize [EINVAL] A supplied parameter is the wrong size.
	ErrJailGetParamWrongSize = einval

	// ErrJailGetUnknownParam [EINVAL] A supplied parameter name does not match
	// any known parameters.
	ErrJailGetUnknownParam = einval
)

// The jail_attach() and jail_remove() system calls will fail with either of the
// below errors
var (
	// ErrJailAttachUnprivilegedUser [EPERM] A user other than the super-user
	// attempted to attach to or remove a jail.
	ErrJailAttachUnprivilegedUser = eperm

	// ErrjailAttachJIDNotExist [EINVAL] The jail specified by jid does not exist.
	ErrjailAttachJIDNotExist = einval
)
