## GOLANG + REDIS + MONGODB + JWT Rest API

full course https://codevoweb.com/api-golang-mongodb-gin-gonic-project-setup/

### Step Create encode publickey and private key

- Go to web https://travistidwell.com/jsencrypt/demo/
- Press Generate New Keys
- Copy Private Key to this https://www.base64encode.org/
- We need convert from base64 to ascii, select ascii and press Encode
- Copy result encode ascii to ACCESS_TOKEN_PRIVATE_KEY
- Copy Public Key from step 3, and repeat process and copy result to ACCESS_TOKEN_PUBLIC_KEY
- Repeat the process for the refresh token. REFRESH_TOKEN_PRIVATE_KEY, REFRESH_TOKEN_PUBLIC_KEY
