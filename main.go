package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/OzkrOssa/mikrotik-go"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("RedPlanet User Service Started")
	c := cron.New()
	_, err := c.AddFunc("0 20 * * *", func() {
		t := time.Now()

		var wg sync.WaitGroup
		mikrotikHost, err := LoadHost()
		if err != nil {
			log.Println(err)
		}
		var fullSecrets []mikrotik.Secret
		for _, host := range mikrotikHost {
			wg.Add(1)
			go func(h string) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						log.Println("Recovered from panic:", r)
					}
				}()
				log.Println("Connecting with: ", h)
				repo, err := mikrotik.NewMikrotikRepository(h, os.Getenv("MIKROTIK_API_USER"), os.Getenv("MIKROTIK_API_PASS"))
				if err != nil {
					log.Println(err, h)
					return
				}

				identity, err := repo.GetIdentity()
				if err != nil {
					log.Println(err)
				}

				secrets, err := repo.GetSecrets(identity, h)
				if err != nil {
					log.Println(err)
				}
				fullSecrets = append(fullSecrets, secrets...)
			}(host)
		}
		wg.Wait()

		saeRepo := NewSaeplusService("ACTIVO", "CORTADO", "POR%20CORTAR", "POR%20INSTALAR")
		users, err := saeRepo.FetchAPI()
		if err != nil {
			log.Println(err)
		}

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017", os.Getenv("MONGO_USER"), os.Getenv("MONGO_PASS"), os.Getenv("MONGO_HOST"))))

		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("redplanet").Collection("users")
		err = collection.Drop(context.Background())
		if err != nil {
			log.Println(err)
		}

		err = client.Database("redplanet").CreateCollection(context.Background(), "users")
		if err != nil {
			log.Fatal(err)
		}

		for _, u := range users {
			for _, ch := range fullSecrets {
				if !strings.Contains(ch.Comment, "_") {
					log.Println(ch)
				}
				Abonado := strings.Split(ch.Comment, "_")[1]
				if u.NroContrato == Abonado {
					row := User{
						Subscriber:   u.NroContrato,
						IDContract:   u.IDContrato,
						Document:     u.Cedula,
						FirstName:    u.Nombre,
						LastName:     u.Apellido,
						Status:       u.StatusContrato,
						Address:      u.Direccion,
						Phone:        u.Telefono,
						Subscription: u.Suscripcion,
						InformationMikrotik: MikrotikInformation{
							Secret:        ch.Name,
							CallerID:      ch.CallerID,
							Comment:       ch.Comment,
							RemoteAddress: ch.RemoteAddress,
							Profile:       ch.Profile,
							Bts:           ch.Bts,
							Host:          ch.Host,
						},
					}
					_, err = collection.InsertOne(context.Background(), row)
					if err != nil {
						log.Println(err)
					}
				}

			}
		}
		e := time.Since(t)
		fmt.Println(e)
	})
	if err != nil {
		return
	}
	c.Start()

	select {}

}
