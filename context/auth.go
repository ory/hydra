package context

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/RangelReale/osin"
	log "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/dgrijalva/jwt-go"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/go-errors/errors"
	. "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/handler"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/ladon/policy"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/net/context"
	hjwt "github.com/ory-am/hydra/jwt"
	"net/http"
)

var (
	AuthKey Key = 0
)

func NewContextFromAuthorization(ctx context.Context, req *http.Request, j *hjwt.JWT, p policy.Storage) context.Context {
	bearer := osin.CheckBearerAuth(req)
	if bearer == nil {
		log.Warn("No authorization bearer given.")
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	t, err := j.VerifyToken([]byte(bearer.Code))
	if err != nil {
		log.Warnf(`Token validation errored: "%v".`, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	} else if !t.Valid {
		log.Warn("Token is invalid.")
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	claims := hjwt.ClaimsCarrier(t.Claims)
	user := claims.GetSubject()
	if user == "" {
		log.Warnf(`Sub claim may not be empty, to: "%v".`, t.Claims)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	policies, err := p.FindPoliciesForSubject(user)
	if err != nil {
		log.Warnf(`Policies for "%s" could not be retrieved: "%v"`, user, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	//	user, err := s.Get(id)
	//	if err != nil {
	//		log.Warnf("Subject not found in store: %v %v", t.Claims, err)
	//		return NewContextFromAuthValues(ctx, nil, nil, nil)
	//	}

	return NewContextFromAuthValues(ctx, claims, t, policies)
}

func TokenFromContext(ctx context.Context) (*jwt.Token, error) {
	args, ok := ctx.Value(AuthKey).(*authorization)
	if !ok {
		return nil, errors.Errorf("Could not assert type for %v", ctx.Value(AuthKey))
	}
	return args.token, nil
}

func SubjectFromContext(ctx context.Context) (string, error) {
	args, ok := ctx.Value(AuthKey).(*authorization)
	if !ok {
		return "", errors.Errorf("Could not assert type for %v", ctx.Value(AuthKey))
	}
	return args.claims.GetSubject(), nil
}

func PoliciesFromContext(ctx context.Context) ([]policy.Policy, error) {
	args, ok := ctx.Value(AuthKey).(*authorization)
	if !ok {
		return nil, errors.Errorf("Could not assert array type for %s", ctx.Value(AuthKey))
	}

	symbols := make([]policy.Policy, len(args.policies))
	for i, arg := range args.policies {
		symbols[i], ok = arg.(*policy.DefaultPolicy)
		if !ok {
			return nil, errors.Errorf("Could not assert policy type for %s", ctx.Value(AuthKey))
		}
	}

	return symbols, nil
}

func IsAuthenticatedFromContext(ctx context.Context) bool {
	a, b := ctx.Value(AuthKey).(*authorization)
	return (b && a.token != nil && a.token.Valid)
}

func NewContextFromAuthValues(ctx context.Context, claims hjwt.ClaimsCarrier, token *jwt.Token, policies []policy.Policy) context.Context {
	return context.WithValue(ctx, AuthKey, &authorization{
		claims:   claims,
		token:    token,
		policies: policies,
	})
}

type authorization struct {
	claims   hjwt.ClaimsCarrier
	token    *jwt.Token
	policies []policy.Policy
}
