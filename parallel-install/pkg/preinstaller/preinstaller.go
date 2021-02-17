package preinstaller

import (
	"fmt"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/config"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/deployment"
	"github.com/kyma-incubator/hydroform/parallel-install/pkg/logger"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
)

// PreInstaller performs CRDs installation.
type PreInstaller struct {
	cfg            config.Config
	kubeClient     kubernetes.Interface
	processUpdates chan<- deployment.ProcessUpdate
}

func NewPreInstaller(cfg config.Config, kubeClient kubernetes.Interface, processUpdates chan<- deployment.ProcessUpdate) *PreInstaller {
	return &PreInstaller{
		cfg:            cfg,
		kubeClient:     kubeClient,
		processUpdates: processUpdates,
	}
}

func (i *PreInstaller) InstallCRDs() error {
	resource := newCrdPreInstallerResource()
	err := i.apply(*resource)
	if err != nil {
		return err
	}

	return nil
}

func (i *PreInstaller) CreateNamespaces() error {
	resource := newNamespacePreInstallerResource()
	err := i.apply(*resource)
	if err != nil {
		return err
	}

	return nil
}

func (i *PreInstaller) apply(dataType preInstallerResource) error {
	installationResourcePath := i.cfg.InstallationResourcePath
	path := fmt.Sprintf("%s/%s", installationResourcePath, dataType.name)

	components, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	if len(components) == 0 {
		i.cfg.Log.Warn("There were no components detected for installation. Skipping.")
		return nil
	}

	for _, component := range components {
		componentName := component.Name()
		i.cfg.Log.Infof("Processing component: %s", componentName)
		pathToComponent := fmt.Sprintf("%s/%s", path, componentName)
		resources, err := ioutil.ReadDir(pathToComponent)
		if err != nil {
			return err
		}

		if len(resources) == 0 {
			i.cfg.Log.Warnf("There were no resources detected for component: ", componentName)
			return nil
		}

		for _, resource := range resources {
			resourceName := resource.Name()
			i.cfg.Log.Infof("Processing file: %s", resourceName)
			pathToResource := fmt.Sprintf("%s/%s", pathToComponent, resourceName)
			resourceData, err := ioutil.ReadFile(pathToResource)
			if err != nil {
				return err
			}

			err = applyResource(resourceName, i.kubeClient, string(resourceData), dataType.validator, i.cfg.Log)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func applyResource(name string, kubeClient kubernetes.Interface, data string, validator func(string) bool, logger logger.Interface) error {
	if !validator(data) {
		logger.Warnf("Validation failed for resource: %s", name)
	}
	return nil
}
