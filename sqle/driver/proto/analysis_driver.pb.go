// Code generated by protoc-gen-go. DO NOT EDIT.
// source: analysis_driver.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	analysis_driver.proto
	driver.proto
	query_driver.proto

It has these top-level messages:
	AnalysisDriverInitRequest
	ListTablesInSchemaRequest
	ListTablesInSchemaResponse
	Table
	GetTableMetaByTableNameRequest
	GetTableMetaByTableNameResponse
	TableItem
	ColumnsInfo
	IndexesInfo
	Row
	AnalysisInfoInTableFormat
	AnalysisInfoHead
	GetTableMetaBySQLRequest
	GetTableMetaBySQLResponse
	TableMetaItemBySQL
	ExplainRequest
	ExplainResponse
	ExplainClassicResult
	DSN
	Rule
	Param
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
	AuditRequest
	AuditResult
	AuditResponse
	GenRollbackSQLRequest
	GenRollbackSQLResponse
	MetasResponse
	QueryPrepareRequest
	QueryPrepareConf
	QueryPrepareResponse
	QueryRequest
	QueryConf
	QueryResponse
	QueryResultRow
	QueryResultValue
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

type AnalysisDriverInitRequest struct {
	Dsn *DSN `protobuf:"bytes,1,opt,name=dsn" json:"dsn,omitempty"`
}

func (m *AnalysisDriverInitRequest) Reset()                    { *m = AnalysisDriverInitRequest{} }
func (m *AnalysisDriverInitRequest) String() string            { return proto1.CompactTextString(m) }
func (*AnalysisDriverInitRequest) ProtoMessage()               {}
func (*AnalysisDriverInitRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *AnalysisDriverInitRequest) GetDsn() *DSN {
	if m != nil {
		return m.Dsn
	}
	return nil
}

type ListTablesInSchemaRequest struct {
	Schema string `protobuf:"bytes,1,opt,name=schema" json:"schema,omitempty"`
}

func (m *ListTablesInSchemaRequest) Reset()                    { *m = ListTablesInSchemaRequest{} }
func (m *ListTablesInSchemaRequest) String() string            { return proto1.CompactTextString(m) }
func (*ListTablesInSchemaRequest) ProtoMessage()               {}
func (*ListTablesInSchemaRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ListTablesInSchemaRequest) GetSchema() string {
	if m != nil {
		return m.Schema
	}
	return ""
}

type ListTablesInSchemaResponse struct {
	Tables []*Table `protobuf:"bytes,1,rep,name=tables" json:"tables,omitempty"`
}

func (m *ListTablesInSchemaResponse) Reset()                    { *m = ListTablesInSchemaResponse{} }
func (m *ListTablesInSchemaResponse) String() string            { return proto1.CompactTextString(m) }
func (*ListTablesInSchemaResponse) ProtoMessage()               {}
func (*ListTablesInSchemaResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *ListTablesInSchemaResponse) GetTables() []*Table {
	if m != nil {
		return m.Tables
	}
	return nil
}

type Table struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *Table) Reset()                    { *m = Table{} }
func (m *Table) String() string            { return proto1.CompactTextString(m) }
func (*Table) ProtoMessage()               {}
func (*Table) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *Table) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetTableMetaByTableNameRequest struct {
	Schema string `protobuf:"bytes,1,opt,name=schema" json:"schema,omitempty"`
	Table  string `protobuf:"bytes,2,opt,name=table" json:"table,omitempty"`
}

func (m *GetTableMetaByTableNameRequest) Reset()                    { *m = GetTableMetaByTableNameRequest{} }
func (m *GetTableMetaByTableNameRequest) String() string            { return proto1.CompactTextString(m) }
func (*GetTableMetaByTableNameRequest) ProtoMessage()               {}
func (*GetTableMetaByTableNameRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *GetTableMetaByTableNameRequest) GetSchema() string {
	if m != nil {
		return m.Schema
	}
	return ""
}

func (m *GetTableMetaByTableNameRequest) GetTable() string {
	if m != nil {
		return m.Table
	}
	return ""
}

type GetTableMetaByTableNameResponse struct {
	TableMeta *TableItem `protobuf:"bytes,1,opt,name=tableMeta" json:"tableMeta,omitempty"`
}

