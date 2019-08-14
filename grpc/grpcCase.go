package grpc

import "github.com/breathbath/uAssert/encode"

type SimulationCases []SimulationCase

type SimulationCasesMap map[string]SimulationCase

type SimulationCase struct {
	Request       interface{}
	Response      interface{}
	GrpcNamespace string
}

func GetSimulationMap(scs SimulationCases) SimulationCasesMap {
	scm := SimulationCasesMap{}
	for _, sc := range scs {
		scm[encode.Md5(sc.Request, sc.GrpcNamespace)] = sc
	}

	return scm
}

func FindSimulatedCaseForRequest(request interface{}, namespace string, simCasesMap SimulationCasesMap) (SimulationCase, bool) {
	requestChecksum := encode.Md5(request, namespace)
	foundCase, isFound := simCasesMap[requestChecksum]
	if !isFound {
		return SimulationCase{}, false
	}

	return foundCase, true
}
