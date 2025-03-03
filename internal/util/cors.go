package util

import "net/http"

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого origin (можно указать конкретный домен вместо "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем определенные методы (GET, POST, PUT, DELETE и т.д.)
		w.Header().Set("Access-Control-Allow-Methods", "POST")

		// Разрешаем определенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Разрешаем браузеру кэшировать предварительные запросы (OPTIONS) на 1 час
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Обработка предварительных запросов (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}