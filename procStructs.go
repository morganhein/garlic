package garlic

import (
	"encoding/binary"

	"github.com/mdlayher/netlink/nlenc"
)

//Preserve a lot of the const/enums from cn_proc.h for the sake of documentation
//No one wants see 0x1 everywhere

const ( //proc_cn_mcast_op
	//ProcCnMcastListen registers a listen event with the kernel
	ProcCnMcastListen = iota + 1
	//ProcCnMcastIgnore registers an ignore event with the kernel
	ProcCnMcastIgnore = iota
)

//CnIdxProc is the Id used for proc/connector
const CnIdxProc = 0x1

//CnValProc is the corrisponding value used by chID
const CnValProc = 0x1

//Various message structs from connector.h

/*
 * idx and val are unique identifiers which
 * are used for message routing and
 * must be registered in connector.h for in-kernel usage.
 */
type cbID struct {
	Idx uint32
	Val uint32
}

//The connector message struct
type cnMsg struct {
	ID    cbID
	Seq   uint32
	Ack   uint32
	Len   uint16
	Flags uint16
}

var cnMsgLen = binary.Size(cnMsg{})

//MarshallBinary converts the entire struct into a slice, along with the proc_cn_mcast_op body
func (hdr cnMsg) MarshalBinaryAndBody(body uint32) []byte {

	bytes := make([]byte, binary.Size(hdr)+binary.Size(body))
	nlenc.PutUint32(bytes[0:4], hdr.ID.Idx)
	nlenc.PutUint32(bytes[4:8], hdr.ID.Val)
	nlenc.PutUint32(bytes[8:12], hdr.Seq)
	nlenc.PutUint32(bytes[12:16], hdr.Ack)
	nlenc.PutUint16(bytes[16:18], hdr.Len)
	nlenc.PutUint16(bytes[18:20], hdr.Flags)
	nlenc.PutUint32(bytes[20:24], body)

	return bytes
}

//converts a binary blob to a msg header
func unmarshalCnMsg(data []byte) cnMsg {

	hdr := cnMsg{}
	hdr.ID.Idx = nlenc.Uint32(data[0:4])
	hdr.ID.Val = nlenc.Uint32(data[4:8])
	hdr.Seq = nlenc.Uint32(data[8:12])
	hdr.Ack = nlenc.Uint32(data[12:16])
	hdr.Len = nlenc.Uint16(data[16:18])
	hdr.Flags = nlenc.Uint16(data[18:20])

	return hdr
}

//This is just an internal  header that allows us to easily cast the raw binary data
type procEventHdr struct {
	What      uint32
	CPU       uint32
	Timestamp uint64
}

var procEventHdrLen = binary.Size(procEventHdr{})

//unmarshal the proc hdr
func unmarshalProcEventHdr(data []byte) procEventHdr {
	hdr := procEventHdr{}

	hdr.What = nlenc.Uint32(data[0:4])
	hdr.CPU = nlenc.Uint32(data[4:8])
	hdr.Timestamp = nlenc.Uint64(data[8:16])

	return hdr
}

//from cn_proc.h
const (

	//ProcEventNone is only used for ACK events
	ProcEventNone = 0x00000000
	//ProcEventFork is a fork event
	ProcEventFork = 0x00000001
	//ProcEventExec is a exec() event
	ProcEventExec = 0x00000002
	//ProcEventUID is a user ID change
	ProcEventUID = 0x00000004
	//ProcEventGID is a group ID change
	ProcEventGID = 0x00000040
	//ProcEventSID is a session ID change
	ProcEventSID = 0x00000080
	//ProcEventSID is a process trace event
	ProcEventPtrace = 0x00000100
	//ProcEventComm is a comm(and) value change. Any value over 16 bytes will be truncated
	ProcEventComm = 0x00000200
	//ProcEventCoredump is a core dump event
	ProcEventCoredump = 0x40000000
	//ProcEventExit is an exit() event
	ProcEventExit = 0x80000000
)

//ProcEvent is the struct representing all the event data.
type ProcEvent struct {
	What        uint32
	CPU         uint32
	TimestampNs uint64
	EventData   EventData
}

//EventData is an interface that encapsulates the union type used in cn_proc
type EventData interface {
	Pid() uint32

	Tgid() uint32
}
