package permissions

import (
	"encoding/json"
	"fmt"
)

// TypeDescriptor represents an ERC-7715 type identifier.
// It can be a plain string ("native-token-stream") or
// an object with a name field ({"name": "native-token-stream", "description": "..."}).
type TypeDescriptor struct {
    name        string
    description string
}

// Name returns the type name regardless of original form.
func (t TypeDescriptor) Name() string {
    return t.name
}

// Description returns the optional description (empty if was a plain string).
func (t TypeDescriptor) Description() string {
    return t.description
}

// MarshalJSON encodes TypeDescriptor back to JSON.
func (t TypeDescriptor) MarshalJSON() ([]byte, error) {
    if t.description == "" {
        return json.Marshal(t.name)
    }
    return json.Marshal(struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }{t.name, t.description})
}

// UnmarshalJSON supports both string and object forms.
func (t *TypeDescriptor) UnmarshalJSON(data []byte) error {
    // Try plain string first
    var s string
    if err := json.Unmarshal(data, &s); err == nil {
        if s == "" {
            return fmt.Errorf("type descriptor string must not be empty")
        }
        t.name = s
        return nil
    }
    // Try object form
    var obj struct {
        Name        string `json:"name"`
        Description string `json:"description,omitempty"`
    }
    if err := json.Unmarshal(data, &obj); err != nil {
        return fmt.Errorf("type descriptor must be a string or object with 'name': %w", err)
    }
    if obj.Name == "" {
        return fmt.Errorf("type descriptor object must have a non-empty 'name' field")
    }
    t.name = obj.Name
    t.description = obj.Description
    return nil
}

// Permission represents an ERC-7715 permission.
type Permission struct {
    Type                TypeDescriptor `json:"type"`
    IsAdjustmentAllowed bool           `json:"isAdjustmentAllowed"`
    Data                map[string]any `json:"data"`
}

// Rule represents a constraint limiting how a permission can be used (e.g., expiry time).
type Rule struct {
    Type TypeDescriptor `json:"type"`
    Data map[string]any `json:"data"`
}

// PermissionRequest is a single ERC-7715 permission request.
type PermissionRequest struct {
    ChainID    string     `json:"chainId"`
    From       *string    `json:"from,omitempty"`
    To         string     `json:"to"`
    Permission Permission `json:"permission"`
    Rules      []Rule     `json:"rules"`
}

// RequestExecutionPermissionsParam is the validated parameter struct
// for the requestExecutionPermissions RPC method.
type RequestExecutionPermissionsParam struct {
    PermissionsRequest []PermissionRequest `json:"permissionsRequest"`
    SiteOrigin         string              `json:"siteOrigin"`
}