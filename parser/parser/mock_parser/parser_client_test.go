package mock_parser

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"parser/parser/parserproto"
	"testing"
	"time"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestParse(t *testing.T) {
	// set up test cases
	tests := []struct {
		filePath      string
		wantTitle     string
		wantThumbnail string
		wantContent   string
	}{
		{
			filePath:      "./test_urls/test_url1.html",
			wantTitle:     "Test Page1!",
			wantThumbnail: "",
			wantContent:   "",
		},
		{
			filePath:      "./test_urls/test_url2.html",
			wantTitle:     "Test Page2!",
			wantThumbnail: "3.jpg!",
			wantContent:   "",
		},
		{
			filePath:      "./test_urls/test_url3.html",
			wantTitle:     "Test Page3!",
			wantThumbnail: "3.jpg!",
			wantContent:   "Stuff to test1\nStuff to test2",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockParserClient := NewMockParserServiceClient(ctrl)
	for _, tt := range tests {
		fmt.Println("Testing " + tt.filePath)
		req := &parser.ParserTestRequest{FilePath: tt.filePath}
		mockParserClient.EXPECT().ParseTest(
			gomock.Any(),
			&rpcMsg{msg: req},
		).Return(&parser.ParserResponse{Title: tt.wantTitle, ThumbnailUrl: tt.wantThumbnail, Content: tt.wantContent}, nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		defer cancel()

		r, err := mockParserClient.ParseTest(ctx, &parser.ParserTestRequest{FilePath: tt.filePath})
		if err != nil || r.Title != tt.wantTitle || r.ThumbnailUrl != tt.wantThumbnail || r.Content != tt.wantContent {
			t.Errorf("mocking failed %s", err)
		}
		t.Log("Reply : ", r.Title, " - ", r.ThumbnailUrl, " - ", r.Content)
	}
}
