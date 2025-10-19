package simplejson

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func Unmarshal[T any](data []byte) (result T, err error) {
	err = json.Unmarshal(data, &result)
	return
}

// UnmarshalResponse returns unmarshaled data.
//
// CAUTION: This function automatically close response body so can't use body later.
func UnmarshalResponse[T any](response *http.Response) (result T, err error) {
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	slog.Debug("Unmarshal data", slog.String("raw", string(b)))
	result, err = Unmarshal[T](b)
	return
}
