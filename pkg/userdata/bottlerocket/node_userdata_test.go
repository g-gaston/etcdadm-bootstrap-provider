package bottlerocket

import (
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/aws/etcdadm-bootstrap-provider/api/v1beta1"
	"github.com/aws/etcdadm-bootstrap-provider/pkg/userdata"
)

const (
	userDataMinimum = `
[settings.host-containers.admin]
enabled = true
superpowered = true
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"

[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"

[settings.network]
hostname = ""`

	userDataWithProxyRegistryBootstrapContainers = `
[settings.host-containers.admin]
enabled = true
superpowered = true
source = "custom-admin-image"
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"
[settings.host-containers.control]
enabled = true
superpowered = false
source = "custom-control-image"

[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"

[settings.network]
hostname = ""
https-proxy = "https-proxy"
no-proxy = ["no-proxy-1","no-proxy-2"]

[settings.bootstrap-containers.custom-bootstrap-1]
essential = true
mode = "always"
source = "custom-bootstrap-image-1"
user-data = "abc"
[settings.bootstrap-containers.custom-bootstrap-2]
essential = false
mode = "once"
source = "custom-bootstrap-image-2"
user-data = "xyz"
[settings.container-registry.mirrors]
"public.ecr.aws" = ["https://registry-endpoint"]
[settings.pki.registry-mirror-ca]
data = "Y2FjZXJ0"
trusted=true`

	userDataWithCustomBootstrapContainer = `
[settings.host-containers.admin]
enabled = true
superpowered = true
source = "custom-admin-image"
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"
[settings.host-containers.control]
enabled = true
superpowered = false
source = "custom-control-image"

[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"

[settings.network]
hostname = ""

[settings.bootstrap-containers.custom-bootstrap-1]
essential = true
mode = "always"
source = "custom-bootstrap-image-1"
user-data = "abc"
[settings.bootstrap-containers.custom-bootstrap-2]
essential = false
mode = "once"
source = "custom-bootstrap-image-2"
user-data = "xyz"`

	userDataWithRegistryAuth = `
[settings.host-containers.admin]
enabled = true
superpowered = true
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"

[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"

[settings.network]
hostname = ""
[settings.container-registry.mirrors]
"public.ecr.aws" = ["https://registry-endpoint"]
[settings.pki.registry-mirror-ca]
data = "Y2FjZXJ0"
trusted=true
[[settings.container-registry.credentials]]
registry = "public.ecr.aws"
username = "username"
password = "password"
[[settings.container-registry.credentials]]
registry = "registry-endpoint"
username = "username"
password = "password"`

	userDataWithNTP = `
[settings.host-containers.admin]
enabled = true
superpowered = true
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"

[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"

[settings.network]
hostname = ""
[settings.ntp]
time-servers = ["1.2.3.4", "time-a.capi.com", "time-b.capi.com"]`

	userDataWithHostname = `
[settings.host-containers.admin]
enabled = true
superpowered = true
source = "custom-admin-image"
user-data = "CnsKCSJzc2giOiB7CgkJImF1dGhvcml6ZWQta2V5cyI6IFsic3NoLWtleSJdCgl9Cn0="
[settings.host-containers.kubeadm-bootstrap]
enabled = true
superpowered = true
source = "kubeadm-bootstrap-image"
user-data = "a3ViZWFkbUJvb3RzdHJhcFVzZXJEYXRh"
[settings.host-containers.control]
enabled = true
superpowered = false
source = "custom-control-image"
[settings.kubernetes]
cluster-domain = "cluster.local"
standalone-mode = true
authentication-mode = "tls"
server-tls-bootstrap = false
pod-infra-container-image = "pause-image"
[settings.network]
hostname = "hostname"
https-proxy = "https-proxy"
no-proxy = ["no-proxy-1","no-proxy-2"]
[settings.bootstrap-containers.custom-bootstrap-1]
essential = true
mode = "always"
source = "custom-bootstrap-image-1"
user-data = "abc"
[settings.bootstrap-containers.custom-bootstrap-2]
essential = false
mode = "once"
source = "custom-bootstrap-image-2"
user-data = "xyz"
[settings.container-registry.mirrors]
"public.ecr.aws" = ["https://registry-endpoint"]
[settings.pki.registry-mirror-ca]
data = "Y2FjZXJ0"
trusted=true`
)

