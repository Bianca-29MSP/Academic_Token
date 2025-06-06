// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package schedule

import (
	fmt "fmt"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_RecommendedSubject                    protoreflect.MessageDescriptor
	fd_RecommendedSubject_subjectId          protoreflect.FieldDescriptor
	fd_RecommendedSubject_recommendationRank protoreflect.FieldDescriptor
	fd_RecommendedSubject_reason             protoreflect.FieldDescriptor
	fd_RecommendedSubject_isRequired         protoreflect.FieldDescriptor
	fd_RecommendedSubject_semesterAlignment  protoreflect.FieldDescriptor
	fd_RecommendedSubject_difficultyLevel    protoreflect.FieldDescriptor
)

func init() {
	file_academictoken_schedule_recommended_subject_proto_init()
	md_RecommendedSubject = File_academictoken_schedule_recommended_subject_proto.Messages().ByName("RecommendedSubject")
	fd_RecommendedSubject_subjectId = md_RecommendedSubject.Fields().ByName("subjectId")
	fd_RecommendedSubject_recommendationRank = md_RecommendedSubject.Fields().ByName("recommendationRank")
	fd_RecommendedSubject_reason = md_RecommendedSubject.Fields().ByName("reason")
	fd_RecommendedSubject_isRequired = md_RecommendedSubject.Fields().ByName("isRequired")
	fd_RecommendedSubject_semesterAlignment = md_RecommendedSubject.Fields().ByName("semesterAlignment")
	fd_RecommendedSubject_difficultyLevel = md_RecommendedSubject.Fields().ByName("difficultyLevel")
}

var _ protoreflect.Message = (*fastReflection_RecommendedSubject)(nil)

type fastReflection_RecommendedSubject RecommendedSubject

func (x *RecommendedSubject) ProtoReflect() protoreflect.Message {
	return (*fastReflection_RecommendedSubject)(x)
}