func (m *GetTableMetaByTableNameResponse) Reset()                    { *m = GetTableMetaByTableNameResponse{} }
func (m *GetTableMetaByTableNameResponse) String() string            { return proto1.CompactTextString(m) }
func (*GetTableMetaByTableNameResponse) ProtoMessage()               {}
func (*GetTableMetaByTableNameResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *GetTableMetaByTableNameResponse) GetTableMeta() *TableItem {
	if m != nil {
		return m.TableMeta
	}
	return nil
}

type TableItem struct {
	Name           string       `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Schema         string       `protobuf:"bytes,2,opt,name=schema" json:"schema,omitempty"`
	ColumnsInfo    *ColumnsInfo `protobuf:"bytes,3,opt,name=columnsInfo" json:"columnsInfo,omitempty"`
	IndexesInfo    *IndexesInfo `protobuf:"bytes,4,opt,name=indexesInfo" json:"indexesInfo,omitempty"`
	CreateTableSQL string       `protobuf:"bytes,5,opt,name=createTableSQL" json:"createTableSQL,omitempty"`
}

func (m *TableItem) Reset()                    { *m = TableItem{} }
func (m *TableItem) String() string            { return proto1.CompactTextString(m) }
func (*TableItem) ProtoMessage()               {}
func (*TableItem) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

func (m *TableItem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TableItem) GetSchema() string {
	if m != nil {
		return m.Schema
	}
	return ""
}

func (m *TableItem) GetColumnsInfo() *ColumnsInfo {
	if m != nil {
		return m.ColumnsInfo
	}
	return nil
}

func (m *TableItem) GetIndexesInfo() *IndexesInfo {
	if m != nil {
		return m.IndexesInfo
	}
	return nil
}

func (m *TableItem) GetCreateTableSQL() string {
	if m != nil {
		return m.CreateTableSQL
	}
	return ""
}

type ColumnsInfo struct {
	AnalysisInfoInTableFormat *AnalysisInfoInTableFormat `protobuf:"bytes,1,opt,name=analysisInfoInTableFormat" json:"analysisInfoInTableFormat,omitempty"`
}

func (m *ColumnsInfo) Reset()                    { *m = ColumnsInfo{} }
func (m *ColumnsInfo) String() string            { return proto1.CompactTextString(m) }
func (*ColumnsInfo) ProtoMessage()               {}
func (*ColumnsInfo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

func (m *ColumnsInfo) GetAnalysisInfoInTableFormat() *AnalysisInfoInTableFormat {
	if m != nil {
		return m.AnalysisInfoInTableFormat
	}
	return nil
}

type IndexesInfo struct {
	AnalysisInfoInTableFormat *AnalysisInfoInTableFormat `protobuf:"bytes,1,opt,name=analysisInfoInTableFormat" json:"analysisInfoInTableFormat,omitempty"`
}

func (m *IndexesInfo) Reset()                    { *m = IndexesInfo{} }
func (m *IndexesInfo) String() string            { return proto1.CompactTextString(m) }
func (*IndexesInfo) ProtoMessage()               {}
func (*IndexesInfo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

func (m *IndexesInfo) GetAnalysisInfoInTableFormat() *AnalysisInfoInTableFormat {
	if m != nil {
		return m.AnalysisInfoInTableFormat
	}
	return nil
}

type Row struct {
	Items []string `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
}

func (m *Row) Reset()                    { *m = Row{} }
func (m *Row) String() string            { return proto1.CompactTextString(m) }
func (*Row) ProtoMessage()               {}
func (*Row) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *Row) GetItems() []string {
	if m != nil {
		return m.Items
	}
	return nil
}

type AnalysisInfoInTableFormat struct {
	Columns []*AnalysisInfoHead `protobuf:"bytes,1,rep,name=columns" json:"columns,omitempty"`
	Rows    []*Row              `protobuf:"bytes,2,rep,name=rows" json:"rows,omitempty"`
}

func (m *AnalysisInfoInTableFormat) Reset()                    { *m = AnalysisInfoInTableFormat{} }
func (m *AnalysisInfoInTableFormat) String() string            { return proto1.CompactTextString(m) }
func (*AnalysisInfoInTableFormat) ProtoMessage()               {}
func (*AnalysisInfoInTableFormat) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *AnalysisInfoInTableFormat) GetColumns() []*AnalysisInfoHead {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *AnalysisInfoInTableFormat) GetRows() []*Row {
	if m != nil {
		return m.Rows
	}
	return nil
}