func TestGenerateBottlerocketNodeUserData(t *testing.T) {
	g := NewWithT(t)
	trueVal := true

	testcases := []struct {
		name                     string
		kubeadmBootstrapUserData string
		hostname                 string
		users                    []bootstrapv1.User
		registryCredentials      userdata.RegistryMirrorCredentials
		etcdConfig               v1beta1.EtcdadmConfigSpec
		output                   string
	}{
		{
			name:                     "minimum setting",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					PauseImage:     "pause-image",
				},
			},
			output: userDataMinimum,
		},
		{
			name:                     "with custom bootstrap container, with admin and control image",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					AdminImage:     "custom-admin-image",
					ControlImage:   "custom-control-image",
					PauseImage:     "pause-image",
					CustomBootstrapContainers: []v1beta1.BottlerocketBootstrapContainer{
						{
							Name:      "custom-bootstrap-1",
							Image:     "custom-bootstrap-image-1",
							Essential: true,
							Mode:      "always",
							UserData:  "abc",
						},
						{
							Name:      "custom-bootstrap-2",
							Image:     "custom-bootstrap-image-2",
							Essential: false,
							Mode:      "once",
							UserData:  "xyz",
						},
					},
				},
			},
			output: userDataWithCustomBootstrapContainer,
		},
		{
			name:                     "with proxy, registry and custom bootstrap containers",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					AdminImage:     "custom-admin-image",
					ControlImage:   "custom-control-image",
					PauseImage:     "pause-image",
					CustomBootstrapContainers: []v1beta1.BottlerocketBootstrapContainer{
						{
							Name:      "custom-bootstrap-1",
							Image:     "custom-bootstrap-image-1",
							Essential: true,
							Mode:      "always",
							UserData:  "abc",
						},
						{
							Name:      "custom-bootstrap-2",
							Image:     "custom-bootstrap-image-2",
							Essential: false,
							Mode:      "once",
							UserData:  "xyz",
						},
					},
				},
				Proxy: &v1beta1.ProxyConfiguration{
					HTTPProxy:  "http-proxy",
					HTTPSProxy: "https-proxy",
					NoProxy: []string{
						"no-proxy-1",
						"no-proxy-2",
					},
				},
				RegistryMirror: &v1beta1.RegistryMirrorConfiguration{
					Endpoint: "registry-endpoint",
					CACert:   "cacert",
				},
			},
			output: userDataWithProxyRegistryBootstrapContainers,
		},
		{
			name:                     "with registry with authentication",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			registryCredentials: userdata.RegistryMirrorCredentials{
				Username: "username",
				Password: "password",
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					PauseImage:     "pause-image",
				},
				RegistryMirror: &v1beta1.RegistryMirrorConfiguration{
					Endpoint: "registry-endpoint",
					CACert:   "cacert",
				},
			},
			output: userDataWithRegistryAuth,
		},
		{
			name:                     "with NTP config",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					PauseImage:     "pause-image",
				},
				NTP: &bootstrapv1.NTP{
					Enabled: &trueVal,
					Servers: []string{
						"1.2.3.4",
						"time-a.capi.com",
						"time-b.capi.com",
					},
				},
			},
			output: userDataWithNTP,
		},
		{
			name:                     "with proxy, hostname",
			kubeadmBootstrapUserData: "kubeadmBootstrapUserData",
			hostname:                 "hostname",
			users: []bootstrapv1.User{
				{
					SSHAuthorizedKeys: []string{
						"ssh-key",
					},
				},
			},
			etcdConfig: v1beta1.EtcdadmConfigSpec{
				BottlerocketConfig: &v1beta1.BottlerocketConfig{
					BootstrapImage: "kubeadm-bootstrap-image",
					AdminImage:     "custom-admin-image",
					ControlImage:   "custom-control-image",
					PauseImage:     "pause-image",
					CustomBootstrapContainers: []v1beta1.BottlerocketBootstrapContainer{
						{
							Name:      "custom-bootstrap-1",
							Image:     "custom-bootstrap-image-1",
							Essential: true,
							Mode:      "always",
							UserData:  "abc",
						},
						{
							Name:      "custom-bootstrap-2",
							Image:     "custom-bootstrap-image-2",
							Essential: false,
							Mode:      "once",
							UserData:  "xyz",
						},
					},
				},
				Proxy: &v1beta1.ProxyConfiguration{
					HTTPProxy:  "http-proxy",
					HTTPSProxy: "https-proxy",
					NoProxy: []string{
						"no-proxy-1",
						"no-proxy-2",
					},
				},
				RegistryMirror: &v1beta1.RegistryMirrorConfiguration{
					Endpoint: "registry-endpoint",
					CACert:   "cacert",
				},
			},
			output: userDataWithHostname,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			b, err := generateBottlerocketNodeUserData([]byte(testcase.kubeadmBootstrapUserData), testcase.users, testcase.registryCredentials, testcase.hostname, testcase.etcdConfig, logr.New(log.NullLogSink{}))
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(b)).To(Equal(testcase.output))
		})
	}
}
