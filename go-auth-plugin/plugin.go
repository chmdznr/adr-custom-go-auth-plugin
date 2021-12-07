package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TykTechnologies/tyk/headers"
	"github.com/TykTechnologies/tyk/user"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/golang-jwt/jwt"
)

var oidcJwtAzp = ""
var kcRealmPublicKey = ""
var logger = log.Get()

// AddFooBarHeader adds custom "Foo: Bar" header to the request
//func AddFooBarHeader(rw http.ResponseWriter, r *http.Request) {
//	r.Header.Add("Foo", "Bar")
//}

func main() {
	//sampleJwt := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJHakJFMlNFMWtPYVNsU2xDaWNWZkNWcm1PQ1FWUUdJdlBfMHZMQ2lRT2pFIn0.eyJleHAiOjE2MzY0NzY1NDMsImlhdCI6MTYzNjQwNDU0MywianRpIjoiNGU5N2RjNzUtYzMzNS00ODg0LWIwZTUtMTVkNWNkMzkwZGMxIiwiaXNzIjoiaHR0cDovLzEwLjIuMjkuNTI6MjA4MC9hdXRoL3JlYWxtcy9tb2JzZWMiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiYmJmMTI1NDgtOWQxNi00OTJmLWE5MTMtZDUwNGE4YjAxMDllIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYWQxZ2F0ZV9tb2JpbGUiLCJzZXNzaW9uX3N0YXRlIjoiMTBjY2YwZTctNGU3MC00ZmVmLWEyZDUtNmM2Mzg0YjFkOTJlIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIiwiaHR0cDovL2h0dHBiaW4ub3JnIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLW1vYnNlYyIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InJlYWQgd3JpdGUiLCJzaWQiOiIxMGNjZjBlNy00ZTcwLTRmZWYtYTJkNS02YzYzODRiMWQ5MmUifQ.FgwAbK3-RMGnwUjXjdNEUp2jbEebd9mnP3Eyc1ju53zZF6jqu14Iu35txpQraEKpV6aqE3sNPc1NzQdqXerdy8J2gxE2Pn_6Y43ePSuN2rLAfwfBOuh6qzgvcSrey351o4M4MhTzR3YdJt40-12E-FDRfsPDzmZ71GB5Fw4DF7jkfbJmuB5RXmm7vybdoAt0Jh2bOjWj9VjmjZx1DWUYRhqyCsZb4S228jyfo3uvLza3w5WFV-A-FEoW9IQu3IkUNVTGiSLBiKcx9PQfL5BHIPWFVy2qKwU0CKgyn7bCBydiGRbKMHjrGKs8oYBmuVfb9Ek7y2retUWkjCQpbqwpLw"

	//sampleJwt := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYzNjMyMDI3NCwiZXhwIjoxNjM2MzIzODc0fQ.gSHUZZ9sYbaMOOPn8iwmpu42MKsH3SK2siVEPEv71O2Kt3ouDx88YG43LEC4ycMez1RPyQ0ExQqlZau-Kn1W67g4WkITuFsT7Lhg-k4DAtl-C2rIbqwxhDhXO-DrECfX41PApWdQxSTtE0IPrHG9IxJerjZ5kSlUoW65IfSaHi8pbuDaXBpZt0umMmL9Ym-gNuEowW7weKMuwi9x33MotI4fsN7cyP6uNu5CYtRIIzVzkOxgx2aDyMoHfBi-YOWIuBu3zw_vQkgAM60nDJoujkDkAnxA9c-yeU1LS_g3ImYnKykhdaRwXq7RxxYm9ujlgo_b78hMgovg-Q5Fzr06FQ"
	//

	//
	//finalPublicKey := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6S7asUuzq5Q/3U9rbs+P\nkDVIdjgmtgWreG5qWPsC9xXZKiMV1AiV9LXyqQsAYpCqEDM3XbfmZqGb48yLhb/X\nqZaKgSYaC/h2DjM7lgrIQAp9902Rr8fUmLN2ivr5tnLxUUOnMOc2SQtr9dgzTONY\nW5Zu3PwyvAWk5D6ueIUhLtYzpcB+etoNdL3Ir2746KIy/VUsDwAM7dhrqSK8U2xF\nCGlau4ikOTtvzDownAMHMrfE7q1B6WZQDAQlBmxRQsyKln5DIsKv6xauNsHRgBAK\nctUxZG8M4QJIx3S6Aughd3RZC4Ca5Ae9fd8L8mlNYBCrQhOZ7dS0f4at4arlLcaj\ntwIDAQAB\n-----END PUBLIC KEY-----")
	//

	//isExpired, err := isExpiredToken(sampleJwt)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if isExpired {
	//	fmt.Println("The token is expired")
	//} else {
	//	fmt.Println("The token is not expired")
	//}
	//
	//sampleConfigData := "{\n  \"allowed_clients\": [\n    {\n      \"name\": \"ad1gate_mobile\",\n      \"public_key\": \"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAltnb4lSe2Y9ia8vfep3pW7mgXb1U8oIs9pVJTiZp0P5xNaPjLAwo2yDpNY4pb4HLndfKBvDvh2e7CYa/BttN+mrd/CKuu8YRi1JeMdt2VMEP45o5xQ5aoP0TWVaQMJIIt+rXgLi/6DPS6HWmooHcj/X36FPpDJSDcvisp3Pr7fCpWoK295lsgVQUFMfDh+HRGPTkWCAC1Qu34SaoIAVDlLfrhCMC6yU48dORt2+8mZZcuRpJyjnJs/epuRpH0MlsNAefWccdSbA37PtPitXbWzGNjvo2W/LNkvz1zorOvoIHNZh1O2OKBdh+v5dhXFlkfMPU4yYoyr4BMGGwzQKgtwIDAQAB\"\n    }\n  ]\n}"
	//
	//configData := make(map[string]interface{})
	//err = json.Unmarshal([]byte(sampleConfigData), &configData)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Println(configData)
	//for _, allowedClient := range configData["allowed_clients"].([]interface{}) {
	//	//fmt.Println(i2)
	//	ac := allowedClient.(map[string]interface{})
	//	//fmt.Println(ac["name"])
	//	if ac["name"].(string) == oidcJwtAzp {
	//		// found matching client_id
	//		kcRealmPublicKey = ac["public_key"].(string)
	//	}
	//}
	//
	//if len(kcRealmPublicKey) <= 0 {
	//	fmt.Println("client_id not allowed to access this")
	//	return
	//}
	//
	//finalPublicKey := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", kcRealmPublicKey)
	//isValid, err := isVerifiedToken(sampleJwt, finalPublicKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if isValid {
	//	fmt.Println("The token is valid")
	//} else {
	//	fmt.Println("The token is invalid")
	//}

	//fmt.Println("Pass")

}

