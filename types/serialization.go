package types

import (
	cmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"

	pb "github.com/rollkit/rollkit/types/pb/rollkit"
)

// MarshalBinary encodes Block into binary form and returns it.
func (b *Block) MarshalBinary() ([]byte, error) {
	bp, err := b.ToProto()
	if err != nil {
		return nil, err
	}
	return bp.Marshal()
}

// UnmarshalBinary decodes binary form of Block into object.
func (b *Block) UnmarshalBinary(data []byte) error {
	var pBlock pb.Block
	err := pBlock.Unmarshal(data)
	if err != nil {
		return err
	}
	err = b.FromProto(&pBlock)
	return err
}

// MarshalBinary encodes Header into binary form and returns it.
func (h *Header) MarshalBinary() ([]byte, error) {
	return h.ToProto().Marshal()
}

// UnmarshalBinary decodes binary form of Header into object.
func (h *Header) UnmarshalBinary(data []byte) error {
	var pHeader pb.Header
	err := pHeader.Unmarshal(data)
	if err != nil {
		return err
	}
	err = h.FromProto(&pHeader)
	return err
}

// MarshalBinary encodes Data into binary form and returns it.
func (d *Data) MarshalBinary() ([]byte, error) {
	return d.ToProto().Marshal()
}

// UnmarshalBinary decodes binary form of Data into object.
func (d *Data) UnmarshalBinary(data []byte) error {
	var pData pb.Data
	err := pData.Unmarshal(data)
	if err != nil {
		return err
	}
	err = d.FromProto(&pData)
	return err
}

// MarshalBinary encodes Commit into binary form and returns it.
func (c *Commit) MarshalBinary() ([]byte, error) {
	return c.ToProto().Marshal()
}

// UnmarshalBinary decodes binary form of Commit into object.
func (c *Commit) UnmarshalBinary(data []byte) error {
	var pCommit pb.Commit
	err := pCommit.Unmarshal(data)
	if err != nil {
		return err
	}
	err = c.FromProto(&pCommit)
	return err
}

// ToProto converts SignedHeader into protobuf representation and returns it.
func (sh *SignedHeader) ToProto() (*pb.SignedHeader, error) {
	vSet, err := sh.Validators.ToProto()
	if err != nil {
		return nil, err
	}
	return &pb.SignedHeader{
		Header:     sh.Header.ToProto(),
		Commit:     sh.Commit.ToProto(),
		Validators: vSet,
	}, nil
}

// FromProto fills SignedHeader with data from protobuf representation.
func (sh *SignedHeader) FromProto(other *pb.SignedHeader) error {
	err := sh.Header.FromProto(other.Header)
	if err != nil {
		return err
	}
	err = sh.Commit.FromProto(other.Commit)
	if err != nil {
		return err
	}

	if other.Validators != nil && other.Validators.GetProposer() != nil {
		validators, err := types.ValidatorSetFromProto(other.Validators)
		if err != nil {
			return err
		}

		sh.Validators = validators
	}
	return nil
}

// MarshalBinary encodes SignedHeader into binary form and returns it.
func (sh *SignedHeader) MarshalBinary() ([]byte, error) {
	hp, err := sh.ToProto()
	if err != nil {
		return nil, err
	}
	return hp.Marshal()
}

// UnmarshalBinary decodes binary form of SignedHeader into object.
func (sh *SignedHeader) UnmarshalBinary(data []byte) error {
	var pHeader pb.SignedHeader
	err := pHeader.Unmarshal(data)
	if err != nil {
		return err
	}
	err = sh.FromProto(&pHeader)
	if err != nil {
		return err
	}
	return nil
}

// ToProto converts Header into protobuf representation and returns it.
func (h *Header) ToProto() *pb.Header {
	return &pb.Header{
		Version: &pb.Version{
			Block: h.Version.Block,
			App:   h.Version.App,
		},
		Height:          h.BaseHeader.Height,
		Time:            h.BaseHeader.Time,
		LastHeaderHash:  h.LastHeaderHash[:],
		LastCommitHash:  h.LastCommitHash[:],
		DataHash:        h.DataHash[:],
		ConsensusHash:   h.ConsensusHash[:],
		AppHash:         h.AppHash[:],
		LastResultsHash: h.LastResultsHash[:],
		ProposerAddress: h.ProposerAddress[:],
		ChainId:         h.BaseHeader.ChainID,
	}
}

// FromProto fills Header with data from its protobuf representation.
func (h *Header) FromProto(other *pb.Header) error {
	h.Version.Block = other.Version.Block
	h.Version.App = other.Version.App
	h.BaseHeader.ChainID = other.ChainId
	h.BaseHeader.Height = other.Height
	h.BaseHeader.Time = other.Time
	h.LastHeaderHash = other.LastHeaderHash
	h.LastCommitHash = other.LastCommitHash
	h.DataHash = other.DataHash
	h.ConsensusHash = other.ConsensusHash
	h.AppHash = other.AppHash
	h.LastResultsHash = other.LastResultsHash
	if len(other.ProposerAddress) > 0 {
		h.ProposerAddress = make([]byte, len(other.ProposerAddress))
		copy(h.ProposerAddress, other.ProposerAddress)
	}

	return nil
}

// ToProto converts Block into protobuf representation and returns it.
func (b *Block) ToProto() (*pb.Block, error) {
	sp, err := b.SignedHeader.ToProto()
	if err != nil {
		return nil, err
	}
	return &pb.Block{
		SignedHeader: sp,
		Data:         b.Data.ToProto(),
	}, nil
}

