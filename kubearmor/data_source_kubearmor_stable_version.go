package kubearmor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubearmorStableVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubearmorStableVersionRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubearmorStableVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	url := "https://raw.githubusercontent.com/kubearmor/KubeArmor/main/STABLE-RELEASE"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching the file:", err)
		return diag.FromErr(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch the file. Status code:", response.StatusCode)
		return diag.FromErr(err)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading the file:", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] content : %#v", string(content))

	d.Set("version", string(content))
	idsum := sha256.New()

	id := fmt.Sprintf("%x", idsum.Sum(nil))
	d.SetId(id)

	return nil
}
