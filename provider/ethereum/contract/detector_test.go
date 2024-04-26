package contract_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/node/provider/ethereum"
	"github.com/rss3-network/node/provider/ethereum/contract"
	"github.com/rss3-network/node/provider/ethereum/endpoint"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/stretchr/testify/require"
)

func TestDetectTokenStandard(t *testing.T) {
	t.Parallel()

	type arguments struct {
		ctx     context.Context
		network network.Network
		address common.Address
	}

	testcases := []struct {
		name      string
		arguments arguments
		want      metadata.Standard
		wantError require.ErrorAssertionFunc
	}{
		{
			name: "RSS3",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0xc98d64da73a6616c42117b582e832812e7b8d57f
				address: common.HexToAddress("0xc98D64DA73a6616c42117b582e832812e7B8D57F"),
			},
			want:      metadata.StandardERC20,
			wantError: require.NoError,
		},
		{
			name: "USD Coin",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48
				address: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			},
			want:      metadata.StandardERC20,
			wantError: require.NoError,
		},
		{
			name: "Maker",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2
				address: common.HexToAddress("0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2"),
			},
			want:      metadata.StandardERC20,
			wantError: require.NoError,
		},
		{
			name: "The Genesis RSS3 Avatar NFT",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				address: common.HexToAddress("0x5452C7fB99D99fAb3Cc1875E9DA9829Cb50F7A13"),
			},
			want:      metadata.StandardERC721,
			wantError: require.NoError,
		},
		{
			name: "ENS",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				address: common.HexToAddress("0x57f1887a8BF19b14fC0dF6Fd9B2acc9Af147eA85"),
			},
			want:      metadata.StandardERC721,
			wantError: require.NoError,
		},
		{
			name: "Proof Of Stake Pages", // SBT
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0x5bf5bcc5362f88721167c1068b58c60cad075aac
				address: common.HexToAddress("0x5bF5BCc5362F88721167C1068b58C60caD075aAc"),
			},
			want:      metadata.StandardERC721,
			wantError: require.NoError,
		},
		{
			name: "Love, Death + Robots",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0xFD43D1dA000558473822302e1d44D81dA2e4cC0d
				address: common.HexToAddress("0xFD43D1dA000558473822302e1d44D81dA2e4cC0d"),
			},
			want:      metadata.StandardERC1155,
			wantError: require.NoError,
		},
		{
			name: "TIME NFT Special Issues",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0x8442864d6AB62a9193be2F16580c08E0D7BCda2f
				address: common.HexToAddress("0x8442864d6AB62a9193be2F16580c08E0D7BCda2f"),
			},
			want:      metadata.StandardERC1155,
			wantError: require.NoError,
		},
		{
			name: "Beacon Deposit Contract",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0x00000000219ab540356cbb839cbe05303d7705fa
				address: common.HexToAddress("0x00000000219ab540356cBB839Cbe05303d7705Fa"),
			},
			want:      metadata.StandardUnknown,
			wantError: require.NoError,
		},
		{
			name: "Arbitrum Bridge",
			arguments: arguments{
				ctx:     context.Background(),
				network: network.Ethereum,
				// https://etherscan.io/address/0x8315177ab297ba92a06054ce80a67ed4dbd7ed3a
				address: common.HexToAddress("0x8315177aB297bA92A06054cE80a67Ed4DBd7ed3a"),
			},
			want:      metadata.StandardUnknown,
			wantError: require.NoError,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			chainID, err := network.EthereumChainIDString(testcase.arguments.network.String())
			require.NoError(t, err)

			ethereumClient, err := ethereum.Dial(testcase.arguments.ctx, endpoint.MustGet(testcase.arguments.network))
			testcase.wantError(t, err)

			result, err := contract.DetectTokenStandard(testcase.arguments.ctx, uint64(chainID), testcase.arguments.address, nil, ethereumClient)
			testcase.wantError(t, err)

			require.Equal(t, result, testcase.want)
		})
	}
}

