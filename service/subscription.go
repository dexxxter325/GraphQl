package service

import (
	"GRAPHQL/graph/model"
	"context"
	"math/rand"
	"strconv"
)

var carPublishedChannel map[string]chan *model.Car

func init() { //инициализируем пустую мапу для дальнейшего использования
	carPublishedChannel = make(map[string]chan *model.Car)
}

func (p *postgresService) CarPublished(ctx context.Context) (<-chan *model.Car, error) {
	//подписка для клиента
	id := rand.Int() //id канала
	idString := strconv.Itoa(id)
	carEvent := make(chan *model.Car, 1) //1-канал может хранить ток 1 эл-т
	go func() {                          //ждет,пока клиент выложит(создаст) машину
		<-ctx.Done()
	}()
	carPublishedChannel[idString] = carEvent //добавляем канал в мапу
	return carEvent, nil
}
