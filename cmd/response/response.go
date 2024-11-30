package resp

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type Service struct {
	Writer    *http.ResponseWriter
	Logger    *logrus.Logger
	Operation string
}

func New(writer *http.ResponseWriter, logger *logrus.Logger, operation string) *Service {
	return &Service{Writer: writer, Logger: logger, Operation: operation}
}

type ErrorResponse struct {
	Response
	Message string `json:"message" `
}

type Response struct {
	Status string   `json:"status"`
	Errors []string `json:"errors" `
}

type OKWithDataResponse struct {
	Response
	Data interface{} `json:"data"`
}

func (er *Service) OkResponse() Response {
	return Response{
		Status: StatusOK,
	}
}

func (er *Service) JsonDecodeError() {
	errorText := "При декодировании данных авторизации произошла ошибка"
	er.Logger.Errorf(errorText)
	http.Error(*er.Writer, errorText, http.StatusInternalServerError)
	return
}

func (er *Service) PasswordIncorrect() {
	er.Logger.Warning("Пользователь не найден")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Response: Response{
			Status: StatusError,
			Errors: []string{"PasswordIsIncorrect"},
		},
		Message: "Password is incorrect",
	})
	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) UserNotFound() {
	er.Logger.Warning("Пользователь не найден")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Response: Response{
			Status: StatusError,
			Errors: []string{"UserNotFound"},
		},
		Message: "User not found",
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) WrongParameter() {
	er.Logger.Warning("Неверные параметры")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Response: Response{
			Status: StatusError,
			Errors: []string{"InvalidParameter"},
		},
		Message: "Invalid parameter",
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) OrganizationNotFound() {
	er.Logger.Warning("Организация не найдена")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Response: Response{
			Status: StatusError,
			Errors: []string{"OrganizationNotFound"},
		},
		Message: "Organization not found",
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) InternalServerError() {
	http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
}

func (er *Service) Unauthorized() {
	http.Error(*er.Writer, "User unauthorized", http.StatusUnauthorized)
}

func (er *Service) OK() {
	err := json.NewEncoder(*er.Writer).Encode(er.OkResponse())

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) OKWithData(data interface{}) {
	err := json.NewEncoder(*er.Writer).Encode(OKWithDataResponse{
		Response: er.OkResponse(),
		Data:     data,
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) EmptyData() {
	errorText := "OneOfSpecifiedParametersIsEmpty"
	er.Logger.Errorf("Произошла ошибка: %s", errorText)
	http.Error(*er.Writer, errorText, http.StatusBadRequest)
}

func (er *Service) NotAllowed() {
	http.Error(*er.Writer, "MethodNotAllowed", http.StatusForbidden)
}

func (er *Service) WithError(error string) {
	err := json.NewEncoder(*er.Writer).Encode(Response{
		Status: StatusError,
		Errors: []string{error},
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
