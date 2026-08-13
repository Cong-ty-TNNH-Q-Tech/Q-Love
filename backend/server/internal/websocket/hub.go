package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/pkg/logger"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const RedisChatStream = "qlove:chat_stream"

type Hub struct {
	// Registered clients map[UserID]map[*Client]bool to support multi-device
	clients map[uuid.UUID]map[*Client]bool
	
	Register   chan *Client
	Unregister chan *Client
	
	redisClient *redis.Client
	mu          sync.RWMutex
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[uuid.UUID]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		redisClient: redisClient,
	}
}

func (h *Hub) Run(ctx context.Context) {
	// Start consuming from Redis Stream
	go h.consumeRedisStream(ctx)

	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			logger.Log.Info("Client registered", zap.String("user_id", client.UserID.String()))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID][client]; ok {
				delete(h.clients[client.UserID], client)
				close(client.Send)
				if len(h.clients[client.UserID]) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			logger.Log.Info("Client unregistered", zap.String("user_id", client.UserID.String()))
			
		case <-ctx.Done():
			return
		}
	}
}

// PublishToRedis takes a message, saves it via HTTP handler, and pushes to Redis Stream
func (h *Hub) PublishToRedis(ctx context.Context, msg *models.ChatMessage) error {
	if h.redisClient == nil {
		return nil
	}
	
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// XADD to Redis stream
	err = h.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: RedisChatStream,
		Values: map[string]interface{}{
			"target_id": msg.ReceiverID.String(),
			"payload":   string(msgBytes),
		},
	}).Err()

	return err
}

func (h *Hub) consumeRedisStream(ctx context.Context) {
	if h.redisClient == nil {
		return
	}

	lastID := "$" // Start listening for new messages only
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Block and read from stream
			streams, err := h.redisClient.XRead(ctx, &redis.XReadArgs{
				Streams: []string{RedisChatStream, lastID},
				Count:   10,
				Block:   0, // Block indefinitely
			}).Result()

			if err != nil {
				if err != redis.Nil && err != context.Canceled {
					logger.Log.Error("Redis XREAD error", zap.Error(err))
					time.Sleep(2 * time.Second) // wait before retry
				}
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					lastID = message.ID
					
					targetIDStr, ok1 := message.Values["target_id"].(string)
					payloadStr, ok2 := message.Values["payload"].(string)
					
					if !ok1 || !ok2 {
						continue
					}
					
					targetID, err := uuid.Parse(targetIDStr)
					if err != nil {
						continue
					}
					
					// If the target user is connected to THIS node, send it
					h.mu.RLock()
					if userConns, exists := h.clients[targetID]; exists {
						for client := range userConns {
							select {
							case client.Send <- []byte(payloadStr):
							default:
								// Cannot send (buffer full), close it
								h.mu.RUnlock()
								h.Unregister <- client
								h.mu.RLock()
							}
						}
					}
					h.mu.RUnlock()
				}
			}
		}
	}
}
