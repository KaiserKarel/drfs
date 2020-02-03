// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/binarylog/grpc_binarylog_v1/binarylog.proto

package grpc_binarylog_v1 // import "google.golang.org/grpc/binarylog/grpc_binarylog_v1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Enumerates the type of event
// Note the terminology is different from the RPC semantics
// definition, but the same meaning is expressed here.
type GrpcLogEntry_EventType int32

const (
	GrpcLogEntry_EVENT_TYPE_UNKNOWN GrpcLogEntry_EventType = 0
	// Header sent from client to server
	GrpcLogEntry_EVENT_TYPE_CLIENT_HEADER GrpcLogEntry_EventType = 1
	// Header sent from server to client
	GrpcLogEntry_EVENT_TYPE_SERVER_HEADER GrpcLogEntry_EventType = 2
	// Message sent from client to server
	GrpcLogEntry_EVENT_TYPE_CLIENT_MESSAGE GrpcLogEntry_EventType = 3
	// Message sent from server to client
	GrpcLogEntry_EVENT_TYPE_SERVER_MESSAGE GrpcLogEntry_EventType = 4
	// A signal that client is done sending
	GrpcLogEntry_EVENT_TYPE_CLIENT_HALF_CLOSE GrpcLogEntry_EventType = 5
	// Trailer indicates the end of the RPC.
	// On client side, this event means a trailer was either received
	// from the network or the gRPC library locally generated a status
	// to inform the application about a failure.
	// On server side, this event means the server application requested
	// to send a trailer. Note: EVENT_TYPE_CANCEL may still arrive after
	// this due to races on server side.
	GrpcLogEntry_EVENT_TYPE_SERVER_TRAILER GrpcLogEntry_EventType = 6
	// A signal that the RPC is cancelled. On client side, this
	// indicates the client application requests a cancellation.
	// On server side, this indicates that cancellation was detected.
	// Note: This marks the end of the RPC. Events may arrive after
	// this due to races. For example, on client side a trailer
	// may arrive even though the application requested to cancel the RPC.
	GrpcLogEntry_EVENT_TYPE_CANCEL GrpcLogEntry_EventType = 7
)

var GrpcLogEntry_EventType_name = map[int32]string{
	0: "EVENT_TYPE_UNKNOWN",
	1: "EVENT_TYPE_CLIENT_HEADER",
	2: "EVENT_TYPE_SERVER_HEADER",
	3: "EVENT_TYPE_CLIENT_MESSAGE",
	4: "EVENT_TYPE_SERVER_MESSAGE",
	5: "EVENT_TYPE_CLIENT_HALF_CLOSE",
	6: "EVENT_TYPE_SERVER_TRAILER",
	7: "EVENT_TYPE_CANCEL",
}
var GrpcLogEntry_EventType_value = map[string]int32{
	"EVENT_TYPE_UNKNOWN":           0,
	"EVENT_TYPE_CLIENT_HEADER":     1,
	"EVENT_TYPE_SERVER_HEADER":     2,
	"EVENT_TYPE_CLIENT_MESSAGE":    3,
	"EVENT_TYPE_SERVER_MESSAGE":    4,
	"EVENT_TYPE_CLIENT_HALF_CLOSE": 5,
	"EVENT_TYPE_SERVER_TRAILER":    6,
	"EVENT_TYPE_CANCEL":            7,
}

func (x GrpcLogEntry_EventType) String() string {
	return proto.EnumName(GrpcLogEntry_EventType_name, int32(x))
}
func (GrpcLogEntry_EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{0, 0}
}

// Enumerates the entity that generates the log entry
type GrpcLogEntry_Logger int32

const (
	GrpcLogEntry_LOGGER_UNKNOWN GrpcLogEntry_Logger = 0
	GrpcLogEntry_LOGGER_CLIENT  GrpcLogEntry_Logger = 1
	GrpcLogEntry_LOGGER_SERVER  GrpcLogEntry_Logger = 2
)

var GrpcLogEntry_Logger_name = map[int32]string{
	0: "LOGGER_UNKNOWN",
	1: "LOGGER_CLIENT",
	2: "LOGGER_SERVER",
}
var GrpcLogEntry_Logger_value = map[string]int32{
	"LOGGER_UNKNOWN": 0,
	"LOGGER_CLIENT":  1,
	"LOGGER_SERVER":  2,
}

