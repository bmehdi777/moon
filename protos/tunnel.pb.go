// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v3.19.6
// source: protos/tunnel.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StreamClient struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Event:
	//
	//	*StreamClient_Credentials
	//	*StreamClient_Logout_
	//	*StreamClient_HttpResponse
	//	*StreamClient_Ping_
	Event         isStreamClient_Event `protobuf_oneof:"event"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamClient) Reset() {
	*x = StreamClient{}
	mi := &file_protos_tunnel_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamClient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamClient) ProtoMessage() {}

func (x *StreamClient) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamClient.ProtoReflect.Descriptor instead.
func (*StreamClient) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{0}
}

func (x *StreamClient) GetEvent() isStreamClient_Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *StreamClient) GetCredentials() *StreamClient_Login {
	if x != nil {
		if x, ok := x.Event.(*StreamClient_Credentials); ok {
			return x.Credentials
		}
	}
	return nil
}

func (x *StreamClient) GetLogout() *StreamClient_Logout {
	if x != nil {
		if x, ok := x.Event.(*StreamClient_Logout_); ok {
			return x.Logout
		}
	}
	return nil
}

func (x *StreamClient) GetHttpResponse() []byte {
	if x != nil {
		if x, ok := x.Event.(*StreamClient_HttpResponse); ok {
			return x.HttpResponse
		}
	}
	return nil
}

func (x *StreamClient) GetPing() *StreamClient_Ping {
	if x != nil {
		if x, ok := x.Event.(*StreamClient_Ping_); ok {
			return x.Ping
		}
	}
	return nil
}

type isStreamClient_Event interface {
	isStreamClient_Event()
}

type StreamClient_Credentials struct {
	Credentials *StreamClient_Login `protobuf:"bytes,1,opt,name=credentials,proto3,oneof"`
}

type StreamClient_Logout_ struct {
	Logout *StreamClient_Logout `protobuf:"bytes,2,opt,name=logout,proto3,oneof"`
}

type StreamClient_HttpResponse struct {
	HttpResponse []byte `protobuf:"bytes,3,opt,name=http_response,json=httpResponse,proto3,oneof"`
}

type StreamClient_Ping_ struct {
	// Heartbeat
	Ping *StreamClient_Ping `protobuf:"bytes,4,opt,name=ping,proto3,oneof"`
}

func (*StreamClient_Credentials) isStreamClient_Event() {}

func (*StreamClient_Logout_) isStreamClient_Event() {}

func (*StreamClient_HttpResponse) isStreamClient_Event() {}

func (*StreamClient_Ping_) isStreamClient_Event() {}

type StreamServer struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Event:
	//
	//	*StreamServer_HttpRequest
	//	*StreamServer_Pong_
	Event         isStreamServer_Event `protobuf_oneof:"event"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamServer) Reset() {
	*x = StreamServer{}
	mi := &file_protos_tunnel_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamServer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamServer) ProtoMessage() {}

func (x *StreamServer) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamServer.ProtoReflect.Descriptor instead.
func (*StreamServer) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{1}
}

func (x *StreamServer) GetEvent() isStreamServer_Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *StreamServer) GetHttpRequest() []byte {
	if x != nil {
		if x, ok := x.Event.(*StreamServer_HttpRequest); ok {
			return x.HttpRequest
		}
	}
	return nil
}

func (x *StreamServer) GetPong() *StreamServer_Pong {
	if x != nil {
		if x, ok := x.Event.(*StreamServer_Pong_); ok {
			return x.Pong
		}
	}
	return nil
}

type isStreamServer_Event interface {
	isStreamServer_Event()
}

type StreamServer_HttpRequest struct {
	HttpRequest []byte `protobuf:"bytes,1,opt,name=http_request,json=httpRequest,proto3,oneof"`
}

type StreamServer_Pong_ struct {
	// Heartbeat
	Pong *StreamServer_Pong `protobuf:"bytes,2,opt,name=pong,proto3,oneof"`
}

func (*StreamServer_HttpRequest) isStreamServer_Event() {}

func (*StreamServer_Pong_) isStreamServer_Event() {}

type StreamClient_Login struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamClient_Login) Reset() {
	*x = StreamClient_Login{}
	mi := &file_protos_tunnel_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamClient_Login) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamClient_Login) ProtoMessage() {}

func (x *StreamClient_Login) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamClient_Login.ProtoReflect.Descriptor instead.
func (*StreamClient_Login) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{0, 0}
}

func (x *StreamClient_Login) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type StreamClient_Logout struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamClient_Logout) Reset() {
	*x = StreamClient_Logout{}
	mi := &file_protos_tunnel_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamClient_Logout) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamClient_Logout) ProtoMessage() {}

func (x *StreamClient_Logout) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamClient_Logout.ProtoReflect.Descriptor instead.
func (*StreamClient_Logout) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{0, 1}
}

