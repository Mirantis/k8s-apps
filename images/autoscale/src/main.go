package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	prometheus "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

// Config : store autoscale configuration
type Config struct {
	PrometheusAddress   string `required:"true" envconfig:"PROMETHEUS_ADDRESS"`
	KubernetesNamespace string `envconfig:"KUBERNETES_NAMESPACE" default:""`
	PollingInterval     int    `envconfig:"POLLING_INTERVAL" default:"30"`
}

// ScaleData : scale information taken from k8s objects annotations
type ScaleData struct {
	ScaleUp     string
	ScaleDown   string
	MinReplicas int
	MaxReplicas int
}

func parseAnnotations(annotations map[string]string) (ScaleData, bool) {
	result := false
	var data ScaleData
	ScaleUp, ok := annotations["autoscale/up"]
	if ok {
		result = true
		data.ScaleUp = ScaleUp
	}
	ScaleDown, ok := annotations["autoscale/down"]
	if ok {
		result = true
		data.ScaleDown = ScaleDown
	}
	MinReplicas, ok := annotations["autoscale/minReplicas"]
	if ok {
		data.MinReplicas, _ = strconv.Atoi(MinReplicas)
	}
	MaxReplicas, ok := annotations["autoscale/maxReplicas"]
	if ok {
		data.MaxReplicas, _ = strconv.Atoi(MaxReplicas)
	}
	return data, result
}

func checkForScale(ctx context.Context, prometheusClient promv1.API,
	query string) bool {
	data, err := prometheusClient.Query(ctx, query, time.Now())
	if err != nil {
		log.Println(err.Error())
	}
	vector, ok := data.(model.Vector)
	if ok {
		if vector.Len() > 0 {
			return true
		}
	}
	return false
}

func getPrometheusClient(cfg Config) promv1.API {
	prometheusClient, err := prometheus.NewClient(prometheus.Config{Address: cfg.PrometheusAddress})
	if err != nil {
		panic(err.Error())
	}
	return promv1.NewAPI(prometheusClient)
}

func getKubernetesClient(cfg Config) *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func main() {
	var cfg Config
	err := envconfig.Process("autoscale", &cfg)
	if err != nil {
		panic(err.Error())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	limit := rate.Every(time.Duration(cfg.PollingInterval) * time.Second)
	limiter := rate.NewLimiter(limit, 1)

	var wg sync.WaitGroup
	ctx := context.Background()
	prometheusClient := getPrometheusClient(cfg)
	kubernetesClient := getKubernetesClient(cfg)
ScaleLoop:
	for {
		limiter.Wait(ctx)
		select {
		case <-signals:
			break ScaleLoop
		default:
		}
		deployments, err := kubernetesClient.ExtensionsV1beta1().Deployments(cfg.KubernetesNamespace).List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, deployment := range deployments.Items {
			wg.Add(1)
			go func(deployment extv1beta1.Deployment) {
				defer wg.Done()
				scaleData, ok := parseAnnotations(deployment.Annotations)
				if ok {
					// TODO: check that we are not in the middle of scale/deploy process
					currentReplicas := deployment.Spec.Replicas
					// Scale up
					if scaleData.MaxReplicas == 0 {
						scaleData.MaxReplicas = int(*currentReplicas) + 1
					}
					if scaleData.ScaleUp != "" && int(*currentReplicas) < scaleData.MaxReplicas {
						if checkForScale(ctx, prometheusClient, scaleData.ScaleUp) {
							log.Printf("[ns: %s][deployment: %s]: scaling up",
								deployment.Namespace, deployment.Name)
							*deployment.Spec.Replicas++
							_, err := kubernetesClient.ExtensionsV1beta1().Deployments(
								cfg.KubernetesNamespace).Update(&deployment)
							if err != nil {
								log.Printf("[ns: %s][deployment: %s]: %s",
									deployment.Namespace, deployment.Name, err.Error())
							}
							return
						}
					}
					// Scale down
					if scaleData.ScaleDown != "" && int(*currentReplicas) > scaleData.MinReplicas {
						if checkForScale(ctx, prometheusClient, scaleData.ScaleDown) {
							log.Printf("[ns: %s][deployment: %s]: scaling down",
								deployment.Namespace, deployment.Name)
							*deployment.Spec.Replicas--
							_, err := kubernetesClient.ExtensionsV1beta1().Deployments(
								cfg.KubernetesNamespace).Update(&deployment)
							if err != nil {
								log.Printf("[ns: %s][deployment: %s]: %s",
									deployment.Namespace, deployment.Name, err.Error())
							}
							return
						}
					}
				} else {
					log.Printf("[ns: %s][deployment: %s]: no annotations, skipping",
						deployment.Namespace, deployment.Name)
				}
			}(deployment)
		}

		statefulsets, err := kubernetesClient.AppsV1beta1().StatefulSets(cfg.KubernetesNamespace).List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, statefulset := range statefulsets.Items {
			wg.Add(1)
			go func(statefulset appsv1beta1.StatefulSet) {
				defer wg.Done()
				scaleData, ok := parseAnnotations(statefulset.Annotations)
				if ok {
					// TODO: check that we are not in the middle of scale/deploy process
					currentReplicas := statefulset.Spec.Replicas
					// Scale up
					if scaleData.MaxReplicas == 0 {
						scaleData.MaxReplicas = int(*currentReplicas) + 1
					}
					if scaleData.ScaleUp != "" && int(*currentReplicas) < scaleData.MaxReplicas {
						if checkForScale(ctx, prometheusClient, scaleData.ScaleUp) {
							log.Printf("[ns: %s][statefulset: %s]: scaling up",
								statefulset.Namespace, statefulset.Name)
							*statefulset.Spec.Replicas++
							_, err := kubernetesClient.AppsV1beta1().StatefulSets(
								cfg.KubernetesNamespace).Update(&statefulset)
							if err != nil {
								log.Printf("[ns: %s][statefulset: %s]: %s",
									statefulset.Namespace, statefulset.Name, err.Error())
							}
							return
						}
					}
					// Scale down
					if scaleData.ScaleDown != "" && int(*currentReplicas) > scaleData.MinReplicas {
						log.Printf("SCALEDOWN")
						if checkForScale(ctx, prometheusClient, scaleData.ScaleDown) {
							log.Printf("[ns: %s][statefulset: %s]: scaling down",
								statefulset.Namespace, statefulset.Name)
							*statefulset.Spec.Replicas--
							_, err := kubernetesClient.AppsV1beta1().StatefulSets(
								cfg.KubernetesNamespace).Update(&statefulset)
							if err != nil {
								log.Printf("[ns: %s][statefulset: %s]: %s",
									statefulset.Namespace, statefulset.Name, err.Error())
							}
							return
						}
					}
				} else {
					log.Printf("[ns: %s][statefulset: %s]: no annotations, skipping",
						statefulset.Namespace, statefulset.Name)
				}
			}(statefulset)
		}
		wg.Wait()
	}
}
