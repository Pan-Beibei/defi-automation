package permissions

// DefaultContracts 是所有支持链上的合约地址（同一套部署）。
// 对应 TS: chainMetadata.ts 中的 contracts 常量。
var DefaultContracts = DelegationContracts{
    DelegationManager:            "0xdb9B1e94B5b69Df7e401DDbedE43491141047dB3",
    NativeTokenStreamingEnforcer: "0xD10b97905a320b13a0608f7E9cC506b56747df19",
    ExactCalldataEnforcer:        "0x99F2e9bF15ce5eC84685604836F71aB835DBBdED",
    TimestampEnforcer:            "0x1046bb45C8d673d4ea75321280DB34899413c069",
    NonceEnforcer:                "0xDE4f2FAC4B3D87A1d9953Ca5FC09FCa7F366254f",
    Erc20PeriodTransferEnforcer:  "0x474e3Ae7E169e940607cC624Da8A15Eb120139aB",
    ValueLteEnforcer:             "0x92Bf12322527cAA612fd31a0e810472BBB106A8F",
}