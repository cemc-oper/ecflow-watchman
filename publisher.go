package ecflow_watchman

import "github.com/go-redis/redis"

type Publisher interface {
	Create()
	Publish(key string, message []byte)
	Close()
}

type RedisPublisher struct {
	Client   *redis.Client
	Pubsubs  []*redis.PubSub
	Address  string
	Password string
	Database int
}

func (p *RedisPublisher) Create() {
	p.Client = redis.NewClient(&redis.Options{
		Addr:     p.Address,
		Password: p.Password,
		DB:       p.Database,
	})
}

func (p *RedisPublisher) CreatePubsub(channelName string) *redis.PubSub {
	pubsub := p.Client.Subscribe(channelName)
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}
	p.Pubsubs = append(p.Pubsubs, pubsub)
	return pubsub
}

func (p *RedisPublisher) Close() {
	for _, pubsub := range p.Pubsubs {
		if pubsub != nil {
			pubsub.Close()
		}
	}
	p.Pubsubs = nil
	if p.Client != nil {
		defer p.Client.Close()
	}
}

func (p *RedisPublisher) Publish(key string, message []byte) error {
	redisCmd := p.Client.Publish(key, message)
	return redisCmd.Err()
}
