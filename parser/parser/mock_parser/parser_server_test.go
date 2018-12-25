package mock_parser

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "parser/parser/parserproto"
)

type parser_server struct{}

func (ps *parser_server) Parse(ctx context.Context, input *pb.ParserRequest) (*pb.ParserResponse, error) {
	title, imgUrl, content, err := processHTML(input.Url)
	fmt.Println(title, "-", imgUrl, "-", content, "-", err)
	return &pb.ParserResponse{Title: title, ThumbnailUrl: imgUrl, Content: content}, err
}

func processHTML(inputUrl string) (string, string, string, error) {
	// Check URL validity
	_, err := url.ParseRequestURI(inputUrl)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// HTTP Request
	response, err := http.Get(inputUrl)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("Error loading HTTP response body", err)
		return "", "", "", err
	}

	title := getTitle(*document)
	fmt.Println("Title: " + title)
	imgUrl := getThumbnailImage(*document)
	fmt.Println("ImageURL: " + imgUrl)
	content := getContent(*document)
	fmt.Println("Content: " + content)
	return title, imgUrl, content, err
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
	var sb strings.Builder

	// Content parser for specific to Medium Blog.
	document.Find(".section-inner.sectionLayout--insetColumn").Each(func(index int, item *goquery.Selection) {
		item.Contents().Each(func(i int, ctx *goquery.Selection) {
			tmp := ctx.Text()
			if tmp != "" && !strings.Contains(tmp, "BlockedUnblockFollowFollowing") {
				sb.WriteString(tmp)
				sb.WriteString(" ")
				fmt.Printf("%d: %s\n", i, tmp)
			}
		})
	})

	if sb.String() == "" {
		// Content parser for specific to BBC News.
		document.Find(".story-body__inner").Each(func(index int, item *goquery.Selection) {
			item.ContentsFiltered("p").Each(func(i int, ctx *goquery.Selection) {
				tmp := ctx.Text()
				if tmp != "" && !strings.Contains(tmp, "\n") {
					sb.WriteString(tmp)
					sb.WriteString(" ")
					fmt.Printf("%d: %s\n", i, tmp)
				}
			})
		})
	}

	if sb.String() == "" {
		// Content parser for specific to Fox News.
		document.Find(".article-body").Each(func(index int, item *goquery.Selection) {
			item.ContentsFiltered("p").Each(func(i int, ctx *goquery.Selection) {
				tmp := ctx.Text()
				if tmp != "" && !strings.Contains(tmp, "\n") {
					sb.WriteString(tmp)
					sb.WriteString(" ")
					fmt.Printf("%d: %s\n", i, tmp)
				}
			})
		})
	}

	if sb.String() == "" {
		// Content parser for general usage.
		document.Find("body p").Each(func(index int, item *goquery.Selection) {
			tmp := item.Text()
			if tmp != "" && !strings.Contains(tmp, "\n") {
				sb.WriteString(tmp)
				sb.WriteString(" ")
				fmt.Printf("%d: %s\n", index, tmp)
			}
		})
		document.Find("body ol").Each(func(index int, item *goquery.Selection) {
			tmp := item.Text()
			if tmp != "" {
				sb.WriteString(tmp)
				sb.WriteString(" ")
				fmt.Printf("%d: %s\n", index, tmp)
			}
		})
		document.Find("body ul").Each(func(index int, item *goquery.Selection) {
			tmp := item.Text()
			if tmp != "" {
				sb.WriteString(tmp)
				sb.WriteString(" ")
				fmt.Printf("%d: %s\n", index, tmp)
			}
		})
	}

	if sb.String() == "" {
		sb.WriteString("Input is either empty webpage or its HTML is not parsable with current state of this code!")
	}

	return strings.TrimSpace(sb.String())
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
		prevImgAlt := ""
		document.Find("img").Each(func(i int, item *goquery.Selection) {
			tag := item
			imageAlt, _ := tag.Attr("alt")
			if len(prevImgAlt) < len(imageAlt) {
				imageUrl, _ = tag.Attr("src")
				prevImgAlt = imageAlt
			}
		})
	}

	if imageUrl == "" {
		imageUrl = "There is no image-related tags found in the given URL!"
	}

	return imageUrl
}

func (ps *parser_server) ParseTest(ctx context.Context, input *pb.ParserTestRequest) (*pb.ParserResponse, error) {
	title, imgUrl, content, err := processFileHTML(input.FilePath)
	return &pb.ParserResponse{Title: title, ThumbnailUrl: imgUrl, Content: content}, err
}

func processFileHTML(inputFileHtml string) (string, string, string, error) {
	f, err := os.Open(inputFileHtml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	document, err := goquery.NewDocumentFromReader(f)

	title := getTitle(*document)
	imgUrl := getThumbnailImage(*document)
	content := getContent(*document)
	return title, imgUrl, content, err
}

func Server() {
	port := ":50050"

	lis, err := net.Listen("tcp", port)
	log.Printf("Listening the port %s", port)
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

func TestMain(m *testing.M) {
	go Server()
	os.Exit(m.Run())
}

func TestParse(t *testing.T) {
	// Set up a connection to the Server.
	const address = "localhost:50050"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewParserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	// Test Cases
	tests := []struct {
		path          string
		wantTitle     string
		wantThumbnail string
		wantContent   string
	}{
		{
			path:          "./test_urls/test_url1.html",
			wantTitle:     "Test Page1!",
			wantThumbnail: "There is no image-related tags found in the given URL!",
			wantContent:   "Input is either empty webpage or its HTML is not parsable with current state of this code!",
		},
		{
			path:          "./test_urls/test_url2.html",
			wantTitle:     "Test Page2!",
			wantThumbnail: "3.jpg",
			wantContent:   "Input is either empty webpage or its HTML is not parsable with current state of this code!",
		},
		{
			path:          "./test_urls/test_url3.html",
			wantTitle:     "Test Page3!",
			wantThumbnail: "3.jpg",
			wantContent:   "Stuff to p1 Stuff to p2",
		},
		{
			path:          "./test_urls/test_url4.html",
			wantTitle:     "Test Page4 in h1 tag!",
			wantThumbnail: "3.jpg",
			wantContent:   "Stuff to p1 Stuff to p2",
		},
		{
			path:          "./test_urls/test_url5.html",
			wantTitle:     "There is no title-related tags found in the given URL!",
			wantThumbnail: "There is no image-related tags found in the given URL!",
			wantContent:   "Input is either empty webpage or its HTML is not parsable with current state of this code!",
		},
		{
			path:          "./test_urls/test_url6.html",
			wantTitle:     "There is no title-related tags found in the given URL!",
			wantThumbnail: "There is no image-related tags found in the given URL!",
			wantContent:   "Input is either empty webpage or its HTML is not parsable with current state of this code!",
		},
	}

	// Test ParseTest
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			path := tt.path
			r, err := c.ParseTest(ctx, &pb.ParserTestRequest{FilePath: path})
			if err != nil {
				t.Fatalf("Could not parse: %v", err)
			}
			t.Log("Title: ", r.Title, " - Thumbnail: ", r.ThumbnailUrl, "- Content: ", r.Content)
			if r.Title != tt.wantTitle {
				t.Errorf("Expected '%s', got %s", tt.wantTitle, r.Title)
			}
			if r.ThumbnailUrl != tt.wantThumbnail {
				t.Errorf("Expected '%s', got %s", tt.wantThumbnail, r.ThumbnailUrl)
			}
			if r.Content != tt.wantContent {
				t.Errorf("Expected '%s', got %s", tt.wantContent, r.Content)
			}
		})
	}
}
