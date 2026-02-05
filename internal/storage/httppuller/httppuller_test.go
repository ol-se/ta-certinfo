package httppuller

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/ol-se/ta-certinfo/internal"
	"github.com/ol-se/ta-certinfo/internal/storage/httppuller/mocks"
)

var (
	input = internal.IDs{
		DaID: "QCDEMO",
		CID:  "3",
	}

	mockHostname = "hostname"
	mockURL      = mockHostname + "/daid/" + input.DaID + "/cid/" + input.CID

	mockData      = []byte("mock data")
	mockReadLimit = len(mockData)

	errMock = errors.New("some error")
)

func assertFailure(t *testing.T, certData []byte, err error) {
	t.Helper()

	if !errors.Is(err, internal.ErrPullingCert) {
		t.Errorf("Expected error: %v, got %v\n", internal.ErrPullingCert, err)
	}

	if len(certData) != 0 {
		t.Errorf("Expected no data, got %v\n", certData)
	}
}

func TestPullCert(t *testing.T) {
	t.Parallel()

	t.Run("PullCert: OK", func(t *testing.T) {
		t.Parallel()

		mockReadCloser := mocks.NewReadCloserMock(t)
		mockReadCloser.EXPECT().Read(make([]byte, mockReadLimit)).RunAndReturn(func(buf []byte) (int, error) {
			copy(buf, mockData)

			return len(mockData), nil
		})
		mockReadCloser.EXPECT().Close().Return(nil)

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       mockReadCloser,
		}

		mockGetter := mocks.NewGetterMock(t)
		mockGetter.EXPECT().Get(mockURL).Return(mockResp, nil)

		mockPuller := &Puller{
			Hostname:  mockHostname,
			ReadLimit: int64(mockReadLimit),
			Getter:    mockGetter,
		}

		certData, err := mockPuller.PullCert(input)
		if err != nil {
			t.Errorf("Expected no error, got %v\n", err)
		}

		if !bytes.Equal(certData, mockData) {
			t.Errorf("Expected %v, got %v\n", mockData, certData)
		}
	})

	t.Run("PullCert: error joining link", func(t *testing.T) {
		t.Parallel()

		invalidLinkSegment := string([]rune{0x7f, 0x7f})

		mockPuller := &Puller{
			Hostname: invalidLinkSegment,
		}

		certData, err := mockPuller.PullCert(internal.IDs{})

		assertFailure(t, certData, err)
	})

	t.Run("PullCert: error getting response", func(t *testing.T) {
		t.Parallel()

		mockGetter := mocks.NewGetterMock(t)
		mockGetter.EXPECT().Get(mockURL).Return(nil, errMock)

		mockPuller := &Puller{
			Hostname:  mockHostname,
			ReadLimit: int64(mockReadLimit),
			Getter:    mockGetter,
		}

		certData, err := mockPuller.PullCert(input)

		assertFailure(t, certData, err)
	})

	t.Run("PullCert: wrong status code", func(t *testing.T) {
		t.Parallel()

		mockReadCloser := mocks.NewReadCloserMock(t)
		mockReadCloser.EXPECT().Close().Return(nil)

		mockResp := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       mockReadCloser,
		}

		mockGetter := mocks.NewGetterMock(t)
		mockGetter.EXPECT().Get(mockURL).Return(mockResp, nil)

		mockPuller := &Puller{
			Hostname: mockHostname,
			Getter:   mockGetter,
		}

		certData, err := mockPuller.PullCert(input)

		assertFailure(t, certData, err)
	})

	t.Run("PullCert: error reading body", func(t *testing.T) {
		t.Parallel()

		mockReadCloser := mocks.NewReadCloserMock(t)
		mockReadCloser.EXPECT().Read(make([]byte, mockReadLimit)).Return(0, errMock)
		mockReadCloser.EXPECT().Close().Return(nil)

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       mockReadCloser,
		}

		mockGetter := mocks.NewGetterMock(t)
		mockGetter.EXPECT().Get(mockURL).Return(mockResp, nil)

		mockPuller := &Puller{
			Hostname:  mockHostname,
			ReadLimit: int64(mockReadLimit),
			Getter:    mockGetter,
		}

		certData, err := mockPuller.PullCert(input)

		assertFailure(t, certData, err)
	})
}
