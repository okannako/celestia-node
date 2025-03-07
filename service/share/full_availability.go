package share

import (
	"context"
	"errors"

	"github.com/ipfs/go-blockservice"
	format "github.com/ipfs/go-ipld-format"

	"github.com/celestiaorg/celestia-node/ipld"
)

// fullAvailability implements Availability using the full data square
// recovery technique. It is considered "full" because it is required
// to download enough shares to fully reconstruct the data square.
type fullAvailability struct {
	rtrv *ipld.Retriever
}

// NewFullAvailability creates a new full Availability.
func NewFullAvailability(bServ blockservice.BlockService) Availability {
	return &fullAvailability{
		rtrv: ipld.NewRetriever(bServ),
	}
}

// SharesAvailable reconstructs the data committed to the given Root by requesting
// enough Shares from the network.
func (fa *fullAvailability) SharesAvailable(ctx context.Context, root *Root) error {
	ctx, cancel := context.WithTimeout(ctx, AvailabilityTimeout)
	defer cancel()
	// we assume the caller of this method has already performed basic validation on the
	// given dah/root. If for some reason this has not happened, the node should panic.
	if err := root.ValidateBasic(); err != nil {
		log.Errorw("Availability validation cannot be performed on a malformed DataAvailabilityHeader",
			"err", err)
		panic(err)
	}

	_, err := fa.rtrv.Retrieve(ctx, root)
	if err != nil {
		log.Errorw("availability validation failed", "root", root.Hash(), "err", err)
		if format.IsNotFound(err) || errors.Is(err, context.DeadlineExceeded) {
			return ErrNotAvailable
		}

		return err
	}
	return err
}
