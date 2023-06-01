package application

import (
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/go-playground/form"
	"github.com/nats-io/stan.go"
	"github.com/sgoldenf/wb_l0/internal/interface/order"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache *template.Template
	FormDecoder   *form.Decoder
	Orders        order.OrderModelInterface
	cache         sync.Map
	stanConn      stan.Conn
	stanSub       stan.Subscription
}

func (app *Application) InitOrdersCache() error {
	orders, err := app.Orders.ReadAllOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		app.cache.Store(order.OrderID, order.Data)
	}
	return nil
}

func (app *Application) InitStanConnection() error {
	sc, err := stan.Connect("test-cluster", "test-sub", stan.NatsURL("localhost:4222"))
	if err != nil {
		return err
	}
	app.stanConn = sc
	return nil
}

func (app *Application) InitStanSubscription() error {
	sub, err := app.stanConn.Subscribe("orders", func(m *stan.Msg) {
		app.addOrder(m.Data)
	})
	if err != nil {
		return err
	}
	app.stanSub = sub
	return nil
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.logRequest(app.home))
	mux.HandleFunc("/order/", app.logRequest(app.viewOrder))
	return mux
}

func (app *Application) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) Shutdown() {
	if app.stanSub != nil {
		app.stanSub.Close()
	}
	if app.stanConn != nil {
		app.stanConn.Close()
	}
	app.Orders.Shutdown()
}

func (app *Application) addOrder(data []byte) {
	if err := validateOrderJSON(data); err != nil {
		app.InfoLog.Println("Got invalid message")
		return
	}
	o, err := parseOrderJSON(data)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
	if err := app.Orders.AddOrder(o); err != nil {
		app.ErrorLog.Println(err)
		return
	}
	app.cache.Store(o.OrderID, o.Data)
}
