module github.com/verza-platform/verza/services/ml-svc

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/verza/pkg/common v0.0.0
	go.uber.org/zap v1.26.0
)

replace github.com/verza/pkg/common => ../../pkg/common