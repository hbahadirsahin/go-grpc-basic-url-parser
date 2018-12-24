package mock_parser

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func testParserServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockParserServer := NewMockParserServiceServer(ctrl)

}
