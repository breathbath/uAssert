package simulation

import "github.com/breathbath/uAssert/encode"

func GetSimulationMap(scs GrpcCases) GrpcCasesMap {
	scm := GrpcCasesMap{}
	for _, sc := range scs {
		scm[encode.Md5(sc.Request, sc.Namespace)] = sc
	}

	return scm
}

func FindSimulatedCaseForRequest(request interface{}, namespace string, simCasesMap GrpcCasesMap) (GrpcCase, bool) {
	requestChecksum := encode.Md5(request, namespace)
	foundCase, isFound := simCasesMap[requestChecksum]
	if !isFound {
		return GrpcCase{}, false
	}

	return foundCase, true
}

