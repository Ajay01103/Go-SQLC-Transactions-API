package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	repo "github.com/Ajay01103/goTransactonsAPI/internal/adapters/postgresql/sqlc"
	"github.com/Ajay01103/goTransactonsAPI/internal/auth"
	"github.com/Ajay01103/goTransactonsAPI/internal/users"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr      string
	db        dbConfig
	jwtSecret string
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	// serve the raw OpenAPI spec so Scalar can fetch it
	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../docs/swagger.json")
	})

	// Scalar interactive API reference
	r.Get("/reference", func(w http.ResponseWriter, r *http.Request) {
		specBytes, err := os.ReadFile("../docs/swagger.json")
		if err != nil {
			http.Error(w, "could not read API spec", http.StatusInternalServerError)
			return
		}

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(specBytes),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Go Transactions API",
			},
			DarkMode: false,
			Theme:    scalar.ThemeNone,
			CustomCss: `
				/* ── Warm Sunny Summer theme ── */

				/* ---------- Light Mode ---------- */
				body { background: #FFF8F0; }

				.light-mode {
					--scalar-color-1:            #3D1A00;
					--scalar-color-2:            #7A3D10;
					--scalar-color-3:            #B06030;
					--scalar-color-accent:       #E8552A;
					--scalar-background-1:       #FFF8F0;
					--scalar-background-2:       #FFF0DC;
					--scalar-background-3:       #FFE5C2;
					--scalar-background-accent:  rgba(232, 85, 42, 0.08);
					--scalar-border-color:       rgba(232, 85, 42, 0.18);
					--scalar-button-1:           #E8552A;
					--scalar-button-1-color:     #fff;
					--scalar-button-1-hover:     #C8421A;
					--scalar-color-green:        #2e7d32;
					--scalar-color-red:          #C0392B;
					--scalar-color-yellow:       #F5A623;
					--scalar-color-blue:         #1565c0;
					--scalar-color-orange:       #E8552A;
					--scalar-color-purple:       #6a1b9a;
					--scalar-scrollbar-color:        rgba(232, 85, 42, 0.20);
					--scalar-scrollbar-color-active: rgba(232, 85, 42, 0.40);
				}

				/* Sidebar */
				.light-mode .t-doc__sidebar {
					--scalar-sidebar-background-1:           #FFE8CC;
					--scalar-sidebar-color-1:                #3D1A00;
					--scalar-sidebar-color-2:                #7A3D10;
					--scalar-sidebar-color-active:           #E8552A;
					--scalar-sidebar-item-hover-color:       #E8552A;
					--scalar-sidebar-item-hover-background:  rgba(232, 85, 42, 0.10);
					--scalar-sidebar-item-active-background: rgba(232, 85, 42, 0.16);
					--scalar-sidebar-border-color:           rgba(232, 85, 42, 0.18);
					--scalar-sidebar-search-background:      #FFF0DC;
					--scalar-sidebar-search-border-color:    rgba(232, 85, 42, 0.22);
					--scalar-sidebar-search-color:           #7A3D10;
				}

				/* Header bar */
				.light-mode .t-doc__header {
					background: linear-gradient(135deg, #FFE8CC 0%, #FFD49A 100%);
					border-bottom: 1px solid rgba(232, 85, 42, 0.22);
				}

				/* Cards */
				.light-mode .scalar-card {
					border-color: rgba(232, 85, 42, 0.14);
					border-radius: 10px;
					background: #FFF4E6;
				}

				/* Code blocks */
				.light-mode .scalar-code-block {
					background: #FFF0DC;
					border: 1px solid rgba(232, 85, 42, 0.14);
					border-radius: 8px;
				}

				/* Links */
				.light-mode a {
					color: #E8552A;
				}
				.light-mode a:hover {
					color: #C8421A;
				}

				/* Response section */
				.light-mode .scalar-response {
					background: #FFF4E6;
					border-radius: 8px;
				}

				/* Search input focus ring */
				.light-mode input:focus {
					outline-color: #E8552A;
					border-color:  #E8552A;
				}

				/* ---------- Dark Mode ---------- */
				.dark-mode {
					--scalar-color-1:            #FFF0D8;
					--scalar-color-2:            #D4A070;
					--scalar-color-3:            #A07048;
					--scalar-color-accent:       #F07860;
					--scalar-background-1:       #1C1208;
					--scalar-background-2:       #2A1C0D;
					--scalar-background-3:       #382610;
					--scalar-background-accent:  rgba(240, 120, 96, 0.10);
					--scalar-border-color:       rgba(240, 120, 96, 0.18);
					--scalar-button-1:           #F07860;
					--scalar-button-1-color:     #1C1208;
					--scalar-button-1-hover:     #E86040;
					--scalar-color-green:        #66bb6a;
					--scalar-color-red:          #ef5350;
					--scalar-color-yellow:       #F5A623;
					--scalar-color-blue:         #42a5f5;
					--scalar-color-orange:       #F07860;
					--scalar-color-purple:       #ab47bc;
					--scalar-scrollbar-color:        rgba(240, 120, 96, 0.22);
					--scalar-scrollbar-color-active: rgba(240, 120, 96, 0.44);
				}

				/* Sidebar dark */
				.dark-mode .t-doc__sidebar {
					--scalar-sidebar-background-1:           #231508;
					--scalar-sidebar-color-1:                #FFF0D8;
					--scalar-sidebar-color-2:                #D4A070;
					--scalar-sidebar-color-active:           #F07860;
					--scalar-sidebar-item-hover-color:       #F07860;
					--scalar-sidebar-item-hover-background:  rgba(240, 120, 96, 0.10);
					--scalar-sidebar-item-active-background: rgba(240, 120, 96, 0.18);
					--scalar-sidebar-border-color:           rgba(240, 120, 96, 0.18);
					--scalar-sidebar-search-background:      #2A1C0D;
					--scalar-sidebar-search-border-color:    rgba(240, 120, 96, 0.22);
					--scalar-sidebar-search-color:           #D4A070;
				}

				/* Header bar dark */
				.dark-mode .t-doc__header {
					background: linear-gradient(135deg, #231508 0%, #2A1C0D 100%);
					border-bottom: 1px solid rgba(240, 120, 96, 0.22);
				}

				/* Cards dark */
				.dark-mode .scalar-card {
					border-color: rgba(240, 120, 96, 0.14);
					border-radius: 10px;
					background: #2A1C0D;
				}

				/* Code blocks dark */
				.dark-mode .scalar-code-block {
					background: #231508;
					border: 1px solid rgba(240, 120, 96, 0.14);
					border-radius: 8px;
				}

				/* Links dark */
				.dark-mode a {
					color: #F07860;
				}
				.dark-mode a:hover {
					color: #F5A623;
				}

				/* Search input focus ring dark */
				.dark-mode input:focus {
					outline-color: #F07860;
					border-color:  #F07860;
				}
			`,
		})
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		fmt.Fprintln(w, htmlContent)
	})

	// auth routes
	authRepo := auth.NewPostgresRepository(repo.New(app.db))
	authService := auth.NewService(authRepo, app.config.jwtSecret)
	authHandler := auth.NewHandler(authService)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// users routes (protected)
	usersRepo := users.NewPostgresRepository(repo.New(app.db))
	usersService := users.NewService(usersRepo)
	usersHandler := users.NewHandler(usersService)
	r.Route("/users", func(r chi.Router) {
		r.Use(auth.RequireAuth(app.config.jwtSecret))
		r.Get("/current-user", usersHandler.GetCurrentUser)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
