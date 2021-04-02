package main

import (
	"flag"
	routing "github.com/jackwhelpton/fasthttp-routing"
	"github.com/jackwhelpton/fasthttp-routing/content"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"proximity-chat-grouper/server"
)

var (
	port      = flag.String("port", ":8081", "Port to listen on")
	compress  = flag.Bool("compress", false, "Enable transparent compress compression")
	manager   = server.NewManager()
	throttler = server.NewThrottler()
)

func main() {
	flag.Parse()
	logrus.Info("Starting proximity chat server...")

	r := routing.New()

	throttler.Run()

	serverApi := r.Group("/server")

	serverApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	serverApi.Get("/<id>", helloWorldId)
	serverApi.Post("/new", newServerEndpoint)
	serverApi.Post("/<id>", handlePacket)

	logrus.Fatal(fasthttp.ListenAndServe(*port, r.HandleRequest))
}

func newServerEndpoint(context *routing.Context) error {
	return nil
}

func handlePacket(context *routing.Context) error {

	serverId := context.Param("id")

	locationPacket := server.NewLocation(context)

	logrus.Info(string(context.RequestURI()))

	isThrottled := throttler.IncreaseThrottle(serverId)

	if isThrottled {
		context.SetStatusCode(403)
		err := context.Write(NewMessageObject("You are currently throttled!"))
		return err
	} else {
		err := context.Write(&locationPacket)
		return err
	}
}

func helloWorldId(context *routing.Context) error {
	return context.Write(NewMessageObject("Hello world"))
}

func NewMessageObject(message string) interface{} {
	return struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
}
