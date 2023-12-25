package auth

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/matrix-org/dendrite/clientapi/auth/authtypes"
	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/clientapi/ratelimit"
	"github.com/matrix-org/dendrite/clientapi/userutil"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/gomatrixserverlib/spec"
	"github.com/matrix-org/util"
)

type GeodomePasswordRequest struct {
	Login

	Password string `json:"password"`
	Email    string `json:"username"`
}

// LoginTypeGeodome implements https://matrix.org/docs/spec/client_server/r0.6.1#password-based
type LoginTypeGeodome struct {
	UserApi       api.ClientUserAPI
	Config        *config.ClientAPI
	Rt            *ratelimit.RtFailedLogin
	InhibitDevice bool
	UserLoginAPI  api.UserLoginAPI
}

func (t *LoginTypeGeodome) Name() string {
	return authtypes.LoginTypeGeodome
}

func (t *LoginTypeGeodome) LoginFromJSON(ctx context.Context, reqBytes []byte) (*Login, LoginCleanupFunc, *util.JSONResponse) {
	var r GeodomePasswordRequest
	if err := httputil.UnmarshalJSON(reqBytes, &r); err != nil {
		return nil, nil, err
	}

	login, err := t.Login(ctx, &r)
	if err != nil {
		return nil, nil, err
	}
	login.InhibitDevice = t.InhibitDevice

	return login, func(context.Context, *util.JSONResponse) {}, nil
}

func (t *LoginTypeGeodome) Login(ctx context.Context, req interface{}) (*Login, *util.JSONResponse) {
	r := req.(*GeodomePasswordRequest)

	if r.Identifier.Address != "" {
		r.Email = r.Identifier.Address
	}
	if r.Medium == "" {
		r.Medium = email
		if r.Identifier.Medium != "" {
			r.Medium = r.Identifier.Medium
		}
	}

	if r.Email == "" {
		return nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.BadJSON("An email must be provided."),
		}
	}

	if len(r.Password) == 0 {
		return nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.BadJSON("A password must be provided."),
		}
	}

	var username string
	r.Email = strings.ToLower(r.Email)
	res := api.QueryLocalpartForThreePIDResponse{}
	err := t.UserApi.QueryLocalpartForThreePID(ctx, &api.QueryLocalpartForThreePIDRequest{
		ThreePID: r.Email,
		Medium:   r.Medium,
	}, &res)
	if err != nil && err != sql.ErrNoRows {
		util.GetLogger(ctx).WithError(err).Error("userApi.QueryLocalpartForThreePID failed")
		return nil, &util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: spec.Unknown(""),
		}
	}
	if res.Localpart != "" {
		username = userutil.MakeUserID(res.Localpart, res.ServerName)
	} else {
		// Case where the user is not registered yet
		nreq := &api.QueryNumericLocalpartRequest{
			ServerName: t.Config.Matrix.ServerName,
		}
		nres := &api.QueryNumericLocalpartResponse{}
		err = t.UserApi.QueryNumericLocalpart(ctx, nreq, nres)
		if err != nil {
			util.GetLogger(ctx).WithError(err).Error("userAPI.QueryNumericLocalpart failed")
			return nil, &util.JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: spec.InternalServerError{},
			}
		}
		username = strconv.FormatInt(nres.ID, 10)
	}

	localpart, domain, err := userutil.ParseUsernameParam(username, t.Config.Matrix)
	if err != nil {
		return nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.InvalidUsername(err.Error()),
		}
	}

	if !t.Config.Matrix.IsLocalServerName(domain) {
		return nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.InvalidUsername("The server name is not known."),
		}
	}

	if t.Rt != nil {
		ok, retryIn := t.Rt.CanAct(localpart)
		if !ok {
			return nil, &util.JSONResponse{
				Code: http.StatusTooManyRequests,
				JSON: spec.LimitExceeded("Too Many Requests", retryIn.Milliseconds()),
			}
		}
	}

	urlEncodedBody := `username=` + r.Email + `&password=` + r.Password + `&_xsrf=`
	resp, authErr := http.Post(t.Config.GeodomeAuthEndpoint, "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewBuffer([]byte(urlEncodedBody)))
	if authErr != nil {
		return nil, &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: spec.UserInUse("Failed to authenticate against external auth provider: " + authErr.Error()),
		}
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.Forbidden("Invalid username or password: " + authErr.Error()),
		}
	} else if resp.StatusCode != http.StatusOK {
		return nil, &util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: spec.Unknown("Unknown error while authenticating against external auth provider: " + authErr.Error()),
		}
	}
	r.ExternalAuthResp = resp

	created, acc, errResp := t.getOrCreateAccount(ctx, localpart, r.Password, domain)
	if errResp != nil {
		return nil, errResp
	}
	// Save 3PID for newly created users
	if created {
		err := t.UserApi.PerformSaveThreePIDAssociation(ctx, &api.PerformSaveThreePIDAssociationRequest{
			Medium:     r.Medium,
			ThreePID:   r.Email,
			Localpart:  acc.Localpart,
			ServerName: acc.ServerName,
		}, &struct{}{})
		if err != nil {
			return nil, &util.JSONResponse{
				Code: http.StatusBadRequest,
				JSON: spec.UserInUse("Failed to save 3PID association: " + err.Error()),
			}
		}
	}

	// Set the user, so login.Username() can do the right thing
	r.Identifier.User = acc.UserID
	r.User = acc.UserID
	r.Login.User = username
	return &r.Login, nil
}

func (t *LoginTypeGeodome) getOrCreateAccount(ctx context.Context, localpart, password string, domain spec.ServerName) (bool, *api.Account, *util.JSONResponse) {
	var existing api.QueryAccountByLocalpartResponse
	err := t.UserLoginAPI.QueryAccountByLocalpart(ctx, &api.QueryAccountByLocalpartRequest{
		Localpart:  localpart,
		ServerName: domain,
	}, &existing)

	if err == nil {
		return false, existing.Account, nil
	}
	if err != sql.ErrNoRows {
		return false, nil, &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: spec.InvalidUsername(err.Error()),
		}
	}

	accountType := api.AccountTypeUser
	var created api.PerformAccountCreationResponse
	err = t.UserLoginAPI.PerformAccountCreation(ctx, &api.PerformAccountCreationRequest{
		AppServiceID: "geodome",
		Localpart:    localpart,
		Password:     password,
		AccountType:  accountType,
		OnConflict:   api.ConflictAbort,
	}, &created)

	if err != nil {
		if _, ok := err.(*api.ErrorConflict); ok {
			return false, nil, &util.JSONResponse{
				Code: http.StatusBadRequest,
				JSON: spec.UserInUse("Desired user ID is already taken."),
			}
		}
		return false, nil, &util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: spec.Unknown("failed to create account: " + err.Error()),
		}
	}
	return true, created.Account, nil
}
