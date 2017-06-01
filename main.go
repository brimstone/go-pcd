package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

type ConfigType struct {
	Hostname string     `json:"hostname"`
	API      string     `json:"api"`
	Docker   dockerType `json:"docker"`
	Files    []struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	} `json:"files"`
	Command string `json:"command"`
}

type dockerType struct {
	APICorsHeader        string        `json:"api-cors-header,omitempty"`
	AuthorizationPlugins []interface{} `json:"authorization-plugins,omitempty"`
	Bip                  string        `json:"bip,omitempty"`
	Bridge               string        `json:"bridge,omitempty"`
	CgroupParent         string        `json:"cgroup-parent,omitempty"`
	ClusterStore         string        `json:"cluster-store,omitempty"`
	ClusterStoreOpts     struct {
	} `json:"cluster-store-opts,omitempty"`
	ClusterAdvertise      string        `json:"cluster-advertise,omitempty"`
	Debug                 bool          `json:"debug,omitempty"`
	DefaultGateway        string        `json:"default-gateway,omitempty"`
	DefaultGatewayV6      string        `json:"default-gateway-v6,omitempty"`
	DefaultRuntime        string        `json:"default-runtime,omitempty"`
	DisableLegacyRegistry bool          `json:"disable-legacy-registry,omitempty"`
	DNS                   []interface{} `json:"dns,omitempty"`
	DNSOpts               []interface{} `json:"dns-opts,omitempty"`
	DNSSearch             []interface{} `json:"dns-search,omitempty"`
	ExecOpts              []interface{} `json:"exec-opts,omitempty"`
	ExecRoot              string        `json:"exec-root,omitempty"`
	FixedCidr             string        `json:"fixed-cidr,omitempty"`
	FixedCidrV6           string        `json:"fixed-cidr-v6,omitempty"`
	Graph                 string        `json:"graph,omitempty"`
	Group                 string        `json:"group,omitempty"`
	Hosts                 []string      `json:"hosts,omitempty"`
	Icc                   bool          `json:"icc,omitempty"`
	InsecureRegistries    []interface{} `json:"insecure-registries,omitempty"`
	IP                    string        `json:"ip,omitempty"`
	Iptables              bool          `json:"iptables,omitempty"`
	Ipv6                  bool          `json:"ipv6,omitempty"`
	IPForward             bool          `json:"ip-forward,omitempty"`
	IPMasq                bool          `json:"ip-masq,omitempty"`
	Labels                []interface{} `json:"labels,omitempty"`
	LiveRestore           bool          `json:"live-restore,omitempty"`
	LogDriver             string        `json:"log-driver,omitempty"`
	LogLevel              string        `json:"log-level,omitempty"`
	LogOpts               struct {
	} `json:"log-opts,omitempty"`
	MaxConcurrentDownloads    int           `json:"max-concurrent-downloads,omitempty"`
	MaxConcurrentUploads      int           `json:"max-concurrent-uploads,omitempty"`
	Mtu                       int           `json:"mtu,omitempty"`
	OomScoreAdjust            int           `json:"oom-score-adjust,omitempty"`
	Pidfile                   string        `json:"pidfile,omitempty"`
	RawLogs                   bool          `json:"raw-logs,omitempty"`
	RegistryMirrors           []interface{} `json:"registry-mirrors,omitempty"`
	SelinuxEnabled            bool          `json:"selinux-enabled,omitempty"`
	StorageDriver             string        `json:"storage-driver,omitempty"`
	StorageOpts               []interface{} `json:"storage-opts,omitempty"`
	SwarmDefaultAdvertiseAddr string        `json:"swarm-default-advertise-addr,omitempty"`
	TLS                       bool          `json:"tls,omitempty"`
	Tlscacert                 string        `json:"tlscacert,omitempty"`
	Tlscert                   string        `json:"tlscert,omitempty"`
	Tlskey                    string        `json:"tlskey,omitempty"`
	Tlsverify                 bool          `json:"tlsverify,omitempty"`
	UserlandProxy             bool          `json:"userland-proxy,omitempty"`
	UsernsRemap               string        `json:"userns-remap,omitempty"`
}

type initFunc struct {
	Status bool
	Func   func() bool
}

