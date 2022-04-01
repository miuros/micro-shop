// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: api/comment/v1/comment.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Comment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId  int64  `protobuf:"varint,2,opt,name=productId,proto3" json:"productId,omitempty"`
	UserUuid   string `protobuf:"bytes,3,opt,name=userUuid,proto3" json:"userUuid,omitempty"`
	ToUserUuid string `protobuf:"bytes,4,opt,name=toUserUuid,proto3" json:"toUserUuid,omitempty"`
	Content    string `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	CreateAt   string `protobuf:"bytes,6,opt,name=createAt,proto3" json:"createAt,omitempty"`
	UpdateAt   string `protobuf:"bytes,7,opt,name=updateAt,proto3" json:"updateAt,omitempty"`
	DeleteAt   string `protobuf:"bytes,8,opt,name=deleteAt,proto3" json:"deleteAt,omitempty"`
	IsDeleted  int64  `protobuf:"varint,9,opt,name=isDeleted,proto3" json:"isDeleted,omitempty"`
}

func (x *Comment) Reset() {
	*x = Comment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Comment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comment) ProtoMessage() {}

func (x *Comment) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comment.ProtoReflect.Descriptor instead.
func (*Comment) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{0}
}

func (x *Comment) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Comment) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *Comment) GetUserUuid() string {
	if x != nil {
		return x.UserUuid
	}
	return ""
}

func (x *Comment) GetToUserUuid() string {
	if x != nil {
		return x.ToUserUuid
	}
	return ""
}

func (x *Comment) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Comment) GetCreateAt() string {
	if x != nil {
		return x.CreateAt
	}
	return ""
}

func (x *Comment) GetUpdateAt() string {
	if x != nil {
		return x.UpdateAt
	}
	return ""
}

func (x *Comment) GetDeleteAt() string {
	if x != nil {
		return x.DeleteAt
	}
	return ""
}

func (x *Comment) GetIsDeleted() int64 {
	if x != nil {
		return x.IsDeleted
	}
	return 0
}

type CreateCmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cm *Comment `protobuf:"bytes,1,opt,name=cm,proto3" json:"cm,omitempty"`
}

func (x *CreateCmReq) Reset() {
	*x = CreateCmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCmReq) ProtoMessage() {}

func (x *CreateCmReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCmReq.ProtoReflect.Descriptor instead.
func (*CreateCmReq) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{1}
}

func (x *CreateCmReq) GetCm() *Comment {
	if x != nil {
		return x.Cm
	}
	return nil
}

type CreateCmReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cm *Comment `protobuf:"bytes,1,opt,name=cm,proto3" json:"cm,omitempty"`
}

func (x *CreateCmReply) Reset() {
	*x = CreateCmReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCmReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCmReply) ProtoMessage() {}

func (x *CreateCmReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCmReply.ProtoReflect.Descriptor instead.
func (*CreateCmReply) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{2}
}

func (x *CreateCmReply) GetCm() *Comment {
	if x != nil {
		return x.Cm
	}
	return nil
}

type UpdateCmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cm *Comment `protobuf:"bytes,1,opt,name=cm,proto3" json:"cm,omitempty"`
}

func (x *UpdateCmReq) Reset() {
	*x = UpdateCmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCmReq) ProtoMessage() {}

func (x *UpdateCmReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCmReq.ProtoReflect.Descriptor instead.
func (*UpdateCmReq) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateCmReq) GetCm() *Comment {
	if x != nil {
		return x.Cm
	}
	return nil
}

type UpdateCmReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cm *Comment `protobuf:"bytes,1,opt,name=cm,proto3" json:"cm,omitempty"`
}

func (x *UpdateCmReply) Reset() {
	*x = UpdateCmReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCmReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCmReply) ProtoMessage() {}

func (x *UpdateCmReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCmReply.ProtoReflect.Descriptor instead.
func (*UpdateCmReply) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateCmReply) GetCm() *Comment {
	if x != nil {
		return x.Cm
	}
	return nil
}

type DeleteCmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserUuid string `protobuf:"bytes,2,opt,name=userUuid,proto3" json:"userUuid,omitempty"`
}

func (x *DeleteCmReq) Reset() {
	*x = DeleteCmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCmReq) ProtoMessage() {}

func (x *DeleteCmReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCmReq.ProtoReflect.Descriptor instead.
func (*DeleteCmReq) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteCmReq) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeleteCmReq) GetUserUuid() string {
	if x != nil {
		return x.UserUuid
	}
	return ""
}

type DeleteCmReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteCmReply) Reset() {
	*x = DeleteCmReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCmReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCmReply) ProtoMessage() {}

func (x *DeleteCmReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCmReply.ProtoReflect.Descriptor instead.
func (*DeleteCmReply) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{6}
}

type GetCmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetCmReq) Reset() {
	*x = GetCmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCmReq) ProtoMessage() {}

func (x *GetCmReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCmReq.ProtoReflect.Descriptor instead.
func (*GetCmReq) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{7}
}

func (x *GetCmReq) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetCmReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cm *Comment `protobuf:"bytes,1,opt,name=cm,proto3" json:"cm,omitempty"`
}

func (x *GetCmReply) Reset() {
	*x = GetCmReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCmReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCmReply) ProtoMessage() {}

func (x *GetCmReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCmReply.ProtoReflect.Descriptor instead.
func (*GetCmReply) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{8}
}

func (x *GetCmReply) GetCm() *Comment {
	if x != nil {
		return x.Cm
	}
	return nil
}

type ListCmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page      int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit     int64 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	ProductId int64 `protobuf:"varint,3,opt,name=productId,proto3" json:"productId,omitempty"`
}

func (x *ListCmReq) Reset() {
	*x = ListCmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCmReq) ProtoMessage() {}

func (x *ListCmReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCmReq.ProtoReflect.Descriptor instead.
func (*ListCmReq) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{9}
}

func (x *ListCmReq) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCmReq) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListCmReq) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

type ListCmReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CmList []*Comment `protobuf:"bytes,1,rep,name=CmList,proto3" json:"CmList,omitempty"`
}

func (x *ListCmReply) Reset() {
	*x = ListCmReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_comment_v1_comment_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCmReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCmReply) ProtoMessage() {}

func (x *ListCmReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_comment_v1_comment_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCmReply.ProtoReflect.Descriptor instead.
func (*ListCmReply) Descriptor() ([]byte, []int) {
	return file_api_comment_v1_comment_proto_rawDescGZIP(), []int{10}
}

func (x *ListCmReply) GetCmList() []*Comment {
	if x != nil {
		return x.CmList
	}
	return nil
}

var File_api_comment_v1_comment_proto protoreflect.FileDescriptor

var file_api_comment_v1_comment_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xff, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x55, 0x73, 0x65, 0x72,
	0x55, 0x75, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x55, 0x73,
	0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x41, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x22, 0x2a, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65,
	0x71, 0x12, 0x1b, 0x0a, 0x02, 0x63, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x02, 0x63, 0x6d, 0x22, 0x2c,
	0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x1b, 0x0a, 0x02, 0x63, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x02, 0x63, 0x6d, 0x22, 0x2a, 0x0a, 0x0b,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x02, 0x63,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x02, 0x63, 0x6d, 0x22, 0x2c, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x02, 0x63, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x02, 0x63, 0x6d, 0x22, 0x39, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x55, 0x75, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x55, 0x75, 0x69,
	0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x1a, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x29,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x02,
	0x63, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x02, 0x63, 0x6d, 0x22, 0x53, 0x0a, 0x09, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x22, 0x32,
	0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a,
	0x06, 0x43, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x43, 0x6d, 0x4c, 0x69,
	0x73, 0x74, 0x32, 0xae, 0x03, 0x0a, 0x09, 0x43, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4e, 0x0a, 0x08, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x12, 0x0f, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x3a, 0x01, 0x2a,
	0x12, 0x4e, 0x0a, 0x08, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x12, 0x0f, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x1a, 0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x01, 0x2a,
	0x12, 0x5b, 0x0a, 0x08, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6d, 0x12, 0x0f, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x2a, 0x23, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x7d, 0x12, 0x44, 0x0a,
	0x05, 0x47, 0x65, 0x74, 0x43, 0x6d, 0x12, 0x0c, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43,
	0x6d, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6d, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x67, 0x65, 0x74, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x12, 0x5e, 0x0a, 0x06, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6d, 0x12, 0x0d, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x34, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x12, 0x2c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x7b, 0x70, 0x61, 0x67, 0x65, 0x7d, 0x2f,
	0x7b, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x49, 0x64, 0x7d, 0x42, 0x10, 0x5a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_comment_v1_comment_proto_rawDescOnce sync.Once
	file_api_comment_v1_comment_proto_rawDescData = file_api_comment_v1_comment_proto_rawDesc
)

func file_api_comment_v1_comment_proto_rawDescGZIP() []byte {
	file_api_comment_v1_comment_proto_rawDescOnce.Do(func() {
		file_api_comment_v1_comment_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_comment_v1_comment_proto_rawDescData)
	})
	return file_api_comment_v1_comment_proto_rawDescData
}

var file_api_comment_v1_comment_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_comment_v1_comment_proto_goTypes = []interface{}{
	(*Comment)(nil),       // 0: v1.Comment
	(*CreateCmReq)(nil),   // 1: v1.CreateCmReq
	(*CreateCmReply)(nil), // 2: v1.CreateCmReply
	(*UpdateCmReq)(nil),   // 3: v1.UpdateCmReq
	(*UpdateCmReply)(nil), // 4: v1.UpdateCmReply
	(*DeleteCmReq)(nil),   // 5: v1.DeleteCmReq
	(*DeleteCmReply)(nil), // 6: v1.DeleteCmReply
	(*GetCmReq)(nil),      // 7: v1.GetCmReq
	(*GetCmReply)(nil),    // 8: v1.GetCmReply
	(*ListCmReq)(nil),     // 9: v1.ListCmReq
	(*ListCmReply)(nil),   // 10: v1.ListCmReply
}
var file_api_comment_v1_comment_proto_depIdxs = []int32{
	0,  // 0: v1.CreateCmReq.cm:type_name -> v1.Comment
	0,  // 1: v1.CreateCmReply.cm:type_name -> v1.Comment
	0,  // 2: v1.UpdateCmReq.cm:type_name -> v1.Comment
	0,  // 3: v1.UpdateCmReply.cm:type_name -> v1.Comment
	0,  // 4: v1.GetCmReply.cm:type_name -> v1.Comment
	0,  // 5: v1.ListCmReply.CmList:type_name -> v1.Comment
	1,  // 6: v1.CmService.CreateCm:input_type -> v1.CreateCmReq
	3,  // 7: v1.CmService.UpdateCm:input_type -> v1.UpdateCmReq
	5,  // 8: v1.CmService.DeleteCm:input_type -> v1.DeleteCmReq
	7,  // 9: v1.CmService.GetCm:input_type -> v1.GetCmReq
	9,  // 10: v1.CmService.ListCm:input_type -> v1.ListCmReq
	2,  // 11: v1.CmService.CreateCm:output_type -> v1.CreateCmReply
	4,  // 12: v1.CmService.UpdateCm:output_type -> v1.UpdateCmReply
	6,  // 13: v1.CmService.DeleteCm:output_type -> v1.DeleteCmReply
	8,  // 14: v1.CmService.GetCm:output_type -> v1.GetCmReply
	10, // 15: v1.CmService.ListCm:output_type -> v1.ListCmReply
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_comment_v1_comment_proto_init() }
func file_api_comment_v1_comment_proto_init() {
	if File_api_comment_v1_comment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_comment_v1_comment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Comment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCmReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCmReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCmReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCmReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_comment_v1_comment_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCmReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_comment_v1_comment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_comment_v1_comment_proto_goTypes,
		DependencyIndexes: file_api_comment_v1_comment_proto_depIdxs,
		MessageInfos:      file_api_comment_v1_comment_proto_msgTypes,
	}.Build()
	File_api_comment_v1_comment_proto = out.File
	file_api_comment_v1_comment_proto_rawDesc = nil
	file_api_comment_v1_comment_proto_goTypes = nil
	file_api_comment_v1_comment_proto_depIdxs = nil
}