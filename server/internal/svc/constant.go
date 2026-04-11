package svc

type PermissionOffer struct {
	Type         string
	ProposedName string
}

type ChainMetadata struct {
	Name        string
	ExplorerUrl string
}

var DefaultGatorPermissionToOffer = []PermissionOffer{
	{
		Type:         "native-token-stream",
		ProposedName: "Native Token Stream",
	},
	{
		Type:         "native-token-periodic",
		ProposedName: "Native Token Periodic Transfer",
	},
	{
		Type:         "erc20-token-stream",
		ProposedName: "ERC20 Token Stream",
	},
	{
		Type:         "erc20-token-periodic",
		ProposedName: "ERC20 Token Periodic Transfer",
	},
	{
		Type:         "erc20-token-revocation",
		ProposedName: "ERC20 Token Revocation",
	},
}

var SupportedRuleTypes = map[string][]string{
	"native-token-stream":    {"expiry"},
	"native-token-periodic":  {"expiry"},
	"erc20-token-stream":     {"expiry"},
	"erc20-token-periodic":   {"expiry"},
	"erc20-token-revocation": {"expiry"},
}

var NameAndExplorerUrlByChainId = map[uint64]ChainMetadata{
	// mainnets
	0x1:     {Name: "Ethereum Mainnet", ExplorerUrl: "https://etherscan.io"},
	0xa:     {Name: "OP Mainnet", ExplorerUrl: "https://optimistic.etherscan.io"},
	0x38:    {Name: "BNB Smart Chain Mainnet", ExplorerUrl: "https://bscscan.com"},
	0x64:    {Name: "Gnosis", ExplorerUrl: "https://gnosisscan.io"},
	0x82:    {Name: "Unichain", ExplorerUrl: "https://uniscan.xyz"},
	0x89:    {Name: "Polygon Mainnet", ExplorerUrl: "https://polygonscan.com"},
	0x2105:  {Name: "Base", ExplorerUrl: "https://basescan.org"},
	0xa4b1:  {Name: "Arbitrum One", ExplorerUrl: "https://arbiscan.io"},
	0xa4ba:  {Name: "Arbitrum Nova", ExplorerUrl: "https://nova.arbiscan.io"},
	0xe708:  {Name: "Linea Mainnet", ExplorerUrl: "https://lineascan.build"},
	0x13882: {Name: "Amoy", ExplorerUrl: "https://amoy.polygonscan.com"},
	0x138de: {Name: "Berachain", ExplorerUrl: "https://berascan.com"},

	// testnets
	0x61:     {Name: "BNB Smart Chain Testnet", ExplorerUrl: "https://testnet.bscscan.com"},
	0x515:    {Name: "Unichain Sepolia Testnet", ExplorerUrl: "https://unichain-sepolia.blockscout.com"},
	0x18c6:   {Name: "MegaETH Testnet", ExplorerUrl: ""},
	0x27d8:   {Name: "Gnosis Chiado Testnet", ExplorerUrl: "https://gnosis-chiado.blockscout.com"},
	0xe705:   {Name: "Linea Sepolia Testnet", ExplorerUrl: "https://sepolia.lineascan.build"},
	0x138c5:  {Name: "Berachain Bepolia Testnet", ExplorerUrl: "https://bepolia.beratrail.io"},
	0x14a34:  {Name: "Base Sepolia Testnet", ExplorerUrl: "https://sepolia.basescan.org"},
	0x66eee:  {Name: "Arbitrum Sepolia Testnet", ExplorerUrl: "https://sepolia.arbiscan.io"},
	0xaa36a7: {Name: "Ethereum Sepolia Testnet", ExplorerUrl: "https://sepolia.etherscan.io"},
	0xaa37dc: {Name: "OP Sepolia Testnet", ExplorerUrl: "https://sepolia-optimism.etherscan.io"},
}