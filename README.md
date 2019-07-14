# Library for testing microservices

[![Travis Status for breathbath/uAssert](https://api.travis-ci.org/breathbath/uAssert.svg?branch=master&label=linux+build)](https://travis-ci.org/breathbath/uAssert)

## Problem
Testing of complex microservice architectures in production environments is a big challenge, which results in lack unified tools for integration, functional, load/performance tests . There are good libraries for unit testing and mocking in many programming languages, but it’s hard or even impossible to test communication of micro-services on a real infrastructure.
Asynchronous nature of intercommunication doesn't allow to assert predictable synchronous expectations for I/O operations on global architecture level. 

Another challenge is to simulate production environments on local development machines to reproduce and fix communication issues or even do load testing. 

Docker provides a great help for virtualisation and orchestration of heterogeneous services. 

However it is hard to get some of the services running in every possible environment, especially when they are related to external APIs (payment gateways, non-cotrollable data providers etc.) or working with specific hardware devices (e.g. in telecommunication). In this situation market lacks good tools for simulation of behaviour for such kind of services to build testing pipelines for end-to-end testing scenarios.

## Motivation

The library should help to solve following problems:
- Challenges in testing micro-services intercommunication, like asserting certain chain of events, execution time expectations, fault tolerance on network problems, unpredictable startup times of different services incl infrastructure, behaviour assertions for failure cases affecting multiple services in different combinations, services stack misconfiguration (e.g. unexpected amount of cloned services) etc.
- Challenges in simulating uncontrolled external APIs or hardware devices
- Challenges in performance and load testing simulating load on the whole services stack rather than individual services
- Challenges in simulating of unpredictable services behaviour - late responses, constant restarting, request blocking, nonsense responses etc. 
- Challenges in collecting data for making assertions and consistent assertion logic agnostic to the source of data (logs, queues, api calls, persistent states)
Having one tool solving all problems above will help to create qualitative, predictable and measurable testing environments for any micro-service driven projects. It would give rich automation opportunities for tests execution as well as easify development of complex service architectures with external API calls and specific hardware requirements. 

## Tools list

Currently the library contains following tools:

### Testing runtime for stateful test cases

`go test` is a perfect tool for creating unit tests but it lacks lifecycle management and state control for more advanced functional tests.

Testing runtime allows you to define a group of related tests with a shared state and configure it's lifecycle in a centralized way.
Consider a situation where you want to test bunch of cases for a group of microservices, e.g. how a proxy server processes requests
against a device management service. 

Your test cases might be like this:

- assert proxy is healthy when device management is healthy too and not-healthy if one of them fails
- assert that api responses of device management are returned by proxy unmodified
- assert that proxy doesn't allow unauthorised access to device management

For all the tests you want to make sure:
- that the proxy and DM services are up and running before tests, they should not start if both are not ready
- tests should not start at all if there are startup problems, the whole test suit should fail early in this case
- before each test we want to have certain data to be returned by device management api, any modifications made to this data should not be visible in following tests
- after each test we want to reset any changes made by them
- if all tests are finished, we want all microservices to be down and are not occupying any network resources (ports/addresses/sockets)

With a standard go test functionality it's quite cumbersome to implement this testing logic. As this usecase has many similar variations,
we crated testing runtime which allows to do this complex tests lifecycle management.

A test might look as following:
    
        var testsRuntime *tests.Runtime
        
        func init() {
        	testsRuntime = tests.NewRuntime()
        
        	testsRuntime.BeforeAll(func(r *tests.Runtime){
        		proxyServer := NewProxy("localhost:2233")
        
        		err := proxyServer.Start()
        		if err == nil {
        			panic(err) //no tests will be executed further
        		}
        
        		deviceManager := NewDeviceManager("localhost:2233")
        
        		pong := deviceManager.Ping() //just to demonstrate a different health check
        		if pong == nil {
        			proxyServer.Stop() //we want a proper cleanup if the other server is not ready
        			panic(errors.New("Cannot bring up device manager")) //no tests will be executed further
        		}
        
        		r.SetState("proxy", proxyServer) //we want both services be available for LCM
        		r.SetState("deviceManager", proxyServer)
        	})
        
        	testsRuntime.BeforeEach(func(r *tests.Runtime) {
        		dm := r.GetStateOrFail("deviceManager").(*DeviceManager) //retrieve device manager
        		dm.SetDevices([]string{"deviceA", "deviceB"}) //we want to make sure that before any test the device manager has the same state
        	})
        
        	testsRuntime.AfterEach(func(r *tests.Runtime) {
        		dm := r.GetStateOrFail("deviceManager").(*DeviceManager) //retrieve device manager
        		dm.SetDevicesStates([]string{}) //whatever tests modified in device states we reset it
        	})
        
        	testsRuntime.AfterAll(func(r *tests.Runtime) { //doing cleanup only after all tests are done
        		p := r.GetStateOrFail("proxy").(*Proxy)
        		p.Stop()
        
        		dm := r.GetStateOrFail("deviceManager").(*DeviceManager)
        		dm.Stop()
        	})
        
        	testsRuntime.TestCase(assertProxyHealth) //our test cases
        	testsRuntime.TestCase(assertProxyForwardingResponsesUnmodified)
        	testsRuntime.TestCase(assertProxyRightsControl)
        }
        
        func TestProxy(t *testing.T) {
        	testsRuntime.Run(t)
        }
        
        func assertProxyHealth(t *testing.T, r *tests.Runtime) {
        	//... test1
        }
        
        func assertProxyForwardingResponsesUnmodified(t *testing.T, r *tests.Runtime) {
        	//... test2
        }
        
        func assertProxyRightsControl(t *testing.T, r *tests.Runtime) {
        	//... test3
        }

### Grpc simulator

In the previous example about proxy and device manager, we showed that it's possible to modify the list of returned
devices during the tests execution. However in the real situation this service probably won't have this functionality, since
in production code it will have a more complex device management logic. 

For microservices it's ofter a case that certain services are quite hard to be brought to a certain expected state, as it might
require to trigger set of events on other services or it's a complex state machine which has complex conditions and roles driven transitions.
Some services might also be a proxy to communicate to specific hardware devices, which are not possible to have in testing environment.

In all those situations a microservice simulator is needed, which allows us to:

- bring it to an expected state without complex logic
- simulate specific hardware devices which can only work in production
- give methods to easily reset state when it's needed
- log communication traffic for test assertions
- simulate hard reproduceable behaviour like special failures, delays, non-sense responses, down times etc.

We created a GRPC simulator to provide this functionality.

#### Current features
- Creation of "fuzzy" grpc servers merging multiple microservices into one for simplification without collisions
- Possibility to define response stubs mapped to predefined requests in a generic way agnostic to specific grpc implementations
- Async startup mode for testing purposes (which is not initially supported by grpc go implementation) with proper startup error handling and blocking waiting startup time.
- Up/down LCM for testing purposes

#### Todo features
- Support for streams simulation
- Unspecific behaviour simulation
- Extended logging
- Stateful behaviour
- Mapping state to responses
- Code generation (in further perspective)

## Examples
We created examples to demonstrate how the features of the library might be used. We used the following opensource project for that:

- Voltha (https://github.com/opencord/voltha-go) - a library to abstract representation of specific network hardwaredevices
- Voltha protos (https://github.com/opencord/voltha-protos) - GRPC models for the above project

The microservices set consists of 2 items:

- voltha service (network devices software abstraction layer)
- access proxy (layer to control access to voltha and abstract it even more by simplifying some interfaces)

We want to write functional tests against the access proxy. For that we want to simulate Voltha behaviour, since it's quite complex and requires hardware devices in the system.

![alt text](https://breathbath.com/files/dUFzc2VydEV4YW1wbGVVc2FnZS5wbmc=_bklmt957f7pistjue2gg.png)

Our test case we define as following:

Voltha cannot provide a Device information by it's serial number as it supports only [GetDevice by id method](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/voltha.proto#L337)
On the other hand it has [ListDevices GRPC method](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/voltha.proto#L312) which the access proxy might faciliate to find a device by its id.

We have defined a GRPC service and [GetDeviceBySn method](https://github.com/breathbath/uAssert/blob/master/projects/accessProxy/protos/access_proxy/accessProxy.proto#L12).

To generate a corresponding go files we used the following command

    make build
    
Our [Devices service](https://github.com/breathbath/uAssert/blob/master/projects/accessProxy/protos/access_proxy/accessProxy.proto#L11) is using the imported [Voltha Device model](https://github.com/opencord/voltha-protos/blob/master/protos/voltha_protos/device.proto).

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

We used [uAssert testing Runtime](https://github.com/breathbath/uAssert/blob/master/test/runtime.go) ensure following LCM:

- before any tests we want to make sure that the voltha simulator service and access proxy are up:

    	func init() {
        	testsRuntime = tests.NewRuntime()
        	testsRuntime.BeforeAll(setup)
        	...
        }
    	
    	...
    	
    	func setup(r *tests.Runtime) {
        	volthaServerSimulator := voltha2.NewVolthaServerSimulator(VOLTHA_SERVER)
        	err := volthaServerSimulator.StartAsync(time.Microsecond * 500)
        	if err != nil {
        		log.Panic(err)
        	}
        	r.SetState("voltha_server", volthaServerSimulator)
        
        	accessProxyServer := NewAccessProxyServer(ACCESS_PROXY_SERVER, VOLTHA_SERVER)
        	err = accessProxyServer.StartAsync(time.Microsecond * 500)
        	if err != nil {
        		log.Panic(err)
        	}
        	r.SetState("access_proxy_server", accessProxyServer)
        }

Please note that the both services are started async (which is not natively supported by grpc implementation).        

- we run our test(s)

        func init() {
        	testsRuntime = tests.NewRuntime()
        	testsRuntime.TestCase(testDeviceIdSnMapping)
        	...
        }
        
        func testDeviceIdSnMapping(t *testing.T, r *tests.Runtime) {
            ...
        	assert.NoError(t, err)
        	...
        }
        
- after all tests we want to make sure that both services are down and not blocking any network resources:

    
        func init() {
            testsRuntime = tests.NewRuntime()
            ...
            testsRuntime.AfterAll(cleanup)
            ...
        }
        ...
        func cleanup(r *tests.Runtime) {
            r.GetStateOrFail("voltha_server").(*simulation.GrpcServer).Stop()
            r.GetStateOrFail("access_proxy_server").(*simulation.GrpcServer).Stop()
        }
    
Please note, that disregard the async mode, the server stopping will be blocking until the servers are down.

If our test execution is green, we successfully evaluated a correct GRPC access proxy work without complex setup
procedure for Voltha service and the need of hardware devices behind. 

