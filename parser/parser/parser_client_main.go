package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"

	pb "parser/parser/parserproto"
)

const (
	address     = "localhost:50051"
	default_url = "https://medium.com/jatana/report-on-text-classification-using-cnn-rnn-han-f0e887214d5f"
	//default_url = "http://www.hurriyet.com.tr/gundem/son-dakika-murat-ozdemir-serbest-birakildi-41055578"
	//default_url = "https://www.sozcu.com.tr/2018/gundem/son-dakika-akpde-isyan-eden-ilce-teskilati-gorevden-alindi-2802280/"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewParserServiceClient(conn)

	// Contact the server and print out its response.
	url := default_url
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Parse(ctx, &pb.ParserRequest{Url: url})
	if err != nil {
		log.Fatalf("could not parse: %v", err)
	}
	log.Printf("Parsed Title: %s", r.Title)
	log.Printf("Parsed Thumbnail Image URL: %s", r.Thumbnail)
}