func (x GrpcLogEntry_Logger) String() string {
	return proto.EnumName(GrpcLogEntry_Logger_name, int32(x))
}
func (GrpcLogEntry_Logger) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{0, 1}
}

type Address_Type int32

const (
	Address_TYPE_UNKNOWN Address_Type = 0
	// address is in 1.2.3.4 form
	Address_TYPE_IPV4 Address_Type = 1
	// address is in IPv6 canonical form (RFC5952 section 4)
	// The scope is NOT included in the address string.
	Address_TYPE_IPV6 Address_Type = 2
	// address is UDS string
	Address_TYPE_UNIX Address_Type = 3
)

var Address_Type_name = map[int32]string{
	0: "TYPE_UNKNOWN",
	1: "TYPE_IPV4",
	2: "TYPE_IPV6",
	3: "TYPE_UNIX",
}
var Address_Type_value = map[string]int32{
	"TYPE_UNKNOWN": 0,
	"TYPE_IPV4":    1,
	"TYPE_IPV6":    2,
	"TYPE_UNIX":    3,
}

func (x Address_Type) String() string {
	return proto.EnumName(Address_Type_name, int32(x))
}
func (Address_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{7, 0}
}

// Log entry we store in binary logs
type GrpcLogEntry struct {
	// The timestamp of the binary log message
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Uniquely identifies a call. The value must not be 0 in order to disambiguate
	// from an unset value.
	// Each call may have several log entries, they will all have the same call_id.
	// Nothing is guaranteed about their value other than they are unique across
	// different RPCs in the same gRPC process.
	CallId uint64 `protobuf:"varint,2,opt,name=call_id,json=callId,proto3" json:"call_id,omitempty"`
	// The entry sequence id for this call. The first GrpcLogEntry has a
	// value of 1, to disambiguate from an unset value. The purpose of
	// this field is to detect missing entries in environments where
	// durability or ordering is not guaranteed.
	SequenceIdWithinCall uint64                 `protobuf:"varint,3,opt,name=sequence_id_within_call,json=sequenceIdWithinCall,proto3" json:"sequence_id_within_call,omitempty"`
	Type                 GrpcLogEntry_EventType `protobuf:"varint,4,opt,name=type,proto3,enum=grpc.binarylog.v1.GrpcLogEntry_EventType" json:"type,omitempty"`
	Logger               GrpcLogEntry_Logger    `protobuf:"varint,5,opt,name=logger,proto3,enum=grpc.binarylog.v1.GrpcLogEntry_Logger" json:"logger,omitempty"`
	// The logger uses one of the following fields to record the payload,
	// according to the type of the log entry.
	//
	// Types that are valid to be assigned to Payload:
	//	*GrpcLogEntry_ClientHeader
	//	*GrpcLogEntry_ServerHeader
	//	*GrpcLogEntry_Message
	//	*GrpcLogEntry_Trailer
	Payload isGrpcLogEntry_Payload `protobuf_oneof:"payload"`
	// true if payload does not represent the full message or metadata.
	PayloadTruncated bool `protobuf:"varint,10,opt,name=payload_truncated,json=payloadTruncated,proto3" json:"payload_truncated,omitempty"`
	// Peer address information, will only be recorded on the first
	// incoming event. On client side, peer is logged on
	// EVENT_TYPE_SERVER_HEADER normally or EVENT_TYPE_SERVER_TRAILER in
	// the case of trailers-only. On server side, peer is always
	// logged on EVENT_TYPE_CLIENT_HEADER.
	Peer                 *Address `protobuf:"bytes,11,opt,name=peer,proto3" json:"peer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrpcLogEntry) Reset()         { *m = GrpcLogEntry{} }
func (m *GrpcLogEntry) String() string { return proto.CompactTextString(m) }
func (*GrpcLogEntry) ProtoMessage()    {}
func (*GrpcLogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{0}
}
func (m *GrpcLogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GrpcLogEntry.Unmarshal(m, b)
}
func (m *GrpcLogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GrpcLogEntry.Marshal(b, m, deterministic)
}
func (dst *GrpcLogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrpcLogEntry.Merge(dst, src)
}
func (m *GrpcLogEntry) XXX_Size() int {
	return xxx_messageInfo_GrpcLogEntry.Size(m)
}
func (m *GrpcLogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_GrpcLogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_GrpcLogEntry proto.InternalMessageInfo

func (m *GrpcLogEntry) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *GrpcLogEntry) GetCallId() uint64 {
	if m != nil {
		return m.CallId
	}
	return 0
}

func (m *GrpcLogEntry) GetSequenceIdWithinCall() uint64 {
	if m != nil {
		return m.SequenceIdWithinCall
	}
	return 0
}

func (m *GrpcLogEntry) GetType() GrpcLogEntry_EventType {
	if m != nil {
		return m.Type
	}
	return GrpcLogEntry_EVENT_TYPE_UNKNOWN
}

func (m *GrpcLogEntry) GetLogger() GrpcLogEntry_Logger {
	if m != nil {
		return m.Logger
	}
	return GrpcLogEntry_LOGGER_UNKNOWN
}

type isGrpcLogEntry_Payload interface {
	isGrpcLogEntry_Payload()
}

type GrpcLogEntry_ClientHeader struct {
	ClientHeader *ClientHeader `protobuf:"bytes,6,opt,name=client_header,json=clientHeader,proto3,oneof"`
}

type GrpcLogEntry_ServerHeader struct {
	ServerHeader *ServerHeader `protobuf:"bytes,7,opt,name=server_header,json=serverHeader,proto3,oneof"`
}

type GrpcLogEntry_Message struct {
	Message *Message `protobuf:"bytes,8,opt,name=message,proto3,oneof"`
}

type GrpcLogEntry_Trailer struct {
	Trailer *Trailer `protobuf:"bytes,9,opt,name=trailer,proto3,oneof"`
}

func (*GrpcLogEntry_ClientHeader) isGrpcLogEntry_Payload() {}

func (*GrpcLogEntry_ServerHeader) isGrpcLogEntry_Payload() {}

func (*GrpcLogEntry_Message) isGrpcLogEntry_Payload() {}

func (*GrpcLogEntry_Trailer) isGrpcLogEntry_Payload() {}

func (m *GrpcLogEntry) GetPayload() isGrpcLogEntry_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *GrpcLogEntry) GetClientHeader() *ClientHeader {
	if x, ok := m.GetPayload().(*GrpcLogEntry_ClientHeader); ok {
		return x.ClientHeader
	}
	return nil
}

func (m *GrpcLogEntry) GetServerHeader() *ServerHeader {
	if x, ok := m.GetPayload().(*GrpcLogEntry_ServerHeader); ok {
		return x.ServerHeader
	}
	return nil
}

func (m *GrpcLogEntry) GetMessage() *Message {
	if x, ok := m.GetPayload().(*GrpcLogEntry_Message); ok {
		return x.Message
	}
	return nil
}

func (m *GrpcLogEntry) GetTrailer() *Trailer {
	if x, ok := m.GetPayload().(*GrpcLogEntry_Trailer); ok {
		return x.Trailer
	}
	return nil
}

func (m *GrpcLogEntry) GetPayloadTruncated() bool {
	if m != nil {
		return m.PayloadTruncated
	}
	return false
}

func (m *GrpcLogEntry) GetPeer() *Address {
	if m != nil {
		return m.Peer
	}
	return nil
}

// XXX_OneofFuncs is for the lib use of the proto package.
func (*GrpcLogEntry) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GrpcLogEntry_OneofMarshaler, _GrpcLogEntry_OneofUnmarshaler, _GrpcLogEntry_OneofSizer, []interface{}{
		(*GrpcLogEntry_ClientHeader)(nil),
		(*GrpcLogEntry_ServerHeader)(nil),
		(*GrpcLogEntry_Message)(nil),
		(*GrpcLogEntry_Trailer)(nil),
	}
}

func _GrpcLogEntry_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GrpcLogEntry)
	// payload
	switch x := m.Payload.(type) {
	case *GrpcLogEntry_ClientHeader:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ClientHeader); err != nil {
			return err
		}
	case *GrpcLogEntry_ServerHeader:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ServerHeader); err != nil {
			return err
		}
	case *GrpcLogEntry_Message:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Message); err != nil {
			return err
		}
	case *GrpcLogEntry_Trailer:
		b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Trailer); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GrpcLogEntry.Payload has unexpected type %T", x)
	}
	return nil
}

func _GrpcLogEntry_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GrpcLogEntry)
	switch tag {
	case 6: // payload.client_header
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ClientHeader)
		err := b.DecodeMessage(msg)
		m.Payload = &GrpcLogEntry_ClientHeader{msg}
		return true, err
	case 7: // payload.server_header
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ServerHeader)
		err := b.DecodeMessage(msg)
		m.Payload = &GrpcLogEntry_ServerHeader{msg}
		return true, err
	case 8: // payload.message
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message)
		err := b.DecodeMessage(msg)
		m.Payload = &GrpcLogEntry_Message{msg}
		return true, err
	case 9: // payload.trailer
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Trailer)
		err := b.DecodeMessage(msg)
		m.Payload = &GrpcLogEntry_Trailer{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GrpcLogEntry_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GrpcLogEntry)
	// payload
	switch x := m.Payload.(type) {
	case *GrpcLogEntry_ClientHeader:
		s := proto.Size(x.ClientHeader)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GrpcLogEntry_ServerHeader:
		s := proto.Size(x.ServerHeader)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GrpcLogEntry_Message:
		s := proto.Size(x.Message)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GrpcLogEntry_Trailer:
		s := proto.Size(x.Trailer)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ClientHeader struct {
	// This contains only the metadata from the application.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The name of the RPC method, which looks something like:
	// /<service>/<method>
	// Note the leading "/" character.
	MethodName string `protobuf:"bytes,2,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	// A single process may be used to run multiple virtual
	// servers with different identities.
	// The authority is the name of such a server identitiy.
	// It is typically a portion of the URI in the form of
	// <host> or <host>:<port> .
	Authority string `protobuf:"bytes,3,opt,name=authority,proto3" json:"authority,omitempty"`
	// the RPC timeout
	Timeout              *duration.Duration `protobuf:"bytes,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ClientHeader) Reset()         { *m = ClientHeader{} }
func (m *ClientHeader) String() string { return proto.CompactTextString(m) }
func (*ClientHeader) ProtoMessage()    {}
func (*ClientHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{1}
}
func (m *ClientHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientHeader.Unmarshal(m, b)
}
func (m *ClientHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientHeader.Marshal(b, m, deterministic)
}
func (dst *ClientHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientHeader.Merge(dst, src)
}
func (m *ClientHeader) XXX_Size() int {
	return xxx_messageInfo_ClientHeader.Size(m)
}
func (m *ClientHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ClientHeader proto.InternalMessageInfo

func (m *ClientHeader) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *ClientHeader) GetMethodName() string {
	if m != nil {
		return m.MethodName
	}
	return ""
}

func (m *ClientHeader) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *ClientHeader) GetTimeout() *duration.Duration {
	if m != nil {
		return m.Timeout
	}
	return nil
}

type ServerHeader struct {
	// This contains only the metadata from the application.
	Metadata             *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ServerHeader) Reset()         { *m = ServerHeader{} }
func (m *ServerHeader) String() string { return proto.CompactTextString(m) }
func (*ServerHeader) ProtoMessage()    {}
func (*ServerHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{2}
}
func (m *ServerHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerHeader.Unmarshal(m, b)
}
func (m *ServerHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerHeader.Marshal(b, m, deterministic)
}
func (dst *ServerHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerHeader.Merge(dst, src)
}
func (m *ServerHeader) XXX_Size() int {
	return xxx_messageInfo_ServerHeader.Size(m)
}
func (m *ServerHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerHeader.DiscardUnknown(m)
}

var xxx_messageInfo_ServerHeader proto.InternalMessageInfo

func (m *ServerHeader) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Trailer struct {
	// This contains only the metadata from the application.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The gRPC status code.
	StatusCode uint32 `protobuf:"varint,2,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	// An original status message before any transport specific
	// encoding.
	StatusMessage string `protobuf:"bytes,3,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
	// The value of the 'grpc-status-details-bin' metadata key. If
	// present, this is always an encoded 'google.rpc.Status' message.
	StatusDetails        []byte   `protobuf:"bytes,4,opt,name=status_details,json=statusDetails,proto3" json:"status_details,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Trailer) Reset()         { *m = Trailer{} }
func (m *Trailer) String() string { return proto.CompactTextString(m) }
func (*Trailer) ProtoMessage()    {}
func (*Trailer) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{3}
}
func (m *Trailer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Trailer.Unmarshal(m, b)
}
func (m *Trailer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Trailer.Marshal(b, m, deterministic)
}
func (dst *Trailer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trailer.Merge(dst, src)
}
func (m *Trailer) XXX_Size() int {
	return xxx_messageInfo_Trailer.Size(m)
}
func (m *Trailer) XXX_DiscardUnknown() {
	xxx_messageInfo_Trailer.DiscardUnknown(m)
}

