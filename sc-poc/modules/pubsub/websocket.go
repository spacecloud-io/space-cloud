package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
)

// TODO: channels may or may not have prefix slash

// GetRoutes returns all the apis that are exposed by this app
func (a *App) GetAPIRoutes() apis.APIs {
	return a.apis
}

// getPublishAPI creates a websocket API for sending messages in the channel
func (a *App) getPublisherAPI(channelPath, channelName string) *apis.API {
	return &apis.API{
		Name: fmt.Sprintf("%s-publisher", channelName),
		Path: fmt.Sprintf("/v1/pubsub/default%s/publisher", channelPath),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create the websocket connection
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				a.logger.Error("could not establish websocket connection", zap.String("channel", channelName), zap.Error(err))
				return
			}
			defer conn.Close()

			for {
				// Get the message from the websocket connection
				data, _, err := wsutil.ReadClientData(conn)
				if err != nil {
					a.logger.Error("could not read client data or the connection is closed", zap.String("channel", channelName), zap.Error(err))
					return
				}

				var message Message
				err = json.Unmarshal(data, &message)
				if err != nil {
					a.logger.Error("could not unmarshal data", zap.String("channel", channelName), zap.Error(err))
					continue
				}

				if message.Event == MessageEvent {
					var pubMsg PublishMessage
					err = mapstructure.Decode(message.Data, &pubMsg)
					if err != nil {
						a.logger.Error("could not decode data", zap.String("channel", channelName), zap.Error(err))
						continue
					}

					if pubMsg.ID == "" {
						pubMsg.ID = ksuid.New().String()
					}

					err = a.Publish(channelName, pubMsg, PublishOptions{})
					if err != nil {
						a.logger.Error("could not publish client message", zap.String("channel", channelName), zap.Error(err))
					}
				}
			}
		}),
	}
}

// getSubscriberAPI creates a websocket API for receiving messages from the channel
func (a *App) getSubscriberAPI(channelPath, channelName string) *apis.API {
	return &apis.API{
		Name: fmt.Sprintf("%s-subscriber", channelName),
		Path: fmt.Sprintf("/v1/pubsub/default%s/subscriber", channelPath),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			isCtxActive := false

			// Create the websocket connection
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				a.logger.Error("could not establish websocket connection", zap.String("channel", channelName), zap.Error(err))
				return
			}

			subscribeStatus := false
			// Go routine to receive subscribe/unsubscribe events
			go func() {
				for {
					data, _, err := wsutil.ReadClientData(conn)
					if err != nil {
						a.logger.Error("could not read client data or the connection is closed", zap.String("channel", channelName), zap.Error(err))
						if isCtxActive {
							cancel()
							isCtxActive = false
						}
						return
					}

					var message Message
					if err = json.Unmarshal(data, &message); err != nil {
						a.logger.Error("error unmarshaling data", zap.String("channel", channelName), zap.Error(err))
						continue
					}

					a.logger.Debug("event received", zap.String("event", string(message.Event)))
					if message.Event == SubscribeEvent {
						// Check if user is already subscribed
						if subscribeStatus {
							err = wsutil.WriteServerText(conn, []byte("You are already subscribed"))
							if err != nil {
								a.logger.Error("could not send message to the websocket", zap.Error(err))
							}
							continue
						}

						ctx, cancel = context.WithCancel(context.Background())
						subscribeStatus = true
						isCtxActive = true

						go func() {
							// Subscribe to the watermill channel for receiving messages and send
							// it over the websocket connection
							msgChan, err := a.Subscribe(ctx, uuid.NewString(), channelName, SubscribeOptions{})
							if err != nil {
								a.logger.Error("could not subscribe to the channel", zap.String("channel", channelName), zap.Error(err))
								return
							}
							a.logger.Debug("subscribed to the channel", zap.String("channel", channelName))

							for {
								select {
								case <-ctx.Done():
									// Unsubscribe event received. Close the go routine.
									return
								case msg := <-msgChan:
									// Write the message to the websocket connection
									var data interface{}
									err := json.Unmarshal(msg.Payload, &data)
									if err != nil {
										a.logger.Error("error unmarshaling data", zap.String("channel", channelName), zap.Error(err))
										continue
									}

									pubMsg := PublishMessage{
										ID:       msg.UUID,
										MetaData: msg.Metadata,
										Payload:  data,
									}

									b, err := json.Marshal(pubMsg)
									if err != nil {
										a.logger.Error("error marshaling data", zap.String("channel", channelName), zap.Error(err))
										continue
									}

									err = wsutil.WriteServerText(conn, b)
									if err != nil {
										a.logger.Error("could not send message to the websocket", zap.Error(err))
									}
									msg.Ack()
								}
							}

						}()
					}
					if message.Event == UnsubscribeEvent {
						a.logger.Debug("unsubscribed to the channel", zap.String("channel", channelName))
						if isCtxActive {
							cancel()
							isCtxActive = false
						}
						subscribeStatus = false
					}
				}
			}()
		}),
	}
}
