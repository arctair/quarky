package v1

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Deployments ...
type Deployments struct {
	clientset *kubernetes.Clientset
}

// NewDeployments ...
func NewDeployments() CreateDelete {
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
	return &Deployments{
		clientset,
	}
}

// Create ...
func (d *Deployments) Create() (string, error) {
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "quarky-test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "quarky-test",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "quarky-test",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "quarky-test",
							Image: "arctair/quarky-test:1.0.11",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 5000,
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentsClient := d.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	_, err := deploymentsClient.Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{},
	)

	if err != nil {
		return "", err
	}

	return "", nil
}

// Delete ...
func (d *Deployments) Delete() (string, error) {
	deploymentsClient := d.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(
		context.TODO(),
		"quarky-test",
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
	if err != nil {
		return "", nil
	}

	return "", nil
}
