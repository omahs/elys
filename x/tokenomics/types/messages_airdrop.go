package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAirdrop = "create_airdrop"
	TypeMsgUpdateAirdrop = "update_airdrop"
	TypeMsgDeleteAirdrop = "delete_airdrop"
)

var _ sdk.Msg = &MsgCreateAirdrop{}

func NewMsgCreateAirdrop(
	authority string,
	intent string,
	amount uint64,

) *MsgCreateAirdrop {
	return &MsgCreateAirdrop{
		Authority: authority,
		Intent:    intent,
		Amount:    amount,
	}
}

func (msg *MsgCreateAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgCreateAirdrop) Type() string {
	return TypeMsgCreateAirdrop
}

func (msg *MsgCreateAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAirdrop{}

func NewMsgUpdateAirdrop(
	authority string,
	intent string,
	amount uint64,

) *MsgUpdateAirdrop {
	return &MsgUpdateAirdrop{
		Authority: authority,
		Intent:    intent,
		Amount:    amount,
	}
}

func (msg *MsgUpdateAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAirdrop) Type() string {
	return TypeMsgUpdateAirdrop
}

func (msg *MsgUpdateAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteAirdrop{}

func NewMsgDeleteAirdrop(
	authority string,
	intent string,

) *MsgDeleteAirdrop {
	return &MsgDeleteAirdrop{
		Authority: authority,
		Intent:    intent,
	}
}
func (msg *MsgDeleteAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAirdrop) Type() string {
	return TypeMsgDeleteAirdrop
}

func (msg *MsgDeleteAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDeleteAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
