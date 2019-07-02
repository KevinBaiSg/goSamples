package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/KevinBaiSg/goSamples/grpc/proto"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{
	savedFeatures 	[]*pb.Feature
	mu         		sync.Mutex
	routeNotes 		map[string][]*pb.RouteNote
}

func (s *server) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Feature{Location: point}, nil
}

func (s *server) ListFeatures(rect *pb.Rectangle, stream pb.Route_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *server) RecordRoute(stream pb.Route_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	for {
		startTime := time.Now()
		point, e := stream.Recv()
		if e == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   	pointCount,
				FeatureCount: 	featureCount,
				Distance:     	distance,
				ElapsedTime:  	int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if e != nil {
			return e
		}

		pointCount ++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature, point) {
				featureCount ++
			}
		}

		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}

		lastPoint = point
	}
}

func (s *server) RouteChat(stream pb.Route_RouteChatServer) error {
	for {
		point, e := stream.Recv()
		if e == io.EOF {
			return nil
		}

		if e != nil {
			return e
		}

		key := serialize(point)

		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}

func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel) //TraceLevel WarnLevel or more
}

func main() {
	log.WithFields(log.Fields{
		"start": "main",
	}).Info("gRPC demo server starting")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.WithFields(log.Fields{
		"start": "listen",
		"port":	port,
	}).Info("gRPC demo server listen port")

	s := grpc.NewServer()
	pb.RegisterRouteServer(s, &server{})

	log.WithFields(log.Fields{
		"start": "gRPC server",
		"port":	port,
	}).Info("gRPC demo server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start serve: %v", err)
		return
	}
}
