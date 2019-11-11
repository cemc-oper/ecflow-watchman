package ecflow_watchman

import "github.com/go-redis/redis"

type Publisher interface {
	Create()
	Publish(message []byte)
	Close()
}

type RedisPublisher struct {
	Client      *redis.Client
	Pubsub      *redis.PubSub
	ChannelName string
	Address     string
	Password    string
	Database    int
}

func (p *RedisPublisher) Create() {
	p.Client = redis.NewClient(&redis.Options{
		Addr:     p.Address,
		Password: p.Password,
		DB:       p.Database,
	})

	p.Pubsub = p.Client.Subscribe(p.ChannelName)
	_, err := p.Pubsub.Receive()
	if err != nil {
		panic(err)
	}
}

func (p *RedisPublisher) Close() {
	if p.Pubsub != nil {
		defer p.Pubsub.Close()
	}
	if p.Client != nil {
		defer p.Client.Close()
	}
}

func (p *RedisPublisher) Publish(message []byte) error {
	redisCmd := p.Client.Publish(p.ChannelName, message)
	return redisCmd.Err()
}
