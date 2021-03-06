package rds

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// UpgradeDBInstanceNetWorkInfo invokes the rds.UpgradeDBInstanceNetWorkInfo API synchronously
// api document: https://help.aliyun.com/api/rds/upgradedbinstancenetworkinfo.html
func (client *Client) UpgradeDBInstanceNetWorkInfo(request *UpgradeDBInstanceNetWorkInfoRequest) (response *UpgradeDBInstanceNetWorkInfoResponse, err error) {
	response = CreateUpgradeDBInstanceNetWorkInfoResponse()
	err = client.DoAction(request, response)
	return
}

// UpgradeDBInstanceNetWorkInfoWithChan invokes the rds.UpgradeDBInstanceNetWorkInfo API asynchronously
// api document: https://help.aliyun.com/api/rds/upgradedbinstancenetworkinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpgradeDBInstanceNetWorkInfoWithChan(request *UpgradeDBInstanceNetWorkInfoRequest) (<-chan *UpgradeDBInstanceNetWorkInfoResponse, <-chan error) {
	responseChan := make(chan *UpgradeDBInstanceNetWorkInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpgradeDBInstanceNetWorkInfo(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// UpgradeDBInstanceNetWorkInfoWithCallback invokes the rds.UpgradeDBInstanceNetWorkInfo API asynchronously
// api document: https://help.aliyun.com/api/rds/upgradedbinstancenetworkinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpgradeDBInstanceNetWorkInfoWithCallback(request *UpgradeDBInstanceNetWorkInfoRequest, callback func(response *UpgradeDBInstanceNetWorkInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpgradeDBInstanceNetWorkInfoResponse
		var err error
		defer close(result)
		response, err = client.UpgradeDBInstanceNetWorkInfo(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// UpgradeDBInstanceNetWorkInfoRequest is the request struct for api UpgradeDBInstanceNetWorkInfo
type UpgradeDBInstanceNetWorkInfoRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	ConnectionString     string           `position:"Query" name:"ConnectionString"`
}

// UpgradeDBInstanceNetWorkInfoResponse is the response struct for api UpgradeDBInstanceNetWorkInfo
type UpgradeDBInstanceNetWorkInfoResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpgradeDBInstanceNetWorkInfoRequest creates a request to invoke UpgradeDBInstanceNetWorkInfo API
func CreateUpgradeDBInstanceNetWorkInfoRequest() (request *UpgradeDBInstanceNetWorkInfoRequest) {
	request = &UpgradeDBInstanceNetWorkInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "UpgradeDBInstanceNetWorkInfo", "rds", "openAPI")
	return
}

// CreateUpgradeDBInstanceNetWorkInfoResponse creates a response to parse from UpgradeDBInstanceNetWorkInfo response
func CreateUpgradeDBInstanceNetWorkInfoResponse() (response *UpgradeDBInstanceNetWorkInfoResponse) {
	response = &UpgradeDBInstanceNetWorkInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
