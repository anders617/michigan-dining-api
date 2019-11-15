package dynamoclient

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/golang/glog"
)

func (d *DynamoClient) StreamHearts() (chan pb.HeartCount, chan struct{}) {
	heartCountChan := make(chan pb.HeartCount)
	recordChan, doneChan := d.streamRecords(HeartsTableName)
	go func(heartCountChan chan pb.HeartCount, recordChan chan dynamodbstreams.Record) {
		for record := range recordChan {
			heartCount := pb.HeartCount{}
			err := dynamodbattribute.UnmarshalMap(record.Dynamodb.NewImage, &heartCount)
			if err != nil {
				glog.Warningf("Could not umarshal heart count: %s", err)
				continue
			}
			heartCountChan <- heartCount
		}
	}(heartCountChan, recordChan)
	return heartCountChan, doneChan
}

func (d *DynamoClient) streamRecords(table string) (chan dynamodbstreams.Record, chan struct{}) {
	recordChan := make(chan dynamodbstreams.Record)
	doneChan := make(chan struct{})
	go func(recordChan chan dynamodbstreams.Record, doneChan chan struct{}, table string) {
		processingShards := map[string]struct{}{}
		ticker := time.NewTicker(5 * time.Minute)
		doneChans := []chan struct{}{}
		done := false
		mu := sync.Mutex{}
		shardPollCount := 0
		shardPollCountMu := sync.Mutex{}
		// Send done message to all shard specific doneChans if we get a done message
		go func(doneChan chan struct{}, recordChan chan dynamodbstreams.Record) {
			<-doneChan
			mu.Lock()
			for _, done := range doneChans {
				done <- struct{}{}
			}
			close(recordChan)
			done = true
			mu.Unlock()
		}(doneChan, recordChan)
		for {
			mu.Lock()
			if done {
				break
			}
			streamArn, err := d.getTableStreamArn(table)
			if err != nil {
				glog.Fatalf("Failed to get stream arn")
			}
			glog.Infof("Stream Arn: %s", *streamArn)
			shards, err := d.getStreamShards(*streamArn)
			if err != nil {
				glog.Fatalf("Failed to get stream shards")
			}
			recordChans := []chan dynamodbstreams.Record{}
			// For each shard, start polling for records
			for _, shard := range *shards {
				_, alreadyProcessing := processingShards[*shard.ShardId]
				if alreadyProcessing {
					continue
				}
				glog.Infof("Processing new shard: %s", *shard.ShardId)
				processingShards[*shard.ShardId] = struct{}{}
				shardIt, err := d.getShardIterator(*streamArn, shard)
				if err != nil {
					glog.Fatalf("Failed to get shard iterator")
				}
				shardPollCountMu.Lock()
				shardPollCount++
				shardPollCountMu.Unlock()
				records, done := d.pollShardForRecords(*shardIt)
				recordChans = append(recordChans, records)
				doneChans = append(doneChans, done)
			}
			// Aggregate each shard specific recordChan into one channel
			for _, records := range recordChans {
				go func(records chan dynamodbstreams.Record) {
					for record := range records {
						recordChan <- record
					}
					shardPollCountMu.Lock()
					shardPollCount--
					shardPollCountMu.Unlock()
					glog.Infof("Closed records")
				}(records)
			}
			shardPollCountMu.Lock()
			glog.Infof("%d shard pollers active for table %s", shardPollCount, table)
			shardPollCountMu.Unlock()
			mu.Unlock()
			// Wait for next tick to check for new shards
			<-ticker.C
		}
	}(recordChan, doneChan, table)
	return recordChan, doneChan
}

func (d *DynamoClient) pollShardForRecords(shardIterator string) (chan dynamodbstreams.Record, chan struct{}) {
	recordChan := make(chan dynamodbstreams.Record)
	doneChan := make(chan struct{})
	go func(recordChan chan dynamodbstreams.Record, doneChan chan struct{}, shardIterator string) {
		shardIt := &shardIterator
		defer close(recordChan)
		for shardIt != nil {
			select {
			case <-doneChan:
				return
			default:
				var records *[]dynamodbstreams.Record
				var err error
				shardIt, records, err = d.getRecords(*shardIt)
				if err != nil {
					glog.Warningf("Failed to get records")
					close(recordChan)
				}
				for _, record := range *records {
					recordChan <- record
				}
				if len(*records) == 0 {
					time.Sleep(time.Second)
				}
			}
		}
	}(recordChan, doneChan, shardIterator)
	return recordChan, doneChan
}

func (d *DynamoClient) getRecords(shardIterator string) (*string, *[]dynamodbstreams.Record, error) {
	params := dynamodbstreams.GetRecordsInput{
		ShardIterator: &shardIterator,
	}
	req := d.streamClient.GetRecordsRequest(&params)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return resp.NextShardIterator, &resp.Records, nil
}

func (d *DynamoClient) getShardIterator(arn string, shard dynamodbstreams.Shard) (*string, error) {
	params := dynamodbstreams.GetShardIteratorInput{
		StreamArn:         &arn,
		ShardId:           shard.ShardId,
		ShardIteratorType: dynamodbstreams.ShardIteratorTypeLatest,
	}

	req := d.streamClient.GetShardIteratorRequest(&params)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.ShardIterator, nil
}

func (d *DynamoClient) getStreamShards(arn string) (*[]dynamodbstreams.Shard, error) {
	params := dynamodbstreams.DescribeStreamInput{
		StreamArn: &arn,
	}
	req := d.streamClient.DescribeStreamRequest(&params)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return &resp.StreamDescription.Shards, nil
}

func (d *DynamoClient) getTableStreamArn(table string) (*string, error) {
	// Example sending a request using the ListStreamsRequest method.
	params := dynamodbstreams.ListStreamsInput{
		TableName: &table,
	}
	req := d.streamClient.ListStreamsRequest(&params)

	resp, err := req.Send(context.Background())
	if err != nil { // resp is now filled
		return nil, err
	}
	if len(resp.Streams) == 0 {
		return nil, errors.New("No streams found")
	}
	return resp.Streams[0].StreamArn, nil
}
