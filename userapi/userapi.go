// Copyright 2020 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package userapi

import (
	"github.com/gorilla/mux"
	keyapi "github.com/matrix-org/dendrite/keyserver/api"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/dendrite/userapi/internal"
	"github.com/matrix-org/dendrite/userapi/inthttp"
	"github.com/matrix-org/dendrite/userapi/mail"
	"github.com/matrix-org/dendrite/userapi/storage/accounts"
	"github.com/matrix-org/dendrite/userapi/storage/devices"
	"github.com/matrix-org/dendrite/userapi/storage/threepid"
	"github.com/sirupsen/logrus"
)

// AddInternalRoutes registers HTTP handlers for the internal API. Invokes functions
// on the given input API.
func AddInternalRoutes(router *mux.Router, intAPI api.UserInternalAPI) {
	inthttp.AddRoutes(router, intAPI)
}

// NewInternalAPI returns a concerete implementation of the internal API. Callers
// can call functions directly on the returned API or via an HTTP interface using AddInternalRoutes.
func NewInternalAPI(
	accountDB accounts.Database, cfg *config.UserAPI, appServices []config.ApplicationService, keyAPI keyapi.KeyInternalAPI,
) *internal.UserInternalAPI {
	deviceDB, err := devices.NewDatabase(&cfg.DeviceDatabase, cfg.Matrix.ServerName)
	if err != nil {
		logrus.WithError(err).Panicf("failed to connect to device db")
	}
	threepidDb, err := threepid.NewDatabase(&cfg.ThreepidDatabase)
	if err != nil {
		logrus.WithError(err).Panicf("failed to connect to threepid db")
	}
	mailer, err := mail.NewMailer(cfg)
	if err != nil {
		logrus.WithError(err).Panicf("failed to crate Mailer")
	}
	return &internal.UserInternalAPI{
		AccountDB:   accountDB,
		DeviceDB:    deviceDB,
		ThreePidDB:  threepidDb,
		ServerName:  cfg.Matrix.ServerName,
		AppServices: appServices,
		KeyAPI:      keyAPI,
		Mail:        mailer,
	}
}
