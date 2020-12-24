package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"time"

	"github.com/shopspring/decimal"
	"github.com/streadway/amqp"
)

var ConsumerHandlers = map[string]ConsumerHandler{
	"trade_executor":  &TradeExecutor{},
	"order_processor": &OrderProcessor{},
}

type LimitOrder struct {
	Id       int             `json:"id"`
	MarketId string          `json:"market"`
	MemberId int             `json:"member_id"`
	Type     string          `json:"type"`
	Price    decimal.Decimal `json:"price"`
	Volume   decimal.Decimal `json:"volume"`
}

type PayloadMatching struct {
	Action string     `json:"action"`
	Order  LimitOrder `json:"order"`
}

type PayloadCreateOrder struct {
	Action   string          `json:"action"`
	MarketId string          `json:"market_id"`
	Ask      string          `json:"ask"`
	Bid      string          `json:"bid"`
	Side     string          `json:"side"`
	MemberId int             `json:"member_id"`
	OrdType  string          `json:"ord_type"`
	Price    decimal.Decimal `json:"price"`
	Volume   decimal.Decimal `json:"volume"`
}

type TradeExecutor struct{}

func (h *TradeExecutor) Init() error {
	return nil
}
func (h *TradeExecutor) Process(delivery amqp.Delivery) error {
	fmt.Println("trade_executor", string(delivery.Body))
	return nil
}

type OrderProcessor struct{}

func (h *OrderProcessor) Init() error {
	return nil
}
func (h *OrderProcessor) Process(delivery amqp.Delivery) error {
	fmt.Println("order_processor", string(delivery.Body))
	return nil
}

// consumer handlers control which consumer from amqp.yml will be included
func TestRunConsumers(t *testing.T) {
	forever := make(chan bool)
	GetAMQPServer().RegisterAndRun(ConsumerHandlers)
	<-forever
}

type AMQPPublisherRegister struct {
	Publisher *AMQPPublisher
}

func (r *AMQPPublisherRegister) Register() {
	r.Publisher = GetAMQPServer().NewAMQPPublisher("order_processor", AMQPPublishOptions{
		RoutingKey: "",
		Mandatory:  false,
		Immediate:  false,
		Publishing: amqp.Publishing{
			ContentType: "application/json",
		},
	})

	return
}

// publishing style is just like peatio and there is no need to specify the routing key
func TestAMQPPublisher(t *testing.T) {
	var r AMQPPublisherRegister
	GetAMQPServer().RegisterAndRun(nil, &r)

	//payload := &PayloadMatching{
	//	Action: "submit",
	//	Order: LimitOrder{
	//		Id:       255694,
	//		MarketId: "wonbtc",
	//		MemberId: 495,
	//		Type:     "bid",
	//		Price:    decimal.NewFromFloat(0.00930142),
	//		Volume:   decimal.NewFromFloat(0.00111633),
	//	},
	//}

	//payload := &PayloadMatching{
	//Action: "submit",
	//Order: LimitOrder{
	//    Id:       1577,
	//    MarketId: "wonusd",
	//    MemberId: 268,
	//    Type:     "bid",
	//    Price:    decimal.NewFromFloat(0.01),
	//    Volume:   decimal.NewFromFloat(0.8),
	//},
	//}

	payload := &PayloadCreateOrder{
		Action:   "create_order",
		MarketId: "wonbtc",
		Ask:      "won",
		Bid:      "btc",
		Side:     "sell",
		MemberId: 1,
		OrdType:  "limit",
		Price:    decimal.NewFromFloat(86.422398888),
		Volume:   decimal.NewFromFloat(100.99999999),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	if r.Publisher != nil {
		err := r.Publisher.Publish(payloadBytes)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type StandardPublisherRegister struct {
	Publisher *AMQPPublisher
}

func (r *StandardPublisherRegister) Register() {
	r.Publisher = GetAMQPServer().NewStandardPublisher("", "trade_preprocess", AMQPPublishOptions{
		RoutingKey: "",
		Mandatory:  false,
		Immediate:  false,
		Publishing: amqp.Publishing{
			ContentType: "text/plain",
		},
	})

	return
}

type PayloadPublish struct {
	MarketId    string          `json:"market_id"`
	AskId       int             `json:"ask_id"`
	BidId       int             `json:"bid_id"`
	StrikePrice decimal.Decimal `json:"strike_price"`
	Volume      decimal.Decimal `json:"volume"`
	Funds       decimal.Decimal `json:"funds"`
	TimeStamp   time.Time       `json:"time_stamp"`
}

// input: exchange key, queue key and publishing options including routing key
func TestStandardPublisher(t *testing.T) {
	var r StandardPublisherRegister
	GetAMQPServer().RegisterAndRun(nil, &r)

	if r.Publisher != nil {
		for i := 0; i < 1; i++ {
			payload := &PayloadPublish{
				MarketId:    "wonbtc",
				AskId:       2 * i,
				BidId:       2*i + 1,
				StrikePrice: decimal.NewFromFloat(0.5),
				Volume:      decimal.NewFromFloat(3.2),
				Funds:       decimal.NewFromFloat(1.6),
				TimeStamp:   time.Now(),
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Fatal(err.Error())
			}

			r.Publisher.Publish(payloadBytes)
		}
	}
}
