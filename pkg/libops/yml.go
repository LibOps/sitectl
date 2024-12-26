package libops

type Yml struct {
	Version       float32             `yaml:"version"`
	Php           float32             `yaml:"php"`
	HttpsFirewall []string            `yaml:"https-firewall"`
	SshFirewall   []string            `yaml:"ssh-firewall"`
	BlockedIps    []string            `yaml:"blocked-ips"`
	Developers    map[string][]string `yaml:"developers"`
	MariaDB       string              `yaml:"mariadb,omitempty"`
	Solr          int                 `yaml:"solr,omitempty"`
	DomainMapping map[string][]string `yaml:"domain-mappings,omitempty"`
}
