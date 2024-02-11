package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	handler_auth "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/auth"
	handler_geo "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/geo"
	handler_user "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/handlers/user"
	reverproxy "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/infrastructure/reverProxy"
	swaggerui "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/infrastructure/swaggerUI"
	auth_token_midw "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/middleware/auth"
	auth_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/auth"
	geo_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/geo"
	user_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/user"
)

func NewApiRouter( /*RPC CLIENT*/ ) http.Handler {
	r := chi.NewRouter()

	proxy := reverproxy.NewReverseProxy("hugo", "1313")

	r.Use(middleware.Logger)
	r.Use(proxy.ReverseProxy)

	//SwaggerUI
	r.Get("/swagger", swaggerui.SwaggerUI)

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./client_app/public"))).ServeHTTP(w, r)
	})

	// handlers

	geo_handler := handler_geo.NewGeoHandler(geo_client.NewGeoClient())
	user_handler := handler_user.NewUserHandler(user_client.NewUserClient())
	auth_handler := handler_auth.NewAuthHandler(auth_client.NewUserClient())



		// API
		r.Route("/api", func(r chi.Router) {
			//	geo := controllers.GeoController


			r.Post("/auth/register", auth_handler.Registeretion)

			r.Post("/auth/login", auth_handler.Login)

			// Group Adress
			r.Route("/address", func(r chi.Router) {

				r.Use(auth_token_midw.TokenAuthMiddleware)
				r.Post("/search", geo_handler.SearchAPI)

				r.Post("/geocode", geo_handler.GeocodeAPI)

			})

			r.Route("/user",func(r chi.Router) {
				r.Use(auth_token_midw.TokenAuthMiddleware)
				r.Get("/profile", user_handler.GetUser)
				r.Get("/list", user_handler.GetUsers)

			})
		})

	return r
}
