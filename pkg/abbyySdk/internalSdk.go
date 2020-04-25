package abbyySdk

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
)

type SdkResponse struct {
	Status        string
	StatusCode    int
	BodyBinary    io.ReadCloser
	Body          string
	ContentLength int64
	WithError     bool
	Errors        []error
	IsBinary      bool
}

const (
	server               string = "https://developers.lingvolive.com"
	authRoute            string = "api/v1.1/authenticate"
	apiToken             string = "YzYyMjc0NTgtNGY0OC00OTg3LWE0OTMtMmI0YmZjZmUwMzI3OjU2MzljYjMwNTEwNzQwZTg5NzczZWY4YjI1NDQ1NmE2"
	errorInternalMessage string = "Internal Sdk error. Broken token bearer or error in handle request"
	errorInternalCode           = 1
	keyBearerToken              = "abbyy:service:key:bearer:token"
)

var bearer string

func auth(withCache bool) error {
	var flgEmptyBearer bool
	route := fmt.Sprintf("%s/%s", server, authRoute)
	if withCache {
		keyInRedisApiDay := getAbbyyApiDayKey()
		_ = slRedis.Get(keyInRedisApiDay, &bearer)
	}
	// сверяем длину ключа bearer
	flgEmptyBearer = len(bearer) < 10
	if flgEmptyBearer {
		request := gorequest.New()
		request.Post(route)
		// заполняются только после объявления метода
		request.Set("Authorization", fmt.Sprintf("Basic %s", apiToken))
		resp, body, errs := request.End()
		if len(errs) > 0 {
			return fmt.Errorf("Auth is failed. Url: %s", route)
		} else if len(body) < 10 {
			return fmt.Errorf(
				"Token is broken. Len < 10. Status: %s. Status Code: %d",
				resp.Status,
				resp.StatusCode)
		}
		bearer = body
		if withCache {
			_ = slRedis.Set(keyBearerToken, bearer)
		}
	}

	return nil
}

func handleRequest(r *gorequest.SuperAgent, isBinary ...bool) *SdkResponse {
	var sdkResponse *SdkResponse
	if len(isBinary) > 0 {
		sdkResponse = makeResponse(0, "unknown", "", isBinary[0])
	} else {
		sdkResponse = makeResponse(0, "unknown", "", false)
	}
	err := auth(true)
	if err != nil {
		sdkResponse.WithError = true
		sdkResponse.Errors = append(sdkResponse.Errors, err)
		sdkResponse.Status = "error"
		sdkResponse.StatusCode = errorInternalCode
		return sdkResponse
	}

	r.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	resp, body, errs := r.End()
	/// проверим 401 ошибку сервера
	if resp.StatusCode == 401 {
		/// пытаемся провести авторизацию
		err := auth(true)
		if err != nil {
			sdkResponse.Errors = append(sdkResponse.Errors, fmt.Errorf("%s", errorInternalMessage))
			sdkResponse.WithError = true
			return sdkResponse
		}
		/// пытаемся повторно получить данные по запросу
		r.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
		resp, body, errs = r.End()
	}

	sdkResponse.Status = resp.Status
	sdkResponse.StatusCode = resp.StatusCode

	if len(errs) > 0 {
		err := fmt.Errorf("Error handle request. Url: %s", r.Url)
		errs = append(errs, err)
		sdkResponse.Errors = errs
		sdkResponse.WithError = true
	}

	if sdkResponse.IsBinary {
		sdkResponse.BodyBinary = resp.Body
	} else {
		sdkResponse.Body = body
	}
	sdkResponse.ContentLength = resp.ContentLength
	return sdkResponse
}

func makeResponse(statusCode int, status string, body string, isBinary bool) *SdkResponse {
	sdkResponse := new(SdkResponse)
	sdkResponse.IsBinary = isBinary
	sdkResponse.WithError = false

	sdkResponse.StatusCode = statusCode
	sdkResponse.Status = status
	sdkResponse.Body = body

	return sdkResponse
}

func getAbbyyApiDayKey() string {
	key, err := slConfig.GetValue(keeper.KeyAbbyyApiDay)
	if err != nil {
		return ""
	}
	return key

}
