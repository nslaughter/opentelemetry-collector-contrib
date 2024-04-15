package azureeventhubreceiver

import (
	"context"
	"errors"
	"fmt"
	"time"

	// Updating imports to use the new SDK
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// Assuming there's a struct managing the processor setup
// type EventHubProcessor struct {
//     Processor *azeventhubs.Processor
// }

// Updated initialization function using the new SDK components
func NewEventHubProcessor(ehConn, ehName, storageConn, storageCnt string) (*EventHubProcessor, error) {
	checkpointingProcessor, err := newCheckpointingProcessor(ehConn, ehName, storageConn, storageCnt)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpointing processor: %w", err)
	}

	// Start processing events
	return &EventHubProcessor{
		Processor: checkpointingProcessor,
	}, nil
}

// Assume there's a function to start processing events
func (e *EventHubProcessor) StartProcessing(ctx context.Context) error {
	// Start the processor
	if err := e.Processor.Run(ctx); err != nil {
		return fmt.Errorf("error running processor: %w", err)
	}
	return nil
}

// Assuming there's a struct managing the processor setup
type EventHubProcessor struct {
	Processor *azeventhubs.Processor
}

// These are config values the processor factory can use to create processors:
//
//	(a) EventHubConnectionString
//	(b) EventHubName
//	(c) StorageConnectionString
//	(d) StorageContainerName
//
// You always need the EventHub variable values.
// And you need all 4 of these to checkpoint.
//
// I think the config values should be managed in the factory struct.
/*
func (pf *processorFactory) CreateProcessor() (*azeventhubs.Processor, error) {
	// Create the consumer client
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(pf.EventHubConnectionString, pf.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
	if err != nil {
		return nil, err
	}

	// Create the blob container client for the checkpoint store
	blobContainerClient, err := container.NewClientFromConnectionString(pf.StorageConnectionString, pf.StorageContainerName, nil)
	if err != nil {
		return nil, err
	}

	// Create the checkpoint store using the blob container client
	checkpointStore, err := azeventhubs.NewBlobCheckpointStore(blobContainerClient, nil)
	// checkpointStore, err := azeventhubs.NewBlobCheckpointStore(blobContainerClient, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// Create the processor with checkpointing
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		return nil, err
	}

	return processor, nil
}
*/

func newCheckpointingProcessor(eventHubConnectionString, eventHubName, storageConnectionString, storageContainerName string) (*azeventhubs.Processor, error) {
	blobContainerClient, err := container.NewClientFromConnectionString(storageConnectionString, storageContainerName, nil)
	if err != nil {
		return nil, err
	}
	checkpointStore, err := checkpoints.NewBlobStore(blobContainerClient, nil)
	if err != nil {
		return nil, err
	}

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(eventHubConnectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	if err != nil {
		return nil, err
	}

	return azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
}

func dispatchPartitionClients(processor *azeventhubs.Processor) {
	for {
		processorPartitionClient := processor.NextPartitionClient(context.TODO())
		if processorPartitionClient == nil {
			break
		}

		go func() {
			if err := processEventsForPartition(processorPartitionClient); err != nil {
				panic(err)
			}
		}()
	}
}

func processEventsForPartition(partitionClient *azeventhubs.ProcessorPartitionClient) error {
	defer shutdownPartitionResources(partitionClient)
	if err := initializePartitionResources(partitionClient.PartitionID()); err != nil {
		return err
	}

	for {
		receiveCtx, cancelReceive := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		cancelReceive()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if len(events) == 0 {
			continue
		}

		fmt.Printf("Received %d event(s)\n", len(events))
		for _, event := range events {
			fmt.Printf("Event received with body %v\n", event.Body)
		}
		if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
			return err
		}
	}
}
