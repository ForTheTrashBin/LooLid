package nutsandbolts

import (
	"context"
	"io"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

type ContextReader struct {
	Context context.Context
	Reader  io.Reader
}

func (ctxRdr *ContextReader) Read(p []byte) (int, error) {

	select {

	case <-ctxRdr.Context.Done():

		return 0, ctxRdr.Context.Err()

	default:
	}

	return ctxRdr.Reader.Read(p)
}