func BenchmarkDetectERC20WithCode(b *testing.B) {
	// RSS3 Token
	address := common.HexToAddress("0xc98D64DA73a6616c42117b582e832812e7B8D57F")
	code, err := hex.DecodeString("608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461012957806370a082311461013c57806395d89b411461014f578063a457c2d714610157578063a9059cbb1461016a578063dd62ed3e1461017d576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100ec57806323b872dd14610101578063313ce56714610114575b600080fd5b6100b6610190565b6040516100c391906106dd565b60405180910390f35b6100df6100da3660046106a9565b610222565b6040516100c391906106d2565b6100f461023f565b6040516100c39190610911565b6100df61010f36600461066e565b610245565b61011c6102de565b6040516100c3919061091a565b6100df6101373660046106a9565b6102e3565b6100f461014a36600461061b565b610337565b6100b6610356565b6100df6101653660046106a9565b610365565b6100df6101783660046106a9565b6103de565b6100f461018b36600461063c565b6103f2565b60606003805461019f9061094c565b80601f01602080910402602001604051908101604052809291908181526020018280546101cb9061094c565b80156102185780601f106101ed57610100808354040283529160200191610218565b820191906000526020600020905b8154815290600101906020018083116101fb57829003601f168201915b5050505050905090565b600061023661022f61041d565b8484610421565b50600192915050565b60025490565b60006102528484846104d5565b6001600160a01b03841660009081526001602052604081208161027361041d565b6001600160a01b03166001600160a01b03168152602001908152602001600020549050828110156102bf5760405162461bcd60e51b81526004016102b6906107fb565b60405180910390fd5b6102d3856102cb61041d565b858403610421565b506001949350505050565b601290565b60006102366102f061041d565b8484600160006102fe61041d565b6001600160a01b03908116825260208083019390935260409182016000908120918b16815292529020546103329190610928565b610421565b6001600160a01b0381166000908152602081905260409020545b919050565b60606004805461019f9061094c565b6000806001600061037461041d565b6001600160a01b03908116825260208083019390935260409182016000908120918816815292529020549050828110156103c05760405162461bcd60e51b81526004016102b6906108cc565b6103d46103cb61041d565b85858403610421565b5060019392505050565b60006102366103eb61041d565b84846104d5565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b3390565b6001600160a01b0383166104475760405162461bcd60e51b81526004016102b690610888565b6001600160a01b03821661046d5760405162461bcd60e51b81526004016102b690610773565b6001600160a01b0380841660008181526001602090815260408083209487168084529490915290819020849055517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906104c8908590610911565b60405180910390a3505050565b6001600160a01b0383166104fb5760405162461bcd60e51b81526004016102b690610843565b6001600160a01b0382166105215760405162461bcd60e51b81526004016102b690610730565b61052c8383836105ff565b6001600160a01b038316600090815260208190526040902054818110156105655760405162461bcd60e51b81526004016102b6906107b5565b6001600160a01b0380851660009081526020819052604080822085850390559185168152908120805484929061059c908490610928565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516105e69190610911565b60405180910390a36105f98484846105ff565b50505050565b505050565b80356001600160a01b038116811461035157600080fd5b60006020828403121561062c578081fd5b61063582610604565b9392505050565b6000806040838503121561064e578081fd5b61065783610604565b915061066560208401610604565b90509250929050565b600080600060608486031215610682578081fd5b61068b84610604565b925061069960208501610604565b9150604084013590509250925092565b600080604083850312156106bb578182fd5b6106c483610604565b946020939093013593505050565b901515815260200190565b6000602080835283518082850152825b81811015610709578581018301518582016040015282016106ed565b8181111561071a5783604083870101525b50601f01601f1916929092016040019392505050565b60208082526023908201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260408201526265737360e81b606082015260800190565b60208082526022908201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604082015261737360f01b606082015260800190565b60208082526026908201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604082015265616c616e636560d01b606082015260800190565b60208082526028908201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616040820152676c6c6f77616e636560c01b606082015260800190565b60208082526025908201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604082015264647265737360d81b606082015260800190565b60208082526024908201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646040820152637265737360e01b606082015260800190565b60208082526025908201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604082015264207a65726f60d81b606082015260800190565b90815260200190565b60ff91909116815260200190565b6000821982111561094757634e487b7160e01b81526011600452602481fd5b500190565b60028104600182168061096057607f821691505b6020821081141561098157634e487b7160e01b600052602260045260246000fd5b5091905056fea2646970667358221220dce3469df9bbc6af8b36a7047024a2509da1d09910cf2bdb5ec57acc4d7031b564736f6c63430008000033")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		require.False(b, contract.DetectERC165WithCode(uint64(network.EthereumChainIDMainnet), address, code))
	}
}
