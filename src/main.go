/*
 * Main entry point for server
 */

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/volatiletech/sqlboiler/boil"
	"google.golang.org/grpc"
	"log"
	"net"
	"nuro-rampup-ticketing/database/models"
	pb "nuro-rampup-ticketing/proto"

	_ "github.com/go-sql-driver/mysql"
)

var (
	port = flag.Int("port", 50051, "The server port")
)
var db_global *sql.DB

type server struct {
	pb.UnimplementedPokemonCRUDServiceServer
}

func (s *server) Test(ctx context.Context, in *pb.TestReq) (*pb.TestRes, error) {
	return &pb.TestRes{Message: "Hello " + in.GetName()}, nil
}

// TODO: step 5. Write controllers for your endpoints that you defined in `pokemon.proto`


// Starts the Server
func main() {

	db, err := sql.Open("mysql", "root@/pokemon")
	db_global = db
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	boil.SetDB(db)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPokemonCRUDServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