var xxx_messageInfo_Trailer proto.InternalMessageInfo

func (m *Trailer) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Trailer) GetStatusCode() uint32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *Trailer) GetStatusMessage() string {
	if m != nil {
		return m.StatusMessage
	}
	return ""
}

func (m *Trailer) GetStatusDetails() []byte {
	if m != nil {
		return m.StatusDetails
	}
	return nil
}

// Message payload, used by CLIENT_MESSAGE and SERVER_MESSAGE
type Message struct {
	// Length of the message. It may not be the same as the length of the
	// data field, as the logging payload can be truncated or omitted.
	Length uint32 `protobuf:"varint,1,opt,name=length,proto3" json:"length,omitempty"`
	// May be truncated or omitted.
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{4}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetLength() uint32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// A list of metadata pairs, used in the payload of client header,
// server header, and server trailer.
// Implementations may omit some entries to honor the header limits
// of GRPC_BINARY_LOG_CONFIG.
//
// Header keys added by gRPC are omitted. To be more specific,
// implementations will not log the following entries, and this is
// not to be treated as a truncation:
// - entries handled by grpc that are not user visible, such as those
//   that begin with 'grpc-' (with exception of grpc-trace-bin)
//   or keys like 'lb-token'
// - transport specific entries, including but not limited to:
//   ':path', ':authority', 'content-encoding', 'user-agent', 'te', etc
// - entries added for call credentials
//
// Implementations must always log grpc-trace-bin if it is present.
// Practically speaking it will only be visible on server side because
// grpc-trace-bin is managed by low level client side mechanisms
// inaccessible from the application level. On server side, the
// header is just a normal metadata key.
// The pair will not count towards the size limit.
type Metadata struct {
	Entry                []*MetadataEntry `protobuf:"bytes,1,rep,name=entry,proto3" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{5}
}
func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (dst *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(dst, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetEntry() []*MetadataEntry {
	if m != nil {
		return m.Entry
	}
	return nil
}

// A metadata key value pair
type MetadataEntry struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetadataEntry) Reset()         { *m = MetadataEntry{} }
func (m *MetadataEntry) String() string { return proto.CompactTextString(m) }
func (*MetadataEntry) ProtoMessage()    {}
func (*MetadataEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{6}
}
func (m *MetadataEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataEntry.Unmarshal(m, b)
}
func (m *MetadataEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataEntry.Marshal(b, m, deterministic)
}
func (dst *MetadataEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataEntry.Merge(dst, src)
}
func (m *MetadataEntry) XXX_Size() int {
	return xxx_messageInfo_MetadataEntry.Size(m)
}
func (m *MetadataEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataEntry.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataEntry proto.InternalMessageInfo

func (m *MetadataEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *MetadataEntry) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// Address information
type Address struct {
	Type    Address_Type `protobuf:"varint,1,opt,name=type,proto3,enum=grpc.binarylog.v1.Address_Type" json:"type,omitempty"`
	Address string       `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// only for TYPE_IPV4 and TYPE_IPV6
	IpPort               uint32   `protobuf:"varint,3,opt,name=ip_port,json=ipPort,proto3" json:"ip_port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_binarylog_264c8c9c551ce911, []int{7}
}
func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (dst *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(dst, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetType() Address_Type {
	if m != nil {
		return m.Type
	}
	return Address_TYPE_UNKNOWN
}

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Address) GetIpPort() uint32 {
	if m != nil {
		return m.IpPort
	}
	return 0
}

func init() {
	proto.RegisterType((*GrpcLogEntry)(nil), "grpc.binarylog.v1.GrpcLogEntry")
	proto.RegisterType((*ClientHeader)(nil), "grpc.binarylog.v1.ClientHeader")
	proto.RegisterType((*ServerHeader)(nil), "grpc.binarylog.v1.ServerHeader")
	proto.RegisterType((*Trailer)(nil), "grpc.binarylog.v1.Trailer")
	proto.RegisterType((*Message)(nil), "grpc.binarylog.v1.Message")
	proto.RegisterType((*Metadata)(nil), "grpc.binarylog.v1.Metadata")
	proto.RegisterType((*MetadataEntry)(nil), "grpc.binarylog.v1.MetadataEntry")
	proto.RegisterType((*Address)(nil), "grpc.binarylog.v1.Address")
	proto.RegisterEnum("grpc.binarylog.v1.GrpcLogEntry_EventType", GrpcLogEntry_EventType_name, GrpcLogEntry_EventType_value)
	proto.RegisterEnum("grpc.binarylog.v1.GrpcLogEntry_Logger", GrpcLogEntry_Logger_name, GrpcLogEntry_Logger_value)
	proto.RegisterEnum("grpc.binarylog.v1.Address_Type", Address_Type_name, Address_Type_value)
}

func init() {
	proto.RegisterFile("grpc/binarylog/grpc_binarylog_v1/binarylog.proto", fileDescriptor_binarylog_264c8c9c551ce911)
}

var fileDescriptor_binarylog_264c8c9c551ce911 = []byte{
	// 900 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0x51, 0x6f, 0xe3, 0x44,
	0x10, 0x3e, 0x37, 0x69, 0xdc, 0x4c, 0x92, 0xca, 0x5d, 0x95, 0x3b, 0x5f, 0x29, 0x34, 0xb2, 0x04,
	0x0a, 0x42, 0x72, 0xb9, 0x94, 0xeb, 0xf1, 0x02, 0x52, 0x92, 0xfa, 0xd2, 0x88, 0x5c, 0x1a, 0x6d,
	0x72, 0x3d, 0x40, 0x48, 0xd6, 0x36, 0x5e, 0x1c, 0x0b, 0xc7, 0x6b, 0xd6, 0x9b, 0xa0, 0xfc, 0x2c,
	0xde, 0x90, 0xee, 0x77, 0xf1, 0x8e, 0xbc, 0x6b, 0x27, 0xa6, 0x69, 0x0f, 0x09, 0xde, 0x3c, 0xdf,
	0x7c, 0xf3, 0xcd, 0xee, 0x78, 0x66, 0x16, 0xbe, 0xf2, 0x79, 0x3c, 0x3b, 0xbf, 0x0b, 0x22, 0xc2,
	0xd7, 0x21, 0xf3, 0xcf, 0x53, 0xd3, 0xdd, 0x98, 0xee, 0xea, 0xc5, 0xd6, 0x67, 0xc7, 0x9c, 0x09,
	0x86, 0x8e, 0x52, 0x8a, 0xbd, 0x45, 0x57, 0x2f, 0x4e, 0x3e, 0xf5, 0x19, 0xf3, 0x43, 0x7a, 0x2e,
	0x09, 0x77, 0xcb, 0x5f, 0xce, 0xbd, 0x25, 0x27, 0x22, 0x60, 0x91, 0x0a, 0x39, 0x39, 0xbb, 0xef,
	0x17, 0xc1, 0x82, 0x26, 0x82, 0x2c, 0x62, 0x45, 0xb0, 0xde, 0xeb, 0x50, 0xef, 0xf3, 0x78, 0x36,
	0x64, 0xbe, 0x13, 0x09, 0xbe, 0x46, 0xdf, 0x40, 0x75, 0xc3, 0x31, 0xb5, 0xa6, 0xd6, 0xaa, 0xb5,
	0x4f, 0x6c, 0xa5, 0x62, 0xe7, 0x2a, 0xf6, 0x34, 0x67, 0xe0, 0x2d, 0x19, 0x3d, 0x03, 0x7d, 0x46,
	0xc2, 0xd0, 0x0d, 0x3c, 0x73, 0xaf, 0xa9, 0xb5, 0xca, 0xb8, 0x92, 0x9a, 0x03, 0x0f, 0xbd, 0x84,
	0x67, 0x09, 0xfd, 0x6d, 0x49, 0xa3, 0x19, 0x75, 0x03, 0xcf, 0xfd, 0x3d, 0x10, 0xf3, 0x20, 0x72,
	0x53, 0xa7, 0x59, 0x92, 0xc4, 0xe3, 0xdc, 0x3d, 0xf0, 0xde, 0x49, 0x67, 0x8f, 0x84, 0x21, 0xfa,
	0x16, 0xca, 0x62, 0x1d, 0x53, 0xb3, 0xdc, 0xd4, 0x5a, 0x87, 0xed, 0x2f, 0xec, 0x9d, 0xdb, 0xdb,
	0xc5, 0x83, 0xdb, 0xce, 0x8a, 0x46, 0x62, 0xba, 0x8e, 0x29, 0x96, 0x61, 0xe8, 0x3b, 0xa8, 0x84,
	0xcc, 0xf7, 0x29, 0x37, 0xf7, 0xa5, 0xc0, 0xe7, 0xff, 0x26, 0x30, 0x94, 0x6c, 0x9c, 0x45, 0xa1,
	0xd7, 0xd0, 0x98, 0x85, 0x01, 0x8d, 0x84, 0x3b, 0xa7, 0xc4, 0xa3, 0xdc, 0xac, 0xc8, 0x62, 0x9c,
	0x3d, 0x20, 0xd3, 0x93, 0xbc, 0x6b, 0x49, 0xbb, 0x7e, 0x82, 0xeb, 0xb3, 0x82, 0x9d, 0xea, 0x24,
	0x94, 0xaf, 0x28, 0xcf, 0x75, 0xf4, 0x47, 0x75, 0x26, 0x92, 0xb7, 0xd5, 0x49, 0x0a, 0x36, 0xba,
	0x04, 0x7d, 0x41, 0x93, 0x84, 0xf8, 0xd4, 0x3c, 0xc8, 0x7f, 0xcb, 0x8e, 0xc2, 0x1b, 0xc5, 0xb8,
	0x7e, 0x82, 0x73, 0x72, 0x1a, 0x27, 0x38, 0x09, 0x42, 0xca, 0xcd, 0xea, 0xa3, 0x71, 0x53, 0xc5,
	0x48, 0xe3, 0x32, 0x32, 0xfa, 0x12, 0x8e, 0x62, 0xb2, 0x0e, 0x19, 0xf1, 0x5c, 0xc1, 0x97, 0xd1,
	0x8c, 0x08, 0xea, 0x99, 0xd0, 0xd4, 0x5a, 0x07, 0xd8, 0xc8, 0x1c, 0xd3, 0x1c, 0x47, 0x36, 0x94,
	0x63, 0x4a, 0xb9, 0x59, 0x7b, 0x34, 0x43, 0xc7, 0xf3, 0x38, 0x4d, 0x12, 0x2c, 0x79, 0xd6, 0x5f,
	0x1a, 0x54, 0x37, 0x3f, 0x0c, 0x3d, 0x05, 0xe4, 0xdc, 0x3a, 0xa3, 0xa9, 0x3b, 0xfd, 0x71, 0xec,
	0xb8, 0x6f, 0x47, 0xdf, 0x8f, 0x6e, 0xde, 0x8d, 0x8c, 0x27, 0xe8, 0x14, 0xcc, 0x02, 0xde, 0x1b,
	0x0e, 0xd2, 0xef, 0x6b, 0xa7, 0x73, 0xe5, 0x60, 0x43, 0xbb, 0xe7, 0x9d, 0x38, 0xf8, 0xd6, 0xc1,
	0xb9, 0x77, 0x0f, 0x7d, 0x02, 0xcf, 0x77, 0x63, 0xdf, 0x38, 0x93, 0x49, 0xa7, 0xef, 0x18, 0xa5,
	0x7b, 0xee, 0x2c, 0x38, 0x77, 0x97, 0x51, 0x13, 0x4e, 0x1f, 0xc8, 0xdc, 0x19, 0xbe, 0x76, 0x7b,
	0xc3, 0x9b, 0x89, 0x63, 0xec, 0x3f, 0x2c, 0x30, 0xc5, 0x9d, 0xc1, 0xd0, 0xc1, 0x46, 0x05, 0x7d,
	0x04, 0x47, 0x45, 0x81, 0xce, 0xa8, 0xe7, 0x0c, 0x0d, 0xdd, 0xea, 0x42, 0x45, 0xb5, 0x19, 0x42,
	0x70, 0x38, 0xbc, 0xe9, 0xf7, 0x1d, 0x5c, 0xb8, 0xef, 0x11, 0x34, 0x32, 0x4c, 0x65, 0x34, 0xb4,
	0x02, 0xa4, 0x52, 0x18, 0x7b, 0xdd, 0x2a, 0xe8, 0x59, 0xfd, 0xad, 0xf7, 0x1a, 0xd4, 0x8b, 0xcd,
	0x87, 0x5e, 0xc1, 0xc1, 0x82, 0x0a, 0xe2, 0x11, 0x41, 0xb2, 0xe1, 0xfd, 0xf8, 0xc1, 0x2e, 0x51,
	0x14, 0xbc, 0x21, 0xa3, 0x33, 0xa8, 0x2d, 0xa8, 0x98, 0x33, 0xcf, 0x8d, 0xc8, 0x82, 0xca, 0x01,
	0xae, 0x62, 0x50, 0xd0, 0x88, 0x2c, 0x28, 0x3a, 0x85, 0x2a, 0x59, 0x8a, 0x39, 0xe3, 0x81, 0x58,
	0xcb, 0xb1, 0xad, 0xe2, 0x2d, 0x80, 0x2e, 0x40, 0x4f, 0x17, 0x01, 0x5b, 0x0a, 0x39, 0xae, 0xb5,
	0xf6, 0xf3, 0x9d, 0x9d, 0x71, 0x95, 0x6d, 0x26, 0x9c, 0x33, 0xad, 0x3e, 0xd4, 0x8b, 0x1d, 0xff,
	0x9f, 0x0f, 0x6f, 0xfd, 0xa1, 0x81, 0x9e, 0x75, 0xf0, 0xff, 0xaa, 0x40, 0x22, 0x88, 0x58, 0x26,
	0xee, 0x8c, 0x79, 0xaa, 0x02, 0x0d, 0x0c, 0x0a, 0xea, 0x31, 0x8f, 0xa2, 0xcf, 0xe0, 0x30, 0x23,
	0xe4, 0x73, 0xa8, 0xca, 0xd0, 0x50, 0x68, 0x36, 0x7a, 0x05, 0x9a, 0x47, 0x05, 0x09, 0xc2, 0x44,
	0x56, 0xa4, 0x9e, 0xd3, 0xae, 0x14, 0x68, 0xbd, 0x04, 0x3d, 0x8f, 0x78, 0x0a, 0x95, 0x90, 0x46,
	0xbe, 0x98, 0xcb, 0x03, 0x37, 0x70, 0x66, 0x21, 0x04, 0x65, 0x79, 0x8d, 0x3d, 0x19, 0x2f, 0xbf,
	0xad, 0x2e, 0x1c, 0xe4, 0x67, 0x47, 0x97, 0xb0, 0x4f, 0xd3, 0xcd, 0x65, 0x6a, 0xcd, 0x52, 0xab,
	0xd6, 0x6e, 0x7e, 0xe0, 0x9e, 0x72, 0xc3, 0x61, 0x45, 0xb7, 0x5e, 0x41, 0xe3, 0x1f, 0x38, 0x32,
	0xa0, 0xf4, 0x2b, 0x5d, 0xcb, 0xec, 0x55, 0x9c, 0x7e, 0xa2, 0x63, 0xd8, 0x5f, 0x91, 0x70, 0x49,
	0xb3, 0xdc, 0xca, 0xb0, 0xfe, 0xd4, 0x40, 0xcf, 0xe6, 0x18, 0x5d, 0x64, 0xdb, 0x59, 0x93, 0xcb,
	0xf5, 0xec, 0xf1, 0x89, 0xb7, 0x0b, 0x3b, 0xd9, 0x04, 0x9d, 0x28, 0x34, 0xeb, 0xb0, 0xdc, 0x4c,
	0x1f, 0x8f, 0x20, 0x76, 0x63, 0xc6, 0x85, 0xac, 0x6a, 0x03, 0x57, 0x82, 0x78, 0xcc, 0xb8, 0xb0,
	0x1c, 0x28, 0xcb, 0x1d, 0x61, 0x40, 0xfd, 0xde, 0x76, 0x68, 0x40, 0x55, 0x22, 0x83, 0xf1, 0xed,
	0xd7, 0x86, 0x56, 0x34, 0x2f, 0x8d, 0xbd, 0x8d, 0xf9, 0x76, 0x34, 0xf8, 0xc1, 0x28, 0x75, 0x7f,
	0x86, 0xe3, 0x80, 0xed, 0x1e, 0xb2, 0x7b, 0xd8, 0x95, 0xd6, 0x90, 0xf9, 0xe3, 0xb4, 0x51, 0xc7,
	0xda, 0x4f, 0xed, 0xac, 0x71, 0x7d, 0x16, 0x92, 0xc8, 0xb7, 0x19, 0x57, 0x4f, 0xf3, 0x87, 0x5e,
	0xea, 0xbb, 0x8a, 0xec, 0xf2, 0x8b, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xe7, 0xf6, 0x4b, 0x50,
	0xd4, 0x07, 0x00, 0x00,
}