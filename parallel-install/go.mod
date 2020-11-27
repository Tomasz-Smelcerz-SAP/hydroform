module github.com/kyma-incubator/hydroform/parallel-install

go 1.14

replace github.com/hashicorp/consul/api v1.3.0 => github.com/hashicorp/consul/api v0.0.0-20191112221531-8742361660b6

//commit 8742361660b63923954d58bd7bfea4e19b0041ad (HEAD, tag: api/v1.3.0)
//Author: Mike Morris <mikemorris@users.noreply.github.com>
//Date:   Tue Nov 12 17:15:31 2019 -0500
//
//    api: bump consul/sdk version

require (
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/kyma-project/kyma/components/kyma-operator v0.0.0-20201125092745-687c943ac940
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	helm.sh/helm/v3 v3.3.4
	k8s.io/api v0.18.9
	k8s.io/apimachinery v0.18.9
	k8s.io/cli-runtime v0.18.9
	k8s.io/client-go v0.18.9
)
