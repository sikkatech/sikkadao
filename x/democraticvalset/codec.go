package democraticvalset

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUpdateDescription{}, "cosmos-sdk/MsgUpdateDescription", nil)
	cdc.RegisterConcrete(MsgUpdateConsPubKey{}, "cosmos-sdk/MsgUpdateConsPubKey", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
