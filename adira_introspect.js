function callKeycloakIntrospect(auth, access_token, server, realm) {
    var keycloakRequest = {
        "Method": "POST",
        "Body": "token=" + access_token + "&token_type_hint=access_token",
        "Headers": {
            "Authorization": auth,
            "Content-Type": "application/x-www-form-urlencoded",
        },
        "Domain": server,
        "Resource": "/auth/realms/" + realm + "/protocol/openid-connect/token/introspect"
    };
    var kcResp = TykMakeHttpRequest(JSON.stringify(keycloakRequest));
    rawlog(keycloakRequest);
    return kcResp;
}

function decodeBase64(s) {
    var e = {}, i, b = 0, c, x, l = 0, a, r = '', w = String.fromCharCode, L = s.length;
    var A = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    for (i = 0; i < 64; i++) { e[A.charAt(i)] = i; }
    for (x = 0; x < L; x++) {
        c = e[s.charAt(x)]; b = (b << 6) + c; l += 6;
        while (l >= 8) { ((a = (b >>> (l -= 8)) & 0xff) || (x < (L - 2))) && (r += w(a)); }
    }
    return r;
}

function introspectToken6780(request, session, config) {
    var auth = request.Headers.Authorization[0];
    var access_token = request.Params.token[0];

    var tmp = auth.split(' ');   // Split on a space,

    // var buf = new Buffer(tmp[1], 'base64'); // create a buffer and tell it the data coming in is base64
    // var plain_auth = buf.toString();        // read it back out as a string
    var plain_auth = decodeBase64(tmp[1]);

    rawlog("Decoded Authorization ", plain_auth);

    // At this point plain_auth = "username:password"

    var creds = plain_auth.split(':');      // split on a ':'
    var client_id = creds[0];
    var client_secret = creds[1];

    var kcResp = "";
    if (client_id.trim() === "revamp_ms2".trim()) {
        kcResp = callKeycloakIntrospect(request.Headers.Authorization[0], access_token, config.config_data.keycloak_server, config.config_data.keycloak_realm_ad1access);
    } else if (client_id.trim() === "ad1gate_mobile".trim()) {
        kcResp = callKeycloakIntrospect(request.Headers.Authorization[0], access_token, config.config_data.keycloak_server, config.config_data.keycloak_realm_ad1gate);
    } else {
        kcResp = callKeycloakIntrospect(request.Headers.Authorization[0], access_token, config.config_data.keycloak_server, config.config_data.keycloak_realm_default);

    }

    var kcRespObj = JSON.parse(kcResp);
    var kcRespBodyObj = JSON.parse(kcRespObj.Body);

    rawlog(kcRespObj);

    var responseObject = {
        Headers: { "Content-Type": "application/json" },
        Body: kcRespObj.Body,
        Code: kcRespObj.Code
    };
    return TykJsResponse(responseObject, session.meta_data);

}