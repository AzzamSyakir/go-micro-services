// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: user.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email     *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password  *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	Balance   *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=balance,proto3" json:"balance,omitempty"`
	CreatedAt *timestamppb.Timestamp  `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp  `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt *timestamppb.Timestamp  `protobuf:"bytes,8,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *User) GetEmail() *wrapperspb.StringValue {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *User) GetPassword() *wrapperspb.StringValue {
	if x != nil {
		return x.Password
	}
	return nil
}

func (x *User) GetBalance() *wrapperspb.StringValue {
	if x != nil {
		return x.Balance
	}
	return nil
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *User) GetDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

type ById struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ById) Reset() {
	*x = ById{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ById) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ById) ProtoMessage() {}

func (x *ById) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ById.ProtoReflect.Descriptor instead.
func (*ById) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *ById) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ByEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *ByEmail) Reset() {
	*x = ByEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ByEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ByEmail) ProtoMessage() {}

func (x *ByEmail) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ByEmail.ProtoReflect.Descriptor instead.
func (*ByEmail) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *ByEmail) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

type UserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *User  `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *UserResponse) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UserResponse) GetData() *User {
	if x != nil {
		return x.Data
	}
	return nil
}

type UserResponseRepeated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int64   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string  `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*User `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *UserResponseRepeated) Reset() {
	*x = UserResponseRepeated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserResponseRepeated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponseRepeated) ProtoMessage() {}

func (x *UserResponseRepeated) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponseRepeated.ProtoReflect.Descriptor instead.
func (*UserResponseRepeated) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *UserResponseRepeated) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UserResponseRepeated) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UserResponseRepeated) GetData() []*User {
	if x != nil {
		return x.Data
	}
	return nil
}

type Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     *string `protobuf:"bytes,2,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Email    *string `protobuf:"bytes,3,opt,name=email,proto3,oneof" json:"email,omitempty"`
	Password *string `protobuf:"bytes,4,opt,name=password,proto3,oneof" json:"password,omitempty"`
	Balance  *string `protobuf:"bytes,5,opt,name=balance,proto3,oneof" json:"balance,omitempty"`
}

func (x *Update) Reset() {
	*x = Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Update) ProtoMessage() {}

func (x *Update) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Update.ProtoReflect.Descriptor instead.
func (*Update) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *Update) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Update) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Update) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *Update) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

func (x *Update) GetBalance() string {
	if x != nil && x.Balance != nil {
		return *x.Balance
	}
	return ""
}

type Create struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Balance  *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Email    *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Create) Reset() {
	*x = Create{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Create) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Create) ProtoMessage() {}

func (x *Create) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Create.ProtoReflect.Descriptor instead.
func (*Create) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{7}
}

func (x *Create) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *Create) GetBalance() *wrapperspb.StringValue {
	if x != nil {
		return x.Balance
	}
	return nil
}

func (x *Create) GetEmail() *wrapperspb.StringValue {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *Create) GetPassword() *wrapperspb.StringValue {
	if x != nil {
		return x.Password
	}
	return nil
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x67, 0x6f,
	0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x9f, 0x03, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x38, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x36, 0x0a, 0x07, 0x62, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x22, 0x16, 0x0a, 0x04, 0x42, 0x79, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1f, 0x0a, 0x07, 0x42, 0x79,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x07, 0x0a, 0x05, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x69, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x71, 0x0a, 0x14, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0xb8, 0x01, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01,
	0x01, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x88, 0x01,
	0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0xe0, 0x01,
	0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x62, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x38, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x32, 0xcf, 0x03, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x46, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x65, 0x42, 0x79, 0x49, 0x64, 0x12, 0x17,
	0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x42, 0x79, 0x49, 0x64, 0x1a, 0x1f, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63,
	0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4f,
	0x6e, 0x65, 0x42, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x2e, 0x67, 0x6f, 0x5f, 0x6d,
	0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x42, 0x79,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x1f, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a,
	0x1f, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x48, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19,
	0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x1f, 0x2e, 0x67, 0x6f, 0x5f, 0x6d,
	0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0a, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x17, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69,
	0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x42, 0x79, 0x49,
	0x64, 0x1a, 0x1f, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4e, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x18, 0x2e, 0x67, 0x6f, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x27, 0x2e, 0x67, 0x6f, 0x5f, 0x6d,
	0x69, 0x63, 0x72, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x42, 0x05, 0x5a, 0x03, 0x70, 0x62, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_user_proto_goTypes = []interface{}{
	(*User)(nil),                   // 0: go_micro_services.User
	(*ById)(nil),                   // 1: go_micro_services.ById
	(*ByEmail)(nil),                // 2: go_micro_services.ByEmail
	(*Empty)(nil),                  // 3: go_micro_services.empty
	(*UserResponse)(nil),           // 4: go_micro_services.UserResponse
	(*UserResponseRepeated)(nil),   // 5: go_micro_services.UserResponseRepeated
	(*Update)(nil),                 // 6: go_micro_services.Update
	(*Create)(nil),                 // 7: go_micro_services.Create
	(*wrapperspb.StringValue)(nil), // 8: google.protobuf.StringValue
	(*timestamppb.Timestamp)(nil),  // 9: google.protobuf.Timestamp
}
var file_user_proto_depIdxs = []int32{
	8,  // 0: go_micro_services.User.name:type_name -> google.protobuf.StringValue
	8,  // 1: go_micro_services.User.email:type_name -> google.protobuf.StringValue
	8,  // 2: go_micro_services.User.password:type_name -> google.protobuf.StringValue
	8,  // 3: go_micro_services.User.balance:type_name -> google.protobuf.StringValue
	9,  // 4: go_micro_services.User.created_at:type_name -> google.protobuf.Timestamp
	9,  // 5: go_micro_services.User.updated_at:type_name -> google.protobuf.Timestamp
	9,  // 6: go_micro_services.User.deleted_at:type_name -> google.protobuf.Timestamp
	0,  // 7: go_micro_services.UserResponse.data:type_name -> go_micro_services.User
	0,  // 8: go_micro_services.UserResponseRepeated.data:type_name -> go_micro_services.User
	8,  // 9: go_micro_services.Create.name:type_name -> google.protobuf.StringValue
	8,  // 10: go_micro_services.Create.balance:type_name -> google.protobuf.StringValue
	8,  // 11: go_micro_services.Create.email:type_name -> google.protobuf.StringValue
	8,  // 12: go_micro_services.Create.password:type_name -> google.protobuf.StringValue
	1,  // 13: go_micro_services.UserService.GetOneById:input_type -> go_micro_services.ById
	2,  // 14: go_micro_services.UserService.GetOneByEmail:input_type -> go_micro_services.ByEmail
	6,  // 15: go_micro_services.UserService.UpdateUser:input_type -> go_micro_services.Update
	7,  // 16: go_micro_services.UserService.CreateUser:input_type -> go_micro_services.Create
	1,  // 17: go_micro_services.UserService.DeleteUser:input_type -> go_micro_services.ById
	3,  // 18: go_micro_services.UserService.ListUsers:input_type -> go_micro_services.empty
	4,  // 19: go_micro_services.UserService.GetOneById:output_type -> go_micro_services.UserResponse
	4,  // 20: go_micro_services.UserService.GetOneByEmail:output_type -> go_micro_services.UserResponse
	4,  // 21: go_micro_services.UserService.UpdateUser:output_type -> go_micro_services.UserResponse
	4,  // 22: go_micro_services.UserService.CreateUser:output_type -> go_micro_services.UserResponse
	4,  // 23: go_micro_services.UserService.DeleteUser:output_type -> go_micro_services.UserResponse
	5,  // 24: go_micro_services.UserService.ListUsers:output_type -> go_micro_services.UserResponseRepeated
	19, // [19:25] is the sub-list for method output_type
	13, // [13:19] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ById); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ByEmail); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserResponse); i {
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
		file_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserResponseRepeated); i {
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
		file_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Update); i {
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
		file_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Create); i {
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
	file_user_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}