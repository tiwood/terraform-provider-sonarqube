package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateProjectAssociationCreate,
		Read:   resourceSonarqubeQualityGateProjectAssociationRead,
		Delete: resourceSonarqubeQualityGateProjectAssociationDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"projectkey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/select"
	sonarQubeURL.RawQuery = url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateProjectAssociationCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("gateid").(string), d.Get("projectkey").(string))
	d.SetId(id)
	return nil
}

func resourceSonarqubeQualityGateProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/search"
	sonarQubeURL.RawQuery = url.Values{
		"gateId": []string{d.Get("gateid").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityGateProjectAssociationRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateAssociationReadResponse := GetQualityGateAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateAssociationReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationRead: Failed to decode json into struct")
	}

	// ID is in format <gateid>/<projectkey>. This splits the id into gateid and projectkey
	// EG: "1/my_project" >> ["1", "my_project"]
	idSlice := strings.Split(d.Id(), "/")

	for _, value := range qualityGateAssociationReadResponse.Results {
		if idSlice[1] == value.Key {
			d.Set("gateid", idSlice[0])
			d.Set("projectkey", value.Key)
		}
	}

	return nil
}

func resourceSonarqubeQualityGateProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/deselect"
	sonarQubeURL.RawQuery = url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateProjectAssociationDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
