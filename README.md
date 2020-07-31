---
id: go-client
title: Go
---
# Getting Started with Camunda Cloud using Go

The [Zeebe Go Client](https://github.com/zeebe-io/zeebe/tree/develop/clients/go) is available for Go applications. 

Watch a [video tutorial on YouTube](https://youtu.be/7vBxJmXD3Js) walking through this Getting Started Guide.

[![](assets/getting-started-java-thumbnail.png)](https://youtu.be/7vBxJmXD3Js)

Estimated time to complete: 60 minutes

## Prerequisites

* [Go 1.13+](https://golang.org/dl/)
* [Zeebe Modeler](https://github.com/zeebe-io/zeebe-modeler/releases)

## Scaffolding the project

* Run the following command to create a new Go project:

```bash
mkdir -p $GOPATH/src/github.com/$USER/cloud-starter
cd $GOPATH/src/github.com/$USER/cloud-starter
```

[Video link](https://youtu.be/7vBxJmXD3Js?t=71)

* Add the [Zeebe Go Client](https://github.com/zeebe-io/spring-zeebe) to the project:

```bash
go get -u github.com/zeebe-io/zeebe/clients/go/...
```
[Video link](https://youtu.be/7vBxJmXD3Js?t=192)

## Create Camunda Cloud cluster

[Video link](https://youtu.be/7vBxJmXD3Js?t=233)

* Log in to [https://camunda.io](https://camunda.io).
* Create a new Zeebe 0.23.3 cluster.
* When the new cluster appears in the console, create a new set of client credentials. 
* Copy the client Connection Info environment variables block.

## Configure connection

[Video link](https://youtu.be/7vBxJmXD3Js?t=368)

We will use [GoDotEnv](https://github.com/joho/godotenv) to environmentalize the client connection credentials.

* Add GoDotEnv to the project: 

```bash
go get github.com/joho/godotenv
```

* Add the client connection credentials for your cluster to the file `.env`:

**Note**: _make sure to remove the `export` keyword from each line_.

```
ZEEBE_ADDRESS='aae86771-0906-4186-8d82-e228097e1ef7.zeebe.camunda.io:443'
ZEEBE_CLIENT_ID='hj9PHRIiRqT0~qHvFeqXZV-J8fLRfifB'
ZEEBE_CLIENT_SECRET='.95Vlv6joiuVR~mJDjGPlyYk5Pz6iIwFYmmQyX8yU3xdB1gezntVMoT1SQTdrCsl'
ZEEBE_AUTHORIZATION_SERVER_URL='https://login.cloud.camunda.io/oauth/token'
```

* Save the file.

## Test Connection with Camunda Cloud

[Video link](https://youtu.be/7vBxJmXD3Js?t=436)

* Paste the following code into the file `main.go`:

```go
package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zeebe-io/zeebe/clients/go/pkg/pb"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"log"
	"os"
)

func main() {
    zbClient := getClient()
    getStatus(zbClient)
}

func getClient() zbc.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gatewayAddress := os.Getenv("ZEEBE_ADDRESS")

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:  gatewayAddress,
	})

	if err != nil {
		panic(err)
	}

    return zbClient
}

func getStatus(zbClient zbc.Client) {
	ctx := context.Background()
	topology, err := zbClient.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}

	for _, broker := range topology.Brokers {
		fmt.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			fmt.Println("  Partition", partition.PartitionId, ":", roleToString(partition.Role))
		}
	}
}

func roleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}
```

* Run the program with the command `go run main.go`. 

* You will see output similar to the following:

```
2020/07/29 06:21:08 Broker zeebe-0.zeebe-broker-service.aae86771-0906-4186-8d82-e228097e1ef7-zeebe.svc.cluster.local : 26501
2020/07/29 06:21:08   Partition 1 : Leader
```

This is the topology response from the cluster.

## Create a BPMN model

[Video link](https://youtu.be/7vBxJmXD3Js?t=900)

* Download and install the [Zeebe Modeler](https://github.com/zeebe-io/zeebe-modeler/releases).
* Open Zeebe Modeler and create a new BPMN Diagram.
* Create a new BPMN diagram.
* Add a StartEvent, an EndEvent, and a Task.
* Click on the Task, click on the little spanner/wrench icon, and select "Service Task".
* Set the _Name_ of the Service Task to `Get Time`, and the _Type_ to `get-time`.

It should look like this:

![](assets/gettingstarted_first-model.png)

* Click on the blank canvas of the diagram, and set the _Id_ to `test-process`, and the _Name_ to "Test Process".
* Save the diagram to `test-process.bpmn` in your project.

## Deploy the BPMN model to Camunda Cloud

[Video Link](https://youtu.be/7vBxJmXD3Js?t=1069)

* Edit the `main.go` file, and add a new function `deploy`:

```go
func deploy (zbClient zbc.Client) {
	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile("test-process.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(response.String())
}
```

* Now update the `main()` function to look like this:

```go
func main() {
	zbClient := getClient()
	getStatus(zbClient)
	deploy(zbClient)
}
```
* Run the program with `go run main.go`.

You will see the deployment response:

```
2020/07/29 06:23:07 Broker zeebe-0.zeebe-broker-service.aae86771-0906-4186-8d82-e228097e1ef7-zeebe.svc.cluster.local : 26501
2020/07/29 06:23:07   Partition 1 : Leader
2020/07/29 06:23:08 key:2251799813685251 workflows:<bpmnProcessId:"test-process" version:1 workflowKey:2251799813685249 resourceName:"test-process.bpmn" > 
```

## Start a Workflow Instance

[Video Link](https://youtu.be/7vBxJmXD3Js?t=1093)

We will add a webserver, and use it to provide a REST interface for our program.
	
* Add `"net/http"` to the imports in `main.go`.

* Edit the `main.go` file, and add a REST handler to start an instance:

```go
type BoundHandler func(w http.ResponseWriter, r * http.Request)

func createStartHandler(client zbc.Client) BoundHandler {
	f := func (w http.ResponseWriter, r * http.Request) {
		ctx := context.Background()
		request, err := client.NewCreateInstanceCommand().BPMNProcessId("test-process").LatestVersion().Send(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, request.String())
	}
	return f
}
```

* Now update the `main()` function to add an HTTP server and the `/start` route:

```go
func main() {
	zbClient := getClient()
	getStatus(zbClient)
	deploy(zbClient)

	http.HandleFunc("/start", createStartHandler(zbClient))

	http.ListenAndServe(":3000", nil)
}
```

* Run the program with the command: `go run main.go`.

* Visit [http://localhost:3000/start](http://localhost:3000/start) in your browser.

You will see output similar to the following: 

```
workflowKey:2251799813685249 bpmnProcessId:"test-process" version:1 workflowInstanceKey:2251799813685257 
``` 

A workflow instance has been started. Let's view it in Operate.

## View a Workflow Instance in Operate

[Video Link](https://youtu.be/7vBxJmXD3Js?t=1223)

* Go to your cluster in the [Camunda Cloud Console](https://camunda.io).
* In the cluster detail view, click on "_View Workflow Instances in Camunda Operate_".
* In the "_Instances by Workflow_" column, click on "_Test Process - 1 Instance in 1 Version_".
* Click the Instance Id to open the instance.
* You will see the token is stopped at the "_Get Time_" task.

Let's create a task worker to serve the job represented by this task.

## Create a Job Worker

[Video Link](https://youtu.be/7vBxJmXD3Js?t=1288)

We will create a worker that logs out the job metadata, and completes the job with success.

* Edit the `main.go` file, and add a handler function for the worker:

```go
func handleGetTime(client worker.JobClient, job entities.Job) {
	log.Println(job)
	ctx := context.Background()
	client.NewCompleteJobCommand().JobKey(job.Key).Send(ctx)
}
```

* Update the `main` function to create a worker in a go routine, which allows it to run in a background thread:

```go
func main() {
	zbClient := getClient()
	getStatus(zbClient)
	deploy(zbClient)
	
	go zbClient.NewJobWorker().JobType("get-time").Handler(handleGetTime).Open()

	http.HandleFunc("/start", createStartHandler(zbClient))
	http.ListenAndServe(":3000", nil)
}
```

* Run the worker program with the command: `go run main.go`.

You will see output similar to: 

```
2020/07/29 23:48:15 {{2251799813685262 get-time 2251799813685257 test-process 1 2251799813685249 Activity_1ucrvca 2251799813685261 {} default 3 1596023595126 {} {} [] 0}}
```

* Go back to Operate. You will see that the workflow instance is gone.
* Click on "Running Instances".
* In the filter on the left, select "_Finished Instances_".

You will see the completed workflow instance.

## Create and Await the Outcome of a Workflow Instance 

[Video Link](https://youtu.be/7vBxJmXD3Js?t=1616)

We will now create the workflow instance, and get the final outcome in the calling code.

* Edit the `createStartHandler` function in the `main.go` file, and make the following change:

```go
// ...
request, err := client.NewCreateInstanceCommand().BPMNProcessId("test-process").LatestVersion().WithResult().Send(ctx)
// ...
```

* Run the program with the command: `go run main.go`.

* Visit [http://localhost:3000/start](http://localhost:3000/start) in your browser.

You will see output similar to the following:

```
workflowKey:2251799813685249 bpmnProcessId:"test-process" version:1 workflowInstanceKey:2251799813689873 variables:"{}" 
```

## Call a REST Service from the Worker 

[Video link](https://youtu.be/7vBxJmXD3Js?t=1845)

* Edit the file `main.go` and add these structs for the REST response payload and the worker variable payload:

```go
type Time struct {
	Time string `json:"time"`
	Hour int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
	Day int `json:"day"`
	Month int `json:"month"`
	Year int `json:"year"`
}

type GetTimeCompleteVariables struct {
	Time Time `json:"time"`
}
```

* Replace the `handleGetTime` function in `main.go` with the following code:

```go
func handleGetTime(client worker.JobClient, job entities.Job) {
	var (
		data []byte
	)
	response, err := http.Get("https://json-api.joshwulf.com/time")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}
	var time Time
	err = json.Unmarshal(data, &time)
	if err != nil {
		log.Fatalln(err)
		return
	}
	payload := &GetTimeCompleteVariables{time}

	ctx := context.Background()
	cmd, _ := client.NewCompleteJobCommand().JobKey(job.Key).VariablesFromObject(payload)
	_, err = cmd.Send(ctx)
	
	if err !=nil {
		log.Fatalln(err)
	}
}
```

* Run the program with the command: `go run main.go`.
* Visit [http://localhost:3000/start](http://localhost:3000/start) in your browser.

You will see output similar to the following:

```
workflowKey:2251799813685249 bpmnProcessId:"test-process" version:1 workflowInstanceKey:2251799813693481 
variables:"{\"time\":{\"time\":\"Thu, 30 Jul 2020 13:33:13 GMT\",\"hour\":13,\"minute\":33,\"second\":13,\"day\":4,\"month\":6,\"year\":2020}}" 
```

## Make a Decision 

[Video link](https://youtu.be/7vBxJmXD3Js?t=2066)

We will edit the model to add a Conditional Gateway.

* Open the BPMN model file `bpmn/test-process.bpmn` in the Zeebe Modeler.
* Drop a Gateway between the Service Task and the End event.
* Add two Service Tasks after the Gateway.
* In one, set the _Name_ to `Before noon` and the _Type_ to `make-greeting`.
* Switch to the _Headers_ tab on that Task, and create a new Key `greeting` with the Value `Good morning`.
* In the second, set the _Name_ to `After noon` and the _Type_ to `make-greeting`.
* Switch to the _Headers_ tab on that Task, and create a new Key `greeting` with the Value `Good afternoon`.
* Click on the arrow connecting the Gateway to the _Before noon_ task.
* Under _Details_ enter the following in _Condition expression_: 

```
=time.hour >=0 and time.hour <=11
```

* Click on the arrow connecting the Gateway to the _After noon_ task. 
* Click the spanner/wrench icon and select "Default Flow".
* Connect both Service Tasks to the End Event.

It should look like this:

![](assets/gettingstarted_second-model.png)  

## Create a Worker that acts based on Custom Headers 

[Video link](https://youtu.be/7vBxJmXD3Js?t=2399)

We will create a second worker that takes the custom header and applies it to the variables in the workflow.

* Edit the `main.go` file, and add structs for the new worker's request and response payloads:

```go
type MakeGreetingCompleteVariables struct {
	Say string 	`json:"say"`
}

type MakeGreetingCustomHeaders struct {
	Greeting string `json:"greeting"`
}

type MakeGreetingJobVariables struct {
	Name string `json:"name"`
}
```

* Create the `handleMakeGreeting` method in `main.go`:

```go
func handleMakeGreeting(client worker.JobClient, job entities.Job) {
	variables, _ := job.GetVariablesAsMap()
	name := variables["name"]
	var headers, _ = job.GetCustomHeadersAsMap()
	greeting := headers["greeting"]
	greetingString := greeting + " " + name.(string)
	say := &MakeGreetingCompleteVariables{greetingString}
	log.Println(say)
	ctx := context.Background()
	response, _ := client.NewCompleteJobCommand().JobKey(job.Key).VariablesFromObject(say)
	response.Send(ctx)
}
```

* Edit the `main` function and create a worker:

```go
// ...
	go zbClient.NewJobWorker().JobType("make-greeting").Handler(handleMakeGreeting).Open()
/ ...
```

* Edit the `startWorkflowInstance` method, and make it look like this:

```go
type InitialVariables struct {
	Name string `json:"name"`
}
func createStartHandler(client zbc.Client) BoundHandler {
	f := func (w http.ResponseWriter, r * http.Request) {
		variables := &InitialVariables{"Josh Wulf"}
		ctx := context.Background()
		request, err := client.NewCreateInstanceCommand().BPMNProcessId("test-process").LatestVersion().VariablesFromObject(variables)
		if err != nil {
			panic(err)
		}
		response, _ := request.WithResult().Send(ctx)
		var result map[string]interface{}
		json.Unmarshal([]byte(response.Variables), &result)
		fmt.Fprint(w, result["say"])
	}
	return f
}
```

You can change the variable value to your own name (or derive it from the url path or a parameter).

* Run the program with the command: `go run main`.
* Visit [http://localhost:3000/start](http://localhost:3000/start) in your browser.

You will see output similar to the following:

```
Good Morning Josh Wulf
```

## Profit!

Congratulations. You've completed the Getting Started Guide for Camunda Cloud using Java with Spring.

