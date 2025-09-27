module github.com/verza

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/verza/pkg/blockchain v0.0.0
	github.com/verza/pkg/database v0.0.0
	github.com/verza/pkg/kms v0.0.0
	github.com/verza/pkg/security v0.0.0
	github.com/verza/pkg/vc v0.0.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/gorm v1.31.0 // indirect
)

replace github.com/verza/pkg/blockchain => ./pkg/blockchain

replace github.com/verza/pkg/database => ./pkg/database

replace github.com/verza/pkg/kms => ./pkg/kms

replace github.com/verza/pkg/security => ./pkg/security

replace github.com/verza/pkg/vc => ./pkg/vc

replace github.com/verza/pkg/common => ./pkg/common
