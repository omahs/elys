package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalCommitmentInfo
// Stores the
type TotalCommitmentInfo struct {
	// Total Elys staked
	TotalElysBonded sdk.Int
	// Total Eden + Eden boost committed
	TotalEdenEdenBoostCommitted sdk.Int
	// Gas fees collected and DEX revenus
	TotalFeesCollected sdk.Coins
	// Total Lp Token committed
	TotalLpTokensCommitted map[string]sdk.Int
	// Revenue tracking per pool, key => (poolId)
	PoolRevenueTrack map[string]sdk.Dec
}

// Returns the pool revenue tracking key.
// Unique per pool per epoch, clean once complete the calculation.
func GetPoolRevenueTrackKey(poolId uint64) string {
	return fmt.Sprintf("pool_revenue_%d", poolId)
}
