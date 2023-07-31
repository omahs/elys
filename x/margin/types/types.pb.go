// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/margin/types.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

type Position int32

const (
	Position_UNSPECIFIED Position = 0
	Position_LONG        Position = 1
	Position_SHORT       Position = 2
)

var Position_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "LONG",
	2: "SHORT",
}

var Position_value = map[string]int32{
	"UNSPECIFIED": 0,
	"LONG":        1,
	"SHORT":       2,
}

func (x Position) String() string {
	return proto.EnumName(Position_name, int32(x))
}

func (Position) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cd1c09c977f732f9, []int{0}
}

type MTP struct {
	Address                  string                                  `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	CollateralAsset          string                                  `protobuf:"bytes,2,opt,name=collateral_asset,json=collateralAsset,proto3" json:"collateral_asset,omitempty"`
	CollateralAmount         github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,3,opt,name=collateral_amount,json=collateralAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"collateral_amount"`
	Liabilities              github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,4,opt,name=liabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"liabilities"`
	InterestPaidCollateral   github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,5,opt,name=interest_paid_collateral,json=interestPaidCollateral,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"interest_paid_collateral"`
	InterestPaidCustody      github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,6,opt,name=interest_paid_custody,json=interestPaidCustody,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"interest_paid_custody"`
	InterestUnpaidCollateral github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,7,opt,name=interest_unpaid_collateral,json=interestUnpaidCollateral,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"interest_unpaid_collateral"`
	CustodyAsset             string                                  `protobuf:"bytes,8,opt,name=custody_asset,json=custodyAsset,proto3" json:"custody_asset,omitempty"`
	CustodyAmount            github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,9,opt,name=custody_amount,json=custodyAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"custody_amount"`
	Leverage                 github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,10,opt,name=leverage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverage"`
	MtpHealth                github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,11,opt,name=mtp_health,json=mtpHealth,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"mtp_health"`
	Position                 Position                                `protobuf:"varint,12,opt,name=position,proto3,enum=elys.margin.Position" json:"position,omitempty"`
	Id                       uint64                                  `protobuf:"varint,13,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *MTP) Reset()         { *m = MTP{} }
func (m *MTP) String() string { return proto.CompactTextString(m) }
func (*MTP) ProtoMessage()    {}
func (*MTP) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd1c09c977f732f9, []int{0}
}
func (m *MTP) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MTP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MTP.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MTP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MTP.Merge(m, src)
}
func (m *MTP) XXX_Size() int {
	return m.Size()
}
func (m *MTP) XXX_DiscardUnknown() {
	xxx_messageInfo_MTP.DiscardUnknown(m)
}

var xxx_messageInfo_MTP proto.InternalMessageInfo

func (m *MTP) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MTP) GetCollateralAsset() string {
	if m != nil {
		return m.CollateralAsset
	}
	return ""
}

func (m *MTP) GetCustodyAsset() string {
	if m != nil {
		return m.CustodyAsset
	}
	return ""
}

func (m *MTP) GetPosition() Position {
	if m != nil {
		return m.Position
	}
	return Position_UNSPECIFIED
}

func (m *MTP) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterEnum("elys.margin.Position", Position_name, Position_value)
	proto.RegisterType((*MTP)(nil), "elys.margin.MTP")
}

func init() { proto.RegisterFile("elys/margin/types.proto", fileDescriptor_cd1c09c977f732f9) }

