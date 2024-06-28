package router

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/rudderlabs/rudder-go-kit/config"

	"github.com/rudderlabs/rudder-server/jobsdb"
	"github.com/rudderlabs/rudder-server/router/types"
)

type workerJobStatus struct {
	userID string
	worker *worker
	job    *jobsdb.JobT
	status *jobsdb.JobStatusT
}

type HandleDestOAuthRespParams struct {
	ctx            context.Context
	destinationJob types.DestinationJobT
	workerID       int
	trRespStCd     int
	trRespBody     string
	secret         json.RawMessage
	contentType    string
}

type Diagnostic struct {
	diagnosisTicker    *time.Ticker
	requestsMetricLock sync.RWMutex
	requestsMetric     []requestMetric
	failureMetricLock  sync.RWMutex
	failuresMetric     map[string]map[string]int
}

type requestMetric struct {
	RequestRetries       int
	RequestAborted       int
	RequestSuccess       int
	RequestCompletedTime time.Duration
}

type JobResponse struct {
	jobID                  int64
	destinationJob         *types.DestinationJobT
	destinationJobMetadata *types.JobMetadataT
	respStatusCode         int
	respBody               string
	errorAt                string
	status                 *jobsdb.JobStatusT
}

func (j JobResponse) GetTransformerApiLogs() []types.ApiLog {
	apiLogs := make([]types.ApiLog, 0)
	for _, metadata := range j.destinationJob.JobMetadataArray {
		apiLogs = append(apiLogs, metadata.ApiLogs...)
	}
	return apiLogs
}

func (j JobResponse) GetRouterAPiLogs() ([]types.ApiLog, error) {
	apiLogs := make([]types.ApiLog, 0)
	reqRaw := j.destinationJob.Message
	var reqParsed []map[string]interface{}
	if err := json.Unmarshal(reqRaw, &reqParsed); err != nil {
		return apiLogs, err
	}
	
	var respParsedArr []map[string]interface{} = make([]map[string]interface{}, 0)
	respArr := strings.Split(j.respBody, " ") 
	for _, resp := range respArr {
		var respParsed map[string]interface{}
		err := json.Unmarshal([]byte(resp), &respParsed)
		if err != nil {
			return apiLogs, err
		}
		respParsedArr = append(respParsedArr, respParsed)
	}

	for i, _ := range respParsedArr {
		apiLogs = append(apiLogs, types.ApiLog{
			Request: reqParsed[i],
			Response: respParsedArr[i],
		})
	}
	return apiLogs, nil
}

func (j JobResponse) GetAllApiLogs() ([]types.ApiLog, error) {
	apiLogs := make([]types.ApiLog, 0)
	apiLogs = append(apiLogs, j.GetTransformerApiLogs()...)
	routerApiLogs, err := j.GetRouterAPiLogs()
	if err != nil {
		return apiLogs, err
	}
	apiLogs = append(apiLogs, routerApiLogs...)
	return apiLogs, nil
}

type reloadableConfig struct {
	jobQueryBatchSize                 config.ValueLoader[int]
	updateStatusBatchSize             config.ValueLoader[int]
	readSleep                         config.ValueLoader[time.Duration]
	maxStatusUpdateWait               config.ValueLoader[time.Duration]
	minRetryBackoff                   config.ValueLoader[time.Duration]
	maxRetryBackoff                   config.ValueLoader[time.Duration]
	jobsBatchTimeout                  config.ValueLoader[time.Duration]
	failingJobsPenaltyThreshold       config.ValueLoader[float64]
	failingJobsPenaltySleep           config.ValueLoader[time.Duration]
	noOfJobsToBatchInAWorker          config.ValueLoader[int]
	jobsDBCommandTimeout              config.ValueLoader[time.Duration]
	jobdDBMaxRetries                  config.ValueLoader[int]
	maxFailedCountForJob              config.ValueLoader[int]
	maxFailedCountForSourcesJob       config.ValueLoader[int]
	payloadLimit                      config.ValueLoader[int64]
	retryTimeWindow                   config.ValueLoader[time.Duration]
	sourcesRetryTimeWindow            config.ValueLoader[time.Duration]
	pickupFlushInterval               config.ValueLoader[time.Duration]
	maxDSQuerySize                    config.ValueLoader[int]
	savePayloadOnError                config.ValueLoader[bool]
	transformerProxy                  config.ValueLoader[bool]
	skipRtAbortAlertForTransformation config.ValueLoader[bool] // represents if event delivery(via transformerProxy) should be alerted via router-aborted-count alert def
	skipRtAbortAlertForDelivery       config.ValueLoader[bool] // represents if transformation(router or batch) should be alerted via router-aborted-count alert def
	oauthV2Enabled                    config.ValueLoader[bool]
	oauthV2ExpirationTimeDiff         config.ValueLoader[time.Duration]
}
