package erc20tokenperiodic

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

var (
    hexStrRegex  = regexp.MustCompile(`^0x[a-fA-F0-9]*$`)
    addressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
)

const (
    zeroAddress     = "0x0000000000000000000000000000000000000000"
    tenYearsSeconds = int64(10 * 365 * 24 * 60 * 60) // 315_360_000
)

// timePeriods are the standard period durations in seconds.
// Matches TIME_PERIOD_TO_SECONDS in time.ts.
var timePeriods = []int64{
    3600,     // hourly
    86400,    // daily
    604800,   // weekly
    1209600,  // biweekly
    2592000,  // monthly
    31536000, // yearly
}

// getClosestTimePeriod returns the standard period (seconds) nearest to the given value.
// Matches getClosestTimePeriod() in time.ts.
func getClosestTimePeriod(seconds int64) int64 {
    closest := timePeriods[0]
    minDiff := abs64(seconds - closest)
    for _, p := range timePeriods[1:] {
        if d := abs64(seconds - p); d < minDiff {
            minDiff = d
            closest = p
        }
    }
    return closest
}

func abs64(n int64) int64 {
    if n < 0 {
        return -n
    }
    return n
}

// getExpiryTimestamp extracts the expiry unix timestamp from rules if present.
func getExpiryTimestamp(rules []Rule) (int64, bool) {
    for _, rule := range rules {
        if rule.Type.Name == "expiry" {
            if v, ok := rule.Data["timestamp"]; ok {
                switch ts := v.(type) {
                case float64:
                    return int64(ts), true
                case int64:
                    return ts, true
                case int:
                    return int64(ts), true
                }
            }
        }
    }
    return 0, false
}

// validateHexInteger validates that value is a non-zero 0x-prefixed hex integer.
// Matches validateHexInteger() in validation.ts with required=true, allowZero=false.
func validateHexInteger(name, value string) error {
    if !hexStrRegex.MatchString(value) {
        return fmt.Errorf("invalid %s: must be a valid hex integer", name)
    }
    n := new(big.Int)
    if _, ok := n.SetString(strings.TrimPrefix(value, "0x"), 16); !ok {
        return fmt.Errorf("invalid %s: must be a valid hex integer", name)
    }
    if n.Sign() == 0 {
        return fmt.Errorf("invalid %s: must be greater than 0", name)
    }
    return nil
}

// validatePermissionSchema checks structural correctness of the permission and
// normalises periodDuration to the closest standard period (like zPeriodDuration transform).
func validatePermissionSchema(perm *Permission) error {
    if perm.Type.Name != "erc20-token-periodic" {
        return fmt.Errorf("invalid permission type: expected erc20-token-periodic, got %q", perm.Type.Name)
    }
    if err := validateHexInteger("periodAmount", perm.Data.PeriodAmount); err != nil {
        return err
    }
    if perm.Data.PeriodDuration <= 0 {
        return fmt.Errorf("invalid periodDuration: must be a positive integer")
    }
    if perm.Data.PeriodDuration > tenYearsSeconds {
        return fmt.Errorf("invalid periodDuration: must be <= %d seconds (10 years)", tenYearsSeconds)
    }
    if perm.Data.StartTime != nil {
        if *perm.Data.StartTime <= 0 {
            return fmt.Errorf("invalid startTime: must be a positive integer")
        }
        // startTime must be >= start of today (local midnight), matching zStartTime.
        now := time.Now()
        startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
        if *perm.Data.StartTime < startOfToday {
            return fmt.Errorf("invalid startTime: must be today or later")
        }
    }
    if !addressRegex.MatchString(perm.Data.TokenAddress) {
        return fmt.Errorf("invalid tokenAddress: must be a valid Ethereum address")
    }
    if strings.EqualFold(perm.Data.TokenAddress, zeroAddress) {
        return fmt.Errorf("invalid tokenAddress: cannot be the zero address")
    }

    // Normalise periodDuration — mirrors zPeriodDuration's .transform().
    perm.Data.PeriodDuration = getClosestTimePeriod(perm.Data.PeriodDuration)
    return nil
}


func validatePermissionData(perm Permission, rules []Rule) error {
    expiry, hasExpiry := getExpiryTimestamp(rules)
    if hasExpiry && perm.Data.StartTime != nil && *perm.Data.StartTime >= expiry {
        return fmt.Errorf("invalid startTime: must be before expiry")
    }
    return nil
}


func ParseAndValidatePermission(req PermissionRequest) (PermissionRequest, error) {
    if err := validatePermissionSchema(&req.Permission); err != nil {
        return PermissionRequest{}, fmt.Errorf("invalid permission: %w", err)
    }
    if err := validatePermissionData(req.Permission, req.Rules); err != nil {
        return PermissionRequest{}, err
    }
    return req, nil
}