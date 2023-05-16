module github.com/AmirMohG/container-runner/runner

replace github.com/AmirMohG/container-runner/services => ../services

go 1.19

require (
	github.com/AmirMohG/container-runner/services v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/streadway/amqp v1.0.0
	gorm.io/driver/mysql v1.5.0
	gorm.io/gorm v1.25.1
)

require (
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.6 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
