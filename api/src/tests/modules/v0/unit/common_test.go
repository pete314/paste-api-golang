//Author: Peter Nagy <https://peternagy.ie>
//Since: 07, 2017
//Description: unit tests for function in common package

package unit

import (
	"testing"

	"../../../../modules/v0/common"
	"bytes"
	"encoding/json"
	"net"
	"strings"
)

var (
	config *common.Config
)

//*********************************
//	configuration.go tests
//*********************************

//TestLoadConfig - testing common/configuration::LoadConfig
func TestLoadConfig(t *testing.T) {
	if config = common.LoadConfig("../../../../../../config/server-config.json"); config == nil {
		t.Error("Config did not initialize")
	}
}

//*********************************
//	crypto.go tests
//*********************************
//TestGenRandomHash - testing common/crypto::GenRandomHash
func TestGenRandomHash(t *testing.T) {
	h1 := common.GenRandomHash()
	h2 := common.GenRandomHash()

	if strings.EqualFold(h1, h2) {
		t.Error("common::GenRandomHash() not random, or unique")
	}

	if len(h1) != len(h2) {
		t.Error("common::GenRandomHash() not constant length")
	}
}

//TestGetStringHash - testing common/crypto::GenConfirmHash
func TestGetStringHash(t *testing.T) {
	uname := "Joe"

	h1 := common.GetStringHash(uname)
	h2 := common.GetStringHash(uname)

	if !strings.EqualFold(h1, h2) {
		t.Error("common::GenConfirmHash() not equal")
	}
}

//TestCreateBcryptPassword - testing common/crypto::CreateBcryptPassword
func TestCreateBcryptPassword(t *testing.T) {
	pass := "P@s5w0rd"
	pb1, e1 := common.CreateBcryptPassword(pass)
	if e1 != nil {
		t.Error("common::CreateBcryptPassword() ", e1)
	}

	pb2, e2 := common.CreateBcryptPassword(pass)
	if e2 != nil {
		t.Error("common::CreateBcryptPassword() ", e2)
	}

	if 0 == bytes.Compare(pb1, pb2) {
		t.Error("common::CreateBcryptPassword() passwords hash equal")
	}
}

//TestValidateBcryptPassword - testing common/crypto::ValidateBcryptPassword
func TestValidateBcryptPassword(t *testing.T) {
	pass := "P@s5w0rd"
	pb1, e1 := common.CreateBcryptPassword(pass)
	if e1 != nil {
		t.Error("common::CreateBcryptPassword() ", e1)
	}

	if common.ValidateBcryptPassword(string(pb1), "invalid password") {
		t.Error("common::CreateBcryptPassword() passwords hash equal")
	}

	if !common.ValidateBcryptPassword(string(pb1), pass) {
		t.Error("common::CreateBcryptPassword() passwords hash equal")
	}
}

//TestGetRandomBytes - testing common/crypto::GetRandomBytes
func TestGetRandomBytes(t *testing.T) {
	b1 := common.GetRandomBytes(12)
	b2 := common.GetRandomBytes(12)
	b3 := common.GetRandomBytes(16)

	if 0 == bytes.Compare(b1, b2) {
		t.Error("common::GetRandomBytes() not random, not uniqe")
	}

	if len(b3) != 16 {
		t.Error("common::GetRandomBytes() not valid length")
	}
}

//*********************************
//	redis_client.go tests
//*********************************

//TestGetRedisClient - testing common/redis_client::GetRedisClient
func TestGetRedisClient(t *testing.T) {
	rc1 := common.GetRedisClient()
	rc2 := common.GetRedisClient()

	if rc1 == nil || rc2 == nil {
		t.Error("RedisClient did not initialize")
	}

	if rc1 != rc2 {
		t.Error("GetRedisClient is not singleton")
	}
}

//TestNewRedisClient - testing common/redis_client::NewRedisClient
func TestNewRedisClient(t *testing.T) {
	if nil == common.NewRedisClient(config) {
		t.Error("RedisClient did not initialize")
	}
}

//TestInfo - testing common/redis_client::Info
func TestInfo(t *testing.T) {
	r := common.NewRedisClient(config)
	if _, err := r.Info("server"); err != nil {
		t.Error("Could not get Redis info message")
	}
}

//TestSetKey - testing common/redis_client::SetKey
func TestSetKey(t *testing.T) {
	r := common.NewRedisClient(config)
	data := map[string]interface{}{
		"key": "value",
	}
	b, _ := json.Marshal(data)
	res := r.SetKey("test", b, 3)
	if !res {
		t.Error("Could not create Redis key")
	}
}

//TestGetKey - testing common/redis_client::GetKey
func TestGetKey(t *testing.T) {
	r := common.NewRedisClient(config)

	b := r.GetKey("test")
	if len(b) == 0 {
		t.Error("Could not get Redis key/value")
	}
}

//TestDeleteKey - testing common/redis_client::DeleteKey
func TestDeleteKey(t *testing.T) {
	r := common.NewRedisClient(config)
	data := map[string]interface{}{
		"key": "value",
	}
	b, _ := json.Marshal(data)
	res := r.SetKey("test-del", b, 3)
	if !res {
		t.Error("Could not create Redis key")
	}

	if 1 != r.DeleteKey("test-del") {
		t.Error("redis_client::DeleteKey Could delete key")
	}
}

