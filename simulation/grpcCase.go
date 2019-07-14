package simulation

type SimulationCases []SimulationCase

type SimulationCasesMap map[string]SimulationCase

type SimulationCase struct {
	Request  interface{}
	Response interface{}
	Namespace  string
}
