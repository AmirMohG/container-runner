package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	services "github.com/AmirMohG/container-runner/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// handle error
	}
	type Container struct {
		ID        uint `gorm:"primaryKey"`
		Image     string
		Tag       string
		Command   []string `gorm:"-"`
		Args      []string `gorm:"-"`
		Cid       string
		Output    string
		Expiry    time.Time
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	db.AutoMigrate(Container{})
	//redis connection

	redis := services.NewRedisClient()
	//rabbitmq connection
	conn, err := services.NewAMQPConnection()
	if err != nil {
		panic(err)
	}
	err = conn.CreateQueue("containers")
	if err != nil {
		panic(err)
	}
	ch, _ := conn.Channel()
	msgs, _ := conn.CreateConsumer(ch, "containers")
	var forever chan struct{}
	go func() {

		for d := range msgs {

			if d.Headers["reqType"] == "add" {
				var msg Container
				json.Unmarshal(d.Body, &msg)
				msg.Cid = d.Headers["cid"].(string)
				msg.Args = append(msg.Args, "--name", d.Headers["cid"].(string))
				stdout, err := runContainer(msg.Image, msg.Tag, msg.Command, msg.Args)
				if err != nil {
					fmt.Println(err.Error())
					err = redis.Set(d.Headers["cid"].(string), err.Error())
					if err != nil {
						msg.Output = err.Error()
						fmt.Println(err)
						panic(err)

					}

				} else {
					msg.Output = stdout
					err = redis.Set(d.Headers["cid"].(string), stdout)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
				}
				msg.Expiry = time.Now().Add(time.Hour * 24 * 30)
				db.Create(&msg)
				err = ch.Ack(d.DeliveryTag, false)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}

			} else if d.Headers["reqType"] == "delete" {
				cid := string(d.Body)
				deleteContainer(cid)
				result := db.Where("cid = ?", cid).Delete(&Container{})
				if result.Error != nil {
					// handle error
				}
				err := redis.Del(cid)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				err = ch.Ack(d.DeliveryTag, false)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}

			}
		}

	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func runContainer(image string, tag string, command []string, args []string) (string, error) {

	cmdArgs := []string{"docker", "run"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, fmt.Sprintf("%s:%s", image, tag))
	cmdArgs = append(cmdArgs, command...)

	// Execute the docker command and capture its output and error
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
func deleteContainer(cid string) error {
	cmd := exec.Command("docker", "rm", "-f", cid)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