var (
	COMMITHASH    = "dev"
	BUILDDATETIME = "today"
	API_VERSION   = "1"
	BASE_URL      = "127.0.0.1:8080"
	MyAPIGet      func(string) string
	MyAPIPost     func(string, string)
	MyExec        func(string, ...string) ([]byte, error)
	MyReadFile    func(string) ([]byte, error)
	MyWriteFile   func(string, []byte, os.FileMode) error
	cmds          []*cobra.Command
	configfile    = "/boot/config.yaml"
	inits         = make(map[string]*initFunc)
	listener      net.Listener
	forever       chan bool
	config        ConfigType
)

func RealReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func RealWriteFile(filename string, contents []byte, mode os.FileMode) error {
	return ioutil.WriteFile(filename, contents, mode)
}

func RealExec(cmd string, arg ...string) ([]byte, error) {
	return exec.Command(cmd, arg...).CombinedOutput()
}

func getUrl(url string) *http.Response {
	for range []int{1, 2, 3} {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil
		}
		req.Header.Set("Metadata-Flavor", "Google")
		resp, err := client.Do(req)
		if err == nil {
			return resp
		}
		log.Println("Error getting config from url:", err)
		log.Println("Waiting to retry")
		time.Sleep(time.Second * 2)
	}
	return nil
}

func readKernelConfig() error {
	processUrl("http://169.254.169.254/latest/user-data", &config)
	processUrl("http://metadata.google.internal/computeMetadata/v1/instance/attributes/startup-script", &config)
	cmdline, err := MyReadFile("/proc/cmdline")
	if err != nil {
		return err
	}
	kernel := make(map[string]string)
	options := strings.Split(strings.TrimSpace(string(cmdline)), " ")
	for _, option := range options {
		kv := strings.SplitN(option, "=", 2)
		if len(kv) < 2 {
			continue
		}

		if kv[0] == "pcd.url" || kv[0] == "url" {
			log.Println("Got a url", kv[1])
			err := processUrl(kv[1], &config)
			if err != nil {
				log.Println(err)
				continue
			}

		} else if kv[0][0:4] == "pcd." {
			kernel[kv[0]] = kv[1]
		} else if kv[0] == "hostname" {
			kernel[kv[0]] = kv[1]
		}
	}
	for k, v := range kernel {
		if k == "hostname" {
			config.Hostname = v
		}
	}
	config.Docker.Hosts = append(config.Docker.Hosts, "unix:///var/run/docker.sock")
	config.Docker.StorageDriver = "overlay2"

	return nil
}

func processUrl(url string, config *ConfigType) error {
	resp := getUrl(url)
	if resp == nil {
		return fmt.Errorf("Unable to fetch URL: %s", url)
	}
	configcontents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error getting config from url: %s", err)
	}
	err = yaml.Unmarshal(configcontents, &config)
	if err != nil {
		return fmt.Errorf("Error parsing config: %s", err)
	}
	return nil
}

func saveConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err != nil {
		return err
	}
	b, _ := yaml.Marshal(config)
	err = MyWriteFile(configfile, b, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Config saved.")
	MyExec("umount", "/boot")
	return nil
}

func readConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err != nil {
		return nil
	}
	// we have our defaults, kernel config, any url config. Now apply what's on disk
	// The file on disk shouldn't overwrite anything we already have.
	configcontents, err := ioutil.ReadFile(configfile)
	err = yaml.Unmarshal(configcontents, &config)
	_, err = MyExec("umount", "/boot")
	if err != nil {
		return err
	}
	saveConfig()
	return nil
}

func init() {
	MyAPIGet = APIGet
	MyAPIPost = APIPost
	MyReadFile = RealReadFile
	MyWriteFile = RealWriteFile
	MyExec = RealExec
}

func main() {
	var rootCmd = &cobra.Command{
		Use:  os.Args[0],
		Long: "Pancake Crop Deli Control Program",
	}

	for cmd := range cmds {
		rootCmd.AddCommand(cmds[cmd])
	}
	rootCmd.PersistentFlags().StringVarP(&BASE_URL, "address", "a", "127.0.0.1:8080", "Address for API server")
	rootCmd.Execute()
}
