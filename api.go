package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"log"

	"github.com/boltdb/bolt"

	"github.com/hashicorp/consul/api"
)

type animals struct {
	name string
	age  int
}

func main() {
	updateServiceInConsul()
	httpRouter := httprouter.New()
	httpRouter.POST("/hello", getHello)
	httpRouter.GET("/hello", postHello)
	log.Fatal(http.ListenAndServe(":8888", httpRouter))

}

func getHello(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	db, err := bolt.Open("test_bolt", 0777, nil)
	fmt.Print("error is", err)
	//defer db.Close()

	requestMap := req.URL.Query()
	nameOfAnimal := requestMap.Get("name")
	ageOfAnimal := requestMap.Get("age")
	resp.WriteHeader(201)
	fmt.Println(nameOfAnimal, ageOfAnimal)
	bucketname := []byte("devops")
	keyName := []byte(nameOfAnimal)
	valueName := []byte(ageOfAnimal)

	db.Update(func(tx *bolt.Tx) error {

		b, errorBucket := tx.CreateBucket(bucketname)
		if errorBucket != nil {

			b = tx.Bucket(bucketname)
		}
		fmt.Println("Error bucket is", errorBucket)
		b.Put(keyName, valueName)

		outVal := b.Get(keyName)
		fmt.Print("Updated key and value", string(outVal))
		fmt.Println("Updated key and value", string(keyName))
		return nil

	})

	db.Close()

}

func postHello(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	db, err := bolt.Open("test_bolt", 0777, nil)
	fmt.Print("error is", err)
	//defer db.Close()

	requestMap := req.URL.Query()
	nameOfAnimal := requestMap.Get("name")
	resp.WriteHeader(201)
	fmt.Println(nameOfAnimal)
	bucketname := []byte("devops")
	keyName := []byte(nameOfAnimal)

	db.Update(func(tx *bolt.Tx) error {

		b, errorBucket := tx.CreateBucket(bucketname)
		if errorBucket != nil {

			b = tx.Bucket(bucketname)
		}
		fmt.Println("Error bucket is", errorBucket)

		outVal := b.Get(keyName)
		fmt.Println("Out value is", string(outVal))
		return nil

	})

	db.Close()

}

func updateServiceInConsul() {
	defaultConfig := api.DefaultConfig()
	client, err := api.NewClient(defaultConfig)
	fmt.Println("Error is", err)
	agent := client.Agent()
	agentServiceCheck := &api.AgentServiceCheck{}
	agentServiceCheck.HTTP = "http://localhost:8888/hello"
	agentServiceCheck.Name = "Hello check"
	agentServiceCheck.Interval = "10s"

	agentServiceRegister := &api.AgentServiceRegistration{}
	agentServiceRegister.Name = "devops-service"
	agentServiceRegister.Port = 8888
	agentServiceRegister.Address = "127.0.0.1"
	agentServiceRegister.Check = agentServiceCheck
	var tags []string
	tags = append(tags, "my_web_service")

	agentServiceRegister.Tags = tags
	returnError := agent.ServiceRegister(agentServiceRegister)
	fmt.Println("Return error is", returnError)

}
