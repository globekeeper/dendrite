// Copyright 2022 The Matrix.org Foundation C.I.C.
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

package perform

import (
	"github.com/matrix-org/dendrite/roomserver/internal/input"
	"github.com/matrix-org/dendrite/roomserver/internal/query"
	"github.com/matrix-org/dendrite/roomserver/storage"
	"github.com/matrix-org/dendrite/setup/config"
)

type DataRetention struct {
	DB      storage.Database
	Cfg     *config.RoomServer
	Queryer *query.Queryer
	Inputer *input.Inputer
	Leaver  *Leaver
}

// PerformRoomDataRetention will data retain a given room.
// func (r *DataRetention) PerformRoomDataRetention(
// 	ctx context.Context,
// 	req *api.PerformDataRetentionRequest,
// 	res *api.PerformDataRetentionResponse,
// ) *util.JSONResponse {
// 	roomInfo, err := r.DB.RoomInfo(ctx, req.RoomID)
// 	if err != nil {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON(fmt.Sprintf("r.DB.RoomInfo: %w", err)),
// 		}
// 	}
// 	if roomInfo == nil || roomInfo.IsStub() {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON(fmt.Sprintf("Room %s not found", req.RoomID)),
// 		}
// 	}

// 	memberNIDs, err := r.DB.GetMembershipEventNIDsForRoom(ctx, roomInfo.RoomNID, true, true)
// 	if err != nil {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON(fmt.Sprintf("r.DB.GetMembershipEventNIDsForRoom: %w", err)),
// 		}
// 	}

// 	memberEvents, err := r.DB.Events(ctx, roomInfo.RoomVersion, memberNIDs)
// 	if err != nil {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON(fmt.Sprintf("r.DB.Events: %w", err)),
// 		}
// 	}

// 	inputEvents := make([]api.InputRoomEvent, 0, len(memberEvents))
// 	res.Affected = make([]string, 0, len(memberEvents))
// 	latestReq := &api.QueryLatestEventsAndStateRequest{
// 		RoomID: req.RoomID,
// 	}
// 	latestRes := &api.QueryLatestEventsAndStateResponse{}
// 	if err = r.Queryer.QueryLatestEventsAndState(ctx, latestReq, latestRes); err != nil {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON(fmt.Sprintf("r.Queryer.QueryLatestEventsAndState: %w", err)),
// 		}
// 	}

// 	prevEvents := latestRes.LatestEvents
// 	for _, memberEvent := range memberEvents {
// 		if memberEvent.StateKey() == nil {
// 			continue
// 		}

// 		var memberContent gomatrixserverlib.MemberContent
// 		if err = json.Unmarshal(memberEvent.Content(), &memberContent); err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON(fmt.Sprintf("json.Unmarshal: %w", err)),
// 			}
// 		}
// 		memberContent.Membership = spec.Leave

// 		stateKey := *memberEvent.StateKey()
// 		fledglingEvent := &gomatrixserverlib.EventBuilder{
// 			RoomID:     req.RoomID,
// 			Type:       spec.MRoomMember,
// 			StateKey:   &stateKey,
// 			SenderID:   stateKey,
// 			PrevEvents: prevEvents,
// 		}

// 		_, senderDomain, err := gomatrixserverlib.SplitID('@', fledglingEvent.SenderID)
// 		if err != nil {
// 			continue
// 		}

// 		if fledglingEvent.Content, err = json.Marshal(memberContent); err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON(fmt.Sprintf("json.Marshal: %s", err)),
// 			}
// 		}

// 		eventsNeeded, err := gomatrixserverlib.StateNeededForProtoEvent(fledglingEvent)
// 		if err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON("spec.StateNeededForEventBuilder: %s", err),
// 			}
// 		}

// 		identity, err := r.Cfg.Matrix.SigningIdentityFor(senderDomain)
// 		if err != nil {
// 			continue
// 		}

// 		event, err := eventutil.BuildEvent(ctx, fledglingEvent, r.Cfg.Matrix, identity, time.Now(), &eventsNeeded, latestRes)
// 		if err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON("eventutil.BuildEvent: %s", err),
// 			}
// 		}

// 		inputEvents = append(inputEvents, api.InputRoomEvent{
// 			Kind:         api.KindNew,
// 			Event:        event,
// 			Origin:       senderDomain,
// 			SendAsServer: string(senderDomain),
// 		})
// 		res.Affected = append(res.Affected, stateKey)
// 		prevEvents = []spec.EventReference{
// 			event.EventReference(),
// 		}
// 	}

