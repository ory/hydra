package context

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ory-am/hydra/account"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

const (
	authKey key = 0
)

func NewContextFromAuthorization(ctx context.Context, req *http.Request, j *hjwt.JWT, s account.Storage, p policy.Storer) context.Context {
	authorization := req.Header.Get("Authorization")
	if len(authorization) <= len("Bearer ") {
		log.Printf("Authorization header has no Bearer: %v", req.Header)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	} else if authorization[:len("Bearer ")] != "Bearer " {
		log.Printf("Authorization header has no Bearer: %v", req.Header)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	data := authorization[len("Bearer "):]
	t, err := j.VerifyToken([]byte(data))
	if err != nil {
		log.Printf("%s token validation errored: %v", authorization, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	} else if !t.Valid {
		log.Printf("%s token invalid: %v", authorization, t.Valid)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	idclaim, ok := t.Claims["subject"]
	if !ok {
		log.Printf("Claim subject not found: %v", t.Claims)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}
	id, ok := idclaim.(string)
	if !ok {
		log.Printf("Claim subject not found: %v", t.Claims)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	user, err := s.Get(id)
	if err != nil {
		log.Printf("Subject not found in store: %v %v", t.Claims, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	policies, err := p.FindPoliciesForSubject(id)
	if err != nil {
		log.Printf("Subject not found in store: %v %v", t.Claims, err)
		return NewContextFromAuthValues(ctx, nil, nil, nil)
	}

	log.Printf("Authentication successfull: %v", t)
	return NewContextFromAuthValues(ctx, user, t, policies)
}

func TokenFromContext(ctx context.Context) (*jwt.Token, error) {
	args, ok := ctx.Value(authKey).(*authorization)
	if !ok {
		return nil, fmt.Errorf("Could not assert type for %v", ctx.Value(authKey))
	}
	return args.token, nil
}

func SubjectFromContext(ctx context.Context) (account.Account, error) {
	args, ok := ctx.Value(authKey).(*authorization)
	if !ok {
		return nil, fmt.Errorf("Could not assert type for %v", ctx.Value(authKey))
	}
	u, ok := args.subject.(*account.DefaultAccount)
	if !ok {
		return nil, fmt.Errorf("Could not assert type for %v", ctx.Value(authKey))
	}
	return u, nil
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
	return b && a.token.Valid
}

func NewContextFromAuthValues(ctx context.Context, subject account.Account, token *jwt.Token, policies []policy.Policy) context.Context {
	if subject == nil {
		subject = &account.DefaultAccount{}
	}
	if token == nil {
		token = &jwt.Token{}
	}
	if policies == nil {
		policies = []policy.Policy{}
	}

	return context.WithValue(ctx, authKey, &authorization{subject, token, policies})
}

type authorization struct {
	subject  account.Account
	token    *jwt.Token
	policies []policy.Policy
}
