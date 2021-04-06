package main

import (
	"flag"
	routing "github.com/jackwhelpton/fasthttp-routing"
	"github.com/jackwhelpton/fasthttp-routing/content"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"proximity-chat-grouper/server"
	"proximity-chat-grouper/util"
)

var (
	port              = flag.String("port", ":8081", "Port to listen on")
	locationThrottler = server.NewThrottler(1000)
	otherThrottler    = server.NewThrottler(100)
	serverService     = server.NewService()
)

func main() {
	flag.Parse()
	logrus.Info("Starting proximity chat server on port " + *port)

	r := routing.New()

	locationThrottler.Run()
	otherThrottler.Run()

	v1 := r.Group("/v1")

	v1.Use(
		content.TypeNegotiator(content.JSON),
	)

	serverApi := v1.Group("/server") // Minecraft
	clientApi := v1.Group("/client") // Website

	serverApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	serverApi.Get("/<id>/", readServer)
	serverApi.Post("/", createServer)
	serverApi.Delete("/<id>", deleteServer)

	joinApi := serverApi.Group("/join")

	joinApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	joinApi.Post("/<id>", joinServer)

	clientApi.Use(
		content.TypeNegotiator(content.JSON),
	)

	logrus.Fatal(fasthttp.ListenAndServe(*port, r.HandleRequest))
}

func readServer(context *routing.Context) error {
	context = util.DecorateContext(context)

	id := context.Param("id")

	if otherThrottler.IncreaseRate(id) {
		return util.Throttled(context)
	}

	serv, err := serverService.ReadServerAsSimplified(id)

	if err != nil {
		return util.ErrorContext(util.NOT_FOUND, context, err)
	}

	return util.SuccessContext(context, serv)
}

func joinServer(context *routing.Context) error {
	context = util.DecorateContext(context)
	id := context.Param("id")

	logrus.Info(string(context.PostBody()))

	if otherThrottler.IncreaseRate(id) {
		return util.Throttled(context)
	}

	logrus.Info("joining server...")

	body := struct {
		Uuid string
	}{}

	if err := context.Read(&body); err != nil {
		logrus.Errorf("Unable to parse joinServer %s", err)
		return util.ErrorContext(util.BAD_REQUEST, context, err)
	}

	serv, err := serverService.ReadServer(id)

	if err != nil {
		return util.ErrorContext(util.NOT_FOUND, context, err)
	}

	if err = context.Read(&body); err != nil {
		return util.ErrorContext(util.BAD_REQUEST, context, err)
	}

	potentialUser, err := serv.CreatePotentialUser(body.Uuid)

	if err != nil {
		return util.ErrorContext(util.INTERNAL_SERVER_ERROR, context, err)
	}

	return context.Write(potentialUser)
}

func deleteServer(context *routing.Context) error {
	context = util.DecorateContext(context)
	id := context.Param("id")

	if otherThrottler.IncreaseRate(id) {
		return util.Throttled(context)
	}

	return serverService.DeleteServer(id)
}

func createServer(context *routing.Context) error {
	context = util.DecorateContext(context)

	body := struct {
		Name string
	}{}

	if err := context.Read(&body); err != nil {
		logrus.Errorf("Unable to parse createServer %s", err)
		return util.ErrorContext(util.BAD_REQUEST, context, err)
	}

	simplifiedServer := serverService.CreateServer(body.Name)

	return util.SuccessContext(context, simplifiedServer)
}

//func handlePacket(context *routing.Context) error {
//
//	serverId := context.Param("id")
//
//	locationPacket := server.NewLocation(context)
//
//	logrus.Info(string(context.RequestURI()))
//
//	isThrottled := locationThrottler.IncreaseRate(serverId)
//
//	if isThrottled {
//		context.SetStatusCode(403)
//		err := context.Write(NewMessageObject("You are currently throttled!"))
//		return err
//	} else {
//		err := context.Write(&locationPacket)
//		return err
//	}
//}
