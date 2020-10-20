package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kyma-incubator/hydroform/installation-poc/pkg/helm"
	"sigs.k8s.io/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var resources = "/Users/i304607/Yaas/go/src/github.com/kyma-project/kyma/resources"
var overridesFile = "/Users/i304607/overrides.yaml"
var kubeconfig = "/Users/i304607/Downloads/mst.yml"
var commonListOpts = metav1.ListOptions{LabelSelector: "installer=overrides"}

type Engine struct {
	componentsProvider ComponentsProvider
	//values.yaml
}

func NewEngine(componentsProvider ComponentsProvider) *Engine {
	return &Engine{
		componentsProvider: componentsProvider,
	}
}

type Components struct{}

type ComponentsProvider interface {
	GetComponents() []Component
}

type Installation interface {
	Install(components []Component) error
}

type Component struct {
	Name              string
	Namespace         string
	OverridesProvider OverridesProvider
	HelmClient        helm.ClientInterface
}

type ComponentInstallation interface {
	InstallComponent() error
}

type Overrides struct{}

type OverridesProvider interface {
	OverridesFor() map[string]interface{}
}



func (c *Component) InstallComponent() error {
	chartDir := path.Join(resources, c.Name)
	log.Printf("MST Installing %s in %s from %s", c.Name, c.Namespace, chartDir)

	overrides := c.OverridesProvider.OverridesFor()

	err := c.HelmClient.InstallRelease(chartDir, c.Namespace, c.Name, overrides)
	if err != nil {
		log.Printf("MST Error installing %s: %v", c.Name, err)
		return err
	}

	log.Printf("MST Installed %s in %s", c.Name, c.Namespace)
	return nil

	return nil
}

func (c *Components) GetComponents() []Component {
	overridesProvider := &Overrides{}
	helmClient := &helm.Client{}
	return []Component{
		Component{
			Name:              "istio-kyma-patch",
			Namespace:         "istio-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "knative-serving",
			Namespace:         "knative-serving",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "knative-eventing",
			Namespace:         "knative-eventing",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "dex",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "ory",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "api-gateway",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "rafter",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "service-catalog",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "service-catalog-addons",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "nats-streaming",
			Namespace:         "natss",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "core",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "cluster-users",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "permission-controller",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "apiserver-proxy",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "iam-kubeconfig-service",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "serverless",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "knative-provisioner-natss",
			Namespace:         "knative-eventing",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "event-sources",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "application-connector",
			Namespace:         "kyma-integration",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
		Component{
			Name:              "console",
			Namespace:         "kyma-system",
			OverridesProvider: overridesProvider,
			HelmClient:helmClient,
		},
	}
}

func (o *Overrides) OverridesFor() map[string]interface{} {
	config, err := getClientConfig(kubeconfig)
	if err != nil {
		log.Fatalf("Unable to build kubernetes configuration. Error: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to create internal client. Error: %v", err)
	}

	configmaps, err := kubeClient.CoreV1().ConfigMaps("kyma-installer").List(context.TODO(), commonListOpts)
	overrides := make(map[string]interface{})

	for _, cm := range configmaps.Items {
		log.Printf("%s data %v", cm.Name, cm.Data)

		yamlData, err := yaml.Marshal(cm.Data)

		//save to .yaml
		tmpFile, err := ioutil.TempFile(os.TempDir(), cm.Name+"-")
		if err != nil {
			log.Fatal("Cannot create temporary file", err)
		}

		fmt.Println("Created File: " + tmpFile.Name())
		defer os.Remove(tmpFile.Name())

		if _, err = tmpFile.Write(yamlData); err != nil {
			log.Fatal("Failed to write to temporary file", err)
		}

		// Close the file
		if err := tmpFile.Close(); err != nil {
			log.Fatal(err)
		}

		//read from file
		var data map[string]interface{}
		bs, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(bs, &data); err != nil {
			panic(err)
		}

		for k, v := range data {
			overrides[k] = v
		}

	}

	unflatten := unflattenToMap(overrides)

	//save to .yaml
	unflattenData, err := yaml.Marshal(unflatten)
	if err := ioutil.WriteFile(overridesFile, unflattenData, 0644); err != nil {
		panic(err)
	}

	return unflatten
}

func installPrerequisites() error {
	overridesProvider := &Overrides{}
	helmClient := &helm.Client{}

	clusterEssentials := &Component{
		Name:      "cluster-essentials",
		Namespace: "kyma-system",
		OverridesProvider:overridesProvider,
		HelmClient:helmClient,
	}
	err := clusterEssentials.InstallComponent()
	if err != nil {
		return err
	}

	istio := &Component{
		Name:      "istio",
		Namespace: "istio-system",
		OverridesProvider:overridesProvider,
		HelmClient:helmClient,
	}
	err = istio.InstallComponent()
	if err != nil {
		return err
	}

	xipPatch := &Component{
		Name:      "xip-patch",
		Namespace: "kyma-installer",
		OverridesProvider:overridesProvider,
		HelmClient:helmClient,
	}
	err = xipPatch.InstallComponent()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Install() error {
	err := installPrerequisites()
	if err != nil {
		return err
	}

	components := e.componentsProvider.GetComponents()

	//Install the rest of the components
	jobChan := make(chan Component, 30)
	for _, comp := range components {
		if !enqueueJob(comp, jobChan) {
			log.Printf("Max capacity reached, component dismissed: %s", comp.Name)
		}
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobChan)
	}

	// to stop the workers, first close the job channel
	close(jobChan)
	wait(&wg, 10*time.Minute)
	cancel()

	return nil
}

func main() {
	componentsProvider := &Components{}

	engine := NewEngine(componentsProvider)

	err := engine.Install()
	if err != nil {
		log.Fatalf("Kyma installation fialed. Error: %v", err)
	}

	fmt.Println("Kyma installed")
}

func unflattenToMap(sourceMap map[string]interface{}) map[string]interface{} {
	mergedMap := map[string]interface{}{}
	if len(sourceMap) == 0 {
		return mergedMap
	}

	for key, value := range sourceMap {
		keys := strings.Split(key, ".")
		mergeIntoMap(keys, value.(string), mergedMap)
	}

	return mergedMap
}

func mergeIntoMap(keys []string, value string, dstMap map[string]interface{}) {
	currentKey := keys[0]
	//Last key points directly to string value
	if len(keys) == 1 {

		//Conversion to boolean to satisfy Helm requirements.yaml: "enable:true/false syntax"
		var vv interface{} = value
		if value == "true" || value == "false" {
			vv, _ = strconv.ParseBool(value)
		}

		dstMap[currentKey] = vv
		return
	}

	//All keys but the last one should point to a nested map
	nestedMap, isMap := dstMap[currentKey].(map[string]interface{})

	if !isMap {
		nestedMap = map[string]interface{}{}
		dstMap[currentKey] = nestedMap
	}

	mergeIntoMap(keys[1:], value, nestedMap)
}

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobChan <-chan Component) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobChan:
			if ctx.Err() != nil || !ok {
				return
			}
			if ok {
				job.InstallComponent()
			}
		}
	}
}

func wait(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		log.Println("Timeout occurred!")
		return false
	}
}

func enqueueJob(job Component, jobChan chan<- Component) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}


