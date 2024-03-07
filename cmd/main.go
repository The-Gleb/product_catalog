package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/The-Gleb/product_catalog/internal/adapter/db"
	"github.com/The-Gleb/product_catalog/internal/adapter/dummyjson"
	"github.com/The-Gleb/product_catalog/internal/config"
	v1 "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler"
	category_handlers "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler/category"
	product_handlers "github.com/The-Gleb/product_catalog/internal/controller/http/v1/handler/product"
	middleware "github.com/The-Gleb/product_catalog/internal/controller/http/v1/middleware"
	"github.com/The-Gleb/product_catalog/internal/domain/service"
	"github.com/The-Gleb/product_catalog/internal/domain/usecase"
	"github.com/The-Gleb/product_catalog/internal/logger"
	"github.com/The-Gleb/product_catalog/pkg/client/postgresql"
	"github.com/go-chi/chi/v5"
)

func main() {

	if err := Run(); err != nil {
		panic(err)
	}

}

func Run() error {
	config := config.MustBuild("")
	logger.Initialize("debug")
	client, err := postgresql.NewClient(context.Background(), config.DB)
	if err != nil {
		return err
	}

	productStorage := db.NewProductStorage(client)
	categoryStorage := db.NewCategoryStorage(client)
	sessionStorage := db.NewSessionStorage(client)
	userStorage := db.NewUserStorage(client)

	productClient := dummyjson.NewProductClient(config.DummyJSONAddress)

	productService := service.NewProductService(productStorage, productClient, config.ProductUpdateInterval)
	categoryService := service.NewCategoryService(categoryStorage)
	sessionService := service.NewSessionService(sessionStorage)
	userService := service.NewUserService(userStorage)

	productUsecase := usecase.NewProductUsecase(productService)
	categoryUsecase := usecase.NewCategoryUsecase(categoryService)
	registerUsecase := usecase.NewRegisterUsecase(userService, sessionService)
	loginUsecase := usecase.NewLoginUsecase(userService, sessionService)
	authUsecase := usecase.NewAuthUsecase(sessionService)

	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	r := chi.NewRouter()

	v1.NewRegisterHandler(registerUsecase).AddToRouter(r)
	v1.NewLoginHandler(loginUsecase).AddToRouter(r)
	product_handlers.NewGetProductsByCategoryHandler(productUsecase).AddToRouter(r)
	product_handlers.NewAddProductHandler(productUsecase).Middlewares(authMiddleware.Do).AddToRouter(r)
	product_handlers.NewDeleteProductHandler(productService).Middlewares(authMiddleware.Do).AddToRouter(r)
	product_handlers.NewUpdateProductNameHandler(productUsecase).Middlewares(authMiddleware.Do).AddToRouter(r)
	product_handlers.NewUpdateProductCategoryHandler(productService).Middlewares(authMiddleware.Do).AddToRouter(r)
	category_handlers.NewGetAllCategoriesHandler(categoryUsecase).AddToRouter(r)

	category_handlers.NewAddCategoryHandler(categoryUsecase).Middlewares(authMiddleware.Do).AddToRouter(r)
	category_handlers.NewDeleteCategoryHandler(categoryUsecase).Middlewares(authMiddleware.Do).AddToRouter(r)
	category_handlers.NewUpdateCategoryNameHandler(categoryUsecase).Middlewares(authMiddleware.Do).AddToRouter(r)

	server := http.Server{
		Addr:    config.RunAddress,
		Handler: r,
	}

	ServerShutdownSignal := make(chan os.Signal, 1)
	signal.Notify(ServerShutdownSignal, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := productService.CheckNewProducts(ctx)
		if err != nil {
			slog.Error("error in updating products from dummyjson")
			ServerShutdownSignal <- syscall.SIGINT
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ServerShutdownSignal
		server.Shutdown(context.Background())
		cancel()
	}()

	slog.Info("config", "struct", config)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	wg.Wait()
	slog.Info("server shutdown")
	return nil
}