// ToProto converts Data into protobuf representation and returns it.
func (d *Data) ToProto() *pb.Data {
	return &pb.Data{
		Txs:                    txsToByteSlices(d.Txs),
		IntermediateStateRoots: d.IntermediateStateRoots.RawRootsList,
		// Note: Temporarily remove Evidence #896
		// Evidence:               evidenceToProto(d.Evidence),
	}
}

// FromProto fills Block with data from its protobuf representation.
func (b *Block) FromProto(other *pb.Block) error {
	err := b.SignedHeader.FromProto(other.SignedHeader)
	if err != nil {
		return err
	}
	err = b.Data.FromProto(other.Data)
	if err != nil {
		return err
	}

	return nil
}

// FromProto fills the Data with data from its protobuf representation
func (d *Data) FromProto(other *pb.Data) error {
	d.Txs = byteSlicesToTxs(other.Txs)
	d.IntermediateStateRoots.RawRootsList = other.IntermediateStateRoots
	// Note: Temporarily remove Evidence #896
	// d.Evidence = evidenceFromProto(other.Evidence)

	return nil
}

// ToProto converts Commit into protobuf representation and returns it.
func (c *Commit) ToProto() *pb.Commit {
	return &pb.Commit{
		Signatures: signaturesToByteSlices(c.Signatures),
	}
}

// FromProto fills Commit with data from its protobuf representation.
func (c *Commit) FromProto(other *pb.Commit) error {
	c.Signatures = byteSlicesToSignatures(other.Signatures)

	return nil
}

// ToProto converts State into protobuf representation and returns it.
func (s *State) ToProto() (*pb.State, error) {

	return &pb.State{
		Version:                          &s.Version,
		ChainId:                          s.ChainID,
		InitialHeight:                    s.InitialHeight,
		LastBlockHeight:                  s.LastBlockHeight,
		LastBlockID:                      s.LastBlockID.ToProto(),
		LastBlockTime:                    s.LastBlockTime,
		DAHeight:                         s.DAHeight,
		ConsensusParams:                  s.ConsensusParams,
		LastHeightConsensusParamsChanged: s.LastHeightConsensusParamsChanged,
		LastResultsHash:                  s.LastResultsHash[:],
		AppHash:                          s.AppHash[:],
	}, nil
}

// FromProto fills State with data from its protobuf representation.
func (s *State) FromProto(other *pb.State) error {
	var err error
	s.Version = *other.Version
	s.ChainID = other.ChainId
	s.InitialHeight = other.InitialHeight
	s.LastBlockHeight = other.LastBlockHeight
	lastBlockID, err := types.BlockIDFromProto(&other.LastBlockID)
	if err != nil {
		return err
	}
	s.LastBlockID = *lastBlockID
	s.LastBlockTime = other.LastBlockTime
	s.DAHeight = other.DAHeight
	s.ConsensusParams = other.ConsensusParams
	s.LastHeightConsensusParamsChanged = other.LastHeightConsensusParamsChanged
	s.LastResultsHash = other.LastResultsHash
	s.AppHash = other.AppHash

	return nil
}

func txsToByteSlices(txs Txs) [][]byte {
	if txs == nil {
		return nil
	}
	bytes := make([][]byte, len(txs))
	for i := range txs {
		bytes[i] = txs[i]
	}
	return bytes
}

func byteSlicesToTxs(bytes [][]byte) Txs {
	if len(bytes) == 0 {
		return nil
	}
	txs := make(Txs, len(bytes))
	for i := range txs {
		txs[i] = bytes[i]
	}
	return txs
}

// Note: Temporarily remove Evidence #896

// func evidenceToProto(evidence EvidenceData) []*abci.Evidence {
// 	var ret []*abci.Evidence
// 	for _, e := range evidence.Evidence {
// 		for i := range e.ABCI() {
// 			ae := e.ABCI()[i]
// 			ret = append(ret, &ae)
// 		}
// 	}
// 	return ret
// }

// func evidenceFromProto(evidence []*abci.Evidence) EvidenceData {
// 	var ret EvidenceData
// 	// TODO(tzdybal): right now Evidence is just an interface without implementations
// 	return ret
// }

func signaturesToByteSlices(sigs []Signature) [][]byte {
	if sigs == nil {
		return nil
	}
	bytes := make([][]byte, len(sigs))
	for i := range sigs {
		bytes[i] = sigs[i]
	}
	return bytes
}

func byteSlicesToSignatures(bytes [][]byte) []Signature {
	if bytes == nil {
		return nil
	}
	sigs := make([]Signature, len(bytes))
	for i := range bytes {
		sigs[i] = bytes[i]
	}
	return sigs
}

// ConsensusParamsFromProto converts protobuf consensus parameters to consensus parameters
func ConsensusParamsFromProto(pbParams cmproto.ConsensusParams) types.ConsensusParams {
	c := types.ConsensusParams{
		Block: types.BlockParams{
			MaxBytes: pbParams.Block.MaxBytes,
			MaxGas:   pbParams.Block.MaxGas,
		},
		Validator: types.ValidatorParams{
			PubKeyTypes: pbParams.Validator.PubKeyTypes,
		},
		Version: types.VersionParams{
			App: pbParams.Version.App,
		},
	}
	if pbParams.Abci != nil {
		c.ABCI.VoteExtensionsEnableHeight = pbParams.Abci.GetVoteExtensionsEnableHeight()
	}
	return c
}
