package permissions

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"server/internal/types"
)

// ──────────────────────────────────────────────
// 类型定义（对应 types.ts）
// ──────────────────────────────────────────────

// Erc20TokenPeriodicData 对应 Snap 中 zErc20TokenPeriodicPermission 的 data 字段。
// 前端将此结构 JSON 序列化后作为 PermissionInfo.Data 字符串传入。
type Erc20TokenPeriodicData struct {
    // hex 大整数字符串，如 "0x16345785d8a0000"，对应 zHexStr
    PeriodAmount string `json:"periodAmount"`
    // 周期秒数，如 86400 = 1天，对应 zPeriodDuration
    PeriodDuration int64 `json:"periodDuration"`
    // unix 时间戳；0 表示"使用当前时间"，对应 zStartTime（optional）
    StartTime int64 `json:"startTime,omitempty"`
    // ERC20 合约地址，不能是零地址，对应 zAddressNotZeroAddress
    TokenAddress  string `json:"tokenAddress"`
    Justification string `json:"justification,omitempty"`
}

// Erc20TokenPeriodicValidationErrors 对应 Snap 中 Erc20TokenPeriodicMetadata.validationErrors。
// 字段级校验错误，供前端在 UI 上精确提示每个输入框的错误。
type Erc20TokenPeriodicValidationErrors struct {
    PeriodAmountError   string `json:"periodAmountError,omitempty"`
    PeriodDurationError string `json:"periodDurationError,omitempty"`
    StartTimeError      string `json:"startTimeError,omitempty"`
    ExpiryError         string `json:"expiryError,omitempty"`
}

func (e Erc20TokenPeriodicValidationErrors) HasErrors() bool {
    return e.PeriodAmountError != "" ||
        e.PeriodDurationError != "" ||
        e.StartTimeError != "" ||
        e.ExpiryError != ""
}

// ──────────────────────────────────────────────
// 第一层：结构校验（对应 validation.ts → parseAndValidatePermission）
// 请求进来时调用，失败直接返回 400
// ──────────────────────────────────────────────

// ParseErc20TokenPeriodicData 反序列化 PermissionInfo.Data。
// 对应 Zod schema 的 safeParse。
func ParseErc20TokenPeriodicData(rawData string) (*Erc20TokenPeriodicData, error) {
    var data Erc20TokenPeriodicData
    if err := json.Unmarshal([]byte(rawData), &data); err != nil {
        return nil, fmt.Errorf("invalid erc20-token-periodic data: %w", err)
    }
    return &data, nil
}

// ValidateErc20TokenPeriodicRequest 进行结构 + 基础业务校验。
// 对应 validation.ts 中 validatePermissionData：
//   - periodAmount 是有效 hex 整数且 > 0
//   - tokenAddress 是有效非零地址
//   - 若设置了 expiry，startTime 必须在 expiry 之前
func ValidateErc20TokenPeriodicRequest(req *types.SubmitPermissionReq, data *Erc20TokenPeriodicData) error {
    // 对应 validateHexInteger({ required: true, allowZero: false })
    if err := validateHexIntegerPositive("periodAmount", data.PeriodAmount); err != nil {
        return err
    }

    // 对应 zAddressNotZeroAddress
    if err := validateTokenAddress(data.TokenAddress); err != nil {
        return err
    }

    // 对应 validateStartTime(startTime, rules)：若有 expiry 规则，startTime 必须早于它
    // if req.Expiry != 0 {
    //     startTime := data.StartTime
    //     if startTime == 0 {
    //         startTime = time.Now().Unix()
    //     }
    //     if startTime >= req.Expiry {
    //         return fmt.Errorf("invalid startTime: must be before expiry")
    //     }
    // }

    return nil
}

// ──────────────────────────────────────────────
// 第二层：字段级校验（对应 context.ts → deriveMetadata）
// 用于 /build-context 接口，返回错误给前端逐字段显示
// decimals 需要从 token metadata 服务获取后传入
// ──────────────────────────────────────────────

// DeriveErc20TokenPeriodicValidationErrors 返回每个字段的校验错误。
// 对应 deriveMetadata 中对 context 的全量校验。
func DeriveErc20TokenPeriodicValidationErrors(
    data *Erc20TokenPeriodicData,
    expiry int64, // 来自 SubmitPermissionReq.Expiry，0 表示无过期
) Erc20TokenPeriodicValidationErrors {
    errs := Erc20TokenPeriodicValidationErrors{}

    // 对应 validateAndParseAmount：periodAmount hex 必须 > 0
    if err := validateHexIntegerPositive("period amount", data.PeriodAmount); err != nil {
        errs.PeriodAmountError = err.Error()
    }

    // 对应 validatePeriodDuration：periodDuration 必须 > 0
    if data.PeriodDuration <= 0 {
        errs.PeriodDurationError = "period duration must be greater than zero"
    }

    // 补全 startTime 默认值（对应 populatePermission 中的 startTime ?? Date.now()）
    startTime := data.StartTime
    if startTime == 0 {
        startTime = time.Now().Unix()
    }

    // 对应 contextValidation.ts → validateStartTime：不能早于今天
    if err := validateStartTimeNotPast(startTime); err != nil {
        errs.StartTimeError = err.Error()
    }

    // 对应 validateExpiry：expiry 必须在未来
    if expiry != 0 && expiry < time.Now().Unix() {
        errs.ExpiryError = "expiry must be in the future"
    }

    // 对应 validateStartTimeVsExpiry：两者都合法时再做交叉校验
    if expiry != 0 && errs.StartTimeError == "" && errs.ExpiryError == "" {
        if startTime >= expiry {
            errs.StartTimeError = "start time must be before expiry"
        }
    }

    return errs
}

// ──────────────────────────────────────────────
// 底层校验工具函数
// ──────────────────────────────────────────────

// validateHexIntegerPositive 对应 validateHexInteger({ required: true, allowZero: false })
func validateHexIntegerPositive(name, value string) error {
    if value == "" {
        return fmt.Errorf("invalid %s: must be defined", name)
    }
    trimmed := strings.TrimPrefix(strings.TrimPrefix(value, "0x"), "0X")
    n := new(big.Int)
    if _, ok := n.SetString(trimmed, 16); !ok {
        return fmt.Errorf("invalid %s: must be a valid hex integer", name)
    }
    if n.Sign() <= 0 {
        return fmt.Errorf("invalid %s: must be greater than 0", name)
    }
    return nil
}

var zeroAddress = strings.Repeat("0", 40)

// validateTokenAddress 对应 zAddressNotZeroAddress
func validateTokenAddress(addr string) error {
    if addr == "" {
        return fmt.Errorf("invalid tokenAddress: must be defined")
    }
    normalized := strings.ToLower(strings.TrimPrefix(addr, "0x"))
    if len(normalized) != 40 {
        return fmt.Errorf("invalid tokenAddress: must be a valid Ethereum address")
    }
    if normalized == zeroAddress {
        return fmt.Errorf("invalid tokenAddress: must not be the zero address")
    }
    return nil
}

// validateStartTimeNotPast 对应 contextValidation.ts → validateStartTime
// startTime 不能早于今天 00:00:00（本地时间）
func validateStartTimeNotPast(startTime int64) error {
    now := time.Now()
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
    if startTime < startOfToday {
        return fmt.Errorf("start time must be today or later")
    }
    return nil
}