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
	// "nuro-rampup-ticketing/database/models"
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

func (s *server) CreatePokemon(ctx context.Context, in *pb.CreatePokemonReq) (*pb.CreatePokemonRes, error) {
	pokemon := in.GetPokemon()
	db_global.Exec("INSERT INTO pokemon (Type, Name, Points, ID) VALUES (?, ?, ?, ?)", pokemon.Type, pokemon.Name, pokemon.Points, pokemon.ID)
	return &pb.CreatePokemonRes{Success: true}, nil
}

func (s *server) ReadPokemon(ctx context.Context, in *pb.ReadPokemonReq) (*pb.ReadPokemonRes, error) {
	id := in.GetID()
	row := db_global.QueryRow("SELECT * FROM pokemon WHERE ID = ?", id)
	var type_ string
	var name string
	var points int32
	var id_ int32
	row.Scan(&type_, &name, &points, &id_)
	return &pb.ReadPokemonRes{Pokemon: &pb.Pokemon{Type: type_, Name: name, Points: points, ID: id_}}, nil
}

func (s *server) UpdatePokemon(ctx context.Context, in *pb.UpdatePokemonReq) (*pb.UpdatePokemonRes, error) {
	ID := in.GetID()
	Points := in.GetPoints()
	db_global.Exec("UPDATE pokemon SET Points = ? WHERE ID = ?", Points, ID)
	return &pb.UpdatePokemonRes{Success: true}, nil
}

func (s *server) DeletePokemon(ctx context.Context, in *pb.DeletePokemonReq) (*pb.DeletePokemonRes, error) {
	ID := in.GetID()
	db_global.Exec("DELETE FROM pokemon WHERE ID = ?", ID)
	return &pb.DeletePokemonRes{Success: true}, nil
}

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
