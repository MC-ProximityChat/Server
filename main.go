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
	logrus.Info("Starting proximity chat server on port " + *port)

	r := routing.New()

	throttler.Run()

	initCli()

	serverApi := r.Group("/server")

	authApi := r.Group("/auth")

	serverApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	authApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	authApi.Get("/token")

	serverApi.Get("/<id>", contextInfo)
	serverApi.Post("/new", newServerEndpoint)
	serverApi.Post("/<id>", handlePacket)

	logrus.Fatal(fasthttp.ListenAndServe(*port, r.HandleRequest))
}

func newServerEndpoint(context *routing.Context) error {

	serverBody := struct {
		Description string `json:"description"`
	}{}

	if err := context.Read(&serverBody); err != nil {
		logrus.Fatalf("Unable to read JSON %s", err)
	}

	serv := server.NewServer(serverBody.Description)
	manager.Add(serv.ID, serv)

	return context.Write(struct {
		ID   string
		Name string
	}{ID: serv.ID, Name: serv.Name})
}

func handlePacket(context *routing.Context) error {

	serverId := context.Param("id")

	locationPacket := server.NewLocation(context)

	logrus.Info(string(context.RequestURI()))

	isThrottled := throttler.IncreaseRate(serverId)

	if isThrottled {
		context.SetStatusCode(403)
		err := context.Write(NewMessageObject("You are currently throttled!"))
		return err
	} else {
		err := context.Write(&locationPacket)
		return err
	}
}

func contextInfo(context *routing.Context) error {
	context.Response.Header.Set("Access-Control-Allow-Origin", "*")
	serv, ok := manager.Get(context.Param("id"))

	if !ok {
		context.SetStatusCode(404)
		return context.Write(NewMessageObject("Server not found!"))
	}
	return context.Write(struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{ID: serv.ID, Name: serv.Name})
}

func NewMessageObject(message string) interface{} {
	return struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
}
