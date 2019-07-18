# Library for testing microservices

[![Travis Status for breathbath/uAssert](https://api.travis-ci.org/breathbath/uAssert.svg?branch=master&label=linux+build)](https://travis-ci.org/breathbath/uAssert)

## Problem
Testing of complex micro-service architectures is a big challenge because of lack of good tools for integration, functional, load/performance testing. In many programming languages there are libraries for unit testing and mocks generation, but it’s hard or sometimes even impossible to test a subset of services in their communication. Wild mixtures of async and sync ways to exchange data, events, commands, method calls in heterogenous micro-services environments make it hard to do predictable, synchronous expectations for I/O operations on global architecture level. 

Another challenge is to reproduce production similar setups in local environments, especially if you have services which communicate to external data/service providers or specific hardware devices. This makes it hard to find communication bugs, replay specific use cases or do load testing. 

In this situation market requires tools for making advanced assertions about async communications, easy usable simulation libraries or special traffic generators for load testing.

## Motivation

The project uAssert is intended to solve following problems:
- Testing of intercommunication in complex micro-service architectures, like asserting certain chain of events, execution time expectations, fault tolerance on network problems, unpredictable startup times of different services, behaviour assertions in case of network failures, lags or misconfiguration.
- Easy simulation of uncontrollable services (e.g. locked for specific hardware or external APIs), or isolation of a subset of services with a full control of states and behaviour for some of them (e.g. to replay user checkout with simulated payment gateway and accounting system)
- Traffic generation for different use cases to do performance and load testing
- Testing of async executions with unpredictable times for performing assertions (e.g. asserting successful device activation which is performed by a background process with unknown activation time or asserting that a certain chain of events happened before activation took place)
- Consistent assertion scenarios agnostic to the source of events (logs, queues, api calls, persistent states)
- Possibility to run complex testing scenarios with easily manageable state and life-cycle management (LCM) of the involved services

### Recommendations for a stateful functional testing

`go test` is a perfect tool for unit tests, for functional testing we can use sub tests for stateful execution and LCM.

Subtests allow you to execute a group of related tests with a shared state. Consider a situation where you want to test bunch of cases for the same micro-services subset. You can destroy and recover the whole setup for each test case, but it takes much time. You also probably would like to play subsequent testing scenarios where the state changes by each test and then evaluate the final result. 
Imagine a proxy service which controls access to some device manager service. 

Your test cases might be like this:

- assert proxy is healthy when device management is healthy and not-healthy if one of them fails
- assert that api responses of device management are forwarded by proxy unmodified
- assert that proxy doesn't allow unauthorised access to device management (DM)

For all the assertions you want to make sure:
- that the proxy and DM services are up and running before testing, tests should not be executed if both services are not ready yet to accept connections
- tests should not start at all if there are startup problems, the whole test suit should fail early in this case
- before each test we want to reset the list of devices, any modifications made to this data should not be visible in following tests
- after each test we want to reset any saved data
- if all tests are finished, we want all micro-services to be down to make sure they are not blocking network ports/IPs for further tests

Let's use sub tests to accomplish those LCM requirements:
    
        var proxyServer *Proxy
        var deviceManager *DeviceManager
        
        func setup() {
        	proxyServer = NewProxy("localhost:2233")
        
            err := proxyServer.Start()
            if err == nil {
                panic(err) //no tests will be executed further
            }
        
        	deviceManager = NewDeviceManager("localhost:2233")
        	pong := deviceManager.Ping() //just to demonstrate a different health check
            if pong == nil {
                proxyServer.Stop() //we want a proper cleanup if the other server is not ready
                panic(errors.New("Cannot bring up device manager")) //no tests will be executed further
            }
        }
        
        func cleanup() {
            proxyServer.Stop()
            deviceManager.Stop() 
        }
        
        func TestProxy(t *testing.T) { //main test to execute subtests
        	setup()
        	defer cleanup()
        	
        	deviceManager.SetDevices([]string{"deviceA", "deviceB"}) //we want to make sure that before any test the device manager has the same state
        	t.Run("assertProxyHealth", testProxyHealth)
        	deviceManager.SetDevicesStates([]string{}) //whatever tests modified in device states we reset it
        	
        	deviceManager.SetDevices([]string{"deviceA", "deviceB"})
        	t.Run("assertProxyForwardingResponsesUnmodified", testProxyForwardingResponsesUnmodified)
        	deviceManager.SetDevicesStates([]string{})
        	
            deviceManager.SetDevices([]string{"deviceA", "deviceB"})
            t.Run("assertProxyRightsControl", testProxyRightsControl)
            deviceManager.SetDevicesStates([]string{})
        }
       
        
        func testProxyHealth(t *testing.T, r *tests.Runtime) {
        	//... sub test1
        }
        
        func testProxyForwardingResponsesUnmodified(t *testing.T, r *tests.Runtime) {
        	//... sub test2
        }
        
        func testProxyRightsControl(t *testing.T, r *tests.Runtime) {
        	//... sub test3
        }

The chain of events for this test will look like this:

- Starting proxy, exit on failure
- Starting device manager, exit on failure
- Setting the list of devices to "deviceA", "deviceB"
- Executing test for health check for both services
- Resetting devices states storage so it's empty for the next test
- Setting the list of devices to "deviceA", "deviceB"
- Executing test for proxy for immutable forwarding of requests
- Resetting devices states storage so it's empty for the next test
- Setting the list of devices to "deviceA", "deviceB"
- Executing tests for forbidding unauthorized access to device manager
- Resetting devices states storage so it's empty for the next test
- Stopping device manager and proxy

## Current status of uAssert library

Currently the library provides following tools:

### GRPC simulator

In the previous example about proxy and device manager, we showed that it's possible to modify the list of returned devices during the tests execution. However in a real application it could be hard (e.g. requiring to clean the database) or very slow. In this situation we urgently need a simulator, which could pretend to be a device manager and keep devices list in memory allowing fast and simple database reset. In some cases it's impossible to have a real service in testing environment, e.g. our device manager talks to a specific hardware, which only exists in production.

The simulator should be able to:

- bring a service to an expected state without complex logic (e.g. setting a devices list without complicated triggers on a hardware)
- easily simulate specific responses for specific requests
- quickly reset states (e.g. devices list stored in a in-memory database) without overhead
- log communication traffic for test assertions (e.g. that the device manager got an expected sequence of inputs)
- simulate non-typical behaviour like failures, delays, non-sense responses, down times etc.

For GRPC serives we created a generic GRPC simulator with following features

#### Features
- Ability to combine multiple micro-services into a single GRPC server for simplification and reduction of startup times.
- Easy way to define response stubs mapped to predefined requests for any specific GRPC implementation
- Easy background execution (which is not supported by go-grpc library out of the box) and sync failures for startup errors. Configurable waiting time for the startup process. 
- Sync termination in async mode (test won't exit until the background service is terminated)

#### Todo features
- Streams simulation
- Unspecific behaviour simulation
- Extended logging
- Predefined dynamic logic (e.g. timestamps/ids/hashes/counters generation)
- Stateful execution (e.g. first request adds a value, second returns it) for more complex scenarios
- Responses dependant on state (e.g. if a service is called for the first time it returns data A, afterwards it returns data B)
- Code generation (in further perspective) for GRPC simulators

## Examples
To demonstrate features of the library we created an example project in projects folder. We took [Voltha](https://github.com/opencord/voltha-go) as an example of micro-services, which we want to simulate (as it requires specific hardware). Generally it represent a set of services to abstract communication to specific network hardware devices. We also took its models in GRPC format (https://github.com/opencord/voltha-protos).

For the simplification we will work with 2 micro-services:

- Device manager simulator (DMS)
- Proxy to control access to DMS and extend its APIs

We want to write functional tests against the proxy which will call device manager simulator. The later one will return predefined responses. We will use GRPC simulator from our library.

![alt text](https://breathbath.com/files/dUFzc2VydEV4YW1wbGVVc2FnZS5wbmc=_bklmt957f7pistjue2gg.png)

Our test case we define as following:

DMS cannot provide Device information by a serial number as it supports only [GetDevice by id method](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/voltha.proto#L337)
On the other hand it has [ListDevices GRPC method](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/voltha.proto#L312), so our proxy may extend the functionality by adding additional mapping logic.

We have defined a GRPC service by adding a [GetDeviceBySn method](https://github.com/breathbath/uAssert/blob/master/projects/accessProxy/protos/access_proxy/accessProxy.proto#L12).

To generate a corresponding go files we used following command

    make build
    
Our [Devices service](https://github.com/breathbath/uAssert/blob/master/projects/accessProxy/protos/access_proxy/accessProxy.proto#L11) is using the imported [DM Device model](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/device.proto).

We used [uAssert Grpc Simulator](https://github.com/breathbath/uAssert/blob/master/projects/voltha/server.go) with [the stubs](https://github.com/breathbath/uAssert/blob/master/projects/voltha/stubs.go) to simulate specific ListDevices response:

    func GetStubs() simulation.SimulationCases {
    	return simulation.SimulationCases{
    		{
    			Request:  &empty.Empty{},
    			Response: &voltha.Devices{
    				Items:[]*voltha.Device{
    					{
    						Id: "id1",
    						Type: "Olt",
    						Root: true,
    						ParentId: "",
    						ParentPortNo: 22,
    						Vendor: "Some",
    						Model: "xyw",
    						HardwareVersion: "333",
    						FirmwareVersion: "333",
    						Address: &voltha.Device_Ipv4Address{"11:111:111:11"},
    						SerialNumber: "sn2",
    					},
    				},
    			},
    			Namespace:  "/voltha.VolthaService/ListDevices",
    		},
    	}
    }

The GetDeviceBySn method is expected to return a device from the result of ListDevices method when it's serial number is matched to the provided one (sn2).

We used [uAssert testing Runtime](https://github.com/breathbath/uAssert/blob/master/test/runtime.go) for LCM:

- before any tests we want to make sure that the DMS service and proxy are up:
        
        
        var (
        	volthaServerSimulator *simulation.GrpcServer
        	accessProxyServer *simulation.GrpcServer
        	accessProxyDevicesClient access_proxy.DevicesClient
        )
        
        func setup() {
            volthaServerSimulator = voltha2.NewVolthaServerSimulator(VOLTHA_SERVER)
            err := volthaServerSimulator.StartAsync(time.Microsecond * 500)
            if err != nil {
                log.Panic(err)
            }
        
            accessProxyServer = NewAccessProxyServer(ACCESS_PROXY_SERVER, VOLTHA_SERVER)
            err = accessProxyServer.StartAsync(time.Microsecond * 500)
            if err != nil {
                log.Panic(err)
            }
        
            grpcConn, err := grpc.Dial(ACCESS_PROXY_SERVER, grpc.WithInsecure())
            if err != nil {
                log.Panic(err)
            }
            accessProxyDevicesClient = access_proxy.NewDevicesClient(grpcConn)
        }

Please note that the both services are started async so they are not blocking the main test process.        

- we run our test(s)


        func testNoDeviceBySnIsFound(t *testing.T) {
        	device, err := accessProxyDevicesClient.GetDeviceBySn(
        		context.Background(),
        		&access_proxy.SerialNumber{Sn: "some_unknown_sn"},
        		grpc.WaitForReady(true),
        	)
        	assert.EqualError(t, err, "rpc error: code = NotFound desc = some_unknown_sn")
        	assert.Nil(t, device)
        }
        
- after all tests we want to make sure that both services are down and not blocking any network resources:

    
        func cleanup() {
            volthaServerSimulator.Stop()
            accessProxyServer.Stop()
        }
        
- our main test function is:


        func TestAccessProxy(t *testing.T) {
        	setup()
        	defer cleanup()
        	t.Run("testNoDeviceBySnIsFound", testNoDeviceBySnIsFound)
        }
    
Please note, that disregard the async mode, the server stopping will be blocking until the servers are down.

If our tests execution is green, we ensured a correct work of the proxy service with predictable response from device manager simulator. We can also additionally check, if proxy won't fail if DMS returns an empty list, is not using id as a search parameter for serial number by using id instead of SN in the proxy input. If the proxy is expected to cache intermediate requests, we could bring down the DMS after the first request and make sure that the behaviour of proxy hasn't change for the second request. 

As you might see, uAssert library provides a big flexibility for testing complex micro-service architectures.

## Further development ideas

- Support of async executions with message buses
- Timeout/duration assertions
- Assertions against logs and queue messages
- Tools for load tests
- Integration with docker and K8s envs

