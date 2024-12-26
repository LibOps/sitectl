package libops

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type TlsService struct {
	Url         string       `json:"url"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

type TcpService struct {
	Host        string       `json:"host"`
	Name        string       `json:"name,omitempty"`
	Port        int          `json:"port"`
	Credentials *Credentials `json:"credentials"`
}

type ConnectionInfo struct {
	Blazegraph *TlsService `json:"blazegraph,omitempty"`
	Database   *TcpService `json:"database"`
	Drupal     *TlsService `json:"drupal,omitempty"`
	Fcrepo     *TlsService `json:"fcrepo,omitempty"`
	Iiif       *TlsService `json:"iiif,omitempty"`
	Matomo     *TlsService `json:"matomo,omitempty"`
	Solr       *TlsService `json:"solr,omitempty"`
	Ssh        *TcpService `json:"ssh,omitempty"`
}
