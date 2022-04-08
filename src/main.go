/*
 * Main entry point for server
 */

package main

import (
  "context"
  "flag"
  "fmt"
  "log"
  "net"
  "database/sql"
  "github.com/volatiletech/sqlboiler/boil"
  "google.golang.org/grpc"
  pb "nuro-rampup-ticketing/proto"
  "nuro-rampup-ticketing/database/models"

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

func (s *server) CreatePokemon(ctx context.Context, in *pb.CreatePokemonReq) (*pb.CreatePokemonRes, error) {

  pokemon := &models.Pokemon{
    ID: int(in.GetPokemonInfo().GetId()),
    Type: in.GetPokemonInfo().GetType(),
    Name: in.GetPokemonInfo().GetName(),
    Points: int(in.GetPokemonInfo().GetPoints()),
  }
  err := pokemon.Insert(ctx, db_global, boil.Infer())
  if err == nil {
    return &pb.CreatePokemonRes{Message: "Success"}, nil
  } else {
    return &pb.CreatePokemonRes{Message: "Failed!"}, err
  }
}


func (s *server) ReadPokemon(ctx context.Context, in *pb.ReadPokemonReq) (*pb.ReadPokemonRes, error) {

  pokemonInfo, err := models.FindPokemon(ctx, db_global, int(in.GetId()))

  if err == nil {
    pokemon := pb.Pokemon{
      Id: int32(pokemonInfo.ID),
      Name: pokemonInfo.Name,
      Type: pokemonInfo.Type,
      Points: int32(pokemonInfo.Points),
    }

    return &pb.ReadPokemonRes{Pokemon: &pokemon}, nil
  } else {
    return &pb.ReadPokemonRes{Pokemon: nil}, err
  }
}


func (s *server) UpdatePokemonPoints(ctx context.Context, in *pb.UpdatePokemonPointsReq) (*pb.UpdatePokemonPointsRes, error) {

  pokemonDB, _ := models.FindPokemon(ctx, db_global, int(in.GetId()))
  pokemonDB.Points = int(in.GetPoints())
  _, err := pokemonDB.Update(ctx, db_global, boil.Infer())

  if err == nil {
    return &pb.UpdatePokemonPointsRes{Message: "Success"}, nil
  } else {
    return &pb.UpdatePokemonPointsRes{Message: "Failed!"}, err
  }
}


func (s *server) DeletePokemon(ctx context.Context, in *pb.DeletePokemonReq) (*pb.DeletePokemonRes, error) {

  pokemonDB, err := models.FindPokemon(ctx, db_global, int(in.GetId()))
  pokemonDB.Delete(ctx, db_global)

  if err == nil {
    return &pb.DeletePokemonRes{Message: "Success"}, nil
  } else {
    return &pb.DeletePokemonRes{Message: "Failed!"}, err
  }
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
