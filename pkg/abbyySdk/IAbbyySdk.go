package abbyySdk

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
)

func WordFormsRu(text string, isBinary bool) *SdkResponse {
	ru := 1049
	route := fmt.Sprintf("%s/%s", server, "api/v1/wordforms")
	request := gorequest.New()
	request.Get(route)
	request.Query(fmt.Sprintf("text=%s", text))
	request.Query(fmt.Sprintf("lang=%d", ru))
	request.Set("Content-Type", gorequest.TypeJSON)

	response := handleRequest(request, isBinary)

	return response
}
