package main

import (
	"context"
	"io"

	pb "github.com/KevinBaiSg/goSamples/grpc/proto"
	log "github.com/sirupsen/logrus"
)

const (
	target = "localhost:50051"
)

func main() {
	log.WithFields(log.Fields{
		"start": "main",
	}).Info("client start")

	// conn, err := grpc.Dial(target, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatal("connect server failed")
	// }
	conn, err := GetConn(target)
	if err != nil {
		log.Fatal("connect server failed")
	}
	defer conn.Close()

	log.Info("create a connect client successful")

	c := pb.NewRouteClient(conn)
	log.Info("new route client successfully")

	context := context.TODO()

	log.Info("start call GetFeature")
	feature, e := c.GetFeature(context, &pb.Point{
		Latitude:  409146138,
		Longitude: -746188906,
	})
	if e != nil {
		log.WithError(e)
		return
	}

	log.Info(feature)

	for {
		feature, e := c.ListFeatures(context, &pb.Rectangle{})
		if e == io.EOF {
			log.Info("ListFeatures finish")
			break
		}
		if e != nil {
			log.WithError(e)
			return
		}
		log.Info(feature)
	}
}
