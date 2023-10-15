package resp

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	Writer    *http.ResponseWriter
	Logger    *logrus.Logger
	Operation string
}

type ErrorResponse struct {
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Message string   `json:"message" `
}

type OKResponse struct {
	Status string `json:"status"`
	//Errors []string `json:"errors"`
}

type OKWithDataResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
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
		Status:  "error",
		Message: "Password is incorrect",
		Errors:  []string{"PasswordIsIncorrect"}})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) UserNotFound() {
	er.Logger.Warning("Пользователь не найден")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Status:  "error",
		Message: "User not found",
		Errors:  []string{"UserNotFound"}})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) WrongParameter() {
	er.Logger.Warning("Неверные параметры")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Status:  "error",
		Message: "Invalid parameter",
		Errors:  []string{"InvalidParameter"}})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) OrganizationNotFound() {
	er.Logger.Warning("Организация не найдена")
	err := json.NewEncoder(*er.Writer).Encode(ErrorResponse{
		Status:  "error",
		Message: "Organization not found",
		Errors:  []string{"OrganizationNotFound"}})

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
	err := json.NewEncoder(*er.Writer).Encode(OKResponse{
		Status: "ok",
	})

	if err != nil {
		er.Logger.Error("Не удалось кодировать JSON: %s", err)
		http.Error(*er.Writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (er *Service) OKWithData(data interface{}) {
	err := json.NewEncoder(*er.Writer).Encode(OKWithDataResponse{
		Status: "ok",
		Data:   data,
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
