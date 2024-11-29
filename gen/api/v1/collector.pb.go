// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/v1/collector.proto

package gen

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Define a message representing an auth information about user and team.
type Auth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                      // Unique identifer for user that is processing the data
	TeamId      string  `protobuf:"bytes,2,opt,name=team_id,json=teamId,proto3" json:"team_id,omitempty"`                      // Unique identifier for users team
	WorkspaceId *string `protobuf:"bytes,3,opt,name=workspace_id,json=workspaceId,proto3,oneof" json:"workspace_id,omitempty"` // Unique identifier of the Workspace that is running the request
	UserEmail   string  `protobuf:"bytes,4,opt,name=user_email,json=userEmail,proto3" json:"user_email,omitempty"`             // Unique identifier of user that is processing the data
}

func (x *Auth) Reset() {
	*x = Auth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_collector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth) ProtoMessage() {}

func (x *Auth) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_collector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth.ProtoReflect.Descriptor instead.
func (*Auth) Descriptor() ([]byte, []int) {
	return file_api_v1_collector_proto_rawDescGZIP(), []int{0}
}

func (x *Auth) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Auth) GetTeamId() string {
	if x != nil {
		return x.TeamId
	}
	return ""
}

func (x *Auth) GetWorkspaceId() string {
	if x != nil && x.WorkspaceId != nil {
		return *x.WorkspaceId
	}
	return ""
}

func (x *Auth) GetUserEmail() string {
	if x != nil {
		return x.UserEmail
	}
	return ""
}

// Define a message representing a command, including its metadata and timing information.
type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                            // Unique identifier for the command.
	Category      string `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`                                 // Category of the command (e.g., system, user).
	Command       string `protobuf:"bytes,3,opt,name=command,proto3" json:"command,omitempty"`                                   // The actual command string.
	User          string `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`                                         // The user who executed the command.
	Directory     string `protobuf:"bytes,5,opt,name=directory,proto3" json:"directory,omitempty"`                               // The directory from which the command was executed.
	ExecutionTime int64  `protobuf:"varint,6,opt,name=execution_time,json=executionTime,proto3" json:"execution_time,omitempty"` // Execution time of the command in milliseconds.
	StartTime     int64  `protobuf:"varint,7,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`             // Start time of the command execution (Unix timestamp).
	EndTime       int64  `protobuf:"varint,8,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`                   // End time of the command execution (Unix timestamp).
	Result        string `protobuf:"bytes,9,opt,name=result,proto3" json:"result,omitempty"`                                     // Result of executed command => success/failure
	Status        string `protobuf:"bytes,10,opt,name=status,proto3" json:"status,omitempty"`                                    // Status of executed command
	Repository    string `protobuf:"bytes,11,opt,name=repository,proto3" json:"repository,omitempty"`                            // Repository is repository where commands are executed
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_collector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_collector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_api_v1_collector_proto_rawDescGZIP(), []int{1}
}

func (x *Command) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Command) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Command) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *Command) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *Command) GetDirectory() string {
	if x != nil {
		return x.Directory
	}
	return ""
}

func (x *Command) GetExecutionTime() int64 {
	if x != nil {
		return x.ExecutionTime
	}
	return 0
}

func (x *Command) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *Command) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *Command) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *Command) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Command) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

