package certinfo

import (
	"errors"
	"testing"
	"time"

	"github.com/ol-se/ta-certinfo/internal"
	"github.com/ol-se/ta-certinfo/internal/app/certinfo/mocks"
)

var (
	input = internal.IDs{
		DaID: "QCDEMO",
		CID:  "3",
	}

	errMock  = errors.New("some error")
	mockData = []byte("mock data")
)

func assertFailure(t *testing.T, data []internal.CertData, err error) {
	t.Helper()

	if data != nil {
		t.Errorf("Expected no data, got %v\n", data)
	}

	if !errors.Is(err, errMock) {
		t.Errorf("Expected error: %v, got %v\n", errMock, err)
	}
}

func TestPullAndParse(t *testing.T) {
	t.Parallel()

	t.Run("PullAndParse: OK", func(t *testing.T) {
		t.Parallel()

		certData := []internal.CertData{
			{
				Sub: "sub1",
				Iss: "iss1",
				Eat: time.Now().Add(time.Hour),
			},
			{
				Sub: "sub2",
				Iss: "iss2",
				Eat: time.Now().Add(-time.Hour),
			},
		}

		mockStorage := mocks.NewStorageMock(t)
		mockStorage.EXPECT().PullCert(input).Return(mockData, nil)

		mockParser := mocks.NewParserMock(t)
		mockParser.EXPECT().Parse(mockData).Return(certData, nil)

		mockApp := &App{
			Storage: mockStorage,
			Parser:  mockParser,
		}

		data, err := mockApp.PullAndParse(input)

		if !internal.CertDataSliceEqual(data, certData) {
			t.Errorf("Expected %v, got %v\n", certData, data)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v\n", err)
		}
	})

	t.Run("PullAndParse: error pulling", func(t *testing.T) {
		t.Parallel()

		mockStorage := mocks.NewStorageMock(t)
		mockStorage.EXPECT().PullCert(input).Return(nil, errMock)

		mockApp := &App{
			Storage: mockStorage,
		}

		data, err := mockApp.PullAndParse(input)

		assertFailure(t, data, err)
	})

	t.Run("PullAndParse: error parsing", func(t *testing.T) {
		t.Parallel()

		mockStorage := mocks.NewStorageMock(t)
		mockStorage.EXPECT().PullCert(input).Return(mockData, nil)

		mockParser := mocks.NewParserMock(t)
		mockParser.EXPECT().Parse(mockData).Return(nil, errMock)

		mockApp := &App{
			Storage: mockStorage,
			Parser:  mockParser,
		}

		data, err := mockApp.PullAndParse(input)

		assertFailure(t, data, err)
	})
}