var fileDescriptor_cd1c09c977f732f9 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0x86, 0xed, 0xf4, 0x12, 0xe7, 0xa4, 0x4d, 0xc3, 0x40, 0x61, 0xd4, 0x85, 0x1b, 0x81, 0x04,
	0x01, 0x54, 0x9b, 0xcb, 0x13, 0xf4, 0x12, 0x68, 0x11, 0x6d, 0x83, 0xdb, 0xb0, 0x40, 0x48, 0xd1,
	0xc4, 0x1e, 0x39, 0xa3, 0xda, 0x1e, 0xcb, 0x33, 0x01, 0xf2, 0x16, 0x6c, 0x78, 0xa7, 0x2e, 0xbb,
	0x44, 0x2c, 0x2a, 0x94, 0xbc, 0x08, 0xf2, 0xd8, 0x71, 0xdd, 0xae, 0x90, 0x57, 0xc9, 0x9c, 0xf3,
	0xfb, 0xfb, 0xcf, 0xb1, 0x7e, 0x0f, 0x3c, 0xa2, 0xc1, 0x54, 0xd8, 0x21, 0x49, 0x7c, 0x16, 0xd9,
	0x72, 0x1a, 0x53, 0x61, 0xc5, 0x09, 0x97, 0x1c, 0x35, 0xd3, 0x86, 0x95, 0x35, 0xb6, 0x1e, 0xf8,
	0xdc, 0xe7, 0xaa, 0x6e, 0xa7, 0xff, 0x32, 0xc9, 0xe3, 0x5f, 0x75, 0x58, 0x3a, 0x3e, 0xef, 0x23,
	0x0c, 0x75, 0xe2, 0x79, 0x09, 0x15, 0x02, 0xeb, 0x1d, 0xbd, 0xdb, 0x70, 0x16, 0x47, 0xf4, 0x1c,
	0xda, 0x2e, 0x0f, 0x02, 0x22, 0x69, 0x42, 0x82, 0x21, 0x11, 0x82, 0x4a, 0x5c, 0x53, 0x92, 0x8d,
	0x9b, 0xfa, 0x6e, 0x5a, 0x46, 0x5f, 0xe1, 0x5e, 0x59, 0x1a, 0xf2, 0x49, 0x24, 0xf1, 0x52, 0xaa,
	0xdd, 0xb3, 0x2f, 0xaf, 0xb7, 0xb5, 0x3f, 0xd7, 0xdb, 0xcf, 0x7c, 0x26, 0xc7, 0x93, 0x91, 0xe5,
	0xf2, 0xd0, 0x76, 0xb9, 0x08, 0xb9, 0xc8, 0x7f, 0x76, 0x84, 0x77, 0x91, 0x0f, 0x3f, 0x60, 0x91,
	0x74, 0x4a, 0xa6, 0xbb, 0x0a, 0x84, 0x3e, 0x41, 0x33, 0x60, 0x64, 0xc4, 0x02, 0x26, 0x19, 0x15,
	0x78, 0xb9, 0x1a, 0xb7, 0xcc, 0x40, 0x0c, 0x30, 0x8b, 0x24, 0x4d, 0xa8, 0x90, 0xc3, 0x98, 0x30,
	0x6f, 0x78, 0x63, 0x8a, 0x57, 0xaa, 0xf1, 0x1f, 0x2e, 0x80, 0x7d, 0xc2, 0xbc, 0xfd, 0x02, 0x87,
	0x5c, 0xd8, 0xbc, 0x63, 0x35, 0x11, 0x92, 0x7b, 0x53, 0xbc, 0x5a, 0xcd, 0xe7, 0xfe, 0x2d, 0x9f,
	0x8c, 0x85, 0x42, 0xd8, 0x2a, 0x4c, 0x26, 0xd1, 0xdd, 0x8d, 0xea, 0xd5, 0x9c, 0x8a, 0x57, 0x34,
	0x50, 0xc4, 0xd2, 0x4e, 0x4f, 0x60, 0x3d, 0xdf, 0x22, 0xcf, 0x85, 0xa1, 0x72, 0xb1, 0x96, 0x17,
	0xb3, 0x50, 0x7c, 0x86, 0x56, 0x21, 0xca, 0x12, 0xd1, 0xa8, 0x36, 0xc7, 0xc2, 0x2b, 0x8f, 0xc3,
	0x07, 0x30, 0x02, 0xfa, 0x8d, 0x26, 0xc4, 0xa7, 0x18, 0x14, 0xd1, 0xca, 0x89, 0x4f, 0xff, 0x83,
	0x78, 0x40, 0x5d, 0xa7, 0x78, 0x1e, 0x1d, 0x03, 0x84, 0x32, 0x1e, 0x8e, 0x29, 0x09, 0xe4, 0x18,
	0x37, 0x2b, 0xd1, 0x1a, 0xa1, 0x8c, 0x0f, 0x15, 0x00, 0xbd, 0x06, 0x23, 0xe6, 0x82, 0x49, 0xc6,
	0x23, 0xbc, 0xd6, 0xd1, 0xbb, 0xad, 0x37, 0x9b, 0x56, 0xe9, 0x53, 0xb4, 0xfa, 0x79, 0xd3, 0x29,
	0x64, 0xa8, 0x05, 0x35, 0xe6, 0xe1, 0xf5, 0x8e, 0xde, 0x5d, 0x76, 0x6a, 0xcc, 0x7b, 0xf1, 0x0a,
	0x8c, 0x85, 0x0a, 0x6d, 0x40, 0x73, 0x70, 0x72, 0xd6, 0xef, 0xed, 0x1f, 0xbd, 0x3b, 0xea, 0x1d,
	0xb4, 0x35, 0x64, 0xc0, 0xf2, 0xc7, 0xd3, 0x93, 0xf7, 0x6d, 0x1d, 0x35, 0x60, 0xe5, 0xec, 0xf0,
	0xd4, 0x39, 0x6f, 0xd7, 0xf6, 0x7a, 0x97, 0x33, 0x53, 0xbf, 0x9a, 0x99, 0xfa, 0xdf, 0x99, 0xa9,
	0xff, 0x9c, 0x9b, 0xda, 0xd5, 0xdc, 0xd4, 0x7e, 0xcf, 0x4d, 0xed, 0xcb, 0xcb, 0xd2, 0x06, 0xe9,
	0x18, 0x3b, 0x11, 0x95, 0xdf, 0x79, 0x72, 0xa1, 0x0e, 0xf6, 0x8f, 0x5b, 0x37, 0xc7, 0x68, 0x55,
	0xdd, 0x0b, 0x6f, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x08, 0x17, 0x42, 0x3e, 0x55, 0x04, 0x00,
	0x00,
}

