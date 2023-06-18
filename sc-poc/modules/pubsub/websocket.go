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
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// TODO: channels may or may not have prefix slash

// GetRoutes returns all the apis that are exposed by this app
func (a *App) GetAPIRoutes() apis.APIs {
	return a.apis
}

// getProducerAPI creates a websocket API for sending messages in the channel
func (a *App) getProducerAPI(channelPath string, channel v1alpha1.PubsubChannelSpec) *apis.API {
	// Create a schema validator for incoming messages
	channelSchema := channel.Payload
	schemaLoader := gojsonschema.NewGoLoader(channelSchema)
	schemaValidator, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		a.logger.Error("could not create schema validator for channel", zap.String("channel", channel.Channel), zap.Error(err))
	}

	// Add the plugins provided in the channel options
	var plugins []v1alpha1.HTTPPlugin
	if channel.ProducerOptions != nil {
		plugins = channel.ProducerOptions.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("%s-publisher", channel.Channel),
		Path:    fmt.Sprintf("/v1/pubsub/default%s/producer", channelPath),
		Plugins: plugins,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create the websocket connection
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				a.logger.Error("could not establish websocket connection", zap.String("channel", channel.Channel), zap.Error(err))
				return
			}
			defer conn.Close()

			for {
				// Get the message from the websocket connection
				data, _, err := wsutil.ReadClientData(conn)
				if err != nil {
					a.logger.Error("could not read client data or the connection is closed", zap.String("channel", channel.Channel), zap.Error(err))
					return
				}

				// Unmarshal message
				var message Message
				err = json.Unmarshal(data, &message)
				if err != nil {
					a.logger.Error("could not unmarshal data", zap.String("channel", channel.Channel), zap.Error(err))
					continue
				}

				// Handle events of type message
				if message.Event == MessageEvent {
					var pubMsg PublishMessage
					err = mapstructure.Decode(message.Data, &pubMsg)
					if err != nil {
						a.logger.Error("could not decode data", zap.String("channel", channel.Channel), zap.Error(err))
						continue
					}

					// Create a ID if not ID is not already present
					if pubMsg.ID == "" {
						pubMsg.ID = ksuid.New().String()
					}

					// Validate schema of the message
					documentLoader := gojsonschema.NewGoLoader(pubMsg.Payload)
					result, err := schemaValidator.Validate(documentLoader)
					if err != nil {
						a.logger.Error("could not validate schema for channel", zap.String("channel", channel.Channel), zap.Error(err))
					}

					if !result.Valid() {
						var errMsgs WebsocketErrorMessage
						errMsgs.Message = "Payload of invalid format provided"
						for _, desc := range result.Errors() {
							errMsgs.Errors = append(errMsgs.Errors, fmt.Sprint(desc))
						}

						b, _ := json.Marshal(errMsgs)
						err = wsutil.WriteServerText(conn, b)
						if err != nil {
							a.logger.Error("could not send message to the websocket", zap.Error(err))
						}
						continue
					}

					if err := a.Publish(channel.Channel, pubMsg, PublishOptions{}); err != nil {
						a.logger.Error("could not publish client message", zap.String("channel", channel.Channel), zap.Error(err))
					}
				}
			}
		}),
	}
}

// getConsumerAPI creates a websocket API for receiving messages from the channel
func (a *App) getConsumerAPI(channelPath string, channel v1alpha1.PubsubChannelSpec) *apis.API {
	// Add the plugins provided in the channel options
	var plugins []v1alpha1.HTTPPlugin
	if channel.ConsumerOptions != nil {
		plugins = channel.ConsumerOptions.Plugins
	}

	return &apis.API{
		Name:    fmt.Sprintf("%s-subscriber", channel.Channel),
		Path:    fmt.Sprintf("/v1/pubsub/default%s/consumer", channelPath),
		Plugins: plugins,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			isCtxActive := false

			// Create the websocket connection
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				a.logger.Error("could not establish websocket connection", zap.String("channel", channel.Channel), zap.Error(err))
				return
			}

			subscribeStatus := false
			// Go routine to receive subscribe/unsubscribe events
			go func() {
				for {
					data, _, err := wsutil.ReadClientData(conn)
					if err != nil {
						a.logger.Error("could not read client data or the connection is closed", zap.String("channel", channel.Channel), zap.Error(err))
						if isCtxActive {
							cancel()
							isCtxActive = false
						}
						return
					}

					var message Message
					if err = json.Unmarshal(data, &message); err != nil {
						a.logger.Error("error unmarshaling data", zap.String("channel", channel.Channel), zap.Error(err))
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
							msgChan, err := a.Subscribe(ctx, uuid.NewString(), channel.Channel, SubscribeOptions{})
							if err != nil {
								a.logger.Error("could not subscribe to the channel", zap.String("channel", channel.Channel), zap.Error(err))
								return
							}
							a.logger.Debug("subscribed to the channel", zap.String("channel", channel.Channel))

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
										a.logger.Error("error unmarshaling data", zap.String("channel", channel.Channel), zap.Error(err))
										continue
									}

									pubMsg := PublishMessage{
										ID:       msg.UUID,
										MetaData: msg.Metadata,
										Payload:  data,
									}

									b, err := json.Marshal(pubMsg)
									if err != nil {
										a.logger.Error("error marshaling data", zap.String("channel", channel.Channel), zap.Error(err))
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
						a.logger.Debug("unsubscribed to the channel", zap.String("channel", channel.Channel))
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
