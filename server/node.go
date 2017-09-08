package server

import (
	"github.com/mh-cbon/rendez-vous/client"
	"github.com/mh-cbon/rendez-vous/store"
)

func CentralPoint(client *client.Client, registrations *store.TSRegistrations) MessageQueryHandler {
	return OneOf(
		HandlePing(),
		HandleRegister(registrations),
		HandleUnregister(registrations),
		HandleFind(registrations),
		HandleRequestKnock(registrations, client),
	)
}

func PeerPoint(client *client.Client, knocks *store.TSPendingKnocks) MessageQueryHandler {
	return OneOf(
		HandlePing(),
		HandleDoKnock(client, knocks),
		HandleKnock(knocks),
	)
}
