package internal

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/matrix-org/dendrite/internal"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/dendrite/userapi/mail"
	"github.com/matrix-org/dendrite/userapi/storage/threepid"
)

const (
	sessionIdByteLength = 32
	tokenByteLength     = 48
)

var ErrBadSession = errors.New("provided sid, client_secret and token does not point to valid session")

func (a *UserInternalAPI) CreateSession(ctx context.Context, req *api.CreateSessionRequest, res *api.CreateSessionResponse) error {
	s, err := a.ThreePidDB.GetSessionByThreePidAndSecret(ctx, req.ThreePid, req.ClientSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			token, err := internal.GenerateBlob(tokenByteLength)
			if err != nil {
				return err
			}
			sid, err := internal.GenerateBlob(sessionIdByteLength)
			if err != nil {
				return err
			}
			s = &api.Session{
				Sid:          sid,
				ClientSecret: req.ClientSecret,
				ThreePid:     req.ThreePid,
				SendAttempt:  req.SendAttempt,
				Token:        token,
				NextLink:     req.NextLink,
			}
			err = a.ThreePidDB.InsertSession(ctx, s)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if req.SendAttempt > s.SendAttempt {
			err = a.ThreePidDB.UpdateSendAttemptNextLink(ctx, s.Sid, req.NextLink)
			if err != nil {
				return err
			}
		} else {
			res.Sid = s.Sid
			return nil
		}
	}
	res.Sid = s.Sid
	// TODO - if we fail sending email, send_attempt for next requests must be bumped,
	// otherwise we will just return nil from this function and not sent email
	return a.Mail.Send(&mail.Mail{
		To:    s.ThreePid,
		Link:  s.NextLink,
		Token: s.Token,
		Extra: req.Extra,
	}, req.SessionType)
}

func (a *UserInternalAPI) ValidateSession(ctx context.Context, req *api.ValidateSessionRequest, res struct{}) error {
	s, err := getSessionByOwnership(ctx, &req.SessionOwnership, a.ThreePidDB)
	if err != nil {
		return err
	}
	if s.Token != req.Token {
		return ErrBadSession
	}
	return a.ThreePidDB.ValidateSession(ctx, s.Sid, int(time.Now().Unix()))
}

func (a *UserInternalAPI) GetThreePidForSession(ctx context.Context, req *api.SessionOwnership, res *api.GetThreePidForSessionResponse) error {
	s, err := getSessionByOwnership(ctx, req, a.ThreePidDB)
	if err != nil {
		return err
	}
	res.ThreePid = s.ThreePid
	return nil
}

func (a *UserInternalAPI) DeleteSession(ctx context.Context, req *api.SessionOwnership, res struct{}) error {
	s, err := getSessionByOwnership(ctx, req, a.ThreePidDB)
	if err != nil {
		return err
	}
	return a.ThreePidDB.DeleteSession(ctx, s.Sid)
}

func (a *UserInternalAPI) IsSessionValidated(ctx context.Context, req *api.SessionOwnership, res *api.IsSessionValidatedResponse) error {
	s, err := getSessionByOwnership(ctx, req, a.ThreePidDB)
	if err != nil {
		return err
	}
	res.Validated = s.Validated
	res.ValidatedAt = s.ValidatedAt
	return nil
}

func getSessionByOwnership(ctx context.Context, so *api.SessionOwnership, d threepid.Database) (*api.Session, error) {
	s, err := d.GetSession(ctx, so.Sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBadSession
		}
		return nil, err
	}
	if s.ClientSecret != so.ClientSecret {
		return nil, ErrBadSession
	}
	return s, err
}
