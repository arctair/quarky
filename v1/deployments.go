package v1

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Deployments ...
type Deployments struct {
	clientset *kubernetes.Clientset
}

// NewDeployments ...
func NewDeployments(c *kubernetes.Clientset) CreateDelete {
	return &Deployments{c}
}

// Create ...
func (d *Deployments) Create() error {
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

	return err
}

// Delete ...
func (d *Deployments) Delete() error {
	deploymentsClient := d.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	return deploymentsClient.Delete(
		context.TODO(),
		"quarky-test",
		metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		},
	)
}
