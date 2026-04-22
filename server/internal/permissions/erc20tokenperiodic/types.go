package erc20tokenperiodic

import "encoding/json"

// TypeDescriptor can be a plain string or {"name": "...", "description": "..."}.
// Matches TypeScript: zTypeDescriptor = z.union([z.string(), z.object({name, description})])
type TypeDescriptor struct {
    Name string
}

func (t *TypeDescriptor) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err == nil {
        t.Name = s
        return nil
    }
    var obj struct {
        Name string `json:"name"`
    }
    if err := json.Unmarshal(data, &obj); err != nil {
        return err
    }
    t.Name = obj.Name
    return nil
}

func (t TypeDescriptor) MarshalJSON() ([]byte, error) {
    return json.Marshal(t.Name)
}

// PermissionData holds the erc20-token-periodic permission fields.
type PermissionData struct {
    Justification  string `json:"justification"`
    PeriodAmount   string `json:"periodAmount"`        // 0x-prefixed hex, e.g. "0x1a2b"
    PeriodDuration int64  `json:"periodDuration"`      // seconds
    StartTime      *int64 `json:"startTime,omitempty"` // unix timestamp (seconds), optional
    TokenAddress   string `json:"tokenAddress"`        // 0x-prefixed EVM address
}

// Permission represents an erc20-token-periodic permission object.
type Permission struct {
    Type                TypeDescriptor `json:"type"`
    IsAdjustmentAllowed bool           `json:"isAdjustmentAllowed"`
    Data                PermissionData `json:"data"`
}

// Rule represents a constraint rule such as expiry.
type Rule struct {
    Type TypeDescriptor         `json:"type"`
    Data map[string]interface{} `json:"data"`
}

// PermissionRequest is the full ERC-7715 permission request.
type PermissionRequest struct {
    ChainID    string     `json:"chainId"`
    From       *string    `json:"from,omitempty"`
    To         string     `json:"to"`
    Permission Permission `json:"permission"`
    Rules      []Rule     `json:"rules,omitempty"`
}

// PopulatedPermissionData is PermissionData with StartTime guaranteed non-nil.
type PopulatedPermissionData struct {
    Justification  string `json:"justification"`
    PeriodAmount   string `json:"periodAmount"`
    PeriodDuration int64  `json:"periodDuration"`
    StartTime      int64  `json:"startTime"` // always set after PopulatePermission
    TokenAddress   string `json:"tokenAddress"`
}

// PopulatedPermission is a Permission with all optional fields filled in.
type PopulatedPermission struct {
    Type                TypeDescriptor          `json:"type"`
    IsAdjustmentAllowed bool                    `json:"isAdjustmentAllowed"`
    Data                PopulatedPermissionData `json:"data"`
}

// Caveat represents a single delegation enforcement condition.
type Caveat struct {
    Enforcer string `json:"enforcer"` // enforcer contract address
    Terms    string `json:"terms"`    // tightly-packed ABI terms, 0x-prefixed hex
    Args     string `json:"args"`     // "0x" (unused by these enforcers)
}

// DelegationContracts holds the enforcer contract addresses for the target chain.
type DelegationContracts struct {
    Erc20PeriodTransferEnforcer string
    ValueLteEnforcer            string
}