package exceptions

import (
	"fmt"
	"time"
)

type AppLayer string
type ExceptionType string

const (
	Server    AppLayer = "Server"    // Сторонняя ошибка на сервере (пятисотка любого рода)
	Database  AppLayer = "Database"  // Ошибка уровня баз данных
	Service   AppLayer = "Service"   // Ошибка на уровне бизнес-логики
	Transport AppLayer = "Transport" // Ошибка транспортного уровня

	Unprocessable ExceptionType = "Unprocessable"
	Untyped       ExceptionType = "Untyped"
	NotFound      ExceptionType = "NotFound"
	Forbidden     ExceptionType = "Forbidden"
	AlreadyExists ExceptionType = "AlreadyExists"
)

type Exception struct {
	When  time.Time
	What  string
	Layer AppLayer
	Type  ExceptionType
}

func (e Exception) Error() string {
	return fmt.Sprintf("%v", e.What)
}

// ----------------------

// AppException Представляет собой Exception с более подробным отображением ошибки
type AppException struct {
	exc Exception
}

func (e AppException) Error() string {
	return fmt.Sprintf("[%v] (%v) %v: %v", e.exc.When, e.exc.Layer, e.exc.Type, e.exc.What)
}
