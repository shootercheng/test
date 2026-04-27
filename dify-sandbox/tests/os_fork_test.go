package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/shootercheng/test/dify-sandbox/internal/run"
)

func TestOsFork(t *testing.T) {

	osForkCode, err := run.ReadCodeFromFile("../python/os_fork.py")
	if err != nil {
		t.Fatalf("read os_fork.py file error")
	}

	sandboxClient := run.NewSandboxClient("http://localhost:8194", "dify-sandbox")
	defer sandboxClient.Close()

	request := &run.RequestBody{
		Language:      "python3",
		Code:          osForkCode,
		Preload:       osForkCode,
		EnableNetwork: true,
	}

	var wg sync.WaitGroup

	maxConcurrency := 4
	invokeTime := 100

	sem := make(chan int, maxConcurrency)

	isOk := true
	for i := range invokeTime {
		wg.Go(func() {

			sem <- i
			defer func() {
				<-sem
			}()

			res, err := sandboxClient.SendCode(request)
			if err != nil {
				t.Logf("request error:%v", err)
				isOk = false
				return
			}
			fmt.Printf("request result:%s\n", string(res))
		})
	}

	wg.Wait()
	if !isOk {
		t.Fatalf("test os fork request error")
	}
}
