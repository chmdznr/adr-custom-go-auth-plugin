### Urutan compile dan bundle
- ubah nama function (jika perlu)
- download library yg dibutuhkan:
```
go mod tidy 
go mod vendor
```
- compile menggunakan docker (on Linux):
```
docker run --rm -v `pwd`:/plugin-source tykio/tyk-plugin-compiler:v3.2.3-rc2 customAuthPlugin-20220105.v3.so
```
- update file `manifest.json`
- bundle dengan docker:
```
docker run --rm -w "/tmp" -v $(pwd):/tmp --entrypoint "/bin/sh" -it tykio/tyk-gateway:v3.2.3-rc2 -c '/opt/tyk-gateway/tyk bundle build -y'
```
- rename `build.zip` (kasih info tanggal jika perlu)