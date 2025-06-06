// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: academictoken/subject/subject_content.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SubjectContent struct {
	Index         string `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	SubjectId     string `protobuf:"bytes,2,opt,name=subjectId,proto3" json:"subjectId,omitempty"`
	Institution   string `protobuf:"bytes,3,opt,name=institution,proto3" json:"institution,omitempty"`
	CourseId      string `protobuf:"bytes,4,opt,name=course_id,json=courseId,proto3" json:"course_id,omitempty"`
	Title         string `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Code          string `protobuf:"bytes,6,opt,name=code,proto3" json:"code,omitempty"`
	WorkloadHours uint64 `protobuf:"varint,7,opt,name=workloadHours,proto3" json:"workloadHours,omitempty"`
	Credits       uint64 `protobuf:"varint,8,opt,name=credits,proto3" json:"credits,omitempty"`
	Description   string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	ContentHash   string `protobuf:"bytes,10,opt,name=contentHash,proto3" json:"contentHash,omitempty"`
	SubjectType   string `protobuf:"bytes,11,opt,name=subjectType,proto3" json:"subjectType,omitempty"`
	KnowledgeArea string `protobuf:"bytes,12,opt,name=knowledgeArea,proto3" json:"knowledgeArea,omitempty"`
	IpfsLink      string `protobuf:"bytes,13,opt,name=ipfsLink,proto3" json:"ipfsLink,omitempty"`
	Creator       string `protobuf:"bytes,14,opt,name=creator,proto3" json:"creator,omitempty"`
}

func (m *SubjectContent) Reset()         { *m = SubjectContent{} }
func (m *SubjectContent) String() string { return proto.CompactTextString(m) }
func (*SubjectContent) ProtoMessage()    {}
func (*SubjectContent) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0bf936fb104efe4, []int{0}
}
func (m *SubjectContent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubjectContent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubjectContent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubjectContent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubjectContent.Merge(m, src)
}
func (m *SubjectContent) XXX_Size() int {
	return m.Size()
}
func (m *SubjectContent) XXX_DiscardUnknown() {
	xxx_messageInfo_SubjectContent.DiscardUnknown(m)
}

var xxx_messageInfo_SubjectContent proto.InternalMessageInfo

func (m *SubjectContent) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *SubjectContent) GetSubjectId() string {
	if m != nil {
		return m.SubjectId
	}
	return ""
}

func (m *SubjectContent) GetInstitution() string {
	if m != nil {
		return m.Institution
	}
	return ""
}

func (m *SubjectContent) GetCourseId() string {
	if m != nil {
		return m.CourseId
	}
	return ""
}

func (m *SubjectContent) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SubjectContent) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *SubjectContent) GetWorkloadHours() uint64 {
	if m != nil {
		return m.WorkloadHours
	}
	return 0
}

func (m *SubjectContent) GetCredits() uint64 {
	if m != nil {
		return m.Credits
	}
	return 0
}

func (m *SubjectContent) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SubjectContent) GetContentHash() string {
	if m != nil {
		return m.ContentHash
	}
	return ""
}

func (m *SubjectContent) GetSubjectType() string {
	if m != nil {
		return m.SubjectType
	}
	return ""
}

func (m *SubjectContent) GetKnowledgeArea() string {
	if m != nil {
		return m.KnowledgeArea
	}
	return ""
}

func (m *SubjectContent) GetIpfsLink() string {
	if m != nil {
		return m.IpfsLink
	}
	return ""
}

func (m *SubjectContent) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func init() {
	proto.RegisterType((*SubjectContent)(nil), "academictoken.subject.SubjectContent")
}

func init() {
	proto.RegisterFile("academictoken/subject/subject_content.proto", fileDescriptor_a0bf936fb104efe4)
}

