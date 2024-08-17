package temporal

import (
	"demo-temporal-project/configs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
)

func NewTemporalClient(logger log.Logger, conf *configs.Temporal, env string) (*client.Client, func(), error) {
	namespace := "LOCAL"
	c, err := client.Dial(client.Options{
		HostPort:  conf.BaseUrl,
		Namespace: namespace,
	})

	if err != nil {
		return &c, nil, err
	}
	return &c, func() {
		c.Close()
		logger.Log(log.LevelDebug, "Closing temporal client")
	}, nil
}

var ProviderSet = wire.NewSet(NewTemporalClient)
