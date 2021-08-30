// Code generated by protoc-gen-go. DO NOT EDIT.
// source: driver.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	driver.proto

It has these top-level messages:
	InstanceMeta
	InitRequest
	Empty
	ExecRequest
	ExecResponse
	TxRequest
	TxResponse
	DatabasesResponse
	ParseRequest
	Node
	ParseResponse
	Rule
	AuditRequest
	AuditResult
	AuditResponse
	GenRollbackSQLRequest
	GenRollbackSQLResponse
	MetasResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type InstanceMeta struct {
	InstanceHost string `protobuf:"bytes,1,opt,name=instanceHost" json:"instanceHost,omitempty"`
	InstancePort string `protobuf:"bytes,2,opt,name=instancePort" json:"instancePort,omitempty"`
	InstanceUser string `protobuf:"bytes,3,opt,name=instanceUser" json:"instanceUser,omitempty"`
	InstancePass string `protobuf:"bytes,4,opt,name=instancePass" json:"instancePass,omitempty"`
	DatabaseOpen string `protobuf:"bytes,5,opt,name=databaseOpen" json:"databaseOpen,omitempty"`
}

func (m *InstanceMeta) Reset()                    { *m = InstanceMeta{} }
func (m *InstanceMeta) String() string            { return proto1.CompactTextString(m) }
func (*InstanceMeta) ProtoMessage()               {}
func (*InstanceMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *InstanceMeta) GetInstanceHost() string {
	if m != nil {
		return m.InstanceHost
	}
	return ""
}

func (m *InstanceMeta) GetInstancePort() string {
	if m != nil {
		return m.InstancePort
	}
	return ""
}

func (m *InstanceMeta) GetInstanceUser() string {
	if m != nil {
		return m.InstanceUser
	}
	return ""
}

func (m *InstanceMeta) GetInstancePass() string {
	if m != nil {
		return m.InstancePass
	}
	return ""
}

func (m *InstanceMeta) GetDatabaseOpen() string {
	if m != nil {
		return m.DatabaseOpen
	}
	return ""
}

type InitRequest struct {
	InstanceMeta *InstanceMeta `protobuf:"bytes,1,opt,name=instanceMeta" json:"instanceMeta,omitempty"`
	IsOffline    bool          `protobuf:"varint,2,opt,name=isOffline" json:"isOffline,omitempty"`
}

func (m *InitRequest) Reset()                    { *m = InitRequest{} }
func (m *InitRequest) String() string            { return proto1.CompactTextString(m) }
func (*InitRequest) ProtoMessage()               {}
func (*InitRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InitRequest) GetInstanceMeta() *InstanceMeta {
	if m != nil {
		return m.InstanceMeta
	}
	return nil
}

