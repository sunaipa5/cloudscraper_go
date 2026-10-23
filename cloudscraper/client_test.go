package cloudscraper

import (
	"fmt"
	"testing"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func TestCloudScrapper_Get(t *testing.T) {
	fmt.Println("\n==========================================")
	fmt.Println(">>> RUNNING TEST: Get")
	fmt.Println("==========================================")

	client, err := Init(false)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	res, err := client.Get("https://httpbin.org/get", nil, "")
	if err != nil {
		t.Fatalf("client Get() error: %v", err)
	}

	fmt.Printf("Get Status: %d\n", res.Status)
	fmt.Println(res.Body)
	fmt.Println(">>> TEST GET PASSED")
}

func TestCloudScrapper_Post(t *testing.T) {
	fmt.Println("\n==========================================")
	fmt.Println(">>> RUNNING TEST: Post")
	fmt.Println("==========================================")

	client, err := Init(false)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	payload := `{"key": "value"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := client.Post("https://httpbin.org/post", headers, payload)
	if err != nil {
		t.Fatalf("client Post() error: %v", err)
	}

	fmt.Printf("Post Status: %d\n", res.Status)
	fmt.Println(res.Body)
	fmt.Println(">>> TEST POST PASSED")
}

func TestCloudScrapper_Do(t *testing.T) {
	fmt.Println("\n==========================================")
	fmt.Println(">>> RUNNING TEST: Do")
	fmt.Println("==========================================")

	client, err := Init(false)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	options := cycletls.Options{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            `{"test": "data"}`,
		Timeout:         10,
		DisableRedirect: true,
	}

	res, err := client.Do("https://httpbin.org/post", options, "POST")
	if err != nil {
		t.Fatalf("client Do() error: %v", err)
	}

	fmt.Printf("Do Status: %d\n", res.Status)
	fmt.Println(res.Body)
	fmt.Println(">>> TEST DO PASSED")
}

func TestCloudScrapper_Queue(t *testing.T) {
	fmt.Println("\n==========================================")
	fmt.Println(">>> RUNNING TEST: Queue")
	fmt.Println("==========================================")

	client, err := Init(false)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	options := cycletls.Options{
		Body:            "",
		Timeout:         10,
		DisableRedirect: true,
	}

	client.Queue("https://httpbin.org/put", options, "PUT")

	select {
	case res := <-client.RespChan():
		fmt.Printf("Get Status: %d\n", res.Status)
		fmt.Println(res.Body)
		fmt.Println(">>> TEST QUEUE PASSED")
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout: No response received from channel")
	}
}