func (x *RecommendedSubject) slowProtoReflect() protoreflect.Message {
	mi := &file_academictoken_schedule_recommended_subject_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_RecommendedSubject_messageType fastReflection_RecommendedSubject_messageType
var _ protoreflect.MessageType = fastReflection_RecommendedSubject_messageType{}

type fastReflection_RecommendedSubject_messageType struct{}

func (x fastReflection_RecommendedSubject_messageType) Zero() protoreflect.Message {
	return (*fastReflection_RecommendedSubject)(nil)
}
func (x fastReflection_RecommendedSubject_messageType) New() protoreflect.Message {
	return new(fastReflection_RecommendedSubject)
}
func (x fastReflection_RecommendedSubject_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_RecommendedSubject
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_RecommendedSubject) Descriptor() protoreflect.MessageDescriptor {
	return md_RecommendedSubject
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_RecommendedSubject) Type() protoreflect.MessageType {
	return _fastReflection_RecommendedSubject_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_RecommendedSubject) New() protoreflect.Message {
	return new(fastReflection_RecommendedSubject)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_RecommendedSubject) Interface() protoreflect.ProtoMessage {
	return (*RecommendedSubject)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_RecommendedSubject) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.SubjectId != "" {
		value := protoreflect.ValueOfString(x.SubjectId)
		if !f(fd_RecommendedSubject_subjectId, value) {
			return
		}
	}
	if x.RecommendationRank != "" {
		value := protoreflect.ValueOfString(x.RecommendationRank)
		if !f(fd_RecommendedSubject_recommendationRank, value) {
			return
		}
	}
	if x.Reason != "" {
		value := protoreflect.ValueOfString(x.Reason)
		if !f(fd_RecommendedSubject_reason, value) {
			return
		}
	}
	if x.IsRequired != "" {
		value := protoreflect.ValueOfString(x.IsRequired)
		if !f(fd_RecommendedSubject_isRequired, value) {
			return
		}
	}
	if x.SemesterAlignment != "" {
		value := protoreflect.ValueOfString(x.SemesterAlignment)
		if !f(fd_RecommendedSubject_semesterAlignment, value) {
			return
		}
	}
	if x.DifficultyLevel != "" {
		value := protoreflect.ValueOfString(x.DifficultyLevel)
		if !f(fd_RecommendedSubject_difficultyLevel, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_RecommendedSubject) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		return x.SubjectId != ""
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		return x.RecommendationRank != ""
	case "academictoken.schedule.RecommendedSubject.reason":
		return x.Reason != ""
	case "academictoken.schedule.RecommendedSubject.isRequired":
		return x.IsRequired != ""
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		return x.SemesterAlignment != ""
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		return x.DifficultyLevel != ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_RecommendedSubject) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		x.SubjectId = ""
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		x.RecommendationRank = ""
	case "academictoken.schedule.RecommendedSubject.reason":
		x.Reason = ""
	case "academictoken.schedule.RecommendedSubject.isRequired":
		x.IsRequired = ""
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		x.SemesterAlignment = ""
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		x.DifficultyLevel = ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_RecommendedSubject) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		value := x.SubjectId
		return protoreflect.ValueOfString(value)
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		value := x.RecommendationRank
		return protoreflect.ValueOfString(value)
	case "academictoken.schedule.RecommendedSubject.reason":
		value := x.Reason
		return protoreflect.ValueOfString(value)
	case "academictoken.schedule.RecommendedSubject.isRequired":
		value := x.IsRequired
		return protoreflect.ValueOfString(value)
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		value := x.SemesterAlignment
		return protoreflect.ValueOfString(value)
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		value := x.DifficultyLevel
		return protoreflect.ValueOfString(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_RecommendedSubject) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		x.SubjectId = value.Interface().(string)
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		x.RecommendationRank = value.Interface().(string)
	case "academictoken.schedule.RecommendedSubject.reason":
		x.Reason = value.Interface().(string)
	case "academictoken.schedule.RecommendedSubject.isRequired":
		x.IsRequired = value.Interface().(string)
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		x.SemesterAlignment = value.Interface().(string)
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		x.DifficultyLevel = value.Interface().(string)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_RecommendedSubject) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		panic(fmt.Errorf("field subjectId of message academictoken.schedule.RecommendedSubject is not mutable"))
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		panic(fmt.Errorf("field recommendationRank of message academictoken.schedule.RecommendedSubject is not mutable"))
	case "academictoken.schedule.RecommendedSubject.reason":
		panic(fmt.Errorf("field reason of message academictoken.schedule.RecommendedSubject is not mutable"))
	case "academictoken.schedule.RecommendedSubject.isRequired":
		panic(fmt.Errorf("field isRequired of message academictoken.schedule.RecommendedSubject is not mutable"))
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		panic(fmt.Errorf("field semesterAlignment of message academictoken.schedule.RecommendedSubject is not mutable"))
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		panic(fmt.Errorf("field difficultyLevel of message academictoken.schedule.RecommendedSubject is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_RecommendedSubject) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "academictoken.schedule.RecommendedSubject.subjectId":
		return protoreflect.ValueOfString("")
	case "academictoken.schedule.RecommendedSubject.recommendationRank":
		return protoreflect.ValueOfString("")
	case "academictoken.schedule.RecommendedSubject.reason":
		return protoreflect.ValueOfString("")
	case "academictoken.schedule.RecommendedSubject.isRequired":
		return protoreflect.ValueOfString("")
	case "academictoken.schedule.RecommendedSubject.semesterAlignment":
		return protoreflect.ValueOfString("")
	case "academictoken.schedule.RecommendedSubject.difficultyLevel":
		return protoreflect.ValueOfString("")
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: academictoken.schedule.RecommendedSubject"))
		}
		panic(fmt.Errorf("message academictoken.schedule.RecommendedSubject does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_RecommendedSubject) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in academictoken.schedule.RecommendedSubject", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_RecommendedSubject) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_RecommendedSubject) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_RecommendedSubject) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_RecommendedSubject) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*RecommendedSubject)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		l = len(x.SubjectId)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.RecommendationRank)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.Reason)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.IsRequired)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.SemesterAlignment)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.DifficultyLevel)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*RecommendedSubject)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if len(x.DifficultyLevel) > 0 {
			i -= len(x.DifficultyLevel)
			copy(dAtA[i:], x.DifficultyLevel)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.DifficultyLevel)))
			i--
			dAtA[i] = 0x32
		}
		if len(x.SemesterAlignment) > 0 {
			i -= len(x.SemesterAlignment)
			copy(dAtA[i:], x.SemesterAlignment)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.SemesterAlignment)))
			i--
			dAtA[i] = 0x2a
		}
		if len(x.IsRequired) > 0 {
			i -= len(x.IsRequired)
			copy(dAtA[i:], x.IsRequired)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.IsRequired)))
			i--
			dAtA[i] = 0x22
		}
		if len(x.Reason) > 0 {
			i -= len(x.Reason)
			copy(dAtA[i:], x.Reason)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Reason)))
			i--
			dAtA[i] = 0x1a
		}
		if len(x.RecommendationRank) > 0 {
			i -= len(x.RecommendationRank)
			copy(dAtA[i:], x.RecommendationRank)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.RecommendationRank)))
			i--
			dAtA[i] = 0x12
		}
		if len(x.SubjectId) > 0 {
			i -= len(x.SubjectId)
			copy(dAtA[i:], x.SubjectId)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.SubjectId)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*RecommendedSubject)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				wire |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			fieldNum := int32(wire >> 3)
			wireType := int(wire & 0x7)
			if wireType == 4 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: RecommendedSubject: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: RecommendedSubject: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field SubjectId", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.SubjectId = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field RecommendationRank", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.RecommendationRank = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Reason", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Reason = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field IsRequired", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.IsRequired = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 5:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field SemesterAlignment", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.SemesterAlignment = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 6:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field DifficultyLevel", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.DifficultyLevel = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: academictoken/schedule/recommended_subject.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RecommendedSubject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubjectId          string `protobuf:"bytes,1,opt,name=subjectId,proto3" json:"subjectId,omitempty"`
	RecommendationRank string `protobuf:"bytes,2,opt,name=recommendationRank,proto3" json:"recommendationRank,omitempty"`
	Reason             string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	IsRequired         string `protobuf:"bytes,4,opt,name=isRequired,proto3" json:"isRequired,omitempty"`
	SemesterAlignment  string `protobuf:"bytes,5,opt,name=semesterAlignment,proto3" json:"semesterAlignment,omitempty"`
	DifficultyLevel    string `protobuf:"bytes,6,opt,name=difficultyLevel,proto3" json:"difficultyLevel,omitempty"`
}

