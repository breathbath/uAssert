package voltha

import (
	"context"
	"fmt"
	"github.com/breathbath/uAssert/encode"
	"github.com/breathbath/uAssert/simulation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

var volthaServer *simulation.GrpcServer

func init() {
	volthaServer = GetVolthaServer()

	err := volthaServer.StartAsync()
	if err != nil {
		log.Panic(err)
	}
}

func TestVoltha(t *testing.T) {
	defer volthaServer.Stop()

	grpcConn, err := grpc.Dial(GRPC_ADDRESS, grpc.WithInsecure())
	assert.NoError(t, err)

	grpcClientCtx, grpcClientCancel := context.WithTimeout(context.Background(), time.Second)
	defer grpcClientCancel()

	for _, sc := range GetStubs() {
		err := grpcConn.Invoke(grpcClientCtx, sc.Namespace, sc.Request, sc.Response)
		assert.NoError(t, err)
		fmt.Println(encode.StringifyGraceful(sc.Response, true))
	}
}
