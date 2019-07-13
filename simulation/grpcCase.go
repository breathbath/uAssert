package simulation

type GrpcCases []GrpcCase

type GrpcCasesMap map[string]GrpcCase

type GrpcCase struct {
	Request  interface{}
	Response interface{}
	Namespace  string
}