type StreamClient_HttpOut struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	HttpMessage   []byte                 `protobuf:"bytes,1,opt,name=http_message,json=httpMessage,proto3" json:"http_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamClient_HttpOut) Reset() {
	*x = StreamClient_HttpOut{}
	mi := &file_protos_tunnel_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamClient_HttpOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamClient_HttpOut) ProtoMessage() {}

func (x *StreamClient_HttpOut) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamClient_HttpOut.ProtoReflect.Descriptor instead.
func (*StreamClient_HttpOut) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{0, 2}
}

func (x *StreamClient_HttpOut) GetHttpMessage() []byte {
	if x != nil {
		return x.HttpMessage
	}
	return nil
}

type StreamClient_Ping struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamClient_Ping) Reset() {
	*x = StreamClient_Ping{}
	mi := &file_protos_tunnel_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamClient_Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamClient_Ping) ProtoMessage() {}

func (x *StreamClient_Ping) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamClient_Ping.ProtoReflect.Descriptor instead.
func (*StreamClient_Ping) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{0, 3}
}

type StreamServer_Pong struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamServer_Pong) Reset() {
	*x = StreamServer_Pong{}
	mi := &file_protos_tunnel_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamServer_Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamServer_Pong) ProtoMessage() {}

func (x *StreamServer_Pong) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tunnel_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamServer_Pong.ProtoReflect.Descriptor instead.
func (*StreamServer_Pong) Descriptor() ([]byte, []int) {
	return file_protos_tunnel_proto_rawDescGZIP(), []int{1, 0}
}

var File_protos_tunnel_proto protoreflect.FileDescriptor

var file_protos_tunnel_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0xd2, 0x02,
	0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x3e,
	0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x48,
	0x00, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x35,
	0x0a, 0x06, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x48, 0x00, 0x52, 0x06, 0x6c,
	0x6f, 0x67, 0x6f, 0x75, 0x74, 0x12, 0x25, 0x0a, 0x0d, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0c,
	0x68, 0x74, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x04,
	0x70, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x1a, 0x2a, 0x0a,
	0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x08, 0x0a, 0x06, 0x4c, 0x6f, 0x67,
	0x6f, 0x75, 0x74, 0x1a, 0x2c, 0x0a, 0x07, 0x48, 0x74, 0x74, 0x70, 0x4f, 0x75, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x68, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x06, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x22, 0x75, 0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x12, 0x23, 0x0a, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0b, 0x68, 0x74, 0x74, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x6f, 0x6e, 0x67,
	0x48, 0x00, 0x52, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x1a, 0x06, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67,
	0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x32, 0x44, 0x0a, 0x06, 0x54, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x12, 0x3a, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x14, 0x2e,
	0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x1a, 0x14, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42,
	0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x6d,
	0x65, 0x68, 0x64, 0x69, 0x37, 0x37, 0x37, 0x2f, 0x6d, 0x6f, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_tunnel_proto_rawDescOnce sync.Once
	file_protos_tunnel_proto_rawDescData = file_protos_tunnel_proto_rawDesc
)

func file_protos_tunnel_proto_rawDescGZIP() []byte {
	file_protos_tunnel_proto_rawDescOnce.Do(func() {
		file_protos_tunnel_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_tunnel_proto_rawDescData)
	})
	return file_protos_tunnel_proto_rawDescData
}

var file_protos_tunnel_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protos_tunnel_proto_goTypes = []any{
	(*StreamClient)(nil),         // 0: tunnel.StreamClient
	(*StreamServer)(nil),         // 1: tunnel.StreamServer
	(*StreamClient_Login)(nil),   // 2: tunnel.StreamClient.Login
	(*StreamClient_Logout)(nil),  // 3: tunnel.StreamClient.Logout
	(*StreamClient_HttpOut)(nil), // 4: tunnel.StreamClient.HttpOut
	(*StreamClient_Ping)(nil),    // 5: tunnel.StreamClient.Ping
	(*StreamServer_Pong)(nil),    // 6: tunnel.StreamServer.Pong
}
var file_protos_tunnel_proto_depIdxs = []int32{
	2, // 0: tunnel.StreamClient.credentials:type_name -> tunnel.StreamClient.Login
	3, // 1: tunnel.StreamClient.logout:type_name -> tunnel.StreamClient.Logout
	5, // 2: tunnel.StreamClient.ping:type_name -> tunnel.StreamClient.Ping
	6, // 3: tunnel.StreamServer.pong:type_name -> tunnel.StreamServer.Pong
	0, // 4: tunnel.Tunnel.Stream:input_type -> tunnel.StreamClient
	1, // 5: tunnel.Tunnel.Stream:output_type -> tunnel.StreamServer
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_protos_tunnel_proto_init() }
func file_protos_tunnel_proto_init() {
	if File_protos_tunnel_proto != nil {
		return
	}
	file_protos_tunnel_proto_msgTypes[0].OneofWrappers = []any{
		(*StreamClient_Credentials)(nil),
		(*StreamClient_Logout_)(nil),
		(*StreamClient_HttpResponse)(nil),
		(*StreamClient_Ping_)(nil),
	}
	file_protos_tunnel_proto_msgTypes[1].OneofWrappers = []any{
		(*StreamServer_HttpRequest)(nil),
		(*StreamServer_Pong_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_tunnel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_tunnel_proto_goTypes,
		DependencyIndexes: file_protos_tunnel_proto_depIdxs,
		MessageInfos:      file_protos_tunnel_proto_msgTypes,
	}.Build()
	File_protos_tunnel_proto = out.File
	file_protos_tunnel_proto_rawDesc = nil
	file_protos_tunnel_proto_goTypes = nil
	file_protos_tunnel_proto_depIdxs = nil
}
