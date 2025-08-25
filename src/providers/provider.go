package providers

import (
	"subscriber-topic-stars/src/handlers"
	"subscriber-topic-stars/src/repositories"
	"subscriber-topic-stars/src/services"
)

type AppProvider struct {
	Repositories repositories.RepositoryCenter
	Services     services.ServiceCenter
	Handlers     handlers.HandlerCenter
}

func Register() AppProvider {
	repos := repositories.InitRepositories()
	allServices := services.InitServices(repos)
	allHandlers := handlers.InitHandlers(allServices)

	return AppProvider{
		Handlers: allHandlers,
	}
}
