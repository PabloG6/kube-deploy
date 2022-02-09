package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	app := fiber.New();
	var kubeconfig *string;
	app.Get("kubernetes/deploy", func(c *fiber.Ctx) error {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		//read from a docker registry and build a kubernetes pod.
		config ,err := clientcmd.BuildConfigFromFlags("", *kubeconfig);
		if err != nil {
			return c.SendString(err.Error());
		}
		k8s, _ := kubernetes.NewForConfig(config);
		podDeployment := k8s.CoreV1().Pods(apiv1.NamespaceDefault)
		pod := &apiv1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pod-name",
			},
			Spec: apiv1.PodSpec{
				
				ImagePullSecrets: []apiv1.LocalObjectReference{
					{
						Name: "regcred",
					},
				},
				Containers: []apiv1.Container{
					{
						Name: "nodejs-test-pod",
						Image: "docker-saasify-dev.hookforms.dev/node-test-1",
					},
				},
			},
		}


		ctx := context.Background();
		_, err = podDeployment.Create(ctx, pod, metav1.CreateOptions{})


		if err != nil {
			return c.SendString(err.Error())
		}

	
		return c.SendString("Done deploying application to kubernetes. ")

	})


	app.Listen(":3000")
}