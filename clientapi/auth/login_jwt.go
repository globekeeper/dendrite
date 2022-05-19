package auth

import (
	"context"

	"github.com/matrix-org/dendrite/clientapi/auth/authtypes"
	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/util"
)

// LoginTypeToken describes how to authenticate with a login token.
type LoginTypeTokenJwt struct {
	// UserAPI uapi.LoginTokenInternalAPI
	Config *config.ClientAPI
}

// Name implements Type.
func (t *LoginTypeTokenJwt) Name() string {
	return authtypes.LoginTypeJwt
}

// LoginFromJSON implements Type. The cleanup function deletes the token from
// the database on success.
func (t *LoginTypeTokenJwt) LoginFromJSON(ctx context.Context, reqBytes []byte) (*Login, LoginCleanupFunc, *util.JSONResponse) {
	var r loginTokenRequest
	if err := httputil.UnmarshalJSON(reqBytes, &r); err != nil {
		return nil, nil, err
	}

	// TODO-entry validate JWT here and extract userID, follow synapse implementation

	r.Login.Identifier.Type = "m.id.user"
	// r.Login.Identifier.User = //res.Data.UserID

	return &r.Login, nil, nil
}
