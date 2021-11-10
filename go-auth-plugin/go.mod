module adira_custom_auth_plugin

go 1.15

replace github.com/jensneuse/graphql-go-tools => github.com/TykTechnologies/graphql-go-tools v1.6.2-0.20210324124350-140640759f4b

require (
	github.com/TykTechnologies/tyk v1.9.2-0.20210930081546-bda54b0f790c
	github.com/TykTechnologies/tyk-cli v0.0.0-20200618203828-2c91f6de17bd // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/gima/govalid v0.0.0-20170508202833-5e9183219fa1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
)
