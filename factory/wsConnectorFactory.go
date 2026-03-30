package factory

import (
	"github.com/TerraDharitri/drt-go-chain-communication/websocket/data"
	factoryHost "github.com/TerraDharitri/drt-go-chain-communication/websocket/factory"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	"github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-ws-connector-template/config"
	"github.com/TerraDharitri/drt-go-chain-ws-connector-template/process"
)

var log = logger.GetOrCreate("ws-connector")

// CreateWSConnector will create a ws connector able to receive and process incoming data
// from a dharitri node
func CreateWSConnector(cfg config.WebSocketConfig) (process.WSConnector, error) {
	marshaller, err := factory.NewMarshalizer(cfg.MarshallerType)
	if err != nil {
		return nil, err
	}

	dataProcessor, err := process.NewLogDataProcessor(marshaller, log)
	if err != nil {
		return nil, err
	}

	wsHost, err := createWsHost(marshaller, cfg)
	if err != nil {
		return nil, err
	}

	err = wsHost.SetPayloadHandler(dataProcessor)
	if err != nil {
		return nil, err
	}

	return wsHost, nil
}

func createWsHost(wsMarshaller marshal.Marshalizer, cfg config.WebSocketConfig) (factoryHost.FullDuplexHost, error) {
	return factoryHost.CreateWebSocketHost(factoryHost.ArgsWebSocketHost{
		WebSocketConfig: data.WebSocketConfig{
			URL:                        cfg.Url,
			WithAcknowledge:            cfg.WithAcknowledge,
			Mode:                       cfg.Mode,
			RetryDurationInSec:         int(cfg.RetryDuration),
			BlockingAckOnError:         cfg.BlockingAckOnError,
			DropMessagesIfNoConnection: false,
		},
		Marshaller: wsMarshaller,
		Log:        log,
	})
}
