package stake

import (
	"encoding/hex"

	legacy "github.com/lbryio/types/v1/go"
	pb "github.com/lbryio/types/v2/go"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/proto"
)

func (c *Helper) serialized() ([]byte, error) {
	serialized := c.Claim.String() + c.Support.String()
	if serialized == "" {
		return nil, errors.WithStack(errors.New("not initialized"))
	}

	if c.LegacyClaim != nil {
		return proto.Marshal(c.getLegacyProtobuf())
	} else if c.IsSupport() {
		return proto.Marshal(c.getSupportProtobuf())
	}

	return proto.Marshal(c.getClaimProtobuf())
}

func (c *Helper) getClaimProtobuf() *pb.Claim {
	claim := &pb.Claim{
		Title:       c.Claim.GetTitle(),
		Description: c.Claim.GetDescription(),
		Thumbnail:   c.Claim.GetThumbnail(),
		Tags:        c.Claim.GetTags(),
		Languages:   c.Claim.GetLanguages(),
		Locations:   c.Claim.GetLocations(),
	}
	if c.Claim.GetChannel() != nil {
		claim.Type = &pb.Claim_Channel{Channel: c.Claim.GetChannel()}
	} else if c.GetStream() != nil {
		claim.Type = &pb.Claim_Stream{Stream: c.GetStream()}
	} else if c.Claim.GetCollection() != nil {
		claim.Type = &pb.Claim_Collection{Collection: c.Claim.GetCollection()}
	} else if c.Claim.GetRepost() != nil {
		claim.Type = &pb.Claim_Repost{Repost: c.Claim.GetRepost()}
	}

	return claim
}

func (c *Helper) getSupportProtobuf() *pb.Support {
	return &pb.Support{
		Emoji:                c.Support.GetEmoji(),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func (c *Helper) getLegacyProtobuf() *legacy.Claim {
	v := c.LegacyClaim.GetVersion()
	t := c.LegacyClaim.GetClaimType()
	return &legacy.Claim{
		Version:            &v,
		ClaimType:          &t,
		Stream:             c.LegacyClaim.GetStream(),
		Certificate:        c.LegacyClaim.GetCertificate(),
		PublisherSignature: c.LegacyClaim.GetPublisherSignature()}
}

func (c *Helper) serializedHexString() (string, error) {
	serialized, err := c.serialized()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return hex.EncodeToString(serialized), nil
}

func (c *Helper) serializedNoSignature() ([]byte, error) {
	if c.Claim.String() == "" && c.Support.String() == "" {
		return nil, errors.WithStack(errors.New("not initialized"))
	}
	if c.Signature == nil {
		serialized, err := c.serialized()
		if err != nil {
			return nil, err
		}
		return serialized, nil
	} else {
		if c.LegacyClaim != nil {
			clone := &legacy.Claim{}
			proto.Merge(clone, c.getLegacyProtobuf())
			proto.ClearAllExtensions(clone.PublisherSignature)
			clone.PublisherSignature = nil
			return proto.Marshal(clone)
		} else if c.IsSupport() {
			clone := &pb.Support{}
			proto.Merge(clone, c.getSupportProtobuf())
			return proto.Marshal(clone)
		}
		clone := &pb.Claim{}
		proto.Merge(clone, c.getClaimProtobuf())
		return proto.Marshal(clone)
	}
}
