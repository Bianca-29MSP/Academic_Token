// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: academictoken/degree/degree_request.proto

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

type DegreeRequest struct {
	Id                     string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Creator                string   `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`
	StudentId              string   `protobuf:"bytes,3,opt,name=student_id,json=studentId,proto3" json:"student_id,omitempty"`
	InstitutionId          string   `protobuf:"bytes,4,opt,name=institution_id,json=institutionId,proto3" json:"institution_id,omitempty"`
	CurriculumId           string   `protobuf:"bytes,5,opt,name=curriculum_id,json=curriculumId,proto3" json:"curriculum_id,omitempty"`
	ExpectedGraduationDate string   `protobuf:"bytes,6,opt,name=expected_graduation_date,json=expectedGraduationDate,proto3" json:"expected_graduation_date,omitempty"`
	RequestDate            string   `protobuf:"bytes,7,opt,name=request_date,json=requestDate,proto3" json:"request_date,omitempty"`
	Status                 string   `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	ValidationScore        string   `protobuf:"bytes,9,opt,name=validation_score,json=validationScore,proto3" json:"validation_score,omitempty"`
	ValidationDetails      string   `protobuf:"bytes,10,opt,name=validation_details,json=validationDetails,proto3" json:"validation_details,omitempty"`
	MissingRequirements    []string `protobuf:"bytes,11,rep,name=missing_requirements,json=missingRequirements,proto3" json:"missing_requirements,omitempty"`
}

func (m *DegreeRequest) Reset()         { *m = DegreeRequest{} }
func (m *DegreeRequest) String() string { return proto.CompactTextString(m) }
func (*DegreeRequest) ProtoMessage()    {}
func (*DegreeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_07daacb6c13072d5, []int{0}
}
func (m *DegreeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DegreeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DegreeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DegreeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DegreeRequest.Merge(m, src)
}
func (m *DegreeRequest) XXX_Size() int {
	return m.Size()
}
func (m *DegreeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DegreeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DegreeRequest proto.InternalMessageInfo

func (m *DegreeRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DegreeRequest) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *DegreeRequest) GetStudentId() string {
	if m != nil {
		return m.StudentId
	}
	return ""
}

func (m *DegreeRequest) GetInstitutionId() string {
	if m != nil {
		return m.InstitutionId
	}
	return ""
}

func (m *DegreeRequest) GetCurriculumId() string {
	if m != nil {
		return m.CurriculumId
	}
	return ""
}

func (m *DegreeRequest) GetExpectedGraduationDate() string {
	if m != nil {
		return m.ExpectedGraduationDate
	}
	return ""
}

func (m *DegreeRequest) GetRequestDate() string {
	if m != nil {
		return m.RequestDate
	}
	return ""
}

func (m *DegreeRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *DegreeRequest) GetValidationScore() string {
	if m != nil {
		return m.ValidationScore
	}
	return ""
}

func (m *DegreeRequest) GetValidationDetails() string {
	if m != nil {
		return m.ValidationDetails
	}
	return ""
}

func (m *DegreeRequest) GetMissingRequirements() []string {
	if m != nil {
		return m.MissingRequirements
	}
	return nil
}

func init() {
	proto.RegisterType((*DegreeRequest)(nil), "academictoken.degree.DegreeRequest")
}

func init() {
	proto.RegisterFile("academictoken/degree/degree_request.proto", fileDescriptor_07daacb6c13072d5)
}

var fileDescriptor_07daacb6c13072d5 = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xcd, 0x4e, 0xea, 0x40,
	0x14, 0x80, 0x29, 0xdc, 0x0b, 0xb7, 0x87, 0x9f, 0x7b, 0xef, 0x48, 0xc8, 0x2c, 0xb4, 0x41, 0x8d,
	0x09, 0x2c, 0xc4, 0x18, 0x13, 0xe3, 0xda, 0x90, 0x18, 0xb6, 0x75, 0xe7, 0x86, 0x8c, 0x33, 0x27,
	0x64, 0x22, 0x6d, 0x71, 0xe6, 0x8c, 0xc1, 0xb7, 0xf0, 0x29, 0x7c, 0x16, 0x97, 0x2c, 0x5d, 0x1a,
	0x78, 0x11, 0xc3, 0xb4, 0x08, 0xac, 0xda, 0xf9, 0xbe, 0x2f, 0x6d, 0x32, 0xe7, 0x40, 0x5f, 0x48,
	0xa1, 0x30, 0xd1, 0x92, 0xb2, 0x27, 0x4c, 0x2f, 0x14, 0x4e, 0x0c, 0x62, 0xf1, 0x18, 0x1b, 0x7c,
	0x76, 0x68, 0x69, 0x30, 0x33, 0x19, 0x65, 0xac, 0xbd, 0x97, 0x0e, 0xf2, 0xe6, 0xe4, 0xbd, 0x02,
	0xcd, 0xa1, 0x7f, 0x8d, 0xf3, 0x9a, 0xb5, 0xa0, 0xac, 0x15, 0x0f, 0xba, 0x41, 0x2f, 0x8c, 0xcb,
	0x5a, 0x31, 0x0e, 0x35, 0x69, 0x50, 0x50, 0x66, 0x78, 0xd9, 0xc3, 0xcd, 0x91, 0x1d, 0x01, 0x58,
	0x72, 0x0a, 0x53, 0x1a, 0x6b, 0xc5, 0x2b, 0x5e, 0x86, 0x05, 0x19, 0x29, 0x76, 0x06, 0x2d, 0x9d,
	0x5a, 0xd2, 0xe4, 0x48, 0x67, 0xe9, 0x3a, 0xf9, 0xe5, 0x93, 0xe6, 0x0e, 0x1d, 0x29, 0x76, 0x0a,
	0x4d, 0xe9, 0x8c, 0xd1, 0xd2, 0x4d, 0x5d, 0xb2, 0xae, 0x7e, 0xfb, 0xaa, 0xb1, 0x85, 0x23, 0xc5,
	0x6e, 0x80, 0xe3, 0x7c, 0x86, 0x92, 0x50, 0x8d, 0x27, 0x46, 0x28, 0x27, 0xfc, 0x37, 0x95, 0x20,
	0xe4, 0x55, 0xdf, 0x77, 0x36, 0xfe, 0xee, 0x47, 0x0f, 0x05, 0x21, 0x3b, 0x86, 0x46, 0x71, 0x0f,
	0x79, 0x5d, 0xf3, 0x75, 0xbd, 0x60, 0x3e, 0xe9, 0x40, 0xd5, 0x92, 0x20, 0x67, 0xf9, 0x1f, 0x2f,
	0x8b, 0x13, 0xeb, 0xc3, 0xbf, 0x17, 0x31, 0xd5, 0x2a, 0xff, 0x97, 0x95, 0x99, 0x41, 0x1e, 0xfa,
	0xe2, 0xef, 0x96, 0xdf, 0xaf, 0x31, 0x3b, 0x07, 0xb6, 0x93, 0x2a, 0x24, 0xa1, 0xa7, 0x96, 0x83,
	0x8f, 0xff, 0x6f, 0xcd, 0x30, 0x17, 0xec, 0x12, 0xda, 0x89, 0xb6, 0x56, 0xa7, 0x13, 0x3f, 0x24,
	0x6d, 0x30, 0xc1, 0x94, 0x2c, 0xaf, 0x77, 0x2b, 0xbd, 0x30, 0x3e, 0x28, 0x5c, 0xbc, 0xa3, 0x6e,
	0xaf, 0x3f, 0x96, 0x51, 0xb0, 0x58, 0x46, 0xc1, 0xd7, 0x32, 0x0a, 0xde, 0x56, 0x51, 0x69, 0xb1,
	0x8a, 0x4a, 0x9f, 0xab, 0xa8, 0xf4, 0x70, 0xb8, 0xbf, 0x03, 0xf3, 0xcd, 0x16, 0xd0, 0xeb, 0x0c,
	0xed, 0x63, 0xd5, 0x4f, 0xff, 0xea, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xa1, 0x8a, 0x34, 0x8f, 0x2a,
	0x02, 0x00, 0x00,
}

func (m *DegreeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DegreeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DegreeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.MissingRequirements) > 0 {
		for iNdEx := len(m.MissingRequirements) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.MissingRequirements[iNdEx])
			copy(dAtA[i:], m.MissingRequirements[iNdEx])
			i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.MissingRequirements[iNdEx])))
			i--
			dAtA[i] = 0x5a
		}
	}
	if len(m.ValidationDetails) > 0 {
		i -= len(m.ValidationDetails)
		copy(dAtA[i:], m.ValidationDetails)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.ValidationDetails)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.ValidationScore) > 0 {
		i -= len(m.ValidationScore)
		copy(dAtA[i:], m.ValidationScore)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.ValidationScore)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.RequestDate) > 0 {
		i -= len(m.RequestDate)
		copy(dAtA[i:], m.RequestDate)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.RequestDate)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.ExpectedGraduationDate) > 0 {
		i -= len(m.ExpectedGraduationDate)
		copy(dAtA[i:], m.ExpectedGraduationDate)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.ExpectedGraduationDate)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.CurriculumId) > 0 {
		i -= len(m.CurriculumId)
		copy(dAtA[i:], m.CurriculumId)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.CurriculumId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.InstitutionId) > 0 {
		i -= len(m.InstitutionId)
		copy(dAtA[i:], m.InstitutionId)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.InstitutionId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.StudentId) > 0 {
		i -= len(m.StudentId)
		copy(dAtA[i:], m.StudentId)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.StudentId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintDegreeRequest(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDegreeRequest(dAtA []byte, offset int, v uint64) int {
	offset -= sovDegreeRequest(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DegreeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.StudentId)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.InstitutionId)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.CurriculumId)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.ExpectedGraduationDate)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.RequestDate)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.ValidationScore)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	l = len(m.ValidationDetails)
	if l > 0 {
		n += 1 + l + sovDegreeRequest(uint64(l))
	}
	if len(m.MissingRequirements) > 0 {
		for _, s := range m.MissingRequirements {
			l = len(s)
			n += 1 + l + sovDegreeRequest(uint64(l))
		}
	}
	return n
}

func sovDegreeRequest(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDegreeRequest(x uint64) (n int) {
	return sovDegreeRequest(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DegreeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDegreeRequest
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
			return fmt.Errorf("proto: DegreeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DegreeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StudentId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StudentId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InstitutionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InstitutionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurriculumId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurriculumId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpectedGraduationDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExpectedGraduationDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestDate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequestDate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidationScore", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidationScore = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidationDetails", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidationDetails = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MissingRequirements", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDegreeRequest
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
				return ErrInvalidLengthDegreeRequest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDegreeRequest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MissingRequirements = append(m.MissingRequirements, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDegreeRequest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDegreeRequest
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
func skipDegreeRequest(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDegreeRequest
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
					return 0, ErrIntOverflowDegreeRequest
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
					return 0, ErrIntOverflowDegreeRequest
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
				return 0, ErrInvalidLengthDegreeRequest
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDegreeRequest
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDegreeRequest
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDegreeRequest        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDegreeRequest          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDegreeRequest = fmt.Errorf("proto: unexpected end of group")
)
