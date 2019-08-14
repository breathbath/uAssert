package grpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimulationCasesMapGeneration(t *testing.T) {
	simulationCases := SimulationCases{
		{
			Request:       "request1",
			Response:      "response1",
			GrpcNamespace: "namespace",
		},
		{
			Request:       "request1",
			Response:      "response1",
			GrpcNamespace: "namespace",
		},
		{
			Request:       "request2",
			Response:      "response2",
			GrpcNamespace: "namespace",
		},
		{
			Request:       "request3",
			Response:      "response3",
			GrpcNamespace: "namespace1",
		},
	}

	expectedMap := SimulationCasesMap{
		"e900c89e27c3ceb71378591da11e9309": { //as md5 on namespace."request1"
			Request:       "request1",
			Response:      "response1",
			GrpcNamespace: "namespace",
		},
		"4658c8c552e44063b6922094f2889701": { //as md5 of namespace."request2"
			Request:       "request2",
			Response:      "response2",
			GrpcNamespace: "namespace",
		},
		"3f98351c3aaaa5b50b611bb2b91e50a8": { //as md5 of namespace1."request3"
			Request:       "request3",
			Response:      "response3",
			GrpcNamespace: "namespace1",
		},
	}

	assert.Equal(t, expectedMap, GetSimulationMap(simulationCases))
}

func TestSimulationCasesSearchInMap(t *testing.T) {
	mapToSearchIn := SimulationCasesMap{
		"e900c89e27c3ceb71378591da11e9309": { //as md5 on namespace."request1"
			Request:       "request1",
			Response:      "response1",
			GrpcNamespace: "namespace",
		},
		"4658c8c552e44063b6922094f2889701": { //as md5 of namespace."request2"
			Request:       "request2",
			Response:      "response2",
			GrpcNamespace: "namespace",
		},
	}

	simulatedCase, found := FindSimulatedCaseForRequest("request1", "namespace", mapToSearchIn)
	assert.True(t, found)
	assert.Equal(
		t,
		SimulationCase{
			Request:       "request1",
			Response:      "response1",
			GrpcNamespace: "namespace",
		},
		simulatedCase,
	)

	_, found = FindSimulatedCaseForRequest("request1", "namespace1", mapToSearchIn)
	assert.False(t, found)

	_, found = FindSimulatedCaseForRequest("request3", "namespace", mapToSearchIn)
	assert.False(t, found)
}
