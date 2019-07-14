package simulation

import "github.com/breathbath/uAssert/encode"

func GetSimulationMap(scs SimulationCases) SimulationCasesMap {
	scm := SimulationCasesMap{}
	for _, sc := range scs {
		scm[encode.Md5(sc.Request, sc.Namespace)] = sc
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

