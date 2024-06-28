package helpers

type ResponseData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseNoPaginate struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseNoData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(params ResponseParams) any {

	var response any
	var status string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	} else {
		status = "failed"
	}

	if params.Data != nil {
		response = &ResponseData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
			Data:    params.Data,
		}
	} else {
		response = &ResponseNoData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
		}
	}

	return response
}
