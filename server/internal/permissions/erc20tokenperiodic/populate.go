package erc20tokenperiodic

import "time"

// PopulatePermission fills in optional fields with defaults.
// If StartTime is nil it is set to time.Now().Unix() (current unix seconds).
// Matches populatePermission() in context.ts.
func PopulatePermission(perm Permission) PopulatedPermission {
    startTime := time.Now().Unix()
    if perm.Data.StartTime != nil {
        startTime = *perm.Data.StartTime
    }
    return PopulatedPermission{
        Type:                perm.Type,
        IsAdjustmentAllowed: perm.IsAdjustmentAllowed,
        Data: PopulatedPermissionData{
            Justification:  perm.Data.Justification,
            PeriodAmount:   perm.Data.PeriodAmount,
            PeriodDuration: perm.Data.PeriodDuration,
            StartTime:      startTime,
            TokenAddress:   perm.Data.TokenAddress,
        },
    }
}