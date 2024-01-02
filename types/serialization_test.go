package types

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto/ed25519"
	cmstate "github.com/cometbft/cometbft/proto/tendermint/state"
	cmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cmtypes "github.com/cometbft/cometbft/types"

	pb "github.com/rollkit/rollkit/types/pb/rollkit"
)

func TestBlockSerializationRoundTrip(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	// create random hashes
	h := []Hash{}
	for i := 0; i < 8; i++ {
		h1 := make(Hash, 32)
		n, err := rand.Read(h1[:])
		require.Equal(32, n)
		require.NoError(err)
		h = append(h, h1)
	}

	h1 := Header{
		Version: Version{
			Block: 1,
			App:   2,
		},
		BaseHeader: BaseHeader{
			Height: 3,
			Time:   4567,
		},
		LastHeaderHash:  h[0],
		LastCommitHash:  h[1],
		DataHash:        h[2],
		ConsensusHash:   h[3],
		AppHash:         h[4],
		LastResultsHash: h[5],
		ProposerAddress: []byte{4, 3, 2, 1},
	}

	pubKey1 := ed25519.GenPrivKey().PubKey()
	pubKey2 := ed25519.GenPrivKey().PubKey()
	validator1 := &cmtypes.Validator{Address: pubKey1.Address(), PubKey: pubKey1}
	validator2 := &cmtypes.Validator{Address: pubKey2.Address(), PubKey: pubKey2}

	cases := []struct {
		name  string
		input *Block
	}{
		{"empty block", &Block{}},
		{"full", &Block{
			SignedHeader: SignedHeader{
				Header: h1,
				Commit: Commit{
					Signatures: []Signature{Signature([]byte{1, 1, 1}), Signature([]byte{2, 2, 2})},
				},
				Validators: &cmtypes.ValidatorSet{
					Validators: []*cmtypes.Validator{
						validator1,
						validator2,
					},
					Proposer: validator1,
				},
			},
			Data: Data{
				Txs:                    nil,
				IntermediateStateRoots: IntermediateStateRoots{RawRootsList: [][]byte{{0x1}}},
				// TODO(tzdybal): update when we have actual evidence types
				// Note: Temporarily remove Evidence #896
				// Evidence: EvidenceData{Evidence: nil},
			},
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)
			blob, err := c.input.MarshalBinary()
			assert.NoError(err)
			assert.NotEmpty(blob)

			deserialized := &Block{}
			err = deserialized.UnmarshalBinary(blob)
			assert.NoError(err)

			assert.Equal(c.input, deserialized)
		})
	}
}

func TestStateRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		state State
	}{
		{
			"with max bytes",
			State{
				ConsensusParams: cmproto.ConsensusParams{
					Block: &cmproto.BlockParams{
						MaxBytes: 123,
						MaxGas:   456,
					},
				},
			},
		},
		{
			name: "with all fields set",
			state: State{
				Version: cmstate.Version{
					Consensus: cmversion.Consensus{
						Block: 123,
						App:   456,
					},
					Software: "rollkit",
				},
				ChainID:         "testchain",
				InitialHeight:   987,
				LastBlockHeight: 987654321,
				LastBlockID: cmtypes.BlockID{
					Hash: nil,
					PartSetHeader: cmtypes.PartSetHeader{
						Total: 0,
						Hash:  nil,
					},
				},
				LastBlockTime: time.Date(2022, 6, 6, 12, 12, 33, 44, time.UTC),
				DAHeight:      3344,
				ConsensusParams: cmproto.ConsensusParams{
					Block: &cmproto.BlockParams{
						MaxBytes: 12345,
						MaxGas:   6543234,
					},
					Evidence: &cmproto.EvidenceParams{
						MaxAgeNumBlocks: 100,
						MaxAgeDuration:  200,
						MaxBytes:        300,
					},
					Validator: &cmproto.ValidatorParams{
						PubKeyTypes: []string{"secure", "more secure"},
					},
					Version: &cmproto.VersionParams{
						App: 42,
					},
				},
				LastHeightConsensusParamsChanged: 12345,
				LastResultsHash:                  Hash{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2},
				AppHash:                          Hash{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			pState, err := c.state.ToProto()
			require.NoError(err)
			require.NotNil(pState)

			bytes, err := pState.Marshal()
			require.NoError(err)
			require.NotEmpty(bytes)

			var newProtoState pb.State
			var newState State
			err = newProtoState.Unmarshal(bytes)
			require.NoError(err)

			err = newState.FromProto(&newProtoState)
			require.NoError(err)

			assert.Equal(c.state, newState)
		})
	}
}

func TestTxsRoundtrip(t *testing.T) {
	// Test the nil case
	var txs Txs
	byteSlices := txsToByteSlices(txs)
	newTxs := byteSlicesToTxs(byteSlices)
	assert.Nil(t, newTxs)

	// Generate 100 random transactions and convert them to byte slices
	txs = make(Txs, 100)
	for i := range txs {
		txs[i] = []byte{byte(i)}
	}
	byteSlices = txsToByteSlices(txs)

	// Convert the byte slices back to transactions
	newTxs = byteSlicesToTxs(byteSlices)

	// Check that the new transactions match the original transactions
	assert.Equal(t, len(txs), len(newTxs))
	for i := range txs {
		assert.Equal(t, txs[i], newTxs[i])
	}
}

func TestSignaturesRoundtrip(t *testing.T) {
	// Test the nil case
	var sigs []Signature
	bytes := signaturesToByteSlices(sigs)
	newSigs := byteSlicesToSignatures(bytes)
	assert.Nil(t, newSigs)

	// Generate 100 random signatures and convert them to byte slices
	sigs = make([]Signature, 100)
	for i := range sigs {
		sigs[i] = []byte{byte(i)}
	}
	bytes = signaturesToByteSlices(sigs)

	// Convert the byte slices back to signatures
	newSigs = byteSlicesToSignatures(bytes)

	// Check that the new signatures match the original signatures
	assert.Equal(t, len(sigs), len(newSigs))
	for i := range sigs {
		assert.Equal(t, newSigs[i], sigs[i])
	}
}

func TestConsensusParamsFromProto(t *testing.T) {
	// Prepare test case
	pbParams := cmproto.ConsensusParams{
		Block: &cmproto.BlockParams{
			MaxBytes: 12345,
			MaxGas:   67890,
		},
		Validator: &cmproto.ValidatorParams{
			PubKeyTypes: []string{cmtypes.ABCIPubKeyTypeEd25519},
		},
		Version: &cmproto.VersionParams{
			App: 42,
		},
	}

	// Call the function to be tested
	params := ConsensusParamsFromProto(pbParams)

	// Check the results
	assert.Equal(t, int64(12345), params.Block.MaxBytes)
	assert.Equal(t, int64(67890), params.Block.MaxGas)
	assert.Equal(t, uint64(42), params.Version.App)
	assert.Equal(t, []string{cmtypes.ABCIPubKeyTypeEd25519}, params.Validator.PubKeyTypes)
}
