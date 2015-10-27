package context

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

const (
	authKey key = 0
)

func NewContextFromAuthorization(ctx context.Context, req *http.Request, j *hjwt.JWT, p policy.Storer) context.Context {
	bearer := osin.CheckBearerAuth(req)
	if bearer == nil {
		log.Printf("No bearer given: %v %v", bearer, req.Header)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	t, err := j.VerifyToken([]byte(bearer.Code))
	if err != nil {
		log.Printf("%s token validation errored: %v", bearer, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	} else if !t.Valid {
		log.Printf("%s token invalid: %v", bearer, t.Valid)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	claims := hjwt.ClaimsCarrier(t.Claims)
	user := claims.Subject()
	if user == "" {
		log.Printf("Subject not claimed: %v", t.Claims)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	policies, err := p.FindPoliciesForSubject(user)
	if err != nil {
		log.Printf("Subject not found in store: %v %v", t.Claims, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	//	user, err := s.Get(id)
	//	if err != nil {
	//		log.Printf("Subject not found in store: %v %v", t.Claims, err)
	//		return NewContextFromAuthValues(ctx, nil, nil, nil)
	//	}

	log.Printf("Authentication successfull: %v", t)
	return NewContextFromAuthValues(ctx, claims, t, policies)
}

func TokenFromContext(ctx context.Context) (*jwt.Token, error) {
	args, ok := ctx.Value(authKey).(*authorization)
	if !ok {
		return nil, fmt.Errorf("Could not assert type for %v", ctx.Value(authKey))
	}
	return args.token, nil
}

func SubjectFromContext(ctx context.Context) (string, error) {
	args, ok := ctx.Value(authKey).(*authorization)
	if !ok {
		return "", fmt.Errorf("Could not assert type for %v", ctx.Value(authKey))
	}
	return args.claims.Subject(), nil
}

func PoliciesFromContext(ctx context.Context) ([]policy.Policy, error) {
	args, ok := ctx.Value(authKey).(*authorization)
	if !ok {
		return nil, fmt.Errorf("Could not assert array type for %s", ctx.Value(authKey))
	}

	symbols := make([]policy.Policy, len(args.policies))
	for i, arg := range args.policies {
		symbols[i], ok = arg.(*policy.DefaultPolicy)
		if !ok {
			return nil, fmt.Errorf("Could not assert policy type for %s", ctx.Value(authKey))
		}
	}

	return symbols, nil
}

func IsAuthenticatedFromContext(ctx context.Context) bool {
	a, b := ctx.Value(authKey).(*authorization)
	return (b && a.token != nil && a.token.Valid)
}

func NewContextFromAuthValues(ctx context.Context, claims hjwt.ClaimsCarrier, token *jwt.Token, policies []policy.Policy) context.Context {
	return context.WithValue(ctx, authKey, &authorization{
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
