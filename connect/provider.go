package connect

import (
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	kc "github.com/jacoelho/go-kafka-connect/lib/connectors"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"insecure_skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},

		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"kafka-connect_connector": kafkaConnectorResource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c := kc.NewClient(d.Get("url").(string))

	skipVerify := d.Get("insecure_skip_verify").(bool)

	if skipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		c = c.WithHTTPClient(client)
	}

	return c, nil
}
