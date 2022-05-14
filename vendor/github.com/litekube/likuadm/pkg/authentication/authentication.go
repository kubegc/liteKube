package authentication

import "context"

type TokenAuthentication struct {
	Token string
}

func (t *TokenAuthentication) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": t.Token,
	}, nil
}

func (t *TokenAuthentication) RequireTransportSecurity() bool {
	return false
}