func (x *RecommendedSubject) Reset() {
	*x = RecommendedSubject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_academictoken_schedule_recommended_subject_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecommendedSubject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecommendedSubject) ProtoMessage() {}

// Deprecated: Use RecommendedSubject.ProtoReflect.Descriptor instead.
func (*RecommendedSubject) Descriptor() ([]byte, []int) {
	return file_academictoken_schedule_recommended_subject_proto_rawDescGZIP(), []int{0}
}

func (x *RecommendedSubject) GetSubjectId() string {
	if x != nil {
		return x.SubjectId
	}
	return ""
}

func (x *RecommendedSubject) GetRecommendationRank() string {
	if x != nil {
		return x.RecommendationRank
	}
	return ""
}

func (x *RecommendedSubject) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *RecommendedSubject) GetIsRequired() string {
	if x != nil {
		return x.IsRequired
	}
	return ""
}

func (x *RecommendedSubject) GetSemesterAlignment() string {
	if x != nil {
		return x.SemesterAlignment
	}
	return ""
}

func (x *RecommendedSubject) GetDifficultyLevel() string {
	if x != nil {
		return x.DifficultyLevel
	}
	return ""
}

var File_academictoken_schedule_recommended_subject_proto protoreflect.FileDescriptor

var file_academictoken_schedule_recommended_subject_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x63, 0x61, 0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x64, 0x65, 0x64, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x16, 0x61, 0x63, 0x61, 0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x22, 0xf2, 0x01, 0x0a, 0x12, 0x52,
	0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x2e, 0x0a, 0x12, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x61, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x6e, 0x6b, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x73, 0x65, 0x6d, 0x65, 0x73,
	0x74, 0x65, 0x72, 0x41, 0x6c, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x11, 0x73, 0x65, 0x6d, 0x65, 0x73, 0x74, 0x65, 0x72, 0x41, 0x6c, 0x69, 0x67,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75,
	0x6c, 0x74, 0x79, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x42,
	0xd8, 0x01, 0x0a, 0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x63, 0x61, 0x64, 0x65, 0x6d, 0x69, 0x63,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x17,
	0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x53, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x61, 0x63, 0x61, 0x64, 0x65,
	0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x63, 0x61,
	0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0xa2, 0x02, 0x03, 0x41, 0x53, 0x58, 0xaa, 0x02, 0x16, 0x41, 0x63, 0x61, 0x64,
	0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0xca, 0x02, 0x16, 0x41, 0x63, 0x61, 0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5c, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0xe2, 0x02, 0x22, 0x41, 0x63,
	0x61, 0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5c, 0x53, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x17, 0x41, 0x63, 0x61, 0x64, 0x65, 0x6d, 0x69, 0x63, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x3a, 0x3a, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_academictoken_schedule_recommended_subject_proto_rawDescOnce sync.Once
	file_academictoken_schedule_recommended_subject_proto_rawDescData = file_academictoken_schedule_recommended_subject_proto_rawDesc
)

func file_academictoken_schedule_recommended_subject_proto_rawDescGZIP() []byte {
	file_academictoken_schedule_recommended_subject_proto_rawDescOnce.Do(func() {
		file_academictoken_schedule_recommended_subject_proto_rawDescData = protoimpl.X.CompressGZIP(file_academictoken_schedule_recommended_subject_proto_rawDescData)
	})
	return file_academictoken_schedule_recommended_subject_proto_rawDescData
}

var file_academictoken_schedule_recommended_subject_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_academictoken_schedule_recommended_subject_proto_goTypes = []interface{}{
	(*RecommendedSubject)(nil), // 0: academictoken.schedule.RecommendedSubject
}
var file_academictoken_schedule_recommended_subject_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_academictoken_schedule_recommended_subject_proto_init() }
func file_academictoken_schedule_recommended_subject_proto_init() {
	if File_academictoken_schedule_recommended_subject_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_academictoken_schedule_recommended_subject_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecommendedSubject); i {
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
			RawDescriptor: file_academictoken_schedule_recommended_subject_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_academictoken_schedule_recommended_subject_proto_goTypes,
		DependencyIndexes: file_academictoken_schedule_recommended_subject_proto_depIdxs,
		MessageInfos:      file_academictoken_schedule_recommended_subject_proto_msgTypes,
	}.Build()
	File_academictoken_schedule_recommended_subject_proto = out.File
	file_academictoken_schedule_recommended_subject_proto_rawDesc = nil
	file_academictoken_schedule_recommended_subject_proto_goTypes = nil
	file_academictoken_schedule_recommended_subject_proto_depIdxs = nil
}
