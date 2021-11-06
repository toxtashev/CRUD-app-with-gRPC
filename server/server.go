package main

import(
	"log"
	"net"
	"context"

	"app/utils"
	"app/proto"
	"app/models"

	"google.golang.org/grpc"
)

type server struct{
	proto.UnimplementedBookProfilesServer
}

func (*server) Create(c context.Context, req *proto.CreateRequest)(*proto.MainResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	log.Println(req)

	new := db.Table("books").
			Create(&req.Book)

	if utils.IsNotFound(new) {
		log.Fatal(new.Error)
	}

	var newBook models.GetBook

	row := db.Table("books").
			Where("name = ?", req.Book.Name).
			Find(&newBook)

	log.Println(newBook.BookId)

	if utils.IsNotFound(row){
		log.Fatal(row.Error)
	}

	res := &proto.MainResponse{
		Book: &proto.BookMainInformation{
			Id: int32(newBook.BookId),
			Author: newBook.Author,
			Name:  newBook.Name,
			Price:  newBook.Price,
		},
	}

	return res, nil
}

func (*server) Get(c context.Context, req *proto.GetBooksRequest)(*proto.ManyResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	var books []models.GetBook

	rows := db.Table(req.Key).
			Find(&books)

	log.Println(books)

	if utils.IsNotFound(rows){
		log.Fatal(rows.Error)
	}

	var result []*proto.BookMainInformation

	for _, objekt := range books{

		res := &proto.BookMainInformation{
				Id: int32(objekt.BookId),
				Author: objekt.Author,
				Name: objekt.Name,
				Price: objekt.Price,
		}

		result = append(result, res)
	}

	return &proto.ManyResponse{
		Books: result, 
	}, nil
}

func (*server) GetById(c context.Context, req *proto.WorkOnlyIdRequest)(*proto.GetByIdResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	var book models.GetBook

	row := db.Table("books").
			Where("book_id = ?", req.GetId()).
			Find(&book)

	if utils.IsNotFound(row){
		log.Fatal(row.Error)
	}

	log.Println(book.BookId)

	res := &proto.GetByIdResponse{
		Book: &proto.Book{
			Author: book.Author,
			Name: book.Name,
			Price: book.Price,
			Genre: book.Genre,
			Cover: book.Cover,
			Page: book.Page,
		},
	}

	return res, nil
}

func (*server) Discount(c context.Context, req *proto.DiscountRequest)(*proto.MainResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	var oldPrice float32

	getPrice := db.Table("books").
			Select("price").
			Where("book_id = ?", req.Id).
			Find(&oldPrice)

	if utils.IsNotFound(getPrice){
		log.Fatal(getPrice.Error)
	}

	var newPrice float32

	newPrice = (float32(1) - float32(req.Percent)/float32(100)) * oldPrice

	update := db.Table("books").
			Where("book_id = ?", req.Id).
			Update("price", newPrice)

	if utils.IsNotFound(update){
		log.Fatal(update.Error)
	}

	var updateBook models.BookMainInformation

	row := db.Table("books").
			Select("book_id, author, name, price").
			Where("book_id = ?", req.Id).
			Find(&updateBook)

	if utils.IsNotFound(row){
		log.Fatal(row.Error)
	}

	res := &proto.MainResponse{
		Book: &proto.BookMainInformation{
			Id: int32(updateBook.BookId),
			Author: updateBook.Author,
			Name:  updateBook.Name,
			Price:  updateBook.Price,
		},
	}

	return res, nil
}

func (*server) Delete(c context.Context, req *proto.WorkOnlyIdRequest)(*proto.MainResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	var deleteBook models.BookMainInformation

	row := db.Table("books").
			Select("book_id, author, name, price").
			Where("book_id = ?", req.Id).
			Find(&deleteBook).
			Delete(&models.GetBook{})

	if utils.IsNotFound(row) {
		log.Fatal(row.Error)
	}

	res := &proto.MainResponse{
		Book: &proto.BookMainInformation{
			Id: int32(deleteBook.BookId),
			Author: deleteBook.Author,
			Name:  deleteBook.Name,
			Price:  deleteBook.Price,
		},
	}

	return res, nil
}

func (*server) Search(c context.Context, req *proto.SearchRequest)(*proto.ManyResponse, error){

	db, err := utils.DBConnection()

	if err != nil{
		log.Fatalf("Didn't open database %v ", err)
	}

	var books []models.BookMainInformation

	rows := db.Table("books").
			Select("book_id, author, name, price").
			Where("author ILIKE ? OR name ILIKE ?", "%" + req.GetKey() + "%", req.GetKey() + "%").
			Find(&books)

	if utils.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	var result []*proto.BookMainInformation

	for _, objekt := range books{

		res := &proto.BookMainInformation{
			Id: int32(objekt.BookId),
			Author: objekt.Author,
			Name: objekt.Name,
			Price: objekt.Price,
		}

		result = append(result, res)
	}

	return &proto.ManyResponse{
		Books: result, 
	}, nil
}

func main() {

	log.Println("Server is ready ...")

	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	proto.RegisterBookProfilesServer(s, &server{})

	err = s.Serve(lis)

	if err != nil {
		log.Fatal(err)
	}
}