type AnalysisInfoHead struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Desc string `protobuf:"bytes,2,opt,name=desc" json:"desc,omitempty"`
}

func (m *AnalysisInfoHead) Reset()                    { *m = AnalysisInfoHead{} }
func (m *AnalysisInfoHead) String() string            { return proto1.CompactTextString(m) }
func (*AnalysisInfoHead) ProtoMessage()               {}
func (*AnalysisInfoHead) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *AnalysisInfoHead) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AnalysisInfoHead) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

type GetTableMetaBySQLRequest struct {
	Sql string `protobuf:"bytes,1,opt,name=sql" json:"sql,omitempty"`
}

func (m *GetTableMetaBySQLRequest) Reset()                    { *m = GetTableMetaBySQLRequest{} }
func (m *GetTableMetaBySQLRequest) String() string            { return proto1.CompactTextString(m) }
func (*GetTableMetaBySQLRequest) ProtoMessage()               {}
func (*GetTableMetaBySQLRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *GetTableMetaBySQLRequest) GetSql() string {
	if m != nil {
		return m.Sql
	}
	return ""
}

type GetTableMetaBySQLResponse struct {
	TableMetas []*TableMetaItemBySQL `protobuf:"bytes,1,rep,name=tableMetas" json:"tableMetas,omitempty"`
}

func (m *GetTableMetaBySQLResponse) Reset()                    { *m = GetTableMetaBySQLResponse{} }
func (m *GetTableMetaBySQLResponse) String() string            { return proto1.CompactTextString(m) }
func (*GetTableMetaBySQLResponse) ProtoMessage()               {}
func (*GetTableMetaBySQLResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *GetTableMetaBySQLResponse) GetTableMetas() []*TableMetaItemBySQL {
	if m != nil {
		return m.TableMetas
	}
	return nil
}

