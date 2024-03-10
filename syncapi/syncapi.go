// Copyright 2017 Vector Creations Ltd
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

package syncapi

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/matrix-org/dendrite/internal/fulltext"
	"github.com/matrix-org/dendrite/internal/httputil"
	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/setup/process"
	"github.com/sirupsen/logrus"

	"github.com/matrix-org/dendrite/internal/caching"

	"github.com/matrix-org/dendrite/roomserver/api"
	"github.com/matrix-org/dendrite/setup/jetstream"
	userapi "github.com/matrix-org/dendrite/userapi/api"

	"github.com/matrix-org/dendrite/syncapi/consumers"
	"github.com/matrix-org/dendrite/syncapi/notifier"
	"github.com/matrix-org/dendrite/syncapi/producers"
	"github.com/matrix-org/dendrite/syncapi/routing"
	"github.com/matrix-org/dendrite/syncapi/storage"
	"github.com/matrix-org/dendrite/syncapi/storage/connnect"
	"github.com/matrix-org/dendrite/syncapi/streams"
	"github.com/matrix-org/dendrite/syncapi/sync"
	"github.com/matrix-org/dendrite/syncapi/types"
)

// AddPublicRoutes sets up and registers HTTP handlers for the SyncAPI
// component.
func AddPublicRoutes(
	processContext *process.ProcessContext,
	routers httputil.Routers,
	dendriteCfg *config.Dendrite,
	cm *sqlutil.Connections,
	natsInstance *jetstream.NATSInstance,
	userAPI userapi.SyncUserAPI,
	rsAPI api.SyncRoomserverAPI,
	caches caching.LazyLoadCache,
	enableMetrics bool,
) {
	js, natsClient := natsInstance.Prepare(processContext, &dendriteCfg.Global.JetStream)

	syncDB, connnectQ, err := storage.NewSyncServerDatasource(processContext.Context(), cm, &dendriteCfg.SyncAPI.Database)
	if err != nil {
		logrus.WithError(err).Panicf("failed to connect to sync db")
	}

	go func() {
		var affected int64
		for {
			affected, err = connnectQ.DeleteMultiRoomVisibilityByExpireTS(context.Background(), time.Now().UnixMilli())
			if err != nil {
				logrus.WithError(err).Error("failed to expire multiroom visibility")
			}
			logrus.WithField("rows", affected).Info("expired multiroom visibility")
			time.Sleep(time.Minute)
		}
	}()

	eduCache := caching.NewTypingCache()
	notifier := notifier.NewNotifier()
	streams := streams.NewSyncStreamProviders(syncDB, userAPI, rsAPI, eduCache, caches, notifier, connnectQ)
	notifier.SetCurrentPosition(streams.Latest(context.Background()))
	if err = notifier.Load(context.Background(), syncDB); err != nil {
		logrus.WithError(err).Panicf("failed to load notifier ")
	}

	var fts *fulltext.Search
	if dendriteCfg.SyncAPI.Fulltext.Enabled {
		fts, err = fulltext.New(processContext, dendriteCfg.SyncAPI.Fulltext)
		if err != nil {
			logrus.WithError(err).Panicf("failed to create full text")
		}
	}

	federationPresenceProducer := &producers.FederationAPIPresenceProducer{
		Topic:     dendriteCfg.Global.JetStream.Prefixed(jetstream.OutputPresenceEvent),
		JetStream: js,
	}
	presenceConsumer := consumers.NewPresenceConsumer(
		processContext, &dendriteCfg.SyncAPI, js, natsClient, syncDB,
		notifier, streams.PresenceStreamProvider,
		userAPI,
	)

	requestPool := sync.NewRequestPool(syncDB, &dendriteCfg.SyncAPI, userAPI, rsAPI, streams, notifier, federationPresenceProducer, presenceConsumer, enableMetrics)

	if err = presenceConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start presence consumer")
	}

	keyChangeConsumer := consumers.NewOutputKeyChangeEventConsumer(
		processContext, &dendriteCfg.SyncAPI, dendriteCfg.Global.JetStream.Prefixed(jetstream.OutputKeyChangeEvent),
		js, rsAPI, syncDB, notifier,
		streams.DeviceListStreamProvider,
	)
	if err = keyChangeConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start key change consumer")
	}

	var asProducer *producers.AppserviceEventProducer
	if len(dendriteCfg.AppServiceAPI.Derived.ApplicationServices) > 0 {
		asProducer = &producers.AppserviceEventProducer{
			JetStream: js, Topic: dendriteCfg.Global.JetStream.Prefixed(jetstream.OutputAppserviceEvent),
		}
	}

	scheduler := InitDataRetentionScheduler(processContext, rsAPI, connnectQ)
	roomConsumer := consumers.NewOutputRoomEventConsumer(
		processContext, &dendriteCfg.SyncAPI, js, syncDB, notifier, streams.PDUStreamProvider,
		streams.InviteStreamProvider, rsAPI, fts, asProducer, scheduler,
	)
	if err = roomConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start room server consumer")
	}

	clientConsumer := consumers.NewOutputClientDataConsumer(
		processContext, &dendriteCfg.SyncAPI, js, natsClient, syncDB, notifier,
		streams.AccountDataStreamProvider, fts,
	)
	if err = clientConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start client data consumer")
	}

	notificationConsumer := consumers.NewOutputNotificationDataConsumer(
		processContext, &dendriteCfg.SyncAPI, js, syncDB, notifier, streams.NotificationDataStreamProvider,
	)
	if err = notificationConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start notification data consumer")
	}

	typingConsumer := consumers.NewOutputTypingEventConsumer(
		processContext, &dendriteCfg.SyncAPI, js, eduCache, notifier, streams.TypingStreamProvider,
	)
	if err = typingConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start typing consumer")
	}

	sendToDeviceConsumer := consumers.NewOutputSendToDeviceEventConsumer(
		processContext, &dendriteCfg.SyncAPI, js, syncDB, userAPI, notifier, streams.SendToDeviceStreamProvider,
	)
	if err = sendToDeviceConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start send-to-device consumer")
	}

	receiptConsumer := consumers.NewOutputReceiptEventConsumer(
		processContext, &dendriteCfg.SyncAPI, js, syncDB, notifier, streams.ReceiptStreamProvider,
	)
	if err = receiptConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start receipts consumer")
	}

	multiRoomConsumer := consumers.NewOutputMultiRoomDataConsumer(
		processContext, &dendriteCfg.SyncAPI, js, connnectQ, notifier, streams.MultiRoomStreamProvider,
	)
	if err = multiRoomConsumer.Start(); err != nil {
		logrus.WithError(err).Panicf("failed to start multiroom consumer")
	}
	rateLimits := httputil.NewRateLimits(&dendriteCfg.ClientAPI.RateLimiting)

	routing.Setup(
		routers.Client, requestPool, syncDB, userAPI,
		rsAPI, &dendriteCfg.SyncAPI, caches, fts,
		rateLimits,
	)

	go func() {
		ctx := context.Background()
		for {
			notify, err := syncDB.ExpirePresence(ctx)
			if err != nil {
				logrus.WithError(err).Error("failed to expire presence")
			}
			for i := range notify {
				requestPool.Presence.Store(notify[i].UserID, types.PresenceInternal{
					Presence: types.PresenceOffline,
				})
				notifier.OnNewPresence(types.StreamingToken{
					PresencePosition: notify[i].StreamPos,
				}, notify[i].UserID)

			}
			time.Sleep(types.PresenceExpireInterval)
		}
	}()
}

func InitDataRetentionScheduler(process *process.ProcessContext, rsAPI api.SyncRoomserverAPI, d *connnect.Queries) (scheduler *gocron.Scheduler) {
	scheduler = gocron.NewScheduler(time.UTC)
	drs, err := d.SelectDataRetentions(process.Context())
	if err != nil {
		logrus.WithError(err).Panicf("failed to select data retentions")
	}
	for _, dr := range drs {
		logrus.WithField("space_id", dr.SpaceID).Info("adding data retention job")
		scheduler.At(dr.At).Every(dr.Timeframe).Milliseconds().Do(func(SpaceID string) {
			logrus.WithField("space_id", dr.SpaceID).Info("running data retention job")
			res := &api.PerformDataRetentionResponse{}
			if err := rsAPI.PerformDataRetention(process.Context(), &api.PerformDataRetentionRequest{RoomID: SpaceID}, res); err != nil {
				logrus.WithError(err).Error("failed to delete room data")
			}
		}, dr.SpaceID)
	}
	scheduler.StartAsync()
	return
}
