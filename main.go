package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/pb"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	zbClient := getClient()
	getStatus(zbClient)
	deploy(zbClient)

	go zbClient.NewJobWorker().JobType("get-time").Handler(handleGetTime).Open()
	go zbClient.NewJobWorker().JobType("make-greeting").Handler(handleMakeGreeting).Open()

	http.HandleFunc("/start", createStartHandler(zbClient))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

type MakeGreetingCompleteVariables struct {
	Say string 	`json:"say"`
}

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

type BoundHandler func(w http.ResponseWriter, r * http.Request)

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

func deploy (zbClient zbc.Client) {
	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile("test-process.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(response.String())
}

func getClient() zbc.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gatewayAddress := os.Getenv("ZEEBE_ADDRESS")

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
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
		log.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			log.Println("  Partition", partition.PartitionId, ":", roleToString(partition.Role))
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