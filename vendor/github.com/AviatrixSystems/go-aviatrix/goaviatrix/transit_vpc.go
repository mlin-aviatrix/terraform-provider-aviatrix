package goaviatrix

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// Gateway simple struct to hold gateway details
type TransitVpc struct {
	AccountName            string `form:"account_name,omitempty" json:"account_name,omitempty"`
	Action                 string `form:"action,omitempty"`
	CID                    string `form:"CID,omitempty"`
	CloudType              int    `form:"cloud_type,omitempty" json:"cloud_type,omitempty"`
	DnsServer              string `form:"dns_server,omitempty" json:"dns_server,omitempty"`
	GwName                 string `form:"gw_name,omitempty" json:"vpc_name,omitempty"`
	GwSize                 string `form:"gw_size,omitempty"`
	VpcID                  string `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	Subnet                 string `form:"public_subnet,omitempty" json:"vpc_net,omitempty"`
	HASubnet               string `form:"ha_subnet,omitempty"`
	PeeringHASubnet        string `json:"public_subnet,omitempty"`
	VpcRegion              string `form:"region,omitempty" json:"vpc_region,omitempty"`
	VpcSize                string `form:"gw_size,omitempty" json:"gw_size,omitempty"`
	EnableNAT              string `form:"nat_enabled,omitempty" json:"enable_nat,omitempty"`
	TagList                string `form:"tags,omitempty"`
	EnableHybridConnection bool   `form:"enable_hybrid_connection" json:"tgw_enabled,omitempty"`
	ConnectedTransit       string `form:"connected_transit" json:"connected_transit,omitempty"`
}

func (c *Client) LaunchTransitVpc(gateway *TransitVpc) error {
	gateway.CID = c.CID
	gateway.Action = "create_transit_gw"
	resp, err := c.Post(c.baseURL, gateway)
	if err != nil {
		return errors.New("HTTP Post create_transit_gw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode create_transit_gw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API create_transit_gw Post failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableHaTransitVpc(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_transit_ha ") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "enable_transit_ha")
	attachSpokeToTransitGw.Add("gw_name", gateway.GwName)
	attachSpokeToTransitGw.Add("public_subnet", gateway.HASubnet)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	//path := c.baseURL + fmt.Sprintf("?CID=%s&action=enable_transit_ha&gw_name=%s&public_subnet=%s", c.CID,
	//	gateway.GwName, gateway.HASubnet)
	//resp, err := c.Get(path, nil)
	if err != nil {
		return errors.New("HTTP Get enable_transit_ha failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_transit_ha failed: " + err.Error())
	}
	if !data.Return {
		if strings.Contains(data.Reason, "HA GW already exists") {
			log.Printf("[INFO] HA is already enabled %s", data.Reason)
			return nil
		}
		log.Printf("[ERROR] Enabling HA failed with error %s", data.Reason)

		return errors.New("Rest API enable_transit_ha Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) AttachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_transit_gateway_interface_to_aws_tgw ") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "enable_transit_gateway_interface_to_aws_tgw")
	attachSpokeToTransitGw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	//path := c.baseURL + fmt.Sprintf("?action=enable_transit_gateway_interface_to_aws_tgw&CID=%s&gateway_name=%s",
	//	c.CID, gateway.GwName)
	//resp, err := c.Get(path, nil)
	if err != nil {
		return errors.New("HTTP Get enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		if strings.Contains(err.Error(), "already enabled tgw interface") {
			return nil
		}
		return errors.New("Json Decode enable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DetachTransitGWForHybrid(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_transit_gateway_interface_to_aws_tgw") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "disable_transit_gateway_interface_to_aws_tgw")
	attachSpokeToTransitGw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	//path := c.baseURL + fmt.Sprintf("?action=disable_transit_gateway_interface_to_aws_tgw&CID=%s&gateway_name=%s",
	//	c.CID, gateway.GwName)
	//resp, err := c.Get(path, nil)
	if err != nil {
		return errors.New("HTTP Get disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_transit_gateway_interface_to_aws_tgw failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_transit_gateway_interface_to_aws_tgw Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) EnableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for enable_connected_transit_on_gateway") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "enable_connected_transit_on_gateway")
	attachSpokeToTransitGw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	//path := c.baseURL + fmt.Sprintf("?CID=%s&action=enable_connected_transit_on_gateway&gateway_name=%s",
	//	c.CID, gateway.GwName)
	//resp, err := c.Get(path, nil)
	if err != nil {
		return errors.New("HTTP Get enable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode enable_connected_transit_on_gateway failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API enable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) DisableConnectedTransit(gateway *TransitVpc) error {
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New(("url Parsing failed for disable_connected_transit_on_gateway") + err.Error())
	}
	attachSpokeToTransitGw := url.Values{}
	attachSpokeToTransitGw.Add("CID", c.CID)
	attachSpokeToTransitGw.Add("action", "disable_connected_transit_on_gateway")
	attachSpokeToTransitGw.Add("gateway_name", gateway.GwName)
	Url.RawQuery = attachSpokeToTransitGw.Encode()
	resp, err := c.Get(Url.String(), nil)

	//path := c.baseURL + fmt.Sprintf("?CID=%s&action=disable_connected_transit_on_gateway&gateway_name=%s",
	//	c.CID, gateway.GwName)
	//resp, err := c.Get(path, nil)
	if err != nil {
		return errors.New("HTTP Get disable_connected_transit_on_gateway failed: " + err.Error())
	}
	var data APIResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.New("Json Decode disable_connected_transit_on_gateway failed: " + err.Error())
	}
	if !data.Return {
		return errors.New("Rest API disable_connected_transit_on_gateway Get failed: " + data.Reason)
	}
	return nil
}
