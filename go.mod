module sub-app-server

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/gorilla/websocket v1.4.0
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	gopkg.in/elazarl/goproxy.v1 v1.0.0-20180725130230-947c36da3153
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	golang.org/x/crypto => C:/Users/CPX/go/src/golang.org/x/crypto
	golang.org/x/net => C:/Users/CPX/go/src/golang.org/x/net
	golang.org/x/sync => C:/Users/CPX/go/src/golang.org/x/sync
	golang.org/x/sys => C:/Users/CPX/go/src/golang.org/x/sys
	golang.org/x/text => C:/Users/CPX/go/src/golang.org/x/text
	golang.org/x/tools => C:/Users/CPX/go/src/golang.org/x/tools
)