// Define a message representing a process, including its metadata and resource usage.
type Process struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                              // Unique identifier for the process.
	Pid            int64   `protobuf:"varint,2,opt,name=pid,proto3" json:"pid,omitempty"`                                            // Process ID.
	Name           string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`                                           // Process name.
	Status         string  `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`                                       // Current status of the process (e.g., running, sleeping).
	CreatedTime    int64   `protobuf:"varint,5,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`         // Creation time of the process (Unix timestamp).
	StoredTime     int64   `protobuf:"varint,6,opt,name=stored_time,json=storedTime,proto3" json:"stored_time,omitempty"`            // Time at which the process information was stored (Unix timestamp).
	Os             string  `protobuf:"bytes,7,opt,name=os,proto3" json:"os,omitempty"`                                               // Operating system the process is running on.
	Platform       string  `protobuf:"bytes,8,opt,name=platform,proto3" json:"platform,omitempty"`                                   // Platform information (e.g., Linux, Windows).
	PlatformFamily string  `protobuf:"bytes,9,opt,name=platform_family,json=platformFamily,proto3" json:"platform_family,omitempty"` // More detailed platform family information.
	CpuUsage       float64 `protobuf:"fixed64,10,opt,name=cpu_usage,json=cpuUsage,proto3" json:"cpu_usage,omitempty"`                // CPU usage percentage by the process.
	MemoryUsage    float64 `protobuf:"fixed64,11,opt,name=memory_usage,json=memoryUsage,proto3" json:"memory_usage,omitempty"`       // Memory usage by the process in megabytes.
	Ppid           int64   `protobuf:"varint,12,opt,name=ppid,proto3" json:"ppid,omitempty"`                                         // Parent process ID.
}

func (x *Process) Reset() {
	*x = Process{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_collector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Process) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Process) ProtoMessage() {}

func (x *Process) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_collector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Process.ProtoReflect.Descriptor instead.
func (*Process) Descriptor() ([]byte, []int) {
	return file_api_v1_collector_proto_rawDescGZIP(), []int{2}
}

func (x *Process) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Process) GetPid() int64 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *Process) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Process) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Process) GetCreatedTime() int64 {
	if x != nil {
		return x.CreatedTime
	}
	return 0
}

func (x *Process) GetStoredTime() int64 {
	if x != nil {
		return x.StoredTime
	}
	return 0
}

func (x *Process) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *Process) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *Process) GetPlatformFamily() string {
	if x != nil {
		return x.PlatformFamily
	}
	return ""
}

func (x *Process) GetCpuUsage() float64 {
	if x != nil {
		return x.CpuUsage
	}
	return 0
}

func (x *Process) GetMemoryUsage() float64 {
	if x != nil {
		return x.MemoryUsage
	}
	return 0
}

func (x *Process) GetPpid() int64 {
	if x != nil {
		return x.Ppid
	}
	return 0
}

// Defines a request for sending a collection of commands.
type SendCommandsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Commands []*Command `protobuf:"bytes,1,rep,name=commands,proto3" json:"commands,omitempty"` // A list of commands.
	Auth     *Auth      `protobuf:"bytes,2,opt,name=auth,proto3,oneof" json:"auth,omitempty"`   // Optional auth configuration
}

func (x *SendCommandsRequest) Reset() {
	*x = SendCommandsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_collector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendCommandsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCommandsRequest) ProtoMessage() {}

func (x *SendCommandsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_collector_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCommandsRequest.ProtoReflect.Descriptor instead.
func (*SendCommandsRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_collector_proto_rawDescGZIP(), []int{3}
}

func (x *SendCommandsRequest) GetCommands() []*Command {
	if x != nil {
		return x.Commands
	}
	return nil
}

func (x *SendCommandsRequest) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

// Defines a request for sending a collection of processes.
type SendProcessesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Processes []*Process `protobuf:"bytes,1,rep,name=processes,proto3" json:"processes,omitempty"` // A list of processes.
	Auth      *Auth      `protobuf:"bytes,2,opt,name=auth,proto3,oneof" json:"auth,omitempty"`     // Optional auth configuration
}

func (x *SendProcessesRequest) Reset() {
	*x = SendProcessesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_collector_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendProcessesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendProcessesRequest) ProtoMessage() {}

