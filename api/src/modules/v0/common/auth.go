//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: --
package common

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net"
	"strings"
)

//TokenModel - token struct
type TokenModel struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Expires      int    `json:"expires"`
	Scope        int    `json:"scope"`
	UID          string `json:"user_id,omitempty"`
	IP           net.IP `json:"ip,omitempty"`
}

//WithAuthentication - create response handler with authentication
func WithAuthentication(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		if uid, err := validateBearer(ctx, false); err != nil {
			JsonResponse(ctx, 403, NewError("api/auth/validate/token", err.Error(), 403))
		} else {
			ctx.Request.Header.Add("uid", fmt.Sprintf("%s", uid))
			h(ctx)
		}
	})
}

//validateBearer - request with Authorization header
func validateBearer(ctx *fasthttp.RequestCtx, isAdmin bool) (string, error) {
	rc := GetRedisClient()

	if ok, uid := rc.ValidateHeaderToken(string(ctx.Request.Header.Peek("Authorization")), ctx.RemoteIP(), isAdmin); ok {
		return uid, nil
	}

	return "", errors.New("Invalid token, or not enough permission")
}

//ValidateHeaderToken - validate access token
func (rc *RedisClient) ValidateHeaderToken(header string, ip net.IP, isAdmin bool) (bool, string) {
	token := ExtractToken(header)
	if len(token) > 0 {
		if tBytes := rc.GetKey(token); len(tBytes) > 0 {
			tm := &TokenModel{}
			if err := Decode(tBytes, tm); err == nil && strings.EqualFold(tm.AccessToken, token) && ip.Equal(tm.IP) {
				if isAdmin && tm.Scope < 10 {
					return false, ""
				}
				return true, tm.UID
			} else {
				CheckError("Authorization: token", err, false)
			}
		}
	}

	return false, ""
}

//ExtractToken - from Authorization header
func ExtractToken(header string) string {
	const prefix = "Bearer "
	if len(header) > 0 && strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}

	return ""
}

//CreateTokens - create auth tokes
func CreateTokens(uid string, ip net.IP, scope int) *TokenModel {
	rc := GetRedisClient()
	at := &TokenModel{
		AccessToken:  GenRandomHash(),
		RefreshToken: GenRandomHash(),
		Expires:      10800,
		Scope:        scope,
		UID:          uid,
		IP:           ip,
	}
	atB, err := Encode(*at)

	if err == nil {
		resAt := rc.SetKey(at.AccessToken, atB, 10800)
		resRt := rc.SetKey(at.RefreshToken, atB, 1219600)
		if resAt && resRt {
			at.IP = nil
			return at
		}
	}

	CheckError("Error while encoding", err, false)
	return nil
}

//CheckRefreshToken - check refresh token
func CheckRefreshToken(accessToken, refreshToken string, ip net.IP) (*TokenModel, error) {
	rc := GetRedisClient()
	if atb := rc.GetKey(accessToken); len(atb) > 0 {
		//This is only possible if did not expire, so let use it
		result := &TokenModel{}
		Decode(atb, result)
		result.IP = nil
		return result, nil
	}

	if rtb := rc.GetKey(refreshToken); len(rtb) > 0 {
		dbToken := &TokenModel{}
		err := Decode(rtb, dbToken)
		if err == nil && strings.EqualFold(dbToken.AccessToken, accessToken) && ip.Equal(dbToken.IP) {
			//Create new token
			result := CreateTokens(dbToken.UID, ip, dbToken.Scope)
			if result != nil {
				//Not really necessary to check the result as ttl will kill it anyway
				rc.DeleteKey(refreshToken)

				return result, nil
			}

			return nil, errors.New("Could not create token")
		}

		CheckError("common::auth::CheckRefreshToken: invalid token", err, false)
	}

	return nil, errors.New("Invalid token sent")
}
