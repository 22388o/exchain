package keeper

import (
	"fmt"

	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

// Keeper of the mint store
type Keeper struct {
	cdc              *codec.Codec
	storeKey         sdk.StoreKey
	paramSpace       params.Subspace
	sk               types.StakingKeeper
	supplyKeeper     types.SupplyKeeper
	feeCollectorName string

	farmModuleName         string
	originalMintedPerBlock sdk.Dec
	paramCaches            *paramCache
}

type paramCache struct {
	param map[string]types.Params

	mintCustom map[string]types.MinterCustom
}

func newConfig() *paramCache {
	return &paramCache{
		param:      map[string]types.Params{},
		mintCustom: map[string]types.MinterCustom{},
	}
}

var (
	paramKey        = "Params"
	minterCustomKey = "MinterCustom"
)

func (p *paramCache) setParam(data types.Params) {
	p.param[paramKey] = data
}

func (p *paramCache) setMinter(data types.MinterCustom) {
	p.mintCustom[minterCustomKey] = data
}

func (p *paramCache) getParams() (types.Params, bool) {
	data, ok := p.param[paramKey]
	return data, ok
}

func (p *paramCache) getMinter() (types.MinterCustom, bool) {
	data, ok := p.mintCustom[minterCustomKey]
	return data, ok
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	sk types.StakingKeeper, supplyKeeper types.SupplyKeeper, feeCollectorName, farmModule string,
) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	return Keeper{
		cdc:                    cdc,
		storeKey:               key,
		paramSpace:             paramSpace.WithKeyTable(types.ParamKeyTable()),
		sk:                     sk,
		supplyKeeper:           supplyKeeper,
		feeCollectorName:       feeCollectorName,
		farmModuleName:         farmModule,
		originalMintedPerBlock: types.DefaultOriginalMintedPerBlock(),
		paramCaches:            newConfig(),
	}
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.MinterCustom) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(types.MinterKey, b)
}

//______________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if data, ok := k.paramCaches.getParams(); ok {
		return data
	}
	k.paramSpace.GetParamSet(ctx, &params)

	k.paramCaches.setParam(params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	k.paramCaches.setParam(params)
}

//______________________________________________________________________

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) sdk.Dec {
	return k.sk.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.sk.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
