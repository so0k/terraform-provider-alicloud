package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/cdn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnDomainCreate,
		Read:   resourceAlicloudCdnDomainRead,
		Update: resourceAlicloudCdnDomainUpdate,
		Delete: resourceAlicloudCdnDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainName,
			},
			"cdn_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCdnType,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCdnSourceType,
			},
			"source_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validateCdnSourcePort,
			},
			"sources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 20,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCdnScope,
			},

			// configs
			"optimize_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCdnEnable,
			},
			"page_compress_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCdnEnable,
			},
			"range_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCdnEnable,
			},
			"video_seek_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCdnEnable,
			},
			"block_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"parameter_filter_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "off",
							ValidateFunc: validateCdnEnable,
						},
						"hash_key_args": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateCdnHashKeyArg,
							},
							MaxItems: 10,
						},
					},
				},
				MaxItems: 1,
			},

			"page_404_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"page_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "default",
							ValidateFunc: validateCdnPage404Type,
						},
						"custom_page_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 1,
			},

			"refer_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"refer_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "block",
							ValidateFunc: validateCdnReferType,
						},
						"refer_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allow_empty": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "on",
							ValidateFunc: validateCdnEnable,
						},
					},
				},
				MaxItems: 1,
			},

			"auth_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "no_auth",
							ValidateFunc: validateCdnAuthType,
						},
						"master_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateCdnAuthKey,
						},
						"slave_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateCdnAuthKey,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1800,
						},
					},
				},
				MaxItems: 1,
			},

			"http_header_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateCdnHttpHeader,
						},
						"header_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"header_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				MaxItems: 10,
			},

			"cache_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"cache_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateCacheType,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validateIntegerInRange(1, 99),
						},
						"cache_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cdn.AddDomainRequest{
		DomainName: d.Get("domain_name").(string),
		CdnType:    d.Get("cdn_type").(string),
		SourcePort: d.Get("source_port").(int),
	}

	if v, ok := d.GetOk("scope"); ok {
		args.Scope = v.(string)
	}

	if args.CdnType != cdn.LiveStream {
		if v, ok := d.GetOk("sources"); ok && v.(*schema.Set).Len() > 0 {
			sources := expandStringList(v.(*schema.Set).List())
			args.Sources = strings.Join(sources, ",")
		} else {
			return fmt.Errorf("Sources is required when 'cdn_type' is not 'liveStream'.")
		}

		if v, ok := d.GetOk("source_type"); ok && v.(string) != "" {
			args.SourceType = v.(string)
		} else {
			return fmt.Errorf("SourceType is required when 'cdn_type' is not 'liveStream'.")
		}
	}
	_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.AddCdnDomain(args)
	})
	if err != nil {
		return fmt.Errorf("AddCdnDomain got an error: %#v", err)
	}

	d.SetId(args.DomainName)

	err = WaitForDomainStatus(d.Id(), Configuring, 60, meta)
	if err != nil {
		return fmt.Errorf("Timeout when Cdn Domain Available")
	}

	return resourceAlicloudCdnDomainUpdate(d, meta)
}

func resourceAlicloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args := cdn.ModifyDomainRequest{
		DomainName: d.Id(),
		SourceType: d.Get("source_type").(string),
	}

	if !d.IsNewResource() {
		attributeUpdate := false
		if d.HasChange("source_type") {
			d.SetPartial("source_type")
			attributeUpdate = true
		}
		if d.HasChange("sources") {
			d.SetPartial("sources")
			sources := expandStringList(d.Get("sources").(*schema.Set).List())
			args.Sources = strings.Join(sources, ",")
			attributeUpdate = true
		}
		if d.HasChange("source_port") {
			d.SetPartial("source_port")
			args.SourcePort = d.Get("source_port").(int)
			attributeUpdate = true
		}
		if attributeUpdate {
			_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
				return cdnClient.ModifyCdnDomain(args)
			})
			if err != nil {
				return fmt.Errorf("ModifyCdnDomain got an error: %#v", err)
			}
		}
	}

	// set optimize_enable 、range_enable、page_compress_enable and video_seek_enable
	if err := enableConfigUpdate(client, d); err != nil {
		return err
	}

	if d.HasChange("block_ips") {
		d.SetPartial("block_ips")
		blockIps := expandStringList(d.Get("block_ips").(*schema.Set).List())
		args := cdn.IpBlackRequest{DomainName: d.Id(), BlockIps: strings.Join(blockIps, ",")}
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetIpBlackListConfig(args)
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("parameter_filter_config") {
		if err := queryStringConfigUpdate(client, d); err != nil {
			return err
		}
	}

	if d.HasChange("page_404_config") {
		if err := page404ConfigUpdate(client, d); err != nil {
			return err
		}
	}

	if d.HasChange("refer_config") {
		if err := referConfigUpdate(client, d); err != nil {
			return err
		}
	}

	if d.HasChange("auth_config") {
		if err := authConfigUpdate(client, d); err != nil {
			return err
		}
	}

	if d.HasChange("http_header_config") {
		if err := httpHeaderConfigUpdate(client, d); err != nil {
			return err
		}
	}

	if d.HasChange("cache_config") {
		if err := cacheConfigUpdate(client, d); err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceAlicloudCdnDomainRead(d, meta)
}

func resourceAlicloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	domain, err := DescribeDomainDetail(d.Id(), meta)
	if err != nil {
		return fmt.Errorf("DescribeDomainDetail got an error: %#v", err)
	}
	d.Set("domain_name", domain.DomainName)
	d.Set("sources", domain.Sources.Source)
	d.Set("cdn_type", domain.CdnType)
	d.Set("source_type", domain.SourceType)
	d.Set("scope", domain.Scope)

	// get domain configs
	describeConfigArgs := cdn.DomainConfigRequest{
		DomainName: d.Id(),
	}
	raw, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.DescribeDomainConfigs(describeConfigArgs)
	})
	if err != nil {
		return fmt.Errorf("DescribeDomainConfigs got an error: %#v", err)
	}
	resp, _ := raw.(cdn.DomainConfigResponse)
	configs := resp.DomainConfigs

	queryStringConfig := configs.IgnoreQueryStringConfig
	if _, ok := d.GetOk("parameter_filter_config"); ok {
		config := make([]map[string]interface{}, 1)
		config[0] = map[string]interface{}{
			"enable":        queryStringConfig.Enable,
			"hash_key_args": strings.Split(queryStringConfig.HashKeyArgs, ","),
		}
		d.Set("parameter_filter_config", config)
	}

	errorPageConfig := configs.ErrorPageConfig
	if _, ok := d.GetOk("page_404_config"); ok {
		config := make([]map[string]interface{}, 1)
		config[0] = map[string]interface{}{
			"page_type":       errorPageConfig.PageType,
			"error_code":      errorPageConfig.ErrorCode,
			"custom_page_url": errorPageConfig.CustomPageUrl,
		}
		if errorPageConfig.PageType == "" {
			config[0]["page_type"] = "default"
		}
		d.Set("page_404_config", config)
	}

	referConfig := configs.RefererConfig
	if _, ok := d.GetOk("refer_config"); ok {
		config := make([]map[string]interface{}, 1)
		config[0] = map[string]interface{}{
			"refer_type":  referConfig.ReferType,
			"refer_list":  strings.Split(referConfig.ReferList, ","),
			"allow_empty": referConfig.AllowEmpty,
		}
		d.Set("refer_config", config)
	}

	authConfig := configs.ReqAuthConfig
	if _, ok := d.GetOk("auth_config"); ok {
		config := make([]map[string]interface{}, 1)
		timeout, _ := strconv.Atoi(authConfig.TimeOut)
		config[0] = map[string]interface{}{
			"auth_type":  authConfig.AuthType,
			"master_key": authConfig.Key1,
			"slave_key":  authConfig.Key2,
			"timeout":    timeout,
		}
		d.Set("auth_config", config)
	}

	headerConfigs := configs.HttpHeaderConfigs.HttpHeaderConfig
	httpHeaderConfigs := make([]map[string]interface{}, 0, len(headerConfigs))
	for _, v := range headerConfigs {
		val := make(map[string]interface{})
		val["header_key"] = v.HeaderKey
		val["header_value"] = v.HeaderValue
		val["header_id"] = v.ConfigId
		httpHeaderConfigs = append(httpHeaderConfigs, val)
	}
	d.Set("http_header_config", httpHeaderConfigs)

	cacheConfigs := configs.CacheExpiredConfigs.CacheExpiredConfig
	cacheExpiredConfigs := make([]map[string]interface{}, 0, len(cacheConfigs))
	for _, v := range cacheConfigs {
		val := make(map[string]interface{})
		ttl, _ := strconv.Atoi(v.TTL)
		weight, _ := strconv.Atoi(v.Weight)
		val["cache_type"] = v.CacheType
		val["cache_content"] = v.CacheContent
		val["cache_id"] = v.ConfigId
		val["weight"] = weight
		val["ttl"] = ttl
		cacheExpiredConfigs = append(cacheExpiredConfigs, val)
	}
	d.Set("cache_config", cacheExpiredConfigs)

	d.Set("optimize_enable", configs.OptimizeConfig.Enable)
	d.Set("page_compress_enable", configs.PageCompressConfig.Enable)
	d.Set("range_enable", configs.RangeConfig.Enable)
	d.Set("video_seek_enable", configs.VideoSeekConfig.Enable)
	blocks := make([]string, 0)
	if len(configs.CcConfig.BlockIps) > 0 {
		blocks = strings.Split(configs.CcConfig.BlockIps, ",")
	}
	d.Set("block_ips", blocks)

	return nil
}

func resourceAlicloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cdn.DescribeDomainRequest{
		DomainName: d.Id(),
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(args)
		})
		if err != nil {
			if IsExceptedError(err, ServiceBusy) {
				return resource.RetryableError(fmt.Errorf("The specified Domain is configuring, please retry later."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting cdn domain %s: %#v.", d.Id(), err))
		}
		return nil
	})
}

func enableConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	type configFunc func(req cdn.ConfigRequest) (cdn.CdnCommonResponse, error)

	raw, _ := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return map[string]configFunc{
			"optimize_enable":      cdnClient.SetOptimizeConfig,
			"range_enable":         cdnClient.SetRangeConfig,
			"page_compress_enable": cdnClient.SetPageCompressConfig,
			"video_seek_enable":    cdnClient.SetVideoSeekConfig,
		}, nil
	})
	relation, _ := raw.(map[string]configFunc)

	for key, fn := range relation {
		if d.HasChange(key) {
			d.SetPartial(key)
			args := cdn.ConfigRequest{
				DomainName: d.Id(),
				Enable:     d.Get(key).(string),
			}
			if _, err := fn(args); err != nil {
				return err
			}
		}
	}
	return nil
}

func queryStringConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	valSet := d.Get("parameter_filter_config").(*schema.Set)
	args := cdn.QueryStringConfigRequest{DomainName: d.Id()}

	if valSet == nil || valSet.Len() == 0 {
		args.Enable = "off"
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetIgnoreQueryStringConfig(args)
		})
		if err != nil {
			return err
		}
		return nil
	}

	val := valSet.List()[0].(map[string]interface{})
	d.SetPartial("parameter_filter_config")
	args.Enable = val["enable"].(string)
	if v, ok := val["hash_key_args"]; ok && len(v.([]interface{})) > 0 {
		hashKeyArgs := expandStringList(v.([]interface{}))
		args.HashKeyArgs = strings.Join(hashKeyArgs, ",")
	}
	_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.SetIgnoreQueryStringConfig(args)
	})
	if err != nil {
		return err
	}
	return nil
}