//TestGetKeyNames - testing common/redis_client::GetKeyNames
func TestGetKeyNames(t *testing.T) {
	r := common.NewRedisClient(config)

	res := r.SetKey("tttest", "val", 3)
	res2 := r.SetKey("tttist", "val", 3)

	if !res || !res2 {
		t.Error("Could not create Redis key")
	}

	keys := r.GetKeyNames("ttt*")

	if keys == nil {
		t.Error("redis_client::GetKeyNames keys regex failed")
	}

	if len(keys) < 2 {
		t.Error("redis_client::GetKeyNames invalid match")
	}
}

//TestPushToList - testing common/redis_client::PushToList
func TestPushToList(t *testing.T) {
	r := common.NewRedisClient(config)

	res := r.PushToList("test-list", "item")
	res1 := r.PushToList("test-list", "item")

	if !res || res != res1 {
		t.Error("redis_client::PushToList failed to push on list")
	}
}

//TestGetListLength - testing common/redis_client::GetListLength
func TestGetListLength(t *testing.T) {
	r := common.NewRedisClient(config)

	res := r.GetListLength("test-list")

	if res < 1 {
		t.Error("redis_client::GetListLength failed to get list length")
	}
}

//TestGetListRange - testing common/redis_client::GetListRange
func TestGetListRange(t *testing.T) {
	r := common.NewRedisClient(config)

	res := r.GetListRange("test-list", 0, -1)

	if len(res) < 2 {
		t.Error("redis_client::GetListRange failed to get list length")
	}

	if !strings.EqualFold(res[0], "item") {
		t.Error("redis_client::GetListRange invalid items in range")
	}

	r.DeleteKey("test-list")
}

//*********************************
//	request.go tests
//*********************************

//TestDecode - testing common/request.go::Decode
func TestDecode(t *testing.T) {
	type TStruct struct {
		Test string
	}

	tmp := map[string]interface{}{"Test": "test"}
	b, _ := json.Marshal(tmp)
	var decTmp TStruct

	if err := common.Decode(b, decTmp); err != nil {
		t.Error("Could not decode json")
	}
}

//*********************************
//	response.go tests
//*********************************

//TestEncode - testing common/response.go::Encode
func TestEncode(t *testing.T) {
	tmp := map[string]interface{}{"Test": "test"}

	if _, err := common.Encode(tmp); err != nil {
		t.Error("Could not encode json")
	}
}

//TestNewError - testing common/response.go::NewError
func TestNewError(t *testing.T) {
	if e := common.NewError("Test/path", "Unit test", 201); e == nil {
		t.Error("Could not create new ErrorBody")
	}
}

//TestNewErrorBuffer - testing common/response.go::NewErrorBuffer
func TestNewErrorBuffer(t *testing.T) {
	if e := common.NewErrorBuffer("Test/path", "Unit test", 201); e == nil {
		t.Error("Could not create new ErrorBody")
	}
}

//TestNewResponseBody - testing common/response.go::NewResponseBody
func TestNewResponseBody(t *testing.T) {
	tmp := map[string]interface{}{"Test": "test"}
	if e := common.NewSuccessBody(true, tmp); e == nil {
		t.Error("Could not create new ResponseBody")
	}
}

//*********************************
//	auth.go tests
//*********************************

//TestValidateHeaderToken - testing common/validate.go::ValidateHeaderToken
func TestValidateHeaderToken(t *testing.T) {
	rc := common.NewRedisClient(config)
	if nil == rc {
		t.Error("RedisClient did not initialize")
	}

	at := common.GenRandomHash()
	token := &common.TokenModel{
		AccessToken:  at,
		RefreshToken: common.GenRandomHash(),
		Expires:      10,
		Scope:        10,
		IP:           net.ParseIP("127.0.0.1"),
	}

	b, _ := common.Encode(token)
	rc.SetKey(at, b, 10)

	header := "Bearer " + at
	if ok, _ := rc.ValidateHeaderToken(header, net.ParseIP("127.0.0.1"), true); !ok {
		t.Error("Token validation error")
	}
}

//TestExtractHeader - testing common/validate.go::ExtractHeader
func TestExtractHeader(t *testing.T) {
	header := "Bearer xx"
	if h := common.ExtractToken(header); len(h) == 0 {
		t.Error("Token validation error")
	}
}

//TestExtractHeader - testing common/auth.go::CheckRefreshToken
func TestRefreshToken(t *testing.T) {
	rc := common.NewRedisClient(config)
	if nil == rc {
		t.Error("RedisClient did not initialize")
	}
	tm := common.CreateTokens("t", net.ParseIP("127.0.0.1"), 0)

	rc.DeleteKey(tm.AccessToken)

	if _, err := common.CheckRefreshToken(tm.AccessToken, tm.RefreshToken, net.ParseIP("127.0.0.1")); err != nil {
		t.Error("Refresh token process error: " + err.Error())
	}

}
