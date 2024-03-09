package exceptions

import (
	"errors"
	"fmt"
	"time"
)

type Layer string
type Type string

const (
	Server    Layer = "Server"    // Сторонняя ошибка на сервере (пятисотка любого рода)
	Database  Layer = "Database"  // Ошибка уровня баз данных
	Service   Layer = "Service"   // Ошибка на уровне бизнес-логики
	Transport Layer = "Transport" // Ошибка транспортного уровня

	Unprocessable Type = "Unprocessable"  // Невозможно обработать полученные данные
	BadRequest    Type = "Bad Request"    // Полученные данные содержат ошибку
	Untyped       Type = "Untyped"        // Нетипизированная ошибка
	NotFound      Type = "Not Found"      // Не найдено
	Forbidden     Type = "Forbidden"      // Нет доступа
	AlreadyExists Type = "Already Exists" // Уже существует
	Internal      Type = "Internal"       // Внутренняя ошибка
)

var (
	ServerErr    = errors.New("ошибка сервера")
	DatabaseErr  = errors.New("ошибка слоя БД")
	ServiceErr   = errors.New("ошибка бизнес-слоя")
	TransportErr = errors.New("ошибка транспортного слоя")

	UnprocessableErr = errors.New("невозможно обработать данные")
	BadRequestErr    = errors.New("некорректные данные")
	UntypedErr       = errors.New("неожиданная ошибка")
	NotFoundErr      = errors.New("не найдено")
	ForbiddenErr     = errors.New("нет доступа")
	AlreadyExistsErr = errors.New("уже существует")
	InternalErr      = errors.New("внутренняя ошибка")
)

var layerMap = map[Layer]error{
	Server:    ServerErr,
	Database:  DatabaseErr,
	Service:   ServiceErr,
	Transport: TransportErr,
}
var typeMap = map[Type]error{
	Unprocessable: UnprocessableErr,
	BadRequest:    BadRequestErr,
	Untyped:       UntypedErr,
	NotFound:      NotFoundErr,
	Forbidden:     ForbiddenErr,
	AlreadyExists: AlreadyExistsErr,
	Internal:      InternalErr,
}

type IError interface {
	error
	New()
}

type Error struct {
	When      time.Time
	Where     error
	Type      error
	What      error
	ClientMsg error
}

func (e Error) Error() string {
	return fmt.Sprintf("[%v] %v: %v: %v: %v", e.When, e.Where, e.Type, e.What, e.ClientMsg)
}

// New создаёт новую ошибку
// l - слой, на котором произошла ошибка
// t - тип ошибки
// m - сообщение ошибки для внутренних нужд
func New(l Layer, t Type, messages ...string) error {
	switch len(messages) {
	case 0:
		return Error{time.Now(), layerMap[l], typeMap[t], errors.New(""), errors.New("")}
	case 1:
		return Error{time.Now(), layerMap[l], typeMap[t], errors.New(""), errors.New(messages[0])}
	default:
		return Error{time.Now(), layerMap[l], typeMap[t], errors.New(messages[1]), errors.New(messages[0])}
	}
}

func Is(err1, err2 error) bool {
	var e Error
	if errors.As(err1, &e) {
		return errors.Is(err2, e.Type)
	}
	return false
}
