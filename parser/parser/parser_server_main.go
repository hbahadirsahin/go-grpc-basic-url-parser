package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "parser/parser/parserproto"
)

const (
	port = ":50051"
)

type parser_server struct{}

func (ps *parser_server) Parse(ctx context.Context, input *pb.ParserRequest) (*pb.ParserResponse, error) {
	title, imgUrl := processHTML(input.Url)
	return &pb.ParserResponse{Title: title, Thumbnail: imgUrl}, nil
}

func processHTML(url string) (string, string) {
	// HTTP Request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body", err)
	}

	title := getTitle(*document)
	fmt.Println("Title: " + title)
	imgUrl := getThumbnailImage(*document)
	fmt.Println("ImageURL: " + imgUrl)
	return title, imgUrl
}

func getTitle(document goquery.Document) string {
	// Get <Title> tag
	title := document.Find("title").Text()

	// Get all titles tagged with <h1>
	titlesH1 := make([]string, 0)
	document.Find("h1").Each(func(index int, item *goquery.Selection) {
		titlesH1 = append(titlesH1, item.Text())
	})

	// Get all titles tagged with <h2>
	titlesH2 := make([]string, 0)
	document.Find("h2").Each(func(index int, item *goquery.Selection) {
		titlesH2 = append(titlesH2, item.Text())
	})

	// Get all titles tagged with <h3>
	titlesH3 := make([]string, 0)
	document.Find("h3").Each(func(index int, item *goquery.Selection) {
		titlesH3 = append(titlesH3, item.Text())
	})

	// Check whether URL contains <title>.
	// If title is empty, get first <h1> title.
	if title == "" {
		if 0 < len(titlesH1) {
			title = titlesH1[0]
		}
	}

	// Check whether URL contains <title> and/or <h1>
	// If title is empty, get first <h2> title.
	if title == "" {
		if 0 < len(titlesH2) {
			title = titlesH2[0]
		}
	}

	// Check whether URL contains <title>, <h1> and/or <h2>
	// If title is empty, get first <h3> title.
	if title == "" {
		if 0 < len(titlesH3) {
			title = titlesH3[0]
		}
	}

	// If all title related tags are empty, provide a warning string as title.
	if title == "" {
		title = "There is no title-related tags found in the given URL!"
	}

	return title
}

func getContent(document goquery.Document) string {
	content := ""

	return content
}

func getThumbnailImage(document goquery.Document) string {
	imageUrl := ""
	document.Find("figure img").Each(func(index int, item *goquery.Selection) {
		tag := item
		imageUrl, _ = tag.Attr("src")
	})

	if imageUrl == "" {
		prevImgAlt := ""
		document.Find("body article section").Find("img").Each(func(index int, item *goquery.Selection) {
			tag := item
			imageAlt, _ := tag.Attr("alt")
			if len(prevImgAlt) < len(imageAlt) {
				imageUrl, _ = tag.Attr("src")
				prevImgAlt = imageAlt
			}
		})
	}

	if imageUrl == "" {
		prevImgAlt := ""
		document.Find("body div").Find("img").Each(func(index int, item *goquery.Selection) {
			tag := item
			imageAlt, _ := tag.Attr("alt")
			if len(prevImgAlt) < len(imageAlt) {
				imageUrl, _ = tag.Attr("src")
				prevImgAlt = imageAlt
			}
		})
	}

	if imageUrl == "" {
		imageUrl = document.Find("img").Text()
	}

	if imageUrl == "" {
		imageUrl = "There is no image-related tags found in the given URL!"
	}

	return imageUrl
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterParserServiceServer(s, &parser_server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
