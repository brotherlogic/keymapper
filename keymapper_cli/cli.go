package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/keymapper/proto"
)

func main() {
	ctx, cancel := utils.BuildContext("keymapper-cli", "keymapper")
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "keymapper")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeymapperServiceClient(conn)

	switch os.Args[1] {
	case "set":
		_, err := client.Set(ctx, &pb.SetRequest{Key: os.Args[2], Value: os.Args[3]})
		if err != nil {
			log.Fatalf("Error on Add Record: %v", err)
		}
	case "get":
		res, err := client.Get(ctx, &pb.GetRequest{Key: os.Args[2]})
		if err != nil {
			log.Fatalf("Error in listing: %v", err)
		}

		fmt.Printf("%v => %v\n", res.GetKey().GetKey(), res.GetKey().GetValue())

	}

}
