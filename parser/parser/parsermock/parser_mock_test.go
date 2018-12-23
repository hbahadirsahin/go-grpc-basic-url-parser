package parsermock

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"

	"parser/parser/parserproto"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockParserClient := NewMockParserServiceClient(ctrl)
	req := &parser.ParserRequest{Url: "unit_test"}
	mockParserClient.EXPECT().Parse(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&parser.ParserResponse{Title: "Mocked Interface", ThumbnailUrl: "Mocked URL", Content: "Mocked Content"}, nil)
	testParse(t, mockParserClient)
}

func testParse(t *testing.T, client parser.ParserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Parse(ctx, &parser.ParserRequest{Url: "unit_test"})
	if err != nil || r.Title != "Mocked Interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.Title)
}