// 	inputReq := &api.InputRoomEventsRequest{
// 		InputRoomEvents: inputEvents,
// 		Asynchronous:    true,
// 	}
// 	inputRes := &api.InputRoomEventsResponse{}
// 	return r.Inputer.InputRoomEvents(ctx, inputReq, inputRes)
// }

// func (r *Admin) PerformAdminEvacuateUser(
// 	ctx context.Context,
// 	req *api.PerformAdminEvacuateUserRequest,
// 	res *api.PerformAdminEvacuateUserResponse,
// ) *util.JSONResponse {
// 	_, domain, err := spec.SplitID('@', req.UserID)
// 	if err != nil {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON("Malformed user ID: %s", err),
// 		}
// 	}
// 	if !r.Cfg.Matrix.IsLocalServerName(domain) {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: "Can only evacuate local users using this endpoint",
// 		}
// 	}

// 	roomIDs, err := r.DB.GetRoomsByMembership(ctx, req.UserID, spec.Join)
// 	if err != nil && err != sql.ErrNoRows {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON("r.DB.GetRoomsByMembership: %s", err),
// 		}
// 	}

// 	inviteRoomIDs, err := r.DB.GetRoomsByMembership(ctx, req.UserID, spec.Invite)
// 	if err != nil && err != sql.ErrNoRows {
// 		return &util.JSONResponse{
// 			Code: http.StatusBadRequest,
// 			JSON: spec.BadJSON("r.DB.GetRoomsByMembership: %s", err),
// 		}
// 	}

// 	for _, roomID := range append(roomIDs, inviteRoomIDs...) {
// 		leaveReq := &api.PerformLeaveRequest{
// 			RoomID: roomID,
// 			Leaver: spec.UserID(req.UserID),
// 		}
// 		leaveRes := &api.PerformLeaveResponse{}
// 		outputEvents, err := r.Leaver.PerformLeave(ctx, leaveReq, leaveRes)
// 		if err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON("r.Leaver.PerformLeave: %s", err),
// 			}
// 			return nil
// 		}
// 		res.Affected = append(res.Affected, roomID)
// 		if len(outputEvents) == 0 {
// 			continue
// 		}
// 		if err := r.Inputer.OutputProducer.ProduceRoomEvents(roomID, outputEvents); err != nil {
// 			return &util.JSONResponse{
// 				Code: http.StatusBadRequest,
// 				JSON: spec.BadJSON("r.Inputer.WriteOutputEvents: %s", err),
// 			}
// 			return nil
// 		}
// 	}
// 	return nil
// }

// // func (r *Admin) PerformAdminPurgeRoom(
// // 	ctx context.Context,
// // 	req *api.PerformAdminPurgeRoomRequest,
// // 	res *api.PerformAdminPurgeRoomResponse,
// // ) *util.JSONResponse {
// // 	// Validate we actually got a room ID and nothing else
// // 	if _, _, err := spec.SplitID('!', req.RoomID); err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("Malformed room ID: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	logrus.WithField("room_id", req.RoomID).Warn("Purging room from roomserver")
// // 	if err := r.DB.PurgeRoom(ctx, req.RoomID); err != nil {
// // 		logrus.WithField("room_id", req.RoomID).WithError(err).Warn("Failed to purge room from roomserver")
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: err.Error(),
// // 		}
// // 		return nil
// // 	}

// // 	logrus.WithField("room_id", req.RoomID).Warn("Room purged from roomserver")

// // 	return
// // }

// // func (r *Admin) PerformAdminDownloadState(
// // 	ctx context.Context,
// // 	req *api.PerformAdminDownloadStateRequest,
// // 	res *api.PerformAdminDownloadStateResponse,
// // ) *util.JSONResponse {
// // 	_, senderDomain, err := r.Cfg.Matrix.SplitLocalID('@', req.UserID)
// // 	if err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("r.Cfg.Matrix.SplitLocalID: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	roomInfo, err := r.DB.RoomInfo(ctx, req.RoomID)
// // 	if err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("r.DB.RoomInfo: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	if roomInfo == nil || roomInfo.IsStub() {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("room %q not found", req.RoomID),
// // 		}
// // 		return nil
// // 	}

// // 	fwdExtremities, _, depth, err := r.DB.LatestEventIDs(ctx, roomInfo.RoomNID)
// // 	if err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("r.DB.LatestEventIDs: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	authEventMap := map[string]*spec.Event{}
// // 	stateEventMap := map[string]*spec.Event{}