func (x *SendProcessesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_collector_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendProcessesRequest.ProtoReflect.Descriptor instead.
func (*SendProcessesRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_collector_proto_rawDescGZIP(), []int{4}
}

func (x *SendProcessesRequest) GetProcesses() []*Process {
	if x != nil {
		return x.Processes
	}
	return nil
}

func (x *SendProcessesRequest) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

var File_api_v1_collector_proto protoreflect.FileDescriptor

var file_api_v1_collector_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01,
	0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x74, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x74, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0c, 0x77, 0x6f, 0x72, 0x6b,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x88, 0x01, 0x01,
	0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x42,
	0x0f, 0x0a, 0x0d, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x22, 0xb2, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e,
	0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6e,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x6f, 0x72, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x6f, 0x72, 0x79, 0x22, 0xc4, 0x02, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03,
	0x70, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x21, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x6f, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12,
	0x27, 0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x66, 0x61, 0x6d, 0x69,
	0x6c, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x70, 0x75, 0x5f,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x63, 0x70, 0x75,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x6d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x70, 0x69, 0x64,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x70, 0x69, 0x64, 0x22, 0x72, 0x0a, 0x13,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73,
	0x12, 0x25, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x48, 0x00, 0x52, 0x04,
	0x61, 0x75, 0x74, 0x68, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x61, 0x75, 0x74, 0x68,
	0x22, 0x75, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x48, 0x00, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x88, 0x01, 0x01, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x32, 0x9e, 0x01, 0x0a, 0x10, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0c,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x1b, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x45, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x12, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x4f, 0x0a, 0x0a, 0x67, 0x65, 0x6e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x76, 0x7a, 0x65, 0x72, 0x6f, 0x2d, 0x69, 0x6e, 0x63,
	0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x64, 0x65, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x72,
	0x2d, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x74, 0x69, 0x63, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_v1_collector_proto_rawDescOnce sync.Once
	file_api_v1_collector_proto_rawDescData = file_api_v1_collector_proto_rawDesc
)

func file_api_v1_collector_proto_rawDescGZIP() []byte {
	file_api_v1_collector_proto_rawDescOnce.Do(func() {
		file_api_v1_collector_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_collector_proto_rawDescData)
	})
	return file_api_v1_collector_proto_rawDescData
}

var file_api_v1_collector_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_v1_collector_proto_goTypes = []interface{}{
	(*Auth)(nil),                 // 0: api.v1.Auth
	(*Command)(nil),              // 1: api.v1.Command
	(*Process)(nil),              // 2: api.v1.Process
	(*SendCommandsRequest)(nil),  // 3: api.v1.SendCommandsRequest
	(*SendProcessesRequest)(nil), // 4: api.v1.SendProcessesRequest
	(*emptypb.Empty)(nil),        // 5: google.protobuf.Empty
}
var file_api_v1_collector_proto_depIdxs = []int32{
	1, // 0: api.v1.SendCommandsRequest.commands:type_name -> api.v1.Command
	0, // 1: api.v1.SendCommandsRequest.auth:type_name -> api.v1.Auth
	2, // 2: api.v1.SendProcessesRequest.processes:type_name -> api.v1.Process
	0, // 3: api.v1.SendProcessesRequest.auth:type_name -> api.v1.Auth
	3, // 4: api.v1.CollectorService.SendCommands:input_type -> api.v1.SendCommandsRequest
	4, // 5: api.v1.CollectorService.SendProcesses:input_type -> api.v1.SendProcessesRequest
	5, // 6: api.v1.CollectorService.SendCommands:output_type -> google.protobuf.Empty
	5, // 7: api.v1.CollectorService.SendProcesses:output_type -> google.protobuf.Empty
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1_collector_proto_init() }
func file_api_v1_collector_proto_init() {
	if File_api_v1_collector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_collector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth); i {
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
		file_api_v1_collector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
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
		file_api_v1_collector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Process); i {
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
		file_api_v1_collector_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendCommandsRequest); i {
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
		file_api_v1_collector_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendProcessesRequest); i {
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
	file_api_v1_collector_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_api_v1_collector_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_api_v1_collector_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_collector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_collector_proto_goTypes,
		DependencyIndexes: file_api_v1_collector_proto_depIdxs,
		MessageInfos:      file_api_v1_collector_proto_msgTypes,
	}.Build()
	File_api_v1_collector_proto = out.File
	file_api_v1_collector_proto_rawDesc = nil
	file_api_v1_collector_proto_goTypes = nil
	file_api_v1_collector_proto_depIdxs = nil
}