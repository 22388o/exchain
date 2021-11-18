package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if gasConsumed := k.configCache.paramsGas; gasConsumed != 0 {
		ctx.GasMeter().ConsumeGas(gasConsumed, "evm.Keeper.GetParams Error")
		return k.configCache.params
	}
	startGas := ctx.GasMeter().GasConsumed()
	k.paramSpace.GetParamSet(ctx, &params)
	k.configCache.setParams(params, ctx.GasMeter().GasConsumed()-startGas)
	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
