package main

import "fmt"

type CellStrategy interface {
	GetPrice(price float64) float64
}

type StrategyA struct{}

func (a *StrategyA) GetPrice(price float64) float64 {
	fmt.Println("use Strategy A")
	return price * 0.8
}

type StrategyB struct{}

func (b *StrategyB) GetPrice(price float64) float64 {
	if price >= 200 {
		price -= 100
	}

	fmt.Println("use Strategy B")

	return price
}

type Goods struct {
	price    float64
	strategy CellStrategy
}

func (g *Goods) SetStrategy(strategy CellStrategy) {
	g.strategy = strategy
}

func (g *Goods) GetPrice() float64 {
	fmt.Println("原价: ", g.price, "现价: ", g.strategy.GetPrice(g.price))
	return g.strategy.GetPrice(g.price)
}

func main() {
	nike := Goods{
		price: 200,
	}
	nike.SetStrategy(new(StrategyA))
	nike.GetPrice()
}
