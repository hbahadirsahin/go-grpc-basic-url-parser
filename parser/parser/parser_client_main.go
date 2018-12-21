package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "parser/parser/parserproto"
)

const (
	address = "localhost:1234"
	//default_url = "https://medium.com/jatana/report-on-text-classification-using-cnn-rnn-han-f0e887214d5f"
	//default_url = "http://www.hurriyet.com.tr/gundem/son-dakika-murat-ozdemir-serbest-birakildi-41055578"
	//default_url = "https://www.sozcu.com.tr/2018/gundem/son-dakika-akpde-isyan-eden-ilce-teskilati-gorevden-alindi-2802280/"
	//default_url = "https://www.bbc.com/news/uk-england-sussex-46623754"
	default_url = "https://www.foxnews.com/politics/veteran-launches-gofundme-campaign-to-help-fund-border-wall"
)

func main() {
	serverAddress := flag.String("address", "localhost:80", "A string argument for IP. Default value is localhost(it directs to 127.0.0.1:80)")
	inputUrl := flag.String("url", "https://medium.com/jatana/report-on-text-classification-using-cnn-rnn-han-f0e887214d5f", "A string argument for the input URL.")
	flag.Parse()

	fmt.Printf("You are connecting to %s\n", *serverAddress)
	fmt.Printf("Input URL: %s\n", *inputUrl)
	// Set up a connection to the server.
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewParserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	r, err := c.Parse(ctx, &pb.ParserRequest{Url: *inputUrl})
	if err != nil {
		log.Fatalf("could not parse: %v", err)
	}
	log.Printf("Parsed Title: %s", r.Title)
	log.Printf("Parsed Thumbnail Image URL: %s", r.ThumbnailUrl)
	log.Printf("Parsed Content: %s", r.Content)
}
