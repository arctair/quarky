package v1

import (
	"context"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// Services ...
type Services struct {
	clientset *kubernetes.Clientset
}

// NewServices ...
func NewServices(c *kubernetes.Clientset) *Services {
	return &Services{c}
}

// Create ...
func (s *Services) Create() error {
	servicesClient := s.clientset.CoreV1().Services(corev1.NamespaceDefault)
	_, err := servicesClient.Create(
		context.TODO(),
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "quarky-test",
				Labels: map[string]string{
					"app": "quarky-test",
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					{
						Port:       8080,
						TargetPort: intstr.FromInt(5000),
					},
				},
				Selector: map[string]string{
					"app": "quarky-test",
				},
			},
		},
		metav1.CreateOptions{},
	)

	return err
}

// Delete ...
func (s *Services) Delete() error {
	servicesClient := s.clientset.CoreV1().Services(corev1.NamespaceDefault)
	return servicesClient.Delete(
		context.TODO(),
		"quarky-test",
		metav1.DeleteOptions{},
	)
}