var fileDescriptor_a0bf936fb104efe4 = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x1b, 0xfa, 0x37, 0x2e, 0xed, 0x60, 0x81, 0x64, 0x15, 0x88, 0x2a, 0xc4, 0x50, 0x09,
	0xa9, 0x0c, 0x0c, 0xcc, 0xc0, 0xd2, 0x4a, 0x4c, 0x85, 0x89, 0xa5, 0x4a, 0xed, 0xa3, 0x98, 0x14,
	0x3b, 0xd8, 0x8e, 0xda, 0xbe, 0x05, 0x6f, 0xc2, 0x6b, 0x30, 0x76, 0x64, 0x44, 0xed, 0x8b, 0xa0,
	0xd8, 0x09, 0x24, 0x12, 0x53, 0xf2, 0xfd, 0xbe, 0xcf, 0xe7, 0xf3, 0x1d, 0x3a, 0x0f, 0x69, 0xc8,
	0xe0, 0x95, 0x53, 0x23, 0x23, 0x10, 0x17, 0x3a, 0x99, 0xbd, 0x00, 0x35, 0xf9, 0x77, 0x4a, 0xa5,
	0x30, 0x20, 0xcc, 0x30, 0x56, 0xd2, 0x48, 0x7c, 0x58, 0x0a, 0x0f, 0xb3, 0x50, 0x6f, 0xf8, 0x7f,
	0x8d, 0x58, 0x81, 0x82, 0xb7, 0x84, 0x6b, 0x6e, 0x60, 0x3a, 0x57, 0x32, 0x89, 0x5d, 0x99, 0xd3,
	0x8f, 0x2a, 0xea, 0xde, 0xbb, 0xd0, 0xad, 0xab, 0x8f, 0x0f, 0x50, 0x9d, 0x0b, 0x06, 0x2b, 0xe2,
	0xf5, 0xbd, 0x81, 0x3f, 0x71, 0x02, 0x1f, 0x23, 0x3f, 0x2b, 0x36, 0x66, 0x64, 0xcf, 0x3a, 0x7f,
	0x00, 0xf7, 0x51, 0x9b, 0x0b, 0x6d, 0xb8, 0x49, 0x0c, 0x97, 0x82, 0x54, 0xad, 0x5f, 0x44, 0xf8,
	0x08, 0xf9, 0x54, 0x26, 0x4a, 0xc3, 0x94, 0x33, 0x52, 0xb3, 0x7e, 0xcb, 0x81, 0x31, 0x4b, 0xaf,
	0x34, 0xdc, 0x2c, 0x80, 0xd4, 0xdd, 0x95, 0x56, 0x60, 0x8c, 0x6a, 0x54, 0x32, 0x20, 0x0d, 0x0b,
	0xed, 0x3f, 0x3e, 0x43, 0x9d, 0xa5, 0x54, 0xd1, 0x42, 0x86, 0x6c, 0x94, 0x9e, 0x26, 0xcd, 0xbe,
	0x37, 0xa8, 0x4d, 0xca, 0x10, 0x13, 0xd4, 0xa4, 0x0a, 0x18, 0x37, 0x9a, 0xb4, 0xac, 0x9f, 0xcb,
	0xb4, 0x51, 0x06, 0x9a, 0x2a, 0x1e, 0xdb, 0x46, 0x7d, 0xd7, 0x68, 0x01, 0xa5, 0x89, 0x6c, 0xd2,
	0xa3, 0x50, 0x3f, 0x13, 0xe4, 0x12, 0x05, 0x94, 0x26, 0xb2, 0x97, 0x3f, 0xac, 0x63, 0x20, 0x6d,
	0x97, 0x28, 0xa0, 0xb4, 0xcb, 0x48, 0xc8, 0xe5, 0x02, 0xd8, 0x1c, 0xae, 0x15, 0x84, 0x64, 0xdf,
	0x66, 0xca, 0x10, 0xf7, 0x50, 0x8b, 0xc7, 0x4f, 0xfa, 0x8e, 0x8b, 0x88, 0x74, 0xdc, 0x44, 0x72,
	0x9d, 0xbd, 0x20, 0x34, 0x52, 0x91, 0xae, 0xb5, 0x72, 0x79, 0x73, 0xf5, 0xb9, 0x0d, 0xbc, 0xcd,
	0x36, 0xf0, 0xbe, 0xb7, 0x81, 0xf7, 0xbe, 0x0b, 0x2a, 0x9b, 0x5d, 0x50, 0xf9, 0xda, 0x05, 0x95,
	0xc7, 0x93, 0xf2, 0xee, 0x57, 0xbf, 0xdb, 0x37, 0xeb, 0x18, 0xf4, 0xac, 0x61, 0x37, 0x7e, 0xf9,
	0x13, 0x00, 0x00, 0xff, 0xff, 0x87, 0xba, 0xd2, 0x12, 0x67, 0x02, 0x00, 0x00,
}

