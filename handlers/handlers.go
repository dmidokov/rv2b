package handlers

import (
	"encoding/json"
	"github.com/dmidokov/rv2/config"
	"github.com/dmidokov/rv2/handlers/auth"
	branchH "github.com/dmidokov/rv2/handlers/branch"
	navigationH "github.com/dmidokov/rv2/handlers/navigation"
	orgH "github.com/dmidokov/rv2/handlers/organization"
	userH "github.com/dmidokov/rv2/handlers/user"
	"github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/session/cookie"
	"github.com/dmidokov/rv2/sse"
	branchS "github.com/dmidokov/rv2/storage/postgres/branch"
	navigationS "github.com/dmidokov/rv2/storage/postgres/navigation"
	orgS "github.com/dmidokov/rv2/storage/postgres/organization"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/dmidokov/rv2/storage/postgres/user"
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

type Service struct {
	DB          *pgxpool.Pool
	Config      *config.Configuration
	CookieStore SessionStorage
	Logger      *logrus.Logger
	SSE         *sse.EventService
}

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	GetByKey(r *http.Request, key string) (interface{}, bool)
	SetMaxAge(maxAge int)
	Get(r *http.Request) (map[interface{}]interface{}, error)
}

func New(db *pgxpool.Pool, cfg *config.Configuration, sessionStore SessionStorage, log *logrus.Logger, sse *sse.EventService) *Service {
	return &Service{
		DB:          db,
		Config:      cfg,
		CookieStore: sessionStore,
		Logger:      log,
		SSE:         sse,
	}
}

func (hm *Service) Router() (*mux.Router, error) {

	log := hm.Logger

	log.Info(hm.Config.RootPathWeb)

	authHandler := auth.New(hm.Logger, hm.DB, hm.CookieStore, hm.Config)
	orgHandler := orgH.New(hm.Logger, hm.DB, hm.Config)
	branchHandler := branchH.New(hm.Logger, hm.DB, hm.Config)
	userHandler := userH.New(hm.Logger, hm.DB, hm.Config)
	navigationHandler := navigationH.New(hm.Logger, hm.DB, hm.Config)

	userService := user.New(hm.DB, hm.CookieStore, hm.Logger)
	orgService := orgS.New(hm.DB, hm.CookieStore, hm.Logger)
	branchService := branchS.New(hm.DB, hm.CookieStore, hm.Logger)
	navigationService := navigationS.New(hm.DB, hm.CookieStore, hm.Logger)
	rightsService := rights.New(hm.DB, hm.Logger)

	router := mux.NewRouter()

	router.HandleFunc("/", hm.handleFileServer(hm.Config.RootPathWeb, "")).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/temp/{fileName}.{ext}", hm.handleFileServer(hm.Config.TempFolder, "/temp/")).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/{folder}/{fileName}.{ext}", hm.handleFileServer(hm.Config.RootPathWeb, "")).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/auth", authHandler.SignIn(userService, orgService)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/authcheck", authHandler.AuthCheck).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/navigation", hm.loggingMiddleware(navigationHandler.GetNavigation(userService, navigationService))).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/organizations", hm.loggingMiddleware(orgHandler.Get(orgService, userService))).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/organizations", hm.loggingMiddleware(orgHandler.Create(orgService, userService))).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/organizations/{id}", hm.loggingMiddleware(orgHandler.DeleteOrganization(orgService))).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/organizations/{id}", orgHandler.GetById(orgService, userService)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/organizations/current", orgHandler.GetById(orgService, userService)).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/branches", hm.loggingMiddleware(branchHandler.Get(branchService, userService))).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/branches", hm.loggingMiddleware(branchHandler.Create(branchService, userService))).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/branches/{id}[0-9]+", hm.loggingMiddleware(branchHandler.DeleteBranch(branchService, userService))).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/branches/active/{branchId}", hm.loggingMiddleware(branchHandler.SetActiveBranch(userService, hm.CookieStore))).Methods(http.MethodPost, http.MethodOptions)

	//router.HandleFunc("/api/users", hm.loggingMiddleware(userHandler.GetUsers(userService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users", hm.loggingMiddleware(userHandler.Create(userService))).Methods(http.MethodPut, http.MethodOptions)
	//router.HandleFunc("/api/users/icon", hm.loggingMiddleware(userHandler.GetUserIcon(userService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/{id:[0-9]+}", hm.loggingMiddleware(userHandler.DeleteUser(userService))).Methods(http.MethodDelete, http.MethodOptions)
	//router.HandleFunc("/api/users/{id:[0-9]+}", hm.loggingMiddleware(userHandler.GetUser(userService, navigationService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/update", hm.loggingMiddleware(userHandler.Update(userService, rightsService, navigationService))).Methods(http.MethodPost, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", hm.loggingMiddleware(userHandler.AddToSwitcher(userService, rightsService))).Methods(http.MethodPut, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", hm.loggingMiddleware(userHandler.RemoveFromSwitcher(userService, rightsService))).Methods(http.MethodDelete, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher", hm.loggingMiddleware(userHandler.GetSwitcher(userService, rightsService))).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/users/switcher/switch", hm.loggingMiddleware(userHandler.GetSwitcher(userService, rightsService))).Methods(http.MethodGet, http.MethodOptions)

	userRouter := router.
		PathPrefix("/api/users").
		Subrouter()

	userRouter.Use(hm.loggingMiddleware1)
	userRouter.HandleFunc("", userHandler.GetUsers(userService)).Methods(http.MethodGet)
	userRouter.HandleFunc("", userHandler.Create(userService)).Methods(http.MethodPut)
	userRouter.HandleFunc("/icon", userHandler.GetUserIcon(userService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser(userService)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id:[0-9]+}", userHandler.GetUser(userService, navigationService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/update", userHandler.Update(userService, rightsService, navigationService)).Methods(http.MethodPost)
	userRouter.HandleFunc("/switcher", userHandler.AddToSwitcher(userService, rightsService)).Methods(http.MethodPut)
	userRouter.HandleFunc("/switcher", userHandler.RemoveFromSwitcher(userService, rightsService)).Methods(http.MethodDelete)
	userRouter.HandleFunc("/switcher", userHandler.GetSwitcher(userService, rightsService)).Methods(http.MethodGet)
	userRouter.HandleFunc("/switcher/switch", userHandler.SwitchUser(userService, rightsService, hm.CookieStore)).Methods(http.MethodGet)

	router.HandleFunc("/sse/{folder}", hm.sseHandler())
	router.HandleFunc("/send/{event}/{client}", hm.sendMessage())

	router.HandleFunc("/upload", hm.uploadImage()).Methods(http.MethodPost, http.MethodOptions)
	if hm.Config.MODE == config.DEV {
		router.Use(mux.CORSMethodMiddleware(router))
		router.Use(corsMiddleware)
	}

	//for()

	return router, nil

}

