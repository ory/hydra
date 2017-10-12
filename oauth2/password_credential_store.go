package oauth2

import (
	"context"
)

type ResourceOwnerPasswordCredentialsGrantStore struct{

}

func (storage *ResourceOwnerPasswordCredentialsGrantStore) Authenticate(ctx context.Context, name string, secret string) error {
	return nil
}