// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Авторизация пользователя. При успешной авторизации отправляет куки с сессией. Если пользователь уже авторизован, то прежний cookies с сессией перезаписывается",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "Данные для входа",
                        "name": "loginData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Регистрация пользователя. Сразу же возвращает сессию в cookies",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "registerData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/collections/compilation/{genre}": {
            "get": {
                "description": "Возвращает актуальные подборки фильмов по указанному жанру. Если передать cookies с сессией, то подборка будет персонализированной",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Collections"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "default": "session=xxx",
                        "description": "session",
                        "name": "Cookie",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "Название жанра",
                        "name": "genre",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список с id фильмов указанного жанра",
                        "schema": {
                            "$ref": "#/definitions/collections.CompilationData"
                        }
                    },
                    "400": {
                        "description": "Требуется указать жанр",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Такой жанр не найден",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/collections/genres": {
            "get": {
                "description": "Возвращает список всех доступных жанров фильмов и сериалов",
                "tags": [
                    "Collections"
                ],
                "responses": {
                    "200": {
                        "description": "Список с id фильмов указанного жанра",
                        "schema": {
                            "$ref": "#/definitions/collections.GenresData"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/content/contentPreview": {
            "get": {
                "description": "Возвращает краткую информацию о фильме или сериале",
                "tags": [
                    "Content"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID искомого контента. Контентом может быть как фильм, так и сериал",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список с id фильмов указанного жанра",
                        "schema": {
                            "$ref": "#/definitions/content.PreviewInfoData"
                        }
                    },
                    "400": {
                        "description": "Требуется указать валидный id контента",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Контент с таким id не найден",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/playground/ping": {
            "get": {
                "description": "Проверка соединения через классический ping pong",
                "tags": [
                    "Playground"
                ],
                "responses": {
                    "200": {
                        "description": "Pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginData": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string",
                    "format": "string",
                    "example": "email@email.com"
                },
                "password": {
                    "type": "string",
                    "format": "string",
                    "example": "SecretPassword1!"
                }
            }
        },
        "auth.RegisterData": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "format": "string",
                    "example": "email@email.com"
                },
                "password": {
                    "type": "string",
                    "format": "string",
                    "example": "SecretPassword1!"
                }
            }
        },
        "collections.CompilationData": {
            "type": "object",
            "properties": {
                "genre": {
                    "type": "string",
                    "example": "action"
                },
                "ids": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        1,
                        2,
                        3
                    ]
                }
            }
        },
        "collections.GenresData": {
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "action",
                        "drama",
                        "comedian"
                    ]
                }
            }
        },
        "content.PreviewInfoData": {
            "type": "object",
            "properties": {
                "actors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Том Хэнкс",
                        "Сергей Бодров"
                    ]
                },
                "country": {
                    "type": "string",
                    "example": "Россия"
                },
                "director": {
                    "type": "string",
                    "example": "Тарантино"
                },
                "duration": {
                    "type": "integer",
                    "example": 134
                },
                "genre": {
                    "type": "string",
                    "example": "Боевик"
                },
                "original_title": {
                    "type": "string",
                    "example": "Batman"
                },
                "poster": {
                    "type": "string",
                    "example": "/static/poster.jpg"
                },
                "rating": {
                    "type": "number",
                    "example": 9.1
                },
                "release_year": {
                    "type": "integer",
                    "example": 2020
                },
                "title": {
                    "type": "string",
                    "example": "Бэтмен"
                }
            }
        },
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API Киноскопа",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
