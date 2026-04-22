package permissions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"server/internal/types"
)

var (
    hexStrRegex  = regexp.MustCompile(`^0x[a-fA-F0-9]*$`)
    addressRegex = regexp.MustCompile(`(?i)^0x[a-fA-F0-9]{40}$`)
)


// 每种权限类型允许使用的 rule 类型
var supportedRuleTypes = map[string][]string{
    "native-token-stream":   {"expiry"},
    "native-token-periodic":  {"expiry"},
    "erc20-token-stream":     {"expiry"},
    "erc20-token-periodic":   {"expiry"},
    "erc20-token-revocation": {"expiry"},
}

// InvalidInputError signals that the caller provided invalid parameters.
type InvalidInputError struct {
    Message string
}

func (e *InvalidInputError) Error() string { return e.Message }

func invalidInput(format string, args ...any) *InvalidInputError {
    return &InvalidInputError{Message: fmt.Sprintf(format, args...)}
}

// ValidateSubmitPermissionReq 校验前端发来的 ERC7715PermissionPayload。
func ValidateSubmitPermissionReq(req *types.SubmitPermissionReq) error {

    // validate chainId
    if !hexStrRegex.MatchString(req.ChainId) {
        return invalidInput("chainId: must be a hex string (0x...)")
    }

    // if from is provided, validate it
    if req.From != "" && !addressRegex.MatchString(req.From) {
        return invalidInput("from: invalid Ethereum address")
    }

    // validate to address
    if !addressRegex.MatchString(req.To) {
        return invalidInput("to: invalid Ethereum address")
    }

    // required permission type
    if req.Permission.Type == "" {
        return invalidInput("permission.type: must not be empty")
    }

    // permission.data 是前端 JSON.stringify 后的字符串，必须是合法 JSON
    if req.Permission.Data == "" {
        return invalidInput("permission.data: must not be empty")
    }
    if !json.Valid([]byte(req.Permission.Data)) {
        return invalidInput("permission.data: must be valid JSON")
    }


    // rules support multiple selection.
    if err := validateRules(req.Rules, req.Permission.Type); err != nil {
        return err
    }

    return nil
}


func validateRules(rules []types.Rule, permissionTypeName string) error {
    if err := validateRuleDuplicates(rules); err != nil {
        return err
    }
    return validateRulesSupportedForPermission(rules, permissionTypeName)
}

func validateRuleDuplicates(rules []types.Rule) error {
    seen := make(map[string]struct{}, len(rules))
    for _, rule := range rules {
        if _, dup := seen[rule.Type]; dup {
            return &InvalidInputError{"rules: duplicate rule types are not allowed"}
        }
        seen[rule.Type] = struct{}{}
    }
    return nil
}

func validateRulesSupportedForPermission(rules []types.Rule, permissionTypeName string) error {
    supported, ok := supportedRuleTypes[permissionTypeName]
    if !ok {
        return nil
    }
    allowed := make(map[string]struct{}, len(supported))
    for _, s := range supported {
        allowed[s] = struct{}{}
    }
    for _, rule := range rules {
        if _, ok := allowed[rule.Type]; !ok {
            return &InvalidInputError{fmt.Sprintf(
                "rules: rule type %q is not supported for permission type %q. Supported: [%s]",
                rule.Type, permissionTypeName, strings.Join(supported, ", "),
            )}
        }
    }
    return nil
}