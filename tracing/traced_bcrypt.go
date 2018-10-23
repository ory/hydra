package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	HashOpName        = "bcrypt.hash"
	CompareOpName     = "bcrypt.compare"
	WorkFactorTagName = "bcrypt.workfactor"
)

// TracedBCrypt implements the Hasher interface
type TracedBCrypt struct {
	WorkFactor int
}

func (b *TracedBCrypt) Hash(ctx context.Context, data []byte) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, HashOpName)
	defer span.Finish()
	span.SetTag(WorkFactorTagName, b.WorkFactor)

	s, err := bcrypt.GenerateFromPassword(data, b.WorkFactor)
	if err != nil {
		ext.Error.Set(span, true)
		return nil, errors.WithStack(err)
	}
	return s, nil
}

func (b *TracedBCrypt) Compare(ctx context.Context, hash, data []byte) error {
	span, _ := opentracing.StartSpanFromContext(ctx, CompareOpName)
	defer span.Finish()
	span.SetTag(WorkFactorTagName, b.WorkFactor)

	if err := bcrypt.CompareHashAndPassword(hash, data); err != nil {
		ext.Error.Set(span, true)
		return errors.WithStack(err)
	}
	return nil
}
