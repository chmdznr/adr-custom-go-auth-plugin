function callKeycloak(auth, client_id, username, password, server, realm) {
    var keycloakRequest = {
        "Method": "POST",
        "Body": "username=" + username + "&password=" + password + "&grant_type=password&client_id=" + client_id,
        "Headers": {
            "Authorization": auth,
            "Content-Type": "application/x-www-form-urlencoded",
        },
        "Domain": server,
        "Resource": "/auth/realms/" + realm + "/protocol/openid-connect/token"
    };
    var kcResp = TykMakeHttpRequest(JSON.stringify(keycloakRequest));
    rawlog(keycloakRequest);
    return kcResp;
}

function oauthTokenXSWedc(request, session, config) {

    var username = request.Params.username[0];
    var password = request.Params.password[0];
    var grant_type = request.Params.grant_type[0];
    var client_id = request.Params.client_id[0];

    var res_body = "";
    var req_body = "";
    var kcResp = "";
    if (client_id.trim() === "revamp_ms2".trim()) {
        // check login to ad1access
        var ad1accessClientRequest = {
            "Method": "POST",
            "Body": "{\"login\":\"" + username + "\",\"password\":\"" + password + "\"}",
            "Headers": { "content-type": "application/json" },
            "Domain": config.config_data.ad1access_url,
            "Resource": config.config_data.ad1access_login_endpoint
        };
        var ad1accessClientRequestStr = JSON.stringify(ad1accessClientRequest)
        log("ad1accessClientRequest object: " + ad1accessClientRequestStr)

        rawlog("--- before get to upstream ---")
        var ad1accessASResp = TykMakeHttpRequest(ad1accessClientRequestStr);
        rawlog("--- After get to upstream ---")

        log('----')
        var ad1accessASRespObj = JSON.parse(ad1accessASResp);
        var ad1accessASRespCode = JSON.parse(ad1accessASRespObj.Code);
        log('oauthASRespCode: ' + ad1accessASRespCode);
        res_body = ad1accessASRespObj.Body;
        req_body = ad1accessClientRequestStr;

        // continue only if OK
        if (res_body.trim() === "OK") {
            kcResp = callKeycloak(request.Headers.Authorization[0], client_id, username, password, config.config_data.keycloak_server, config.config_data.keycloak_realm_ad1access);
            var kcRespObj = JSON.parse(kcResp);
            var kcRespBodyObj = JSON.parse(kcRespObj.Body);

            rawlog(config);
            // rawlog(keycloakRequest);
            rawlog(kcRespObj);

            if (kcRespObj.Code == 200) {
                var ret = {
                    // "auth": request.Headers.Authorization[0],
                    // "username": username,
                    // "password": password,
                    // "grant_type": grant_type,
                    // "client_id": client_id,
                    // "res_body": res_body,
                    // "req_body": req_body,
                    // "config": JSON.stringify(config),
                    // "kc_resp_code": kcRespObj.Code,
                    // "kc_resp_body": JSON.stringify(kcRespObj.Body),
                    // "kc_req": JSON.stringify(keycloakRequest),
                    "access_token": kcRespBodyObj.access_token,
                    // "refresh_token": kcRespBodyObj.refresh_token,
                    "token_type": kcRespBodyObj.token_type,
                    "expires_in": kcRespBodyObj.expires_in,
                    "scope": kcRespBodyObj.scope
                };
            } else {
                // var ret = {
                //     // "auth": request.Headers.Authorization[0],
                //     // "username": username,
                //     // "password": password,
                //     // "grant_type": grant_type,
                //     // "client_id": client_id,
                //     "res_body": res_body,
                //     "req_body": req_body,
                //     "config": JSON.stringify(config),
                //     "kc_resp_code": kcRespObj.Code,
                //     "kc_resp_body": kcRespObj.Body
                //     // "kc_req": JSON.stringify(keycloakRequest)
                // };

                var kcRespBodyObj = JSON.parse(kcRespObj.Body);

                var ret = {
                    "error": kcRespBodyObj.error,
                    "error_description": kcRespBodyObj.error_description
                }
            }

            var responseObject = {
                Headers: { "Content-Type": "application/json" },
                Body: JSON.stringify(ret),
                Code: kcRespObj.Code
            };
            return TykJsResponse(responseObject, session.meta_data);

        } else {
            var ret = {
                "error": res_body.trim(),
                "error_description": res_body.trim()
            }

            var responseObject = {
                Headers: { "Content-Type": "application/json" },
                Body: JSON.stringify(ret),
                Code: 400
            };
            return TykJsResponse(responseObject, session.meta_data);
        }
    } else if (client_id.trim() === "ad1gate_mobile".trim()) {
        kcResp = callKeycloak(request.Headers.Authorization[0], client_id, username, password, config.config_data.keycloak_server, config.config_data.keycloak_realm_ad1gate);
        var kcRespObj = JSON.parse(kcResp);
        var kcRespBodyObj = JSON.parse(kcRespObj.Body);

        rawlog(config);
        // rawlog(keycloakRequest);
        rawlog(kcRespObj);

        if (kcRespObj.Code == 200) {
            // panggil ad1gate
            var ad1gateClientRequest = {
                "Method": "POST",
                "Body": "{\"uid\":\"" + username + "\",\"pwd\":\"" + password + "\",\"ip\":\"\"}",
                "Headers": { "content-type": "application/json" },
                "Domain": config.config_data.ad1gate_url,
                "Resource": config.config_data.ad1gate_login_endpoint
            };
            var ad1gateClientRequestStr = JSON.stringify(ad1gateClientRequest)
            log("ad1gateClientRequest object: " + ad1gateClientRequestStr)

            rawlog("--- before get to upstream ---")
            var ad1gateASResp = TykMakeHttpRequest(ad1gateClientRequestStr);
            rawlog("--- After get to upstream ---")

            log('----')
            var ad1gateASRespObj = JSON.parse(ad1gateASResp);
            var ad1gateASRespCode = JSON.parse(ad1gateASRespObj.Code);
            log('oauthASRespCode: ' + ad1gateASRespCode);
            res_body = ad1gateASRespObj.Body;
            req_body = ad1gateClientRequestStr;
            res_body_obj = JSON.parse(res_body);

            if (res_body_obj.Status == 0) {
                // sukses
                var ret = {
                    // "res_body": res_body_obj,
                    "access_token": kcRespBodyObj.access_token,
                    // "refresh_token": kcRespBodyObj.refresh_token,
                    "token_type": kcRespBodyObj.token_type,
                    "expires_in": kcRespBodyObj.expires_in,
                    "scope": kcRespBodyObj.scope,
                    "UserDetail": res_body_obj.Data.UserDetail,
                    "DLCName": res_body_obj.Data.DLCName,
                    "LastLogin": res_body_obj.Data.LastLogin
                };
            } else {
                // gagal
                var ret = {
                    "error": kcRespBodyObj.error,
                    "error_description": kcRespBodyObj.error_description
                }
            }


        } else {
            var kcRespBodyObj = JSON.parse(kcRespObj.Body);

            var ret = {
                "error": kcRespBodyObj.error,
                "error_description": kcRespBodyObj.error_description
            }

        }

        var responseObject = {
            Headers: { "Content-Type": "application/json" },
            Body: JSON.stringify(ret),
            Code: kcRespObj.Code
        };
        return TykJsResponse(responseObject, session.meta_data);
    } else {
        kcResp = callKeycloak(request.Headers.Authorization[0], client_id, username, password, config.config_data.keycloak_server, config.config_data.keycloak_realm_default);
        var kcRespObj = JSON.parse(kcResp);
        var kcRespBodyObj = JSON.parse(kcRespObj.Body);

        rawlog(config);
        // rawlog(keycloakRequest);
        rawlog(kcRespObj);

        if (kcRespObj.Code == 200) {
            var ret = {
                // "auth": request.Headers.Authorization[0],
                // "username": username,
                // "password": password,
                // "grant_type": grant_type,
                // "client_id": client_id,
                // "res_body": res_body,
                // "req_body": req_body,
                // "config": JSON.stringify(config),
                // "kc_resp_code": kcRespObj.Code,
                // "kc_resp_body": JSON.stringify(kcRespObj.Body),
                // "kc_req": JSON.stringify(keycloakRequest),
                "access_token": kcRespBodyObj.access_token,
                // "refresh_token": kcRespBodyObj.refresh_token,
                "token_type": kcRespBodyObj.token_type,
                "expires_in": kcRespBodyObj.expires_in,
                "scope": kcRespBodyObj.scope
            };
        } else {
            var kcRespBodyObj = JSON.parse(kcRespObj.Body);

            var ret = {
                "error": kcRespBodyObj.error,
                "error_description": kcRespBodyObj.error_description
            }

        }

        var responseObject = {
            Headers: { "Content-Type": "application/json" },
            Body: JSON.stringify(ret),
            Code: kcRespObj.Code
        };
        return TykJsResponse(responseObject, session.meta_data);
    }

    // var keycloakRequest = {
    //     "Method": "POST",
    //     "Body": "username=" + username + "&password=" + password + "&grant_type=password&client_id=" + client_id,
    //     "Headers": {
    //         "Authorization": request.Headers.Authorization[0],
    //         "Content-Type": "application/x-www-form-urlencoded",
    //     },
    //     "Domain": config.config_data.keycloak_server,
    //     "Resource": "/auth/realms/"+config.config_data.keycloak_realm+"/protocol/openid-connect/token"
    // };
    // var kcResp = TykMakeHttpRequest(JSON.stringify(keycloakRequest));


}