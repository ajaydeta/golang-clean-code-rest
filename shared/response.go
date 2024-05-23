package shared

import "github.com/gofiber/fiber/v2"

type JSONResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
	Reason  string `json:"reason,omitempty"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{}
}

func (r *JSONResponse) SetCount(i int64) *JSONResponse {
	r.Data = fiber.Map{
		"count": i,
	}
	return r
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetCode(code int) *JSONResponse {
	r.Code = code
	return r
}

func (r *JSONResponse) SetReason(reason any) *JSONResponse {

	if reason == nil {
		return r
	}

	switch reason.(type) {
	case error:
		r.Reason = reason.(error).Error()
	case string:
		r.Reason = reason.(string)
	}

	return r
}

func (r *JSONResponse) APIStatusSuccess() *JSONResponse {
	r.Code = fiber.StatusOK
	r.Message = "Success"
	return r
}

func (r *JSONResponse) APIStatusNotFound() *JSONResponse {
	r.Code = fiber.StatusNotFound
	r.Message = "Data Not Found"
	return r
}

func (r *JSONResponse) Send(c *fiber.Ctx) error {
	return c.Status(r.Code).JSON(r)
}
