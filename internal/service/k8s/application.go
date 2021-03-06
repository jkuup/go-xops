package k8ss

import (
	"context"
	"fmt"
	"go-xops/assets/k8s"
	"go-xops/pkg/common"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetApplication(ctx context.Context, namespace, id, name, format string) (k8s.Application, error) {
	result := k8s.Application{}
	data, err := common.ClientSet.AppsV1().Deployments(namespace).List(ctx, v1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", "app", "nginx"),
	})
	if err != nil {
		return result, err
	}

	result.ID = id
	result.Name = name
	result.Format = format
	result.Containers = []k8s.Container{}
	for _, deployment := range data.Items {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			result.Containers = append(result.Containers, k8s.Container{
				Name:  container.Name,
				Image: container.Image,
				Version: strings.Replace(
					container.Image,
					strings.Replace(format, "[.Release]", "", -1),
					"",
					-1,
				),
				Deployment: k8s.Deployment{
					Name: deployment.ObjectMeta.Name,
					UID:  string(deployment.ObjectMeta.UID),
				},
			})
		}
	}
	return result, nil
}

func GetApplications(ctx context.Context) ([]k8s.Application, error) {
	applications := make([]k8s.Application, 0)
	result := k8s.Application{}
	// 获取所有的namespaces
	res, err := GetNameSpaces(ctx)
	if err != nil {
		return applications, err
	}
	// 根据每个namespace来获取所有的 applications
	for _, application := range res {
		data, err := common.ClientSet.AppsV1().Deployments(application.Name).List(ctx, v1.ListOptions{})
		if err != nil {
			// return applications, err
			continue
		}
		for _, deployment := range data.Items {
			for _, container := range deployment.Spec.Template.Spec.Containers {
				result.Containers = append(result.Containers, k8s.Container{
					Name:  container.Name,
					Image: container.Image,
					Version: strings.Replace(
						container.Image,
						strings.Replace("", "[.Release]", "", -1),
						"",
						-1,
					),
					Deployment: k8s.Deployment{
						Name: deployment.ObjectMeta.Name,
						UID:  string(deployment.ObjectMeta.UID),
					},
				})
			}
		}
	}
	applications = append(applications, result)

	return applications, nil
}