func (m *SubjectContent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubjectContent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubjectContent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.IpfsLink) > 0 {
		i -= len(m.IpfsLink)
		copy(dAtA[i:], m.IpfsLink)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.IpfsLink)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.KnowledgeArea) > 0 {
		i -= len(m.KnowledgeArea)
		copy(dAtA[i:], m.KnowledgeArea)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.KnowledgeArea)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.SubjectType) > 0 {
		i -= len(m.SubjectType)
		copy(dAtA[i:], m.SubjectType)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.SubjectType)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.ContentHash) > 0 {
		i -= len(m.ContentHash)
		copy(dAtA[i:], m.ContentHash)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.ContentHash)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x4a
	}
	if m.Credits != 0 {
		i = encodeVarintSubjectContent(dAtA, i, uint64(m.Credits))
		i--
		dAtA[i] = 0x40
	}
	if m.WorkloadHours != 0 {
		i = encodeVarintSubjectContent(dAtA, i, uint64(m.WorkloadHours))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Code) > 0 {
		i -= len(m.Code)
		copy(dAtA[i:], m.Code)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Code)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.CourseId) > 0 {
		i -= len(m.CourseId)
		copy(dAtA[i:], m.CourseId)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.CourseId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Institution) > 0 {
		i -= len(m.Institution)
		copy(dAtA[i:], m.Institution)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Institution)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SubjectId) > 0 {
		i -= len(m.SubjectId)
		copy(dAtA[i:], m.SubjectId)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.SubjectId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintSubjectContent(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSubjectContent(dAtA []byte, offset int, v uint64) int {
	offset -= sovSubjectContent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SubjectContent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.SubjectId)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.Institution)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.CourseId)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	if m.WorkloadHours != 0 {
		n += 1 + sovSubjectContent(uint64(m.WorkloadHours))
	}
	if m.Credits != 0 {
		n += 1 + sovSubjectContent(uint64(m.Credits))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.ContentHash)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.SubjectType)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.KnowledgeArea)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.IpfsLink)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovSubjectContent(uint64(l))
	}
	return n
}

func sovSubjectContent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSubjectContent(x uint64) (n int) {
	return sovSubjectContent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SubjectContent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubjectContent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
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
			return fmt.Errorf("proto: SubjectContent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubjectContent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubjectId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubjectId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Institution", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Institution = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CourseId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CourseId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WorkloadHours", wireType)
			}
			m.WorkloadHours = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WorkloadHours |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Credits", wireType)
			}
			m.Credits = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Credits |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContentHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContentHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubjectType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubjectType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KnowledgeArea", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KnowledgeArea = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IpfsLink", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IpfsLink = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
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
				return ErrInvalidLengthSubjectContent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSubjectContent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubjectContent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSubjectContent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSubjectContent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSubjectContent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSubjectContent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSubjectContent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSubjectContent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSubjectContent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSubjectContent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSubjectContent = fmt.Errorf("proto: unexpected end of group")
)
