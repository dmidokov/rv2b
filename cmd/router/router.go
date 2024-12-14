package router

import (
	"context"
	"encoding/json"
	"github.com/dmidokov/rv2/config"
	"github.com/dmidokov/rv2/handlers/auth"
	branchH "github.com/dmidokov/rv2/handlers/branch"
	"github.com/dmidokov/rv2/handlers/group"
	navigationH "github.com/dmidokov/rv2/handlers/navigation"
	orgH "github.com/dmidokov/rv2/handlers/organization"
	"github.com/dmidokov/rv2/handlers/sse"
	userH "github.com/dmidokov/rv2/handlers/user"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/session/cookie"
	branchS "github.com/dmidokov/rv2/storage/postgres/branch"
	navigationS "github.com/dmidokov/rv2/storage/postgres/navigation"
	orgS "github.com/dmidokov/rv2/storage/postgres/organization"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/dmidokov/rv2/storage/postgres/user"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type SessionStorage interface {
	Save(req *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	GetByKey(req *http.Request, key string) (interface{}, bool)
	SetMaxAge(maxAge int)
	Get(req *http.Request) (map[interface{}]interface{}, error)
}

type Router struct {
	entitie.App
}

func New(ctx context.Context, db *pgxpool.Pool, cfg *config.Configuration, sessionStore *cookie.Service, log *logrus.Logger, sse *sse.EventService) *Router {
	return &Router{
		App: entitie.App{
			Ctx:         ctx,
			DB:          db,
			Config:      cfg,
			CookieStore: sessionStore,
			Logger:      log,
			SSE:         sse,
		},
	}
}

func (r *Router) Router() (*mux.Router, error) {

	authHandler := auth.New(r.Logger, r.DB, r.CookieStore, r.Config)
	orgHandler := orgH.New(r.Logger, r.DB, r.Config)
	branchHandler := branchH.New(r.Logger, r.DB, r.Config)
	userHandler := userH.New(r.Logger, r.DB, r.Config)
	navigationHandler := navigationH.New(r.Logger, r.DB, r.Config)
	groupHandler := group.New(r.Logger, r.DB, r.Config)

	userService := user.New(r.DB, r.CookieStore, r.Logger)
	orgService := orgS.New(r.DB, r.CookieStore, r.Logger)
	branchService := branchS.New(r.DB, r.CookieStore, r.Logger)
	navigationService := navigationS.New(r.DB, r.CookieStore, r.Logger)
	rightsService := rights.New(r.DB, r.Logger)
	sseService := sse.New(r.Logger, r.DB, r.CookieStore, r.Config, r.SSE)

	router := mux.NewRouter()

	router.HandleFunc("/", r.handleFileServer(r.Config.RootPathWeb, "")).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/temp/{fileName}.{ext}", r.handleFileServer(r.Config.TempFolder, "/temp/")).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/{folder}/{fileName}.{ext}", r.handleFileServer(r.Config.RootPathWeb, "")).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/auth", authHandler.SignIn(userService, orgService, sseService)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/authcheck", authHandler.AuthCheck).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/navigation", r.loggingMiddleware(navigationHandler.GetNavigation(userService, navigationService))).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/organizations", r.loggingMiddleware(orgHandler.Get(orgService, userService))).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/organizations", r.loggingMiddleware(orgHandler.Create(orgService, userService))).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/organizations/{id}", r.loggingMiddleware(orgHandler.DeleteOrganization(orgService))).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/organizations/{id}", orgHandler.GetById(orgService, userService)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/organizations/current", orgHandler.GetById(orgService, userService)).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/branches", r.loggingMiddleware(branchHandler.Get(branchService, userService))).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/branches", r.loggingMiddleware(branchHandler.Create(branchService, userService))).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/branches/{id}[0-9]+", r.loggingMiddleware(branchHandler.DeleteBranch(branchService, userService))).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/branches/active/{branchId}", r.loggingMiddleware(branchHandler.SetActiveBranch(userService, r.CookieStore))).Methods(http.MethodPost, http.MethodOptions)

	//router.HandleFunc("/api/users", r.loggingMiddleware(userHandler.GetUsers(userService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users", r.loggingMiddleware(userHandler.Create(userService))).Methods(http.MethodPut, http.MethodOptions)
	//router.HandleFunc("/api/users/icon", r.loggingMiddleware(userHandler.GetUserIcon(userService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/{id:[0-9]+}", r.loggingMiddleware(userHandler.DeleteUser(userService))).Methods(http.MethodDelete, http.MethodOptions)
	//router.HandleFunc("/api/users/{id:[0-9]+}", r.loggingMiddleware(userHandler.GetUser(userService, navigationService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/update", r.loggingMiddleware(userHandler.Update(userService, rightsService, navigationService))).Methods(http.MethodPost, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", r.loggingMiddleware(userHandler.AddToSwitcher(userService, rightsService))).Methods(http.MethodPut, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", r.loggingMiddleware(userHandler.RemoveFromSwitcher(userService, rightsService))).Methods(http.MethodDelete, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", r.loggingMiddleware(userHandler.GetSwitcher(userService, rightsService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher/switch", r.loggingMiddleware(userHandler.GetSwitcher(userService, rightsService))).Methods(http.MethodGet, http.MethodOptions)

	userRouter := router.
		PathPrefix("/api/users").
		Subrouter()

	userRouter.Use(r.loggingMiddleware1)
	userRouter.HandleFunc("", userHandler.GetUsers(userService, rightsService)).Methods(http.MethodGet)
	userRouter.HandleFunc("", userHandler.Create(userService)).Methods(http.MethodPut)
	userRouter.HandleFunc("/icon", userHandler.GetUserIcon(userService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser(userService)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id:[0-9]+}", userHandler.GetUser(userService, navigationService, rightsService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/update", userHandler.Update(userService, rightsService, navigationService)).Methods(http.MethodPost)
	userRouter.HandleFunc("/switcher", userHandler.AddToSwitcher(userService, rightsService)).Methods(http.MethodPut)
	userRouter.HandleFunc("/group", userHandler.AddGroup(userService, rightsService)).Methods(http.MethodPut)
	userRouter.HandleFunc("/group", userHandler.DeleteGroup(userService, rightsService)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/switcher", userHandler.RemoveFromSwitcher(userService, rightsService)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/switcher", userHandler.GetSwitcher(userService, rightsService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/switcher/switch", userHandler.SwitchUser(userService, rightsService, r.CookieStore)).Methods(http.MethodGet)
	//userRouter.HandleFunc("/rights/user", userHandler.GetAvailableRights(userService, rightsService, r.CookieStore)).Methods(http.MethodGet)

	groupRouter := router.PathPrefix("/api/groups").Subrouter()
	groupRouter.Use(r.loggingMiddleware1)
	groupRouter.HandleFunc("", groupHandler.GetGroups(userService, rightsService)).Methods(http.MethodGet)
	groupRouter.HandleFunc("", groupHandler.AddGroup(userService, rightsService)).Methods(http.MethodPut)
	groupRouter.HandleFunc("", groupHandler.DeleteGroup(userService, rightsService)).Methods(http.MethodDelete)
	groupRouter.HandleFunc("/rights", groupHandler.GetAvailableRights(userService, rightsService)).Methods(http.MethodGet)

	router.HandleFunc("/sse/{folder}", sseService.SseHandler())
	router.HandleFunc("/send/{event}/{client}", r.sendMessage())

	router.HandleFunc("/upload", r.uploadImage()).Methods(http.MethodPost, http.MethodOptions)
	if r.Config.MODE == config.DEV {
		router.Use(mux.CORSMethodMiddleware(router))
		router.Use(corsMiddleware)
	}

	return router, nil
}

func (r *Router) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			return
		}

		method := "Login middleware"
		log := r.Logger

		responses := resp.New(&w, log, method)

		if authenticated, ok := r.CookieStore.GetByKey(req, cookie.Authenticated); ok && authenticated.(bool) {
			r.CookieStore.Save(req, w, make(map[string]interface{}))
			if req.URL.String() == "/" {
				w.Header().Set("cache-control", "no-cache")
			}
		} else {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}
		next.ServeHTTP(w, req)
	}
}