type TableMetaItemBySQL struct {
	Name           string       `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Schema         string       `protobuf:"bytes,2,opt,name=schema" json:"schema,omitempty"`
	ColumnsInfo    *ColumnsInfo `protobuf:"bytes,3,opt,name=columnsInfo" json:"columnsInfo,omitempty"`
	IndexesInfo    *IndexesInfo `protobuf:"bytes,4,opt,name=indexesInfo" json:"indexesInfo,omitempty"`
	CreateTableSQL string       `protobuf:"bytes,5,opt,name=createTableSQL" json:"createTableSQL,omitempty"`
	ErrMessage     string       `protobuf:"bytes,6,opt,name=errMessage" json:"errMessage,omitempty"`
}

func (m *TableMetaItemBySQL) Reset()                    { *m = TableMetaItemBySQL{} }
func (m *TableMetaItemBySQL) String() string            { return proto1.CompactTextString(m) }
func (*TableMetaItemBySQL) ProtoMessage()               {}
func (*TableMetaItemBySQL) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

func (m *TableMetaItemBySQL) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TableMetaItemBySQL) GetSchema() string {
	if m != nil {
		return m.Schema
	}
	return ""
}

func (m *TableMetaItemBySQL) GetColumnsInfo() *ColumnsInfo {
	if m != nil {
		return m.ColumnsInfo
	}
	return nil
}

func (m *TableMetaItemBySQL) GetIndexesInfo() *IndexesInfo {
	if m != nil {
		return m.IndexesInfo
	}
	return nil
}

func (m *TableMetaItemBySQL) GetCreateTableSQL() string {
	if m != nil {
		return m.CreateTableSQL
	}
	return ""
}

func (m *TableMetaItemBySQL) GetErrMessage() string {
	if m != nil {
		return m.ErrMessage
	}
	return ""
}

type ExplainRequest struct {
	Sql string `protobuf:"bytes,1,opt,name=sql" json:"sql,omitempty"`
}

func (m *ExplainRequest) Reset()                    { *m = ExplainRequest{} }
func (m *ExplainRequest) String() string            { return proto1.CompactTextString(m) }
func (*ExplainRequest) ProtoMessage()               {}
func (*ExplainRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{15} }

func (m *ExplainRequest) GetSql() string {
	if m != nil {
		return m.Sql
	}
	return ""
}

type ExplainResponse struct {
	ClassicResult *ExplainClassicResult `protobuf:"bytes,1,opt,name=classicResult" json:"classicResult,omitempty"`
}

func (m *ExplainResponse) Reset()                    { *m = ExplainResponse{} }
func (m *ExplainResponse) String() string            { return proto1.CompactTextString(m) }
func (*ExplainResponse) ProtoMessage()               {}
func (*ExplainResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{16} }

func (m *ExplainResponse) GetClassicResult() *ExplainClassicResult {
	if m != nil {
		return m.ClassicResult
	}
	return nil
}

type ExplainClassicResult struct {
	AnalysisInfoInTableFormat *AnalysisInfoInTableFormat `protobuf:"bytes,1,opt,name=analysisInfoInTableFormat" json:"analysisInfoInTableFormat,omitempty"`
}

func (m *ExplainClassicResult) Reset()                    { *m = ExplainClassicResult{} }
func (m *ExplainClassicResult) String() string            { return proto1.CompactTextString(m) }
func (*ExplainClassicResult) ProtoMessage()               {}
func (*ExplainClassicResult) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{17} }

func (m *ExplainClassicResult) GetAnalysisInfoInTableFormat() *AnalysisInfoInTableFormat {
	if m != nil {
		return m.AnalysisInfoInTableFormat
	}
	return nil
}

func init() {
	proto1.RegisterType((*AnalysisDriverInitRequest)(nil), "proto.AnalysisDriverInitRequest")
	proto1.RegisterType((*ListTablesInSchemaRequest)(nil), "proto.ListTablesInSchemaRequest")
	proto1.RegisterType((*ListTablesInSchemaResponse)(nil), "proto.ListTablesInSchemaResponse")
	proto1.RegisterType((*Table)(nil), "proto.Table")
	proto1.RegisterType((*GetTableMetaByTableNameRequest)(nil), "proto.GetTableMetaByTableNameRequest")
	proto1.RegisterType((*GetTableMetaByTableNameResponse)(nil), "proto.GetTableMetaByTableNameResponse")
	proto1.RegisterType((*TableItem)(nil), "proto.TableItem")
	proto1.RegisterType((*ColumnsInfo)(nil), "proto.ColumnsInfo")
	proto1.RegisterType((*IndexesInfo)(nil), "proto.IndexesInfo")
	proto1.RegisterType((*Row)(nil), "proto.Row")
	proto1.RegisterType((*AnalysisInfoInTableFormat)(nil), "proto.AnalysisInfoInTableFormat")
	proto1.RegisterType((*AnalysisInfoHead)(nil), "proto.AnalysisInfoHead")
	proto1.RegisterType((*GetTableMetaBySQLRequest)(nil), "proto.GetTableMetaBySQLRequest")
	proto1.RegisterType((*GetTableMetaBySQLResponse)(nil), "proto.GetTableMetaBySQLResponse")
	proto1.RegisterType((*TableMetaItemBySQL)(nil), "proto.TableMetaItemBySQL")
	proto1.RegisterType((*ExplainRequest)(nil), "proto.ExplainRequest")
	proto1.RegisterType((*ExplainResponse)(nil), "proto.ExplainResponse")
	proto1.RegisterType((*ExplainClassicResult)(nil), "proto.ExplainClassicResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AnalysisDriver service

type AnalysisDriverClient interface {
	// Init should be called at first before calling following methods.
	// It will pass some necessary info to plugin server. In the beginning,
	// we consider that put this info to the executable binary environment.
	// We put all communication on gRPC for unification in the end.
	Init(ctx context.Context, in *AnalysisDriverInitRequest, opts ...grpc.CallOption) (*Empty, error)
	ListTablesInSchema(ctx context.Context, in *ListTablesInSchemaRequest, opts ...grpc.CallOption) (*ListTablesInSchemaResponse, error)
	GetTableMetaByTableName(ctx context.Context, in *GetTableMetaByTableNameRequest, opts ...grpc.CallOption) (*GetTableMetaByTableNameResponse, error)
	GetTableMetaBySQL(ctx context.Context, in *GetTableMetaBySQLRequest, opts ...grpc.CallOption) (*GetTableMetaBySQLResponse, error)
	Explain(ctx context.Context, in *ExplainRequest, opts ...grpc.CallOption) (*ExplainResponse, error)
}

type analysisDriverClient struct {
	cc *grpc.ClientConn
}

func NewAnalysisDriverClient(cc *grpc.ClientConn) AnalysisDriverClient {
	return &analysisDriverClient{cc}
}

func (c *analysisDriverClient) Init(ctx context.Context, in *AnalysisDriverInitRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.AnalysisDriver/Init", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisDriverClient) ListTablesInSchema(ctx context.Context, in *ListTablesInSchemaRequest, opts ...grpc.CallOption) (*ListTablesInSchemaResponse, error) {
	out := new(ListTablesInSchemaResponse)
	err := grpc.Invoke(ctx, "/proto.AnalysisDriver/ListTablesInSchema", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisDriverClient) GetTableMetaByTableName(ctx context.Context, in *GetTableMetaByTableNameRequest, opts ...grpc.CallOption) (*GetTableMetaByTableNameResponse, error) {
	out := new(GetTableMetaByTableNameResponse)
	err := grpc.Invoke(ctx, "/proto.AnalysisDriver/GetTableMetaByTableName", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisDriverClient) GetTableMetaBySQL(ctx context.Context, in *GetTableMetaBySQLRequest, opts ...grpc.CallOption) (*GetTableMetaBySQLResponse, error) {
	out := new(GetTableMetaBySQLResponse)
	err := grpc.Invoke(ctx, "/proto.AnalysisDriver/GetTableMetaBySQL", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisDriverClient) Explain(ctx context.Context, in *ExplainRequest, opts ...grpc.CallOption) (*ExplainResponse, error) {
	out := new(ExplainResponse)
	err := grpc.Invoke(ctx, "/proto.AnalysisDriver/Explain", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AnalysisDriver service

type AnalysisDriverServer interface {
	// Init should be called at first before calling following methods.
	// It will pass some necessary info to plugin server. In the beginning,
	// we consider that put this info to the executable binary environment.
	// We put all communication on gRPC for unification in the end.
	Init(context.Context, *AnalysisDriverInitRequest) (*Empty, error)
	ListTablesInSchema(context.Context, *ListTablesInSchemaRequest) (*ListTablesInSchemaResponse, error)
	GetTableMetaByTableName(context.Context, *GetTableMetaByTableNameRequest) (*GetTableMetaByTableNameResponse, error)
	GetTableMetaBySQL(context.Context, *GetTableMetaBySQLRequest) (*GetTableMetaBySQLResponse, error)
	Explain(context.Context, *ExplainRequest) (*ExplainResponse, error)
}

func RegisterAnalysisDriverServer(s *grpc.Server, srv AnalysisDriverServer) {
	s.RegisterService(&_AnalysisDriver_serviceDesc, srv)
}

func _AnalysisDriver_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalysisDriverInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisDriverServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AnalysisDriver/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisDriverServer).Init(ctx, req.(*AnalysisDriverInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisDriver_ListTablesInSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTablesInSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisDriverServer).ListTablesInSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AnalysisDriver/ListTablesInSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisDriverServer).ListTablesInSchema(ctx, req.(*ListTablesInSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisDriver_GetTableMetaByTableName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTableMetaByTableNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisDriverServer).GetTableMetaByTableName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AnalysisDriver/GetTableMetaByTableName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisDriverServer).GetTableMetaByTableName(ctx, req.(*GetTableMetaByTableNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisDriver_GetTableMetaBySQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTableMetaBySQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisDriverServer).GetTableMetaBySQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AnalysisDriver/GetTableMetaBySQL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisDriverServer).GetTableMetaBySQL(ctx, req.(*GetTableMetaBySQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisDriver_Explain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExplainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisDriverServer).Explain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AnalysisDriver/Explain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisDriverServer).Explain(ctx, req.(*ExplainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AnalysisDriver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AnalysisDriver",
	HandlerType: (*AnalysisDriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _AnalysisDriver_Init_Handler,
		},
		{
			MethodName: "ListTablesInSchema",
			Handler:    _AnalysisDriver_ListTablesInSchema_Handler,
		},
		{
			MethodName: "GetTableMetaByTableName",
			Handler:    _AnalysisDriver_GetTableMetaByTableName_Handler,
		},
		{
			MethodName: "GetTableMetaBySQL",
			Handler:    _AnalysisDriver_GetTableMetaBySQL_Handler,
		},
		{
			MethodName: "Explain",
			Handler:    _AnalysisDriver_Explain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "analysis_driver.proto",
}

func init() { proto1.RegisterFile("analysis_driver.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 661 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x55, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x56, 0xea, 0x24, 0x25, 0x93, 0x52, 0xca, 0xaa, 0x3f, 0x8e, 0x8b, 0xda, 0xb0, 0x82, 0xaa,
	0x07, 0x14, 0x44, 0x8b, 0x10, 0x85, 0x53, 0xff, 0x00, 0x4b, 0x69, 0xa4, 0x6c, 0xaa, 0x4a, 0x70,
	0x00, 0x6d, 0x9d, 0x85, 0x5a, 0xb2, 0xd7, 0xa9, 0x77, 0xd3, 0x34, 0x4f, 0xc3, 0x03, 0xf1, 0x3a,
	0x3c, 0x00, 0xf2, 0x7a, 0x9d, 0xd8, 0x89, 0xdd, 0x9e, 0x7a, 0xe0, 0xe4, 0xdd, 0x99, 0xef, 0x9b,
	0x19, 0x8f, 0xbf, 0xf1, 0xc0, 0x1a, 0xe5, 0xd4, 0x1b, 0x0b, 0x57, 0xfc, 0xe8, 0x87, 0xee, 0x0d,
	0x0b, 0x5b, 0x83, 0x30, 0x90, 0x01, 0xaa, 0xa8, 0x87, 0xb5, 0x94, 0x36, 0xe2, 0x03, 0x68, 0x1c,
	0x6a, 0xf4, 0x89, 0xb2, 0xdb, 0xdc, 0x95, 0x84, 0x5d, 0x0f, 0x99, 0x90, 0xe8, 0x19, 0x18, 0x7d,
	0xc1, 0xcd, 0x52, 0xb3, 0xb4, 0x5b, 0xdf, 0x83, 0x98, 0xd1, 0x3a, 0xe9, 0x75, 0x48, 0x64, 0xc6,
	0xfb, 0xd0, 0x68, 0xbb, 0x42, 0x9e, 0xd3, 0x4b, 0x8f, 0x09, 0x9b, 0xf7, 0x9c, 0x2b, 0xe6, 0xd3,
	0x84, 0xba, 0x0e, 0x55, 0xa1, 0x0c, 0x8a, 0x5d, 0x23, 0xfa, 0x86, 0x8f, 0xc0, 0xca, 0x23, 0x89,
	0x41, 0xc0, 0x05, 0x43, 0x2f, 0xa0, 0x2a, 0x95, 0xc7, 0x2c, 0x35, 0x8d, 0xdd, 0xfa, 0xde, 0x92,
	0xce, 0xa9, 0xe0, 0x44, 0xfb, 0xf0, 0x26, 0x54, 0x94, 0x01, 0x21, 0x28, 0x73, 0xea, 0x33, 0x9d,
	0x42, 0x9d, 0x71, 0x07, 0xb6, 0x3e, 0xb3, 0x38, 0xfe, 0x19, 0x93, 0xf4, 0x68, 0xac, 0x8e, 0x1d,
	0xea, 0xb3, 0x7b, 0x4a, 0x43, 0xab, 0x50, 0x51, 0x09, 0xcc, 0x05, 0x65, 0x8e, 0x2f, 0xb8, 0x0b,
	0xdb, 0x85, 0xf1, 0x74, 0xd5, 0x2d, 0xa8, 0xc9, 0xc4, 0xaf, 0x9b, 0xb5, 0x92, 0x2e, 0xdc, 0x96,
	0xcc, 0x27, 0x53, 0x08, 0xfe, 0x53, 0x82, 0xda, 0xc4, 0x91, 0xf7, 0x12, 0xa9, 0x12, 0x17, 0x32,
	0x25, 0xbe, 0x85, 0xba, 0x13, 0x78, 0x43, 0x9f, 0x0b, 0x9b, 0xff, 0x0c, 0x4c, 0x43, 0xe5, 0x42,
	0x3a, 0xd7, 0xf1, 0xd4, 0x43, 0xd2, 0xb0, 0x88, 0xe5, 0xf2, 0x3e, 0xbb, 0x65, 0x31, 0xab, 0x9c,
	0x61, 0xd9, 0x53, 0x0f, 0x49, 0xc3, 0xd0, 0x0e, 0x2c, 0x3b, 0x21, 0xa3, 0x92, 0xa9, 0x52, 0x7b,
	0xdd, 0xb6, 0x59, 0x51, 0xb5, 0xcc, 0x58, 0xb1, 0x0f, 0xf5, 0x54, 0x66, 0xf4, 0x1d, 0x1a, 0x89,
	0xfc, 0xa2, 0xbb, 0xcd, 0x15, 0xf0, 0x53, 0x10, 0xfa, 0x54, 0xea, 0xe6, 0x34, 0x75, 0xea, 0xc3,
	0x22, 0x1c, 0x29, 0x0e, 0x11, 0xa5, 0x4b, 0x95, 0xfc, 0xe0, 0xe9, 0x36, 0xc1, 0x20, 0xc1, 0x28,
	0xd2, 0x86, 0x2b, 0x99, 0x1f, 0xeb, 0xb2, 0x46, 0xe2, 0x0b, 0xe6, 0xd3, 0xe1, 0x99, 0x63, 0xa2,
	0x37, 0xb0, 0xa8, 0x3f, 0x82, 0x16, 0xf3, 0x46, 0x4e, 0x1d, 0x5f, 0x18, 0xed, 0x93, 0x04, 0x87,
	0xb6, 0xa0, 0x1c, 0x06, 0x23, 0x61, 0x2e, 0x28, 0x7c, 0x32, 0x70, 0x24, 0x18, 0x11, 0x65, 0xc7,
	0x1f, 0x60, 0x65, 0x96, 0x9c, 0x2b, 0x1f, 0x04, 0xe5, 0x3e, 0x13, 0x8e, 0x16, 0x8f, 0x3a, 0xe3,
	0x57, 0x60, 0x66, 0x75, 0xdc, 0xeb, 0xb6, 0x93, 0x89, 0x58, 0x01, 0x43, 0x5c, 0x7b, 0x3a, 0x44,
	0x74, 0xc4, 0x17, 0xd0, 0xc8, 0x41, 0x6b, 0xbd, 0x1f, 0x00, 0x4c, 0xc4, 0x9c, 0xbc, 0x5c, 0x23,
	0x2d, 0xf8, 0xc8, 0x11, 0x69, 0x3b, 0xa6, 0xa5, 0xc0, 0xf8, 0x6f, 0x09, 0xd0, 0x3c, 0xe4, 0x7f,
	0x9f, 0x01, 0xb4, 0x05, 0xc0, 0xc2, 0xf0, 0x8c, 0x09, 0x41, 0x7f, 0x31, 0xb3, 0xaa, 0x30, 0x29,
	0x0b, 0xc6, 0xb0, 0x7c, 0x7a, 0x3b, 0xf0, 0xa8, 0xcb, 0x8b, 0x5b, 0x7e, 0x0e, 0x4f, 0x26, 0x18,
	0xdd, 0xe8, 0x43, 0x78, 0xec, 0x78, 0x54, 0x08, 0xd7, 0x21, 0x4c, 0x0c, 0xbd, 0x44, 0xd0, 0x9b,
	0xba, 0x6c, 0x0d, 0x3f, 0x4e, 0x43, 0x48, 0x96, 0x81, 0x6f, 0x60, 0x35, 0x0f, 0xf6, 0xd0, 0x73,
	0xb3, 0xf7, 0xdb, 0x80, 0xe5, 0xec, 0x62, 0x41, 0xef, 0xa0, 0x1c, 0x2d, 0x17, 0x34, 0x1b, 0x77,
	0x6e, 0xef, 0x58, 0xc9, 0x6f, 0xff, 0xd4, 0x1f, 0xc8, 0x31, 0xfa, 0x0a, 0x68, 0x7e, 0x65, 0x4c,
	0xa2, 0x14, 0xae, 0x20, 0xeb, 0xf9, 0x1d, 0x08, 0xdd, 0xe0, 0x2b, 0xd8, 0x28, 0xf8, 0xb9, 0xa3,
	0x97, 0x9a, 0x7d, 0xf7, 0x32, 0xb1, 0x76, 0xee, 0x83, 0xe9, 0x4c, 0x17, 0xf0, 0x74, 0x6e, 0xa0,
	0xd0, 0x76, 0x2e, 0x79, 0x3a, 0x98, 0x56, 0xb3, 0x18, 0xa0, 0xe3, 0xbe, 0x87, 0x45, 0xfd, 0x7d,
	0xd1, 0x5a, 0x56, 0x16, 0x49, 0x8c, 0xf5, 0x59, 0x73, 0xcc, 0x3c, 0x82, 0x6f, 0x8f, 0x5a, 0xaf,
	0x3f, 0x2a, 0xdf, 0x65, 0x55, 0x3d, 0xf6, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x59, 0xc3, 0x7f,
	0xd8, 0x3a, 0x08, 0x00, 0x00,
}
