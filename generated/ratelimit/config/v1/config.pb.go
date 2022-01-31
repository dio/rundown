// Copyright 2022 Dhi Aurrahman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: ratelimit/config/v1/config.proto

package v1

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

// This is the proto representaion of https://github.com/envoyproxy/ratelimit/blob/8d6488ead8618ce49a492858321dae946f2d97bc/src/config/config_impl.go#L16-L33.
// This is manually syncronized.
// TODO(dio): Parse https://github.com/envoyproxy/ratelimit/blob/main/src/config/config_impl.go and
// generate this message (always to sync with main).
type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domain      string        `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	Descriptors []*Descriptor `protobuf:"bytes,2,rep,name=descriptors,proto3" json:"descriptors,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ratelimit_config_v1_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_ratelimit_config_v1_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_ratelimit_config_v1_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Config) GetDescriptors() []*Descriptor {
	if x != nil {
		return x.Descriptors
	}
	return nil
}

type Descriptor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key         string        `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value       string        `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	RateLimit   *RateLimit    `protobuf:"bytes,3,opt,name=rate_limit,json=rateLimit,proto3" json:"rate_limit,omitempty"`
	Descriptors []*Descriptor `protobuf:"bytes,4,rep,name=descriptors,proto3" json:"descriptors,omitempty"`
	ShadowMode  bool          `protobuf:"varint,5,opt,name=shadow_mode,json=shadowMode,proto3" json:"shadow_mode,omitempty"`
}

func (x *Descriptor) Reset() {
	*x = Descriptor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ratelimit_config_v1_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Descriptor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Descriptor) ProtoMessage() {}

func (x *Descriptor) ProtoReflect() protoreflect.Message {
	mi := &file_ratelimit_config_v1_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Descriptor.ProtoReflect.Descriptor instead.
func (*Descriptor) Descriptor() ([]byte, []int) {
	return file_ratelimit_config_v1_config_proto_rawDescGZIP(), []int{1}
}

func (x *Descriptor) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Descriptor) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Descriptor) GetRateLimit() *RateLimit {
	if x != nil {
		return x.RateLimit
	}
	return nil
}

func (x *Descriptor) GetDescriptors() []*Descriptor {
	if x != nil {
		return x.Descriptors
	}
	return nil
}

func (x *Descriptor) GetShadowMode() bool {
	if x != nil {
		return x.ShadowMode
	}
	return false
}

type RateLimit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestsPerUnit uint32 `protobuf:"varint,1,opt,name=requests_per_unit,json=requestsPerUnit,proto3" json:"requests_per_unit,omitempty"`
	Unit            string `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
	Unlimited       bool   `protobuf:"varint,3,opt,name=unlimited,proto3" json:"unlimited,omitempty"`
}

func (x *RateLimit) Reset() {
	*x = RateLimit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ratelimit_config_v1_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateLimit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateLimit) ProtoMessage() {}

func (x *RateLimit) ProtoReflect() protoreflect.Message {
	mi := &file_ratelimit_config_v1_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateLimit.ProtoReflect.Descriptor instead.
func (*RateLimit) Descriptor() ([]byte, []int) {
	return file_ratelimit_config_v1_config_proto_rawDescGZIP(), []int{2}
}

func (x *RateLimit) GetRequestsPerUnit() uint32 {
	if x != nil {
		return x.RequestsPerUnit
	}
	return 0
}

func (x *RateLimit) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *RateLimit) GetUnlimited() bool {
	if x != nil {
		return x.Unlimited
	}
	return false
}

var File_ratelimit_config_v1_config_proto protoreflect.FileDescriptor

var file_ratelimit_config_v1_config_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x22, 0x63, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x41, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x73, 0x22, 0xd7, 0x01, 0x0a,
	0x0a, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x72, 0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x52, 0x09, 0x72, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x41, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x64, 0x6f, 0x77, 0x5f,
	0x6d, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x64,
	0x6f, 0x77, 0x4d, 0x6f, 0x64, 0x65, 0x22, 0x69, 0x0a, 0x09, 0x52, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x5f,
	0x70, 0x65, 0x72, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x50, 0x65, 0x72, 0x55, 0x6e, 0x69, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x6e, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x75, 0x6e, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65,
	0x64, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x64, 0x69, 0x6f, 0x2f, 0x72, 0x75, 0x6e, 0x64, 0x6f, 0x77, 0x6e, 0x2f, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x72, 0x61, 0x74, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_ratelimit_config_v1_config_proto_rawDescOnce sync.Once
	file_ratelimit_config_v1_config_proto_rawDescData = file_ratelimit_config_v1_config_proto_rawDesc
)

func file_ratelimit_config_v1_config_proto_rawDescGZIP() []byte {
	file_ratelimit_config_v1_config_proto_rawDescOnce.Do(func() {
		file_ratelimit_config_v1_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_ratelimit_config_v1_config_proto_rawDescData)
	})
	return file_ratelimit_config_v1_config_proto_rawDescData
}

var file_ratelimit_config_v1_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_ratelimit_config_v1_config_proto_goTypes = []interface{}{
	(*Config)(nil),     // 0: ratelimit.config.v1.Config
	(*Descriptor)(nil), // 1: ratelimit.config.v1.Descriptor
	(*RateLimit)(nil),  // 2: ratelimit.config.v1.RateLimit
}
var file_ratelimit_config_v1_config_proto_depIdxs = []int32{
	1, // 0: ratelimit.config.v1.Config.descriptors:type_name -> ratelimit.config.v1.Descriptor
	2, // 1: ratelimit.config.v1.Descriptor.rate_limit:type_name -> ratelimit.config.v1.RateLimit
	1, // 2: ratelimit.config.v1.Descriptor.descriptors:type_name -> ratelimit.config.v1.Descriptor
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_ratelimit_config_v1_config_proto_init() }
func file_ratelimit_config_v1_config_proto_init() {
	if File_ratelimit_config_v1_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ratelimit_config_v1_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_ratelimit_config_v1_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Descriptor); i {
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
		file_ratelimit_config_v1_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateLimit); i {
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
			RawDescriptor: file_ratelimit_config_v1_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ratelimit_config_v1_config_proto_goTypes,
		DependencyIndexes: file_ratelimit_config_v1_config_proto_depIdxs,
		MessageInfos:      file_ratelimit_config_v1_config_proto_msgTypes,
	}.Build()
	File_ratelimit_config_v1_config_proto = out.File
	file_ratelimit_config_v1_config_proto_rawDesc = nil
	file_ratelimit_config_v1_config_proto_goTypes = nil
	file_ratelimit_config_v1_config_proto_depIdxs = nil
}