func page404ConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	valSet := d.Get("page_404_config").(*schema.Set)
	args := cdn.ErrorPageConfigRequest{DomainName: d.Id()}

	if valSet == nil || valSet.Len() == 0 {
		args.PageType = "default"
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetErrorPageConfig(args)
		})
		if err != nil {
			return err
		}
		return nil
	}

	val := valSet.List()[0].(map[string]interface{})
	d.SetPartial("page_404_config")
	args.PageType = val["page_type"].(string)
	customPageUrl, ok := val["custom_page_url"]
	if ok {
		args.CustomPageUrl = customPageUrl.(string)
	}

	if args.PageType == "charity" && (ok && customPageUrl.(string) != CharityPageUrl || !ok) {
		return fmt.Errorf("If 'page_type' value is 'charity', you must set 'custom_page_url' with '%s'.", CharityPageUrl)
	}
	if args.PageType == "default" && ok && customPageUrl.(string) != "" {
		return fmt.Errorf("If 'page_type' value is 'default', you can not set 'custom_page_url'.")
	}
	if args.PageType == "other" && (!ok || customPageUrl.(string) == "") {
		return fmt.Errorf("If 'page_type' value is 'other', you must set the value of 'custom_page_url'.")
	}

	_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.SetErrorPageConfig(args)
	})
	if err != nil {
		return err
	}
	return nil
}

func referConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	valSet := d.Get("refer_config").(*schema.Set)
	args := cdn.ReferConfigRequest{DomainName: d.Id()}

	if valSet == nil || valSet.Len() == 0 {
		args.ReferType = "block"
		args.AllowEmpty = "on"
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetRefererConfig(args)
		})
		if err != nil {
			return err
		}
		return nil
	}

	val := valSet.List()[0].(map[string]interface{})
	d.SetPartial("refer_config")
	args.ReferType = val["refer_type"].(string)
	args.AllowEmpty = val["allow_empty"].(string)
	if v, ok := val["refer_list"]; ok && len(v.([]interface{})) > 0 {
		referList := expandStringList(v.([]interface{}))
		args.ReferList = strings.Join(referList, ",")
	}
	_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.SetRefererConfig(args)
	})
	if err != nil {
		return err
	}
	return nil
}

func authConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	ov, nv := d.GetChange("auth_config")
	oldConfig, newConfig := ov.(*schema.Set), nv.(*schema.Set)
	args := cdn.ReqAuthConfigRequest{DomainName: d.Id()}

	if newConfig == nil || newConfig.Len() == 0 {
		args.AuthType = "no_auth"
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetReqAuthConfig(args)
		})
		if err != nil {
			return err
		}
		return nil
	}

	val := newConfig.List()[0].(map[string]interface{})
	d.SetPartial("auth_config")
	args.AuthType = val["auth_type"].(string)
	args.Timeout = strconv.Itoa(val["timeout"].(int))

	masterKey, okMasterKey := val["master_key"]
	slaveKey, okSlaveKey := val["slave_key"]
	if okMasterKey {
		args.Key1 = masterKey.(string)
	}
	if okSlaveKey {
		args.Key2 = slaveKey.(string)
	}

	if args.AuthType == "no_auth" {
		if oldConfig == nil || oldConfig.Len() == 0 {
			if okMasterKey || okSlaveKey {
				return fmt.Errorf("If 'auth_type' value is 'no_auth', you can not set the value of 'master_key' and 'slave_key'.")
			}
		} else {
			oldVal := oldConfig.List()[0].(map[string]interface{})
			if oldVal["master_key"] != val["master_key"] || oldVal["slave_key"] != val["slave_key"] {
				return fmt.Errorf("If 'auth_type' value is 'no_auth', you can not change the value of 'master_key' and 'slave_key'.")
			}
		}
	} else {
		if !okMasterKey || !okSlaveKey {
			return fmt.Errorf("If 'auth_type' value is one of ['type_a', 'type_b', 'type_c'], you must set 'master_key' and 'slave_key' at one time.")
		}
	}

	_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.SetReqAuthConfig(args)
	})
	if err != nil {
		return err
	}
	return nil
}

func httpHeaderConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	ov, nv := d.GetChange("http_header_config")
	oldConfigs := ov.(*schema.Set).List()
	newConfigs := nv.(*schema.Set).List()

	for _, v := range oldConfigs {
		configId := v.(map[string]interface{})["header_id"].(string)
		args := cdn.DeleteHttpHeaderConfigRequest{
			DomainName: d.Id(),
			ConfigID:   configId,
		}
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DeleteHttpHeaderConfig(args)
		})
		if err != nil {
			return err
		}
	}

	if len(newConfigs) == 0 {
		return nil
	}

	for _, v := range newConfigs {
		args := cdn.HttpHeaderConfigRequest{
			DomainName:  d.Id(),
			HeaderKey:   v.(map[string]interface{})["header_key"].(string),
			HeaderValue: v.(map[string]interface{})["header_value"].(string),
		}
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetHttpHeaderConfig(args)
		})
		if err != nil {
			return fmt.Errorf("SetHttpHeaderConfig got an error: %#v", err)
		}
	}

	return nil
}

func cacheConfigUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	ov, nv := d.GetChange("cache_config")
	oldConfigs := ov.(*schema.Set).List()
	newConfigs := nv.(*schema.Set).List()

	for _, v := range oldConfigs {
		val := v.(map[string]interface{})
		configId := val["cache_id"].(string)
		args := cdn.DeleteCacheConfigRequest{
			DomainName: d.Id(),
			ConfigID:   configId,
			CacheType:  val["cache_type"].(string),
		}
		_, err := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.DeleteCacheExpiredConfig(args)
		})
		if err != nil {
			return fmt.Errorf("DeleteCacheExpiredConfig got an error: %#v", err)
		}
	}

	if len(newConfigs) == 0 {
		return nil
	}

	for _, v := range newConfigs {
		val := v.(map[string]interface{})
		args := cdn.CacheConfigRequest{
			DomainName:   d.Id(),
			CacheContent: val["cache_content"].(string),
			TTL:          strconv.Itoa(val["ttl"].(int)),
			Weight:       strconv.Itoa(val["weight"].(int)),
		}
		if err := setCacheExpiredConfig(args, val["cache_type"].(string), client); err != nil {
			return err
		}
	}

	return nil
}

func setCacheExpiredConfig(req cdn.CacheConfigRequest, cacheType string, client *connectivity.AliyunClient) (err error) {
	if cacheType == "suffix" {
		_, err = client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetFileCacheExpiredConfig(req)
		})
	} else {
		_, err = client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
			return cdnClient.SetPathCacheExpiredConfig(req)
		})
	}
	return
}

func DescribeDomainDetail(Id string, meta interface{}) (domain cdn.DomainDetail, err error) {
	client := meta.(*connectivity.AliyunClient)

	args := cdn.DescribeDomainRequest{
		DomainName: Id,
	}
	raw, e := client.WithCdnClient(func(cdnClient *cdn.CdnClient) (interface{}, error) {
		return cdnClient.DescribeCdnDomainDetail(args)
	})
	if e != nil {
		err = fmt.Errorf("DescribeCdnDomainDetail got an error: %#v", e)
		return
	}
	response, _ := raw.(cdn.DomainResponse)
	domain = response.GetDomainDetailModel
	return
}

func WaitForDomainStatus(Id string, status Status, timeout int, meta interface{}) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		domain, err := DescribeDomainDetail(Id, meta)
		if err != nil {
			return err
		}
		if domain.DomainStatus == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Domain", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
