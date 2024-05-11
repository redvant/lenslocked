package main

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/gorilla/csrf"
	"github.com/redvant/lenslocked/controllers"
	"github.com/redvant/lenslocked/middleware"
	"github.com/redvant/lenslocked/migrations"
	"github.com/redvant/lenslocked/models"
	"github.com/redvant/lenslocked/templates"
	"github.com/redvant/lenslocked/views"
)

type config struct {
	PSQL   models.PostgresConfig
	SMTP   models.SMTPConfig
	CSRF   CSRFConfig
	Server ServerConfig
}
type CSRFConfig struct {
	Key    string `env:"CSRF_KEY,required"`
	Secure bool   `env:"CSRF_SECURE" envDefault:"true"`
}
type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS" envDefault:":3000"`
}

func main() {
	// Parse config from env var
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup middleware
	usersMw := middleware.Users{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
		ServerAddress:        cfg.Server.Address,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "signup.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "signin.gohtml",
	))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "forgot-pw.gohtml",
	))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "check-your-email.gohtml",
	))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "reset-pw.gohtml",
	))
	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "galleries/new.gohtml",
	))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "galleries/edit.gohtml",
	))
	galleriesC.Templates.Index = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "galleries/index.gohtml",
	))
	galleriesC.Templates.Show = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "galleries/show.gohtml",
	))
	galleriesC.Templates.ShowPublished = views.Must(views.ParseFS(
		templates.FS, "public.gohtml", "galleries/showPublished.gohtml",
	))

	// Setup router and routes
	router := http.NewServeMux()
	router.HandleFunc("/", controllers.Home(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))))
	router.HandleFunc("GET /contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))))
	router.HandleFunc("GET /faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))))
	router.HandleFunc("GET /signup", usersC.New)
	router.HandleFunc("POST /users", usersC.Create)
	router.HandleFunc("GET /signin", usersC.SignIn)
	router.HandleFunc("POST /signin", usersC.Authenticate)
	router.HandleFunc("POST /signout", usersC.SignOut)
	router.HandleFunc("GET /forgot-pw", usersC.ForgotPassword)
	router.HandleFunc("POST /forgot-pw", usersC.ProcessForgotPassword)
	router.HandleFunc("GET /reset-pw", usersC.ResetPassword)
	router.HandleFunc("POST /reset-pw", usersC.ProcessResetPassword)

	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET /", usersC.CurrentUser)
	userRouter.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})
	router.Handle("/users/me/", http.StripPrefix("/users/me", usersMw.RequireUser(userRouter)))

	galleriesRouter := http.NewServeMux()
	galleriesRouter.HandleFunc("GET /new", galleriesC.New)
	galleriesRouter.HandleFunc("POST /{$}", galleriesC.Create)
	galleriesRouter.HandleFunc("GET /{id}/edit", galleriesC.Edit)
	galleriesRouter.HandleFunc("POST /{id}", galleriesC.Update)
	galleriesRouter.HandleFunc("GET /{$}", galleriesC.Index)
	galleriesRouter.HandleFunc("GET /{id}", galleriesC.Show)
	galleriesRouter.HandleFunc("POST /{id}/delete", galleriesC.Delete)
	galleriesRouter.HandleFunc("POST /{id}/publish", galleriesC.Publish)
	galleriesRouter.HandleFunc("POST /{id}/unpublish", galleriesC.Unpublish)
	galleriesRouter.HandleFunc("GET /{id}/images/{filename}", galleriesC.Image)
	galleriesRouter.HandleFunc("POST /{id}/images", galleriesC.UploadImage)
	galleriesRouter.HandleFunc("POST /{id}/images/{filename}/delete", galleriesC.DeleteImage)
	router.Handle("/galleries/", http.StripPrefix("/galleries", usersMw.RequireUser(galleriesRouter)))

	router.HandleFunc("GET /g/{id}", galleriesC.ShowPublished)
	router.HandleFunc("GET /g/{id}/images/{filename}", galleriesC.PublishedImage)

	// Setup general middleware chain stack
	mwStack := middleware.CreateStack(
		middleware.Logging,
		csrfMw,
		usersMw.SetUser,
	)

	// Setup server
	server := http.Server{
		Addr:    cfg.Server.Address,
		Handler: mwStack(router),
	}

	// Start server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