func (m *InitRequest) GetIsOffline() bool {
	if m != nil {
		return m.IsOffline
	}
	return false
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type ExecRequest struct {
	Query string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
}

func (m *ExecRequest) Reset()                    { *m = ExecRequest{} }
func (m *ExecRequest) String() string            { return proto1.CompactTextString(m) }
func (*ExecRequest) ProtoMessage()               {}
func (*ExecRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ExecRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type ExecResponse struct {
	LastInsertId      int64  `protobuf:"varint,1,opt,name=lastInsertId" json:"lastInsertId,omitempty"`
	LastInsertIdError string `protobuf:"bytes,2,opt,name=lastInsertIdError" json:"lastInsertIdError,omitempty"`
	RowsAffected      int64  `protobuf:"varint,3,opt,name=rowsAffected" json:"rowsAffected,omitempty"`
	RowsAffectedError string `protobuf:"bytes,4,opt,name=rowsAffectedError" json:"rowsAffectedError,omitempty"`
}

func (m *ExecResponse) Reset()                    { *m = ExecResponse{} }
func (m *ExecResponse) String() string            { return proto1.CompactTextString(m) }
func (*ExecResponse) ProtoMessage()               {}
func (*ExecResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ExecResponse) GetLastInsertId() int64 {
	if m != nil {
		return m.LastInsertId
	}
	return 0
}

func (m *ExecResponse) GetLastInsertIdError() string {
	if m != nil {
		return m.LastInsertIdError
	}
	return ""
}

func (m *ExecResponse) GetRowsAffected() int64 {
	if m != nil {
		return m.RowsAffected
	}
	return 0
}

func (m *ExecResponse) GetRowsAffectedError() string {
	if m != nil {
		return m.RowsAffectedError
	}
	return ""
}

type TxRequest struct {
	Queries []string `protobuf:"bytes,1,rep,name=queries" json:"queries,omitempty"`
}

func (m *TxRequest) Reset()                    { *m = TxRequest{} }
func (m *TxRequest) String() string            { return proto1.CompactTextString(m) }
func (*TxRequest) ProtoMessage()               {}
func (*TxRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *TxRequest) GetQueries() []string {
	if m != nil {
		return m.Queries
	}
	return nil
}

type TxResponse struct {
	Resluts []*ExecResponse `protobuf:"bytes,1,rep,name=resluts" json:"resluts,omitempty"`
}

func (m *TxResponse) Reset()                    { *m = TxResponse{} }
func (m *TxResponse) String() string            { return proto1.CompactTextString(m) }
func (*TxResponse) ProtoMessage()               {}
func (*TxResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *TxResponse) GetResluts() []*ExecResponse {
	if m != nil {
		return m.Resluts
	}
	return nil
}

type DatabasesResponse struct {
	Databases []string `protobuf:"bytes,1,rep,name=databases" json:"databases,omitempty"`
}

func (m *DatabasesResponse) Reset()                    { *m = DatabasesResponse{} }
func (m *DatabasesResponse) String() string            { return proto1.CompactTextString(m) }
func (*DatabasesResponse) ProtoMessage()               {}
func (*DatabasesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *DatabasesResponse) GetDatabases() []string {
	if m != nil {
		return m.Databases
	}
	return nil
}

type ParseRequest struct {
	SqlText string `protobuf:"bytes,1,opt,name=sqlText" json:"sqlText,omitempty"`
}

func (m *ParseRequest) Reset()                    { *m = ParseRequest{} }
func (m *ParseRequest) String() string            { return proto1.CompactTextString(m) }
func (*ParseRequest) ProtoMessage()               {}
func (*ParseRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ParseRequest) GetSqlText() string {
	if m != nil {
		return m.SqlText
	}
	return ""
}

type Node struct {
	Text        string `protobuf:"bytes,1,opt,name=text" json:"text,omitempty"`
	Type        string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Fingerprint string `protobuf:"bytes,3,opt,name=fingerprint" json:"fingerprint,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto1.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Node) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Node) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Node) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

type ParseResponse struct {
	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes" json:"nodes,omitempty"`
}

func (m *ParseResponse) Reset()                    { *m = ParseResponse{} }
func (m *ParseResponse) String() string            { return proto1.CompactTextString(m) }
func (*ParseResponse) ProtoMessage()               {}
func (*ParseResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ParseResponse) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type Rule struct {
	Name      string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Desc      string `protobuf:"bytes,2,opt,name=desc" json:"desc,omitempty"`
	Value     string `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
	Level     string `protobuf:"bytes,4,opt,name=level" json:"level,omitempty"`
	Typ       string `protobuf:"bytes,5,opt,name=typ" json:"typ,omitempty"`
	IsDefault bool   `protobuf:"varint,6,opt,name=isDefault" json:"isDefault,omitempty"`
}

func (m *Rule) Reset()                    { *m = Rule{} }
func (m *Rule) String() string            { return proto1.CompactTextString(m) }
func (*Rule) ProtoMessage()               {}
func (*Rule) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Rule) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Rule) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *Rule) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Rule) GetLevel() string {
	if m != nil {
		return m.Level
	}
	return ""
}

func (m *Rule) GetTyp() string {
	if m != nil {
		return m.Typ
	}
	return ""
}

func (m *Rule) GetIsDefault() bool {
	if m != nil {
		return m.IsDefault
	}
	return false
}

type AuditRequest struct {
	Rules []*Rule `protobuf:"bytes,1,rep,name=rules" json:"rules,omitempty"`
	Sql   string  `protobuf:"bytes,2,opt,name=sql" json:"sql,omitempty"`
}

func (m *AuditRequest) Reset()                    { *m = AuditRequest{} }
func (m *AuditRequest) String() string            { return proto1.CompactTextString(m) }
func (*AuditRequest) ProtoMessage()               {}
func (*AuditRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *AuditRequest) GetRules() []*Rule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func (m *AuditRequest) GetSql() string {
	if m != nil {
		return m.Sql
	}
	return ""
}

type AuditResult struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	Level   string `protobuf:"bytes,2,opt,name=level" json:"level,omitempty"`
}

func (m *AuditResult) Reset()                    { *m = AuditResult{} }
func (m *AuditResult) String() string            { return proto1.CompactTextString(m) }
func (*AuditResult) ProtoMessage()               {}
func (*AuditResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *AuditResult) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AuditResult) GetLevel() string {
	if m != nil {
		return m.Level
	}
	return ""
}

type AuditResponse struct {
	Results []*AuditResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *AuditResponse) Reset()                    { *m = AuditResponse{} }
func (m *AuditResponse) String() string            { return proto1.CompactTextString(m) }
func (*AuditResponse) ProtoMessage()               {}
func (*AuditResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *AuditResponse) GetResults() []*AuditResult {
	if m != nil {
		return m.Results
	}
	return nil
}

type GenRollbackSQLRequest struct {
	Sql string `protobuf:"bytes,1,opt,name=sql" json:"sql,omitempty"`
}

func (m *GenRollbackSQLRequest) Reset()                    { *m = GenRollbackSQLRequest{} }
func (m *GenRollbackSQLRequest) String() string            { return proto1.CompactTextString(m) }
func (*GenRollbackSQLRequest) ProtoMessage()               {}
func (*GenRollbackSQLRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *GenRollbackSQLRequest) GetSql() string {
	if m != nil {
		return m.Sql
	}
	return ""
}

type GenRollbackSQLResponse struct {
	Sql    string `protobuf:"bytes,1,opt,name=sql" json:"sql,omitempty"`
	Reason string `protobuf:"bytes,2,opt,name=reason" json:"reason,omitempty"`
}

func (m *GenRollbackSQLResponse) Reset()                    { *m = GenRollbackSQLResponse{} }
func (m *GenRollbackSQLResponse) String() string            { return proto1.CompactTextString(m) }
func (*GenRollbackSQLResponse) ProtoMessage()               {}
func (*GenRollbackSQLResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *GenRollbackSQLResponse) GetSql() string {
	if m != nil {
		return m.Sql
	}
	return ""
}

func (m *GenRollbackSQLResponse) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type MetasResponse struct {
	Name  string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Rules []*Rule `protobuf:"bytes,2,rep,name=rules" json:"rules,omitempty"`
}

func (m *MetasResponse) Reset()                    { *m = MetasResponse{} }
func (m *MetasResponse) String() string            { return proto1.CompactTextString(m) }
func (*MetasResponse) ProtoMessage()               {}
func (*MetasResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *MetasResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MetasResponse) GetRules() []*Rule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func init() {
	proto1.RegisterType((*InstanceMeta)(nil), "proto.InstanceMeta")
	proto1.RegisterType((*InitRequest)(nil), "proto.InitRequest")
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*ExecRequest)(nil), "proto.ExecRequest")
	proto1.RegisterType((*ExecResponse)(nil), "proto.ExecResponse")
	proto1.RegisterType((*TxRequest)(nil), "proto.TxRequest")
	proto1.RegisterType((*TxResponse)(nil), "proto.TxResponse")
	proto1.RegisterType((*DatabasesResponse)(nil), "proto.DatabasesResponse")
	proto1.RegisterType((*ParseRequest)(nil), "proto.ParseRequest")
	proto1.RegisterType((*Node)(nil), "proto.Node")
	proto1.RegisterType((*ParseResponse)(nil), "proto.ParseResponse")
	proto1.RegisterType((*Rule)(nil), "proto.Rule")
	proto1.RegisterType((*AuditRequest)(nil), "proto.AuditRequest")
	proto1.RegisterType((*AuditResult)(nil), "proto.AuditResult")
	proto1.RegisterType((*AuditResponse)(nil), "proto.AuditResponse")
	proto1.RegisterType((*GenRollbackSQLRequest)(nil), "proto.GenRollbackSQLRequest")
	proto1.RegisterType((*GenRollbackSQLResponse)(nil), "proto.GenRollbackSQLResponse")
	proto1.RegisterType((*MetasResponse)(nil), "proto.MetasResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Driver service

type DriverClient interface {
	// Metas returns some base info from plugin server.
	Metas(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MetasResponse, error)
	// Init will should be called at first before calling following methods.
	// It will pass some necessary info to plugin server. In the begginning,
	// we consider that put this info to the executable binary environment.
	// We put all communication on gRPC for unification in the end.
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error)
	Close(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResponse, error)
	Tx(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (*TxResponse, error)
	Databases(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DatabasesResponse, error)
	Parse(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error)
	Audit(ctx context.Context, in *AuditRequest, opts ...grpc.CallOption) (*AuditResponse, error)
	GenRollbackSQL(ctx context.Context, in *GenRollbackSQLRequest, opts ...grpc.CallOption) (*GenRollbackSQLResponse, error)
}

type driverClient struct {
	cc *grpc.ClientConn
}

func NewDriverClient(cc *grpc.ClientConn) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) Metas(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MetasResponse, error) {
	out := new(MetasResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Metas", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.Driver/Init", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Close(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.Driver/Close", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.Driver/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResponse, error) {
	out := new(ExecResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Exec", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Tx(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (*TxResponse, error) {
	out := new(TxResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Tx", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Databases(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DatabasesResponse, error) {
	out := new(DatabasesResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Databases", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Parse(ctx context.Context, in *ParseRequest, opts ...grpc.CallOption) (*ParseResponse, error) {
	out := new(ParseResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Parse", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Audit(ctx context.Context, in *AuditRequest, opts ...grpc.CallOption) (*AuditResponse, error) {
	out := new(AuditResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/Audit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) GenRollbackSQL(ctx context.Context, in *GenRollbackSQLRequest, opts ...grpc.CallOption) (*GenRollbackSQLResponse, error) {
	out := new(GenRollbackSQLResponse)
	err := grpc.Invoke(ctx, "/proto.Driver/GenRollbackSQL", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Driver service

type DriverServer interface {
	// Metas returns some base info from plugin server.
	Metas(context.Context, *Empty) (*MetasResponse, error)
	// Init will should be called at first before calling following methods.
	// It will pass some necessary info to plugin server. In the begginning,
	// we consider that put this info to the executable binary environment.
	// We put all communication on gRPC for unification in the end.
	Init(context.Context, *InitRequest) (*Empty, error)
	Close(context.Context, *Empty) (*Empty, error)
	Ping(context.Context, *Empty) (*Empty, error)
	Exec(context.Context, *ExecRequest) (*ExecResponse, error)
	Tx(context.Context, *TxRequest) (*TxResponse, error)
	Databases(context.Context, *Empty) (*DatabasesResponse, error)
	Parse(context.Context, *ParseRequest) (*ParseResponse, error)
	Audit(context.Context, *AuditRequest) (*AuditResponse, error)
	GenRollbackSQL(context.Context, *GenRollbackSQLRequest) (*GenRollbackSQLResponse, error)
}

func RegisterDriverServer(s *grpc.Server, srv DriverServer) {
	s.RegisterService(&_Driver_serviceDesc, srv)
}

func _Driver_Metas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Metas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Metas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Metas(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Close(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Exec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Exec(ctx, req.(*ExecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Tx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Tx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Tx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Tx(ctx, req.(*TxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Databases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Databases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Databases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Databases(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Parse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Parse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Parse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Parse(ctx, req.(*ParseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Audit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuditRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Audit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/Audit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Audit(ctx, req.(*AuditRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_GenRollbackSQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenRollbackSQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).GenRollbackSQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Driver/GenRollbackSQL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).GenRollbackSQL(ctx, req.(*GenRollbackSQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Driver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Metas",
			Handler:    _Driver_Metas_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _Driver_Init_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Driver_Close_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Driver_Ping_Handler,
		},
		{
			MethodName: "Exec",
			Handler:    _Driver_Exec_Handler,
		},
		{
			MethodName: "Tx",
			Handler:    _Driver_Tx_Handler,
		},
		{
			MethodName: "Databases",
			Handler:    _Driver_Databases_Handler,
		},
		{
			MethodName: "Parse",
			Handler:    _Driver_Parse_Handler,
		},
		{
			MethodName: "Audit",
			Handler:    _Driver_Audit_Handler,
		},
		{
			MethodName: "GenRollbackSQL",
			Handler:    _Driver_GenRollbackSQL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "driver.proto",
}

func init() { proto1.RegisterFile("driver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 777 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xcd, 0x4e, 0xe3, 0x48,
	0x10, 0x56, 0x7e, 0x1c, 0x48, 0x25, 0xac, 0xa0, 0x97, 0x45, 0x56, 0xc4, 0x4a, 0xd0, 0x68, 0xa5,
	0xa0, 0x65, 0x41, 0x1b, 0x0e, 0x7b, 0x40, 0x1c, 0xf8, 0xdb, 0xdd, 0x48, 0xc3, 0x90, 0xf1, 0x30,
	0x97, 0xb9, 0x99, 0xb8, 0x82, 0xac, 0x69, 0xec, 0xa4, 0xbb, 0xcd, 0x24, 0xaf, 0x30, 0x8f, 0x33,
	0x8f, 0x30, 0x4f, 0x34, 0x8f, 0x30, 0xea, 0x76, 0xb7, 0xdd, 0x4e, 0x60, 0x4e, 0xee, 0xaa, 0xfe,
	0x5c, 0xf5, 0x7d, 0xd5, 0x55, 0x05, 0xdd, 0x88, 0xc7, 0xcf, 0xc8, 0x8f, 0xa7, 0x3c, 0x95, 0x29,
	0xf1, 0xf4, 0x87, 0x7e, 0xab, 0x41, 0x77, 0x98, 0x08, 0x19, 0x26, 0x63, 0xbc, 0x45, 0x19, 0x12,
	0x0a, 0xdd, 0xd8, 0xd8, 0xff, 0xa7, 0x42, 0xfa, 0xb5, 0xbd, 0x5a, 0xbf, 0x1d, 0x54, 0x7c, 0x2e,
	0x66, 0x94, 0x72, 0xe9, 0xd7, 0xab, 0x18, 0xe5, 0x73, 0x31, 0x1f, 0x04, 0x72, 0xbf, 0x51, 0xc5,
	0x28, 0x5f, 0x25, 0x4e, 0x28, 0x84, 0xdf, 0x5c, 0x8a, 0x13, 0x0a, 0xa1, 0x30, 0x51, 0x28, 0xc3,
	0x87, 0x50, 0xe0, 0xdd, 0x14, 0x13, 0xdf, 0xcb, 0x31, 0xae, 0x8f, 0x46, 0xd0, 0x19, 0x26, 0xb1,
	0x0c, 0x70, 0x96, 0xa1, 0x90, 0xe4, 0x9f, 0x32, 0xac, 0x92, 0xa4, 0x25, 0x74, 0x06, 0xbf, 0xe6,
	0xc2, 0x8f, 0x5d, 0xb5, 0x41, 0x05, 0x48, 0x76, 0xa1, 0x1d, 0x8b, 0xbb, 0xc9, 0x84, 0xc5, 0x09,
	0x6a, 0x51, 0xeb, 0x41, 0xe9, 0xa0, 0x6b, 0xe0, 0xdd, 0x3c, 0x4d, 0xe5, 0x82, 0x1e, 0x40, 0xe7,
	0x66, 0x8e, 0x63, 0x9b, 0x6e, 0x1b, 0xbc, 0x59, 0x86, 0x7c, 0x61, 0x4a, 0x95, 0x1b, 0xf4, 0x6b,
	0x0d, 0xba, 0x39, 0x4a, 0x4c, 0xd3, 0x44, 0xa0, 0x12, 0xc2, 0x42, 0x21, 0x87, 0x89, 0x40, 0x2e,
	0x87, 0x91, 0x46, 0x37, 0x82, 0x8a, 0x8f, 0x1c, 0xc1, 0x96, 0x6b, 0xdf, 0x70, 0x9e, 0x72, 0x53,
	0xdd, 0xd5, 0x0b, 0x15, 0x91, 0xa7, 0x9f, 0xc5, 0xc5, 0x64, 0x82, 0x63, 0x89, 0x91, 0x2e, 0x71,
	0x23, 0xa8, 0xf8, 0x54, 0x44, 0xd7, 0xce, 0x23, 0xe6, 0x75, 0x5e, 0xbd, 0xa0, 0x7f, 0x40, 0xfb,
	0x7e, 0x6e, 0x75, 0xf9, 0xb0, 0xa6, 0xa4, 0xc4, 0x28, 0xfc, 0xda, 0x5e, 0xa3, 0xdf, 0x0e, 0xac,
	0x49, 0xcf, 0x00, 0x14, 0xcc, 0x08, 0xfb, 0x0b, 0xd6, 0x38, 0x0a, 0x96, 0xc9, 0x1c, 0x57, 0x56,
	0xda, 0x95, 0x1f, 0x58, 0x0c, 0xfd, 0x1b, 0xb6, 0xae, 0xcd, 0xe3, 0x89, 0x22, 0xc6, 0x2e, 0xb4,
	0xed, 0x8b, 0xda, 0x6c, 0xa5, 0x83, 0xf6, 0xa1, 0x3b, 0x0a, 0xb9, 0x40, 0x87, 0x99, 0x98, 0xb1,
	0x7b, 0x9c, 0xdb, 0xf6, 0xb4, 0x26, 0x1d, 0x41, 0xf3, 0x6d, 0x1a, 0x21, 0x21, 0xd0, 0x94, 0xe5,
	0xb5, 0x3e, 0x6b, 0xdf, 0x62, 0x8a, 0xa6, 0x9e, 0xfa, 0x4c, 0xf6, 0xa0, 0x33, 0x89, 0x93, 0x47,
	0xe4, 0x53, 0x1e, 0x27, 0xd2, 0x34, 0xa9, 0xeb, 0xa2, 0x03, 0xd8, 0x30, 0xb9, 0x0d, 0xd5, 0x7d,
	0xf0, 0x92, 0x34, 0x42, 0x2b, 0xb6, 0x63, 0xc4, 0xaa, 0xb4, 0x41, 0x7e, 0x43, 0xbf, 0xd4, 0xa0,
	0x19, 0x64, 0x4c, 0xd3, 0x48, 0xc2, 0x27, 0xb4, 0x34, 0xd4, 0x59, 0xf9, 0x22, 0x14, 0x63, 0x4b,
	0x43, 0x9d, 0x55, 0x0b, 0x3d, 0x87, 0x2c, 0x43, 0x43, 0x20, 0x37, 0x94, 0x97, 0xe1, 0x33, 0x32,
	0xf3, 0x5e, 0xb9, 0x41, 0x36, 0xa1, 0x21, 0x17, 0x53, 0x33, 0x07, 0xea, 0x98, 0xb7, 0xed, 0x35,
	0x4e, 0xc2, 0x8c, 0x49, 0xbf, 0x65, 0xdb, 0xd6, 0x38, 0xe8, 0x15, 0x74, 0x2f, 0xb2, 0xa8, 0x9c,
	0x8e, 0x7d, 0xf0, 0x78, 0xc6, 0x56, 0xf8, 0x2b, 0xbe, 0x41, 0x7e, 0xa3, 0x52, 0x88, 0x19, 0x33,
	0x0c, 0xd5, 0x91, 0x9e, 0x43, 0xc7, 0x04, 0x11, 0x19, 0xd3, 0x0f, 0xf0, 0x84, 0x42, 0x84, 0x8f,
	0x56, 0x9a, 0x35, 0x4b, 0xce, 0x75, 0x87, 0x33, 0x3d, 0x87, 0x0d, 0xfb, 0x7b, 0x5e, 0xc4, 0x23,
	0xdd, 0x33, 0x19, 0x2b, 0x7a, 0x86, 0x18, 0x1a, 0x4e, 0x96, 0xc0, 0x42, 0xe8, 0x21, 0xfc, 0xf6,
	0x1f, 0x26, 0x41, 0xca, 0xd8, 0x43, 0x38, 0xfe, 0xf4, 0xfe, 0xdd, 0x1b, 0xab, 0xc5, 0x10, 0xad,
	0x95, 0x44, 0x2f, 0x61, 0x67, 0x19, 0x6a, 0x52, 0xae, 0x60, 0xc9, 0x0e, 0xb4, 0x38, 0x86, 0x22,
	0x4d, 0x0c, 0x59, 0x63, 0xd1, 0x7f, 0x61, 0x43, 0xad, 0x83, 0xb2, 0x3b, 0x5f, 0x7a, 0xc6, 0xa2,
	0x8c, 0xf5, 0xd7, 0xca, 0x38, 0xf8, 0xde, 0x80, 0xd6, 0xb5, 0xde, 0xb9, 0xe4, 0x4f, 0xf0, 0x74,
	0x48, 0xd2, 0xb5, 0xb3, 0xa1, 0x36, 0x49, 0x6f, 0xdb, 0x58, 0xd5, 0x74, 0x7d, 0x68, 0xaa, 0x75,
	0x46, 0x48, 0xb1, 0xb1, 0x8a, 0xd7, 0xeb, 0x55, 0xfe, 0x27, 0x07, 0xe0, 0x5d, 0xb1, 0x54, 0xe0,
	0x52, 0xd8, 0x2a, 0x88, 0x42, 0x73, 0x14, 0x27, 0x8f, 0x3f, 0xc5, 0x9c, 0x40, 0x53, 0x4d, 0x6b,
	0x91, 0xd2, 0xd9, 0x6f, 0xbd, 0x97, 0xc6, 0x99, 0x1c, 0x42, 0xfd, 0x7e, 0x4e, 0x36, 0xcd, 0x55,
	0xb1, 0x34, 0x7a, 0x5b, 0x8e, 0xc7, 0x40, 0x4f, 0xa1, 0x5d, 0x0c, 0xfc, 0x12, 0x09, 0xdf, 0x58,
	0xab, 0x0b, 0x61, 0x00, 0x9e, 0x1e, 0x3b, 0x62, 0xb3, 0xbb, 0x0b, 0xa0, 0xa8, 0x5b, 0x75, 0x32,
	0x07, 0xe0, 0xe9, 0xf6, 0x29, 0xfe, 0x71, 0xfb, 0xbe, 0xf8, 0xa7, 0xda, 0x88, 0xb7, 0xf0, 0x4b,
	0xb5, 0x5f, 0xc8, 0xae, 0xc1, 0xbd, 0xd8, 0x71, 0xbd, 0xdf, 0x5f, 0xb9, 0xcd, 0xc3, 0x5d, 0xc2,
	0xc7, 0xf5, 0xe3, 0x93, 0x33, 0x0d, 0x79, 0x68, 0xe9, 0xcf, 0xe9, 0x8f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x08, 0x61, 0xd9, 0x51, 0x78, 0x07, 0x00, 0x00,
}
