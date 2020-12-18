package metaforce

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)

const (
	DefaultApiVersion = "50.0"
	DefaultLoginUrl   = "login.salesforce.com"
)

type Client struct {
	ApiVersion string
	ServerUrl  string
	LoginUrl   string
	portType    *MetadataPortType
	loginResult *LoginResult
	debug       bool
}

func NewClient() *Client {
	portType := NewMetadataPortType("", true, nil)
	return &Client{
		portType: portType,
		LoginUrl: DefaultLoginUrl,
		ApiVersion: DefaultApiVersion,
	}
}

func NewDebugClient() *Client {
	portType := NewMetadataPortType("", true, nil)
	return &Client{
		portType: portType,
		LoginUrl: DefaultLoginUrl,
		ApiVersion: DefaultApiVersion,
		debug: true,
	}
}

func (c *Client) SetDebug(debug bool) {
	c.portType.SetDebug(debug)
}

func (c *Client) SetApiVersion(v string) {
	c.ApiVersion = v
	c.setLoginUrl()
}

func (c *Client) SetAccessToken(sid string) {
	sessionHeader := &SessionHeader{
		SessionId: sid,
	}
	c.portType.SetHeader(sessionHeader)
}

func (c *Client) SetLoginUrl(url string) {
	c.LoginUrl = url
	c.setLoginUrl()
}

func (c *Client) setLoginUrl() {
	url := fmt.Sprintf("https://%s/services/Soap/u/%s", c.LoginUrl, c.ApiVersion)
	c.portType.SetServerUrl(url)
}

func (c *Client) SetLogger(logger io.Writer) {
	c.portType.SetLogger(logger)
}

func (c *Client) SetGzip(gz bool) {
	c.portType.SetGzip(gz)
}

//func (c *Client) Logout() error {
//	_, err := c.portType.Logout(&soapforce.Logout{})
//	if err != nil {
//		return err
//	}
//	c.ServerUrl = ""
//	c.setLoginUrl()
//	c.portType.ClearHeader()
//	return nil
//}

func (c *Client) Login(username string, password string) error {
	loginRequest := LoginRequest{Username: username, Password: password}
	loginResponse, err := c.portType.Login(&loginRequest)
	if err != nil {
		return err
	}
	c.loginResult = &loginResponse.LoginResult
	sessionHeader := SessionHeader{
		SessionId: c.loginResult.SessionId,
	}
	c.portType.SetHeader(&sessionHeader)
	c.portType.SetServerUrl(c.loginResult.MetadataServerUrl)
	return nil
}

func (c *Client) UseExistingSession(SessionId, MetadataServerUrl string) {
	sessionHeader := SessionHeader{
		SessionId: SessionId,
	}
	c.portType.SetHeader(&sessionHeader)
	c.portType.SetServerUrl(MetadataServerUrl)
}

func (c *Client) Deploy(buf []byte, options *DeployOptions) (*DeployResponse, error) {
	request := Deploy{
		ZipFile:       base64.StdEncoding.EncodeToString(buf),
		DeployOptions: options,
	}
	return c.portType.Deploy(&request)
}

func (c *Client) CheckDeployStatus(resultId string, includeDetails bool) (*CheckDeployStatusResponse, error) {
	request := CheckDeployStatus{AsyncProcessId: ID(resultId), IncludeDetails: includeDetails}
	return c.portType.CheckDeployStatus(&request)
}

func (c *Client) CancelDeploy(processId string) (*CancelDeployResponse, error) {
	request := CancelDeploy{AsyncProcessId: ID(processId)}
	return c.portType.CancelDeploy(&request)
}

func (c *Client) DescribeMetadata() (*DescribeMetadataResponse, error) {
	f, err := strconv.ParseFloat(c.ApiVersion, 32)
	if err != nil {
		f = 37.0
	}

	request := DescribeMetadata{AsOfVersion: f}
	return c.portType.DescribeMetadata(&request)
}

func (c *Client) DescribeValueType(desc_type string) (*DescribeValueTypeResponse, error) {
	request := DescribeValueType{
		Type: desc_type,
	}
	return c.portType.DescribeValueType(&request)
}

func (c *Client) ListMetadata(listMetadataQuery []*ListMetadataQuery) (*ListMetadataResponse, error) {
	f, err := strconv.ParseFloat(c.ApiVersion, 32)
	if err != nil {
		f = 37.0
	}

	request := ListMetadata{
		Queries: listMetadataQuery,
		AsOfVersion: f,
	}
	return c.portType.ListMetadata(&request)
}

func (c *Client) CreateMetadata(metadata []MetadataInterface) (*CreateMetadataResponse, error) {
	request := CreateMetadata{
		Metadata: metadata,
	}
	return c.portType.CreateMetadata(&request)
}

func (c *Client) DeleteMetadata(typeName string, fullNames []string) (*DeleteMetadataResponse, error) {
	request := DeleteMetadata{
		FullNames: fullNames,
		Type: typeName,
	}
	return c.portType.DeleteMetadata(&request)
}

func (c *Client) ReadMetadata(typeName string, fullNames []string) (*ReadMetadataResponse, error) {
	request := ReadMetadata{
		FullNames: fullNames,
		Type: typeName,
	}
	return c.portType.ReadMetadata(&request)
}

func (c *Client) Retrieve(retrieveRequest *RetrieveRequest) (*RetrieveResponse, error) {
	r := &Retrieve{
		RetrieveRequest: retrieveRequest,
	}
	return c.portType.Retrieve(r)
}

func (c *Client) RenameMetadata(r *RenameMetadata) (*RenameMetadataResponse, error) {
	return c.portType.RenameMetadata(r)
}

func (c *Client) UpdateMetadata(metadata []MetadataInterface) (*UpdateMetadataResponse, error) {
	return c.portType.UpdateMetadata(&UpdateMetadata{Metadata: metadata})
}

func (c *Client) UpsertMetadata(metadata []MetadataInterface) (*UpsertMetadataResponse, error) {
	return c.portType.UpsertMetadata(&UpsertMetadata{Metadata: metadata})
}

func (c *Client) DeployRecentValidation(validationId string) (*DeployRecentValidationResponse, error) {
	return c.portType.DeployRecentValidation(&DeployRecentValidation{
		ValidationId: ID(validationId),
	})
}