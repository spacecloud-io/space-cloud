package syncman

import (
	"context"
	"encoding/json"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *Manager) SetPubSubRoutines(nodeID string) error {
	ch, err := s.pubsubClient.Subscribe(context.Background(), generatePubSubTopic(nodeID, pubSubOperationUpgrade))
	if err != nil {
		return err
	}

	go func() {
		for msg := range ch {
			helpers.Logger.LogDebug("pub-sub-upgrade-process", "Received message", nil)
			pubsubMsg := new(model.PubSubMessage)
			if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
				_ = helpers.Logger.LogError("pub-sub-upgrade-process", "Unable to un marshal incoming license request", err, map[string]interface{}{"payload": msg.Payload})
				continue
			}

			s.handlePubSubUpgradeMessage(pubsubMsg)
		}
	}()

	ch1, err := s.pubsubClient.Subscribe(context.Background(), generatePubSubTopic(nodeID, pubSubOperationRenew))
	if err != nil {
		return helpers.Logger.LogError("syncman-new", "Unable to initialize pub sub client required for sync module", err, nil)
	}

	go func() {
		for msg := range ch1 {
			helpers.Logger.LogDebug("pub-sub-upgrade-process", "Received message", nil)
			pubsubMsg := new(model.PubSubMessage)
			if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
				_ = helpers.Logger.LogError("pub-sub-upgrade-process", "Unable to un marshal incoming license request", err, map[string]interface{}{"payload": msg.Payload})
				continue
			}

			s.handlePubSubUpgradeMessage(pubsubMsg)
		}
	}()
	return nil
}

func (s *Manager) handlePubSubUpgradeMessage(msg *model.PubSubMessage) {
	// TODO: Do i require a lock here ?
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	isLeader, err := s.leader.IsLeader(ctx, s.nodeID)
	if err != nil {
		_ = helpers.Logger.LogError("pub-sub-upgrade-process", "Only leader can process an license upgrade request", err, map[string]interface{}{"nodeId": s.nodeID})
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}
	if !isLeader {
		helpers.Logger.LogDebug("pub-sub-upgrade-process", "Not a leader", nil)
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	// Unmarshal the incoming message
	upgradeRequest := new(model.LicenseUpgradeRequest)
	if err := msg.Unmarshal(upgradeRequest); err != nil {
		_ = helpers.Logger.LogError("pub-sub-upgrade-process", "Unable to un marshal pub sub upgrade request", err, map[string]interface{}{"payload": msg.Payload})
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	helpers.Logger.LogDebug("pub-sub-upgrade-process", "Upgrading license", nil)
	if err := s.ConvertToEnterprise(ctx, upgradeRequest); err != nil {
		_ = helpers.Logger.LogError("pub-sub-upgrade-process", "Unable to upgrade license", err, map[string]interface{}{"payload": msg.Payload})
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}
	helpers.Logger.LogDebug("pub-sub-upgrade-process", "Sending positive acknowledgement", nil)
	_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, true)
}

func (s *Manager) handlePubSubRenewMessage(msg *model.PubSubMessage) {
	// TODO: Do i require a lock here ?
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	isLeader, err := s.leader.IsLeader(ctx, s.nodeID)
	if err != nil {
		_ = helpers.Logger.LogError("pub-sub-renew-process", "Only leader can process an license renew request", err, map[string]interface{}{"nodeId": s.nodeID})
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}
	if !isLeader {
		helpers.Logger.LogDebug("pub-sub-upgrade-process", "Not a leader", nil)
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	helpers.Logger.LogDebug("pub-sub-upgrade-process", "Renewing license", nil)
	if err := s.RenewLicense(ctx); err != nil {
		_ = helpers.Logger.LogError("pub-sub-renew-process", "Unable to renew license", err, map[string]interface{}{"payload": msg.Payload})
		_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}
	helpers.Logger.LogDebug("pub-sub-renew-process", "Sending positive acknowledgement", nil)
	_ = s.pubsubClient.SendAck(ctx, msg.ReplyTo, true)
}
