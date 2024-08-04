package types

import (
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"go.dedis.ch/kyber/v3/suites"
)

// SecretCommits to SecretCommits
func SecretCommitsToProtocol(sc *rabin.SecretCommits) (*SecretCommitJson, error) {
	points := make([][]byte, len(sc.Commitments))
	for i, c := range sc.Commitments {
		cBytes, err := c.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("marshal commitment: %w", err)
		}
		points[i] = cBytes
	}

	return &SecretCommitJson{
		Index: sc.Index,
		// TargetIndex: sc.TargetIndex,
		Commitments: points,
		SessionID:   sc.SessionID,
		Signature:   sc.Signature,
	}, nil
}

func SecretCommitsFromProtocol(suite suites.Suite, sc *SecretCommitJson) (*rabin.SecretCommits, error) {
	// convert kyber points
	points := make([]kyber.Point, len(sc.Commitments))
	for i, c := range sc.Commitments {
		commitPoint := suite.Point()
		err := commitPoint.UnmarshalBinary(c)
		if err != nil {
			return nil, fmt.Errorf("unmarshal commitment: %w", err)
		}
		points[i] = commitPoint
	}

	return &rabin.SecretCommits{
		Index:       sc.Index,
		Commitments: points,
		SessionID:   sc.SessionID,
		Signature:   sc.Signature,
	}, nil
}

// SecretCommits
type SecretCommitJson struct {
	// Index of the Dealer in the list of participants
	Index uint32
	// Commitments generated by the Dealer
	Commitments [][]byte
	// SessionID generated by the Dealer tied to the Deal
	SessionID []byte
	// Signature from the Dealer
	Signature []byte
}

// DistKeyShare
type DistKeyShare struct {
	// Coefficients of the public polynomial holding the public key
	Commits []kyber.Point

	// PriShare of the distributed secret
	PriShare *share.PriShare
}

type PriShareJson struct {
	I int
	V []byte
}

type DistKeyShareJson struct {
	PriShare *PriShareJson
	Commits  [][]byte
}

// DistKeyShare to []byte
func DistKeyShareToProtocol(d *DistKeyShare) ([]byte, error) {
	share := d.PriShare

	// convert kyber point
	sbuf, err := share.V.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal private share: %w", err)
	}

	// convert kyber points
	var shareCommits [][]byte
	commits := d.Commits
	if commits != nil {
		shareCommits = make([][]byte, len(commits))
		for i, cmt := range commits {
			buf, err := cmt.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("couldn't marshal share commitment: %w", err)
			}
			shareCommits[i] = buf
		}
	}

	data := DistKeyShareJson{
		PriShare: &PriShareJson{
			I: share.I,
			V: sbuf,
		},
		Commits: shareCommits,
	}

	return json.Marshal(&data)
}

// []byte to DistKeyShare
func DistKeyShareFromProtocol(suite suites.Suite, data []byte) (*DistKeyShare, error) {
	jsonData := DistKeyShareJson{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal DistKeyShare: %w", err)
	}

	// convert kyber points
	points := make([]kyber.Point, len(jsonData.Commits))
	for i, c := range jsonData.Commits {
		commitPoint := suite.Point()
		err := commitPoint.UnmarshalBinary(c)
		if err != nil {
			return nil, fmt.Errorf("unmarshal commitment: %w", err)
		}
		points[i] = commitPoint
	}

	s := suite.Scalar()
	s.UnmarshalBinary(jsonData.PriShare.V)
	return &DistKeyShare{
		Commits: points,
		PriShare: &share.PriShare{
			I: jsonData.PriShare.I,
			V: s,
		},
	}, nil
}

type Secret struct {
	EncCmt  []byte   `json:"enc_cmt,omitempty"`  // encryption commitment
	EncScrt [][]byte `json:"enc_scrt,omitempty"` // enncrypted secret
}

func CidFromBytes(b []byte) (cid.Cid, error) {
	h, err := mh.Sum(b, mh.SHA2_256, -1)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.Raw, h), nil
}