func (m *MTP) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MTP) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MTP) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x68
	}
	if m.Position != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Position))
		i--
		dAtA[i] = 0x60
	}
	{
		size := m.MtpHealth.Size()
		i -= size
		if _, err := m.MtpHealth.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	{
		size := m.Leverage.Size()
		i -= size
		if _, err := m.Leverage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.CustodyAmount.Size()
		i -= size
		if _, err := m.CustodyAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if len(m.CustodyAsset) > 0 {
		i -= len(m.CustodyAsset)
		copy(dAtA[i:], m.CustodyAsset)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.CustodyAsset)))
		i--
		dAtA[i] = 0x42
	}
	{
		size := m.InterestUnpaidCollateral.Size()
		i -= size
		if _, err := m.InterestUnpaidCollateral.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.InterestPaidCustody.Size()
		i -= size
		if _, err := m.InterestPaidCustody.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.InterestPaidCollateral.Size()
		i -= size
		if _, err := m.InterestPaidCollateral.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Liabilities.Size()
		i -= size
		if _, err := m.Liabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.CollateralAmount.Size()
		i -= size
		if _, err := m.CollateralAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.CollateralAsset) > 0 {
		i -= len(m.CollateralAsset)
		copy(dAtA[i:], m.CollateralAsset)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.CollateralAsset)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MTP) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.CollateralAsset)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.CollateralAmount.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Liabilities.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.InterestPaidCollateral.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.InterestPaidCustody.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.InterestUnpaidCollateral.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.CustodyAsset)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.CustodyAmount.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Leverage.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.MtpHealth.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.Position != 0 {
		n += 1 + sovTypes(uint64(m.Position))
	}
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MTP) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MTP: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MTP: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAsset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CollateralAsset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liabilities", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaidCollateral", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InterestPaidCollateral.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaidCustody", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InterestPaidCustody.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestUnpaidCollateral", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InterestUnpaidCollateral.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CustodyAsset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CustodyAsset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CustodyAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CustodyAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leverage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Leverage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MtpHealth", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MtpHealth.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Position", wireType)
			}
			m.Position = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Position |= Position(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)