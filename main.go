package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	v1 "arctair.com/quarky/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	sha1    string
	version string
)

// StartHTTPServer ...
func StartHTTPServer(wg *sync.WaitGroup) *http.Server {
	clientset := NewClientset()
	server := &http.Server{
		Addr: ":5000",
		Handler: v1.NewRouter(
			v1.NewRolloutController(
				v1.NewRollouts(
					v1.NewDeployments(
						clientset,
					),
					v1.NewServices(
						clientset,
					),
				),
				&v1.LoggerConsole{},
			),
			v1.NewVersionController(
				v1.NewBuild(sha1, version),
			),
		),
	}

	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return server
}

func main() {
	serverExit := &sync.WaitGroup{}
	serverExit.Add(1)
	StartHTTPServer(serverExit)
	serverExit.Wait()
}

func NewClientset() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "absolute path to kubeconfig")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	return clientset
}
