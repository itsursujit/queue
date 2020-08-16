package workers

import (
	"crypto/tls"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"queue/workers/storage"
)

// Options contains the set of configuration options for a manager and/or producer
type Options struct {
	ProcessID      string
	Namespace      string
	PollInterval   int
	Database       int
	Password       string
	PoolSize       int
	PersistentAddr string
	PersistentDB   string
	// Provide one of ServerAddr or (SentinelAddrs + RedisMasterName)
	ServerAddr      string
	SentinelAddrs   string
	RedisMasterName string
	RedisTLSConfig  *tls.Config

	// Optional display name used when displaying manager stats
	ManagerDisplayName string
	store              storage.Store
	client             *redis.Client
	persistentClient   *MongoInstance
}

func processOptions(options Options) (Options, error) {
	options, err := validateGeneralOptions(options)
	if err != nil {
		return Options{}, err
	}

	//redis options
	if options.PoolSize == 0 {
		options.PoolSize = 1
	}
	redisIdleTimeout := 240 * time.Second
	if options.PersistentDB == "" {
		options.PersistentDB = "test_db"
	}
	if options.PersistentAddr == "" {
		options.PersistentAddr = "mongodb://localhost:27017/" + options.PersistentDB
	}
	// Connect to the database
	if err := Connect(options.PersistentAddr, options.PersistentDB); err != nil {
		log.Fatal(err)
	}

	options.persistentClient = MG
	if options.ServerAddr != "" {
		options.client = redis.NewClient(&redis.Options{
			IdleTimeout: redisIdleTimeout,
			Password:    options.Password,
			DB:          options.Database,
			PoolSize:    options.PoolSize,
			Addr:        options.ServerAddr,
			TLSConfig:   options.RedisTLSConfig,
		})
	} else if options.SentinelAddrs != "" {
		if options.RedisMasterName == "" {
			return Options{}, errors.New("Sentinel configuration requires a master name")
		}

		options.client = redis.NewFailoverClient(&redis.FailoverOptions{
			IdleTimeout:   redisIdleTimeout,
			Password:      options.Password,
			DB:            options.Database,
			PoolSize:      options.PoolSize,
			SentinelAddrs: strings.Split(options.SentinelAddrs, ","),
			MasterName:    options.RedisMasterName,
			TLSConfig:     options.RedisTLSConfig,
		})
	} else {
		return Options{}, errors.New("Options requires either the Server or Sentinels option")
	}

	redisStore := storage.NewRedisStore(options.Namespace, options.client)
	options.store = redisStore

	return options, nil
}

func processOptionsWithRedisClient(options Options, client *redis.Client) (Options, error) {
	options, err := validateGeneralOptions(options)
	if err != nil {
		return Options{}, err
	}

	if client == nil {
		return Options{}, errors.New("Redis client is nil; Redis client is not configured")
	}

	options.client = client

	redisStore := storage.NewRedisStore(options.Namespace, options.client)
	options.store = redisStore

	return options, nil
}

func validateGeneralOptions(options Options) (Options, error) {
	if options.ProcessID == "" {
		return Options{}, errors.New("Options requires a ProcessID, which uniquely identifies this instance")
	}

	if options.Namespace != "" {
		options.Namespace += ":"
	}

	if options.PollInterval <= 0 {
		options.PollInterval = 15
	}

	return options, nil
}
