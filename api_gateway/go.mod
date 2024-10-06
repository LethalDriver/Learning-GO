module example.com/chat_app/api_gateway

go 1.22.2

require (
	example.com/chat_app/common v0.0.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
)

replace example.com/chat_app/common => ../common