func AdiraCustomGoAuthPlugin1626(w http.ResponseWriter, r *http.Request) {
	authHeaders := r.Header.Get(headers.Authorization)

	//fmt.Printf("%v", authHeaders)
	logger.Infof("%v", authHeaders)

	// check header sanity
	authParts := strings.Split(authHeaders, " ")
	if len(authParts) < 2 || !strings.Contains(authHeaders, "earer") {
		replyData := map[string]interface{}{
			"error":             "invalid authorization",
			"error_description": "invalid authorization header",
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonReplyData)
		return
	}

	sampleJwt := authParts[1]

	//fmt.Println("Mulai cek expired token")
	logger.Println("Mulai cek expired token")
	isExpired, err := isExpiredToken(sampleJwt)
	if err != nil {
		replyData := map[string]interface{}{
			"error":             "invalid token expiration check",
			"error_description": err,
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write(jsonReplyData)
		return
		//log.Fatal(err)
	}
	logger.Info("Lolos cek expired token")

	if isExpired {
		//fmt.Println("The token is expired")
		logger.Error("The token is expired")
		replyData := map[string]interface{}{
			"error":             "token expired",
			"error_description": "Token is expired",
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonReplyData)
		return
	}

	// get config_data from API definition
	apidef := ctx.GetDefinition(r)
	//fmt.Println("API name is", apidef.Name)
	logger.Info("API name is", apidef.Name)
	configData := apidef.ConfigData

	for _, allowedClient := range configData["allowed_clients"].([]interface{}) {
		//fmt.Println(i2)
		ac := allowedClient.(map[string]interface{})
		logger.Info(ac)
		//fmt.Println(ac["name"])
		if ac["name"].(string) == oidcJwtAzp {
			// found matching client_id
			kcRealmPublicKey = ac["public_key"].(string)
		}
	}

	if len(kcRealmPublicKey) <= 0 {
		//fmt.Println("client_id not allowed to access this")
		logger.Info("client_id not allowed to access this")
		replyData := map[string]interface{}{
			"error":             "invalid client_id",
			"error_description": "client_id not allowed to access this",
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonReplyData)
		return
	}

	finalPublicKey := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", kcRealmPublicKey)
	isValid, err := isVerifiedToken(sampleJwt, finalPublicKey)
	if err != nil {
		replyData := map[string]interface{}{
			"error":             "Token validation error",
			"error_description": err,
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonReplyData)
		return
		//log.Fatal(err)
	}

	if !isValid {
		//fmt.Println("The token is invalid")
		logger.Info("The token is invalid")
		replyData := map[string]interface{}{
			"error":             "invalid token",
			"error_description": "Token is invalid",
		}

		jsonReplyData, _ := json.Marshal(replyData)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonReplyData)
		return
	}

	//session:=getDefaultSession()

	// auth was successful, add session and key to request's context so other middlewares can use it
	//ctx.SetSession(r, session ,true)
	// Let the request continue
	//fmt.Println("Auth passed")
	logger.Info("Auth passed")
}

func getDefaultSession() *user.SessionState {
	// return session
	return &user.SessionState{
		OrgID: "default",
		Alias: "custom-auth-session",
	}
}

func isExpiredToken(tokenStr string) (bool, error) {
	logger.Info("Masuk isExpiredToken")
	var token *jwt.Token
	parser := jwt.Parser{UseJSONNumber: true}
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	logger.Info("Lolos ParseUnverified")

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//fmt.Println(claims["exp"])
		if claims["azp"] != nil {
			oidcJwtAzp = claims["azp"].(string)
		}

		var expiredAt int64
		expiredAt, err = claims["exp"].(json.Number).Int64()
		if err != nil {
			//fmt.Println("failed to convert expiration time")
			//fmt.Println(err)
			return true, err
		} else {
			if expiredAt < time.Now().Unix() {
				// expired
				//fmt.Println("Token expired")
				return true, nil
			} else {
				//fmt.Println("Token not expired")
				return false, nil
			}
		}

	} else {
		//fmt.Println(err)
		return true, err
	}
}

func isVerifiedToken(token, publicKey string) (bool, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return false, err
	}

	parts := strings.Split(token, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		fmt.Println("Verify error\n", err)
		return false, nil
	}
	return true, nil
}