// // 	for _, fwdExtremity := range fwdExtremities {
// // 		var state spec.StateResponse
// // 		state, err = r.Inputer.FSAPI.LookupState(ctx, r.Inputer.ServerName, req.ServerName, req.RoomID, fwdExtremity.EventID, roomInfo.RoomVersion)
// // 		if err != nil {
// // 			return &util.JSONResponse{
// // 				Code: http.StatusBadRequest,
// // 				JSON: spec.BadJSON("r.Inputer.FSAPI.LookupState (%q): %s", fwdExtremity.EventID, err),
// // 			}
// // 			return nil
// // 		}
// // 		for _, authEvent := range state.GetAuthEvents().UntrustedEvents(roomInfo.RoomVersion) {
// // 			if err = authEvent.VerifyEventSignatures(ctx, r.Inputer.KeyRing); err != nil {
// // 				continue
// // 			}
// // 			authEventMap[authEvent.EventID()] = authEvent
// // 		}
// // 		for _, stateEvent := range state.GetStateEvents().UntrustedEvents(roomInfo.RoomVersion) {
// // 			if err = stateEvent.VerifyEventSignatures(ctx, r.Inputer.KeyRing); err != nil {
// // 				continue
// // 			}
// // 			stateEventMap[stateEvent.EventID()] = stateEvent
// // 		}
// // 	}

// // 	authEvents := make([]*spec.HeaderedEvent, 0, len(authEventMap))
// // 	stateEvents := make([]*spec.HeaderedEvent, 0, len(stateEventMap))
// // 	stateIDs := make([]string, 0, len(stateEventMap))

// // 	for _, authEvent := range authEventMap {
// // 		authEvents = append(authEvents, authEvent.Headered(roomInfo.RoomVersion))
// // 	}
// // 	for _, stateEvent := range stateEventMap {
// // 		stateEvents = append(stateEvents, stateEvent.Headered(roomInfo.RoomVersion))
// // 		stateIDs = append(stateIDs, stateEvent.EventID())
// // 	}

// // 	builder := &spec.EventBuilder{
// // 		Type:    "org.matrix.dendrite.state_download",
// // 		Sender:  req.UserID,
// // 		RoomID:  req.RoomID,
// // 		Content: spec.RawJSON("{}"),
// // 	}

// // 	eventsNeeded, err := spec.StateNeededForEventBuilder(builder)
// // 	if err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("spec.StateNeededForEventBuilder: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	queryRes := &api.QueryLatestEventsAndStateResponse{
// // 		RoomExists:   true,
// // 		RoomVersion:  roomInfo.RoomVersion,
// // 		LatestEvents: fwdExtremities,
// // 		StateEvents:  stateEvents,
// // 		Depth:        depth,
// // 	}

// // 	identity, err := r.Cfg.Matrix.SigningIdentityFor(senderDomain)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	ev, err := eventutil.BuildEvent(ctx, builder, r.Cfg.Matrix, identity, time.Now(), &eventsNeeded, queryRes)
// // 	if err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("eventutil.BuildEvent: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	inputReq := &api.InputRoomEventsRequest{
// // 		Asynchronous: false,
// // 	}
// // 	inputRes := &api.InputRoomEventsResponse{}

// // 	for _, authEvent := range append(authEvents, stateEvents...) {
// // 		inputReq.InputRoomEvents = append(inputReq.InputRoomEvents, api.InputRoomEvent{
// // 			Kind:  api.KindOutlier,
// // 			Event: authEvent,
// // 		})
// // 	}

// // 	inputReq.InputRoomEvents = append(inputReq.InputRoomEvents, api.InputRoomEvent{
// // 		Kind:          api.KindNew,
// // 		Event:         ev,
// // 		Origin:        r.Cfg.Matrix.ServerName,
// // 		HasState:      true,
// // 		StateEventIDs: stateIDs,
// // 		SendAsServer:  string(r.Cfg.Matrix.ServerName),
// // 	})

// // 	if err := r.Inputer.InputRoomEvents(ctx, inputReq, inputRes); err != nil {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: spec.BadJSON("r.Inputer.InputRoomEvents: %s", err),
// // 		}
// // 		return nil
// // 	}

// // 	if inputRes.ErrMsg != "" {
// // 		return &util.JSONResponse{
// // 			Code: http.StatusBadRequest,
// // 			JSON: inputRes.ErrMsg,
// // 		}
// // 	}

// // 	return nil
// // }
