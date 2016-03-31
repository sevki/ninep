// Copyright 2009 The Ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// 9P2000 message types
const (
	Tversion MType = 100 + iota
	Rversion
	Tauth
	Rauth
	Tattach
	Rattach
	Terror
	Rerror
	Tflush
	Rflush
	Twalk
	Rwalk
	Topen
	Ropen
	Tcreate
	Rcreate
	Tread
	Rread
	Twrite
	Rwrite
	Tclunk
	Rclunk
	Tremove
	Rremove
	Tstat
	Rstat
	Twstat
	Rwstat
	Tlast
)

const (
	MSIZE   = 2*1048576 + IOHDRSZ // default message size (1048576+IOHdrSz)
	IOHDRSZ = 24                  // the non-data size of the Twrite messages
	PORT    = 564                 // default port for 9P file servers
)

// Qid types
const (
	QTDIR     = 0x80 // directories
	QTAPPEND  = 0x40 // append only files
	QTEXCL    = 0x20 // exclusive use files
	QTMOUNT   = 0x10 // mounted channel
	QTAUTH    = 0x08 // authentication file
	QTTMP     = 0x04 // non-backed-up file
	QTSYMLINK = 0x02 // symbolic link (Unix, 9P2000.u)
	QTLINK    = 0x01 // hard link (Unix, 9P2000.u)
	QTFILE    = 0x00
)

// Flags for the mode field in Topen and Tcreate messages
const (
	OREAD   = 0x0    // open read-only
	OWRITE  = 0x1    // open write-only
	ORDWR   = 0x2    // open read-write
	OEXEC   = 0x3    // execute (== read but check execute permission)
	OTRUNC  = 0x10   // or'ed in (except for exec), truncate file first
	OCEXEC  = 0x20   // or'ed in, close on exec
	ORCLOSE = 0x40   // or'ed in, remove on close
	OAPPEND = 0x80   // or'ed in, append only
	OEXCL   = 0x1000 // or'ed in, exclusive client use
)

// File modes
const (
	DMDIR       = 0x80000000 // mode bit for directories
	DMAPPEND    = 0x40000000 // mode bit for append only files
	DMEXCL      = 0x20000000 // mode bit for exclusive use files
	DMMOUNT     = 0x10000000 // mode bit for mounted channel
	DMAUTH      = 0x08000000 // mode bit for authentication file
	DMTMP       = 0x04000000 // mode bit for non-backed-up file
	DMSYMLINK   = 0x02000000 // mode bit for symbolic link (Unix, 9P2000.u)
	DMLINK      = 0x01000000 // mode bit for hard link (Unix, 9P2000.u)
	DMDEVICE    = 0x00800000 // mode bit for device file (Unix, 9P2000.u)
	DMNAMEDPIPE = 0x00200000 // mode bit for named pipe (Unix, 9P2000.u)
	DMSOCKET    = 0x00100000 // mode bit for socket (Unix, 9P2000.u)
	DMSETUID    = 0x00080000 // mode bit for setuid (Unix, 9P2000.u)
	DMSETGID    = 0x00040000 // mode bit for setgid (Unix, 9P2000.u)
	DMREAD      = 0x4        // mode bit for read permission
	DMWRITE     = 0x2        // mode bit for write permission
	DMEXEC      = 0x1        // mode bit for execute permission
)

const (
	NOTAG uint16 = 0xFFFF     // no tag specified
	NOFID uint32 = 0xFFFFFFFF // no fid specified
	NOUID uint32 = 0xFFFFFFFF // no uid specified
)

// Error values
const (
	EPERM   = 1
	ENOENT  = 2
	EIO     = 5
	EACCES  = 13
	EEXIST  = 17
	ENOTDIR = 20
	EINVAL  = 22
)

// Types contained in 9p messages.
type (
	MType      uint8
	Mode       uint8
	NumEntries uint16
	Tag        uint16
	FID        uint32
	Count      int32
	Perm       int32
	Offset     uint64
	Data       []byte
)

// Error represents a 9P2000 (and 9P2000.u) error
type Error struct {
	Err      string // textual representation of the error
	Errornum uint32 // numeric representation of the error (9P2000.u)
}

// File identifier
type Qid struct {
	Type    uint8  // type of the file (high 8 bits of the mode)
	Version uint32 // version number for the path
	Path    uint64 // server's unique identification of the file
}

// Dir describes a file
type Dir struct {
	Size   uint16 // size-2 of the Dir on the wire
	Type   uint16
	Dev    uint32
	Qid           // file's Qid
	Mode   uint32 // permissions and flags
	Atime  uint32 // last access time in seconds
	Mtime  uint32 // last modified time in seconds
	Length uint64 // file length in bytes
	Name   string // file name
	Uid    string // owner name
	Gid    string // group name
	Muid   string // name of the last user that modified the file
}

// N.B. In all packets, the wire order is assumed to be the order in which you
// put struct members.

type TversionPkt struct {
	Tag     uint16
	Msize   uint32
	Version string
}
