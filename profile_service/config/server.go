package config

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"profile_service/config/config"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

//func (server *Server) Start() {
//	mongoClient := server.initMongoClient()
//	defer func(mongoClient *mongo.Client, ctx context.Context) {
//		err := mongoClient.Disconnect(ctx)
//		if err != nil {
//			log.Printf("error closing db: %s\n", err)
//		}
//	}(mongoClient, context.Background())
//
//	userStore := server.initUserStore(mongoClient)
//
//	userService := server.initUserService(userStore)
//
//	userHandler := server.initUserHandler(userService)
//
//	server.start(userHandler)
//}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := config.GetClient(server.config.AuthDBHost, server.config.AuthDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

//func (server *Server) initUserStore(client *mongo.Client) model.UserStore {
//	return repository.NewUserMongoDBStore(client)
//}
//
//func (server *Server) initUserService(store model.UserStore) *service.UserService {
//	return service.NewUserService(store)
//}
//
//func (server *Server) initUserHandler(service *service.UserService) *handlers.UserHandler {
//	return handlers.NewUserHandler(service)
//}
//
//func (server *Server) start(orderHandler *handlers.UserHandler) {
//	r := mux.NewRouter()
//	orderHandler.Init(r)
//
//	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))
//
//	srv := &http.Server{
//		Addr:    fmt.Sprintf(":%s", server.config.Port),
//		Handler: cors(r),
//	}
//	log.Println("Povezivanje na servis " + server.config.Port)
//	wait := time.Second * 15
//	go func() {
//		if err := srv.ListenAndServe(); err != nil {
//			log.Println("Greska pri povezivanju")
//			log.Println(err)
//		}
//	}()
//
//	c := make(chan os.Signal, 1)
//
//	signal.Notify(c, os.Interrupt)
//	signal.Notify(c, syscall.SIGTERM)
//
//	<-c
//
//	ctx, cancel := context.WithTimeout(context.Background(), wait)
//	defer cancel()
//
//	if err := srv.Shutdown(ctx); err != nil {
//		log.Fatalf("error shutting down server %s", err)
//	}
//	log.Println("server gracefully stopped")
//}
