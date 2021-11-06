package main

import(
	"log"
	"context"
	"strconv"

	"app/proto"
	"app/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"	
)

var cl proto.BookProfilesClient

func Create(ctx *gin.Context){

	var	book models.Book

	ctx.ShouldBindJSON(&book)

	log.Println(book)

	req := &proto.CreateRequest{
		Book: &proto.Book{
			Author: book.Author,
			Name: book.Name,
			Price: book.Price,
			Genre: book.Genre,
			Cover: book.Cover,
			Page: book.Page,
		},
	}

	res, err := cl.Create(context.Background(), req)
	
	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func GetBooks(ctx *gin.Context){

	req := &proto.GetBooksRequest{
		Key: "books",
	}

	res, err := cl.Get(context.Background(), req)

	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func GetBook(ctx *gin.Context){

	idS := ctx.Param("id")

	id, err := strconv.Atoi(idS)

	if err != nil{
		log.Printf("Error with Atoi %v", err)
	}

	req := &proto.WorkOnlyIdRequest{
		Id: int32(id),
	}

	res, err := cl.GetById(context.Background(), req)

	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func Discount(ctx *gin.Context){

	var discount models.Discount

	ctx.BindJSON(&discount)

	req := &proto.DiscountRequest{
		Id: discount.Id,
		Percent: discount.Percent,
	}

	res, err := cl.Discount(context.Background(), req)

	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func Delete(ctx *gin.Context){

	var id models.Id 

	ctx.BindJSON(&id)

	req := &proto.WorkOnlyIdRequest{
		Id : id.Id,
	}

	res, err := cl.Delete(context.Background(), req)

	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func GetWithSearch(ctx *gin.Context){

	key := ctx.Param("key")

	req := &proto.SearchRequest{
		Key: key,
	}

	res, err := cl.Search(context.Background(), req)

	if err != nil {
		log.Fatalf("Error client with server %v", err)
	}

	ctx.JSON(200, res)
}

func main() {

	log.Println("Welcome client ...")

	Router := gin.Default()

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error with connection grpc %v", err)
	}

	defer conn.Close()

	cl = proto.NewBookProfilesClient(conn)

	Router.POST("/create", Create)
	Router.GET("/getbooks", GetBooks)
	Router.GET("/getbook/:id", GetBook)
	Router.PUT("/discount", Discount)
	Router.DELETE("/delete", Delete)
	Router.GET("/search/:key", GetWithSearch)

	Router.Run("localhost:9090")
}