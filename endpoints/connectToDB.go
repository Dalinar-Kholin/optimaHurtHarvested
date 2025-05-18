package endpoints

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"optimaHurt/constAndVars"
)

func ConnectToDB(connectionString string) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the serer
	client, err := mongo.Connect(constAndVars.ContextBackground, opts)
	if err != nil {
		panic(err)
	}
	constAndVars.DbConnect = client.Database(constAndVars.DbName)
	constAndVars.DbClient = client
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(constAndVars.ContextBackground, bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}