func (r *Router) loggingMiddleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			return
		}

		method := "Login middleware"
		log := r.Logger
		log.Info(req.RequestURI)
		responses := resp.Service{Writer: &w, Logger: log, Operation: method}

		if authenticated, ok := r.CookieStore.GetByKey(req, cookie.Authenticated); ok && authenticated.(bool) {
			r.CookieStore.Save(req, w, make(map[string]interface{}))
			if req.URL.String() == "/" {
				w.Header().Set("cache-control", "no-cache")
			}
		} else {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}
		next.ServeHTTP(w, req)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://control.remontti.site:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		next.ServeHTTP(w, req)
	})
}

func (r *Router) Redirect(w http.ResponseWriter, req *http.Request) {
	host := strings.Split(req.Host, ":")[0]
	target := "https://" + host + ":" + r.Config.SSLPort + req.URL.Path
	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

func (r *Router) handleFileServer(dir, prefix string) http.HandlerFunc {
	r.Logger.Info(dir)
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		realHandler(w, req)
	}
}

func (r *Router) sendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		event := vars["event"]
		client, _ := strconv.Atoi(vars["client"])

		if r.SSE.Chanel != nil {
			r.SSE.Chanel <- sse.Event{Name: sse.EventName(event), Value: "data: bye-bye\n\n", UserId: client}
		}

	}
}

func (r *Router) uploadImage() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log := r.Logger
		log.Info("File uploading")
		err := req.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error("Error form parsing")
			log.Error(err)
			return
		}

		file, handler, err := req.FormFile("myFile")
		if err != nil {
			log.Error("Error Retrieving the File")
			log.Error(err)
			return
		}

		isIcon := req.Form.Get("isIcon")

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Error("Close file error")
			}
		}(file)

		log.Info("Uploaded File: %+v\n", handler.Filename)
		log.Info("File Size: %+v\n", handler.Size)
		log.Info("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern

		var tempFile *os.File
		var imageBasePath string

		if isIcon != "" {
			tempFile, err = os.CreateTemp(r.Config.TempFolder, "upload-*.png")
			imageBasePath = "/temp/"
		} else {
			tempFile, err = os.CreateTemp(r.Config.TempFolder, "upload-*.png")
			imageBasePath = "/temp/"
		}

		if err != nil {
			log.Error(err)
		}
		defer func(tempFile *os.File) {
			err := tempFile.Close()
			if err != nil {
				log.Error("Close file error")
			}
		}(tempFile)

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Error(err)
		}

		_, err = tempFile.Write(fileBytes)
		if err != nil {
			return
		}
		log.Info(w, "Successfully Uploaded File\n")

		if isIcon != "" {
			userService := user.New(r.DB, r.CookieStore, r.Logger)
			userService.GetUserIdFromSession(req)

			currentUserId := userService.GetUserIdFromSession(req)
			if currentUserId == 0 {
				log.Warning("Пользователь не найден в сессии")
				return
			}

			err := userService.SetIcon(currentUserId, imageBasePath+filepath.Base(tempFile.Name()))
			if err != nil {
				return
			}
		}
		a := map[string]string{"image-name": imageBasePath + filepath.Base(tempFile.Name())}

		err = json.NewEncoder(w).Encode(a)
		if err != nil {
			log.Errorf("Encode result error %s", err.Error())
		}
	}
}