func (hm *Service) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "Login middleware"
		log := hm.Logger

		responses := resp.Service{Writer: &w, Logger: log, Operation: method}

		if authenticated, ok := hm.CookieStore.GetByKey(r, cookie.Authenticated); ok && authenticated.(bool) {
			hm.CookieStore.Save(r, w, make(map[string]interface{}))
			if r.URL.String() == "/" {
				w.Header().Set("cache-control", "no-cache")
			}
		} else {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (hm *Service) loggingMiddleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "Login middleware"
		log := hm.Logger
		log.Info(r.RequestURI)
		responses := resp.Service{Writer: &w, Logger: log, Operation: method}

		if authenticated, ok := hm.CookieStore.GetByKey(r, cookie.Authenticated); ok && authenticated.(bool) {
			hm.CookieStore.Save(r, w, make(map[string]interface{}))
			if r.URL.String() == "/" {
				w.Header().Set("cache-control", "no-cache")
			}
		} else {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://control.remontti.site:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		next.ServeHTTP(w, r)
	})
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func (hm *Service) handleFileServer(dir, prefix string) http.HandlerFunc {
	hm.Logger.Info(dir)
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		realHandler(w, req)
	}
}

func (hm *Service) sendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		event := vars["event"]
		client, _ := strconv.Atoi(vars["client"])

		if hm.SSE.Chanel != nil {
			hm.SSE.Chanel <- sse.Event{Name: sse.EventName(event), Value: "data: bye-bye\n\n", UserId: client}
		}

	}
}

func (hm *Service) uploadImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := hm.Logger
		log.Info("File uploading")
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error("Error form parsing")
			log.Error(err)
			return
		}

		file, handler, err := r.FormFile("myFile")
		if err != nil {
			log.Error("Error Retrieving the File")
			log.Error(err)
			return
		}

		isIcon := r.Form.Get("isIcon")

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
			tempFile, err = os.CreateTemp(hm.Config.TempFolder, "upload-*.png")
			imageBasePath = "/temp/"
		} else {
			tempFile, err = os.CreateTemp(hm.Config.TempFolder, "upload-*.png")
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
			userService := user.New(hm.DB, hm.CookieStore, hm.Logger)
			userService.GetUserIdFromSession(r)

			currentUserId := userService.GetUserIdFromSession(r)
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
