package models

type Book struct{
	Author string	`json:"author"`
	Name string		`json:"name"`
	Price float32	`json:"price"`
	Genre string	`json:"genre"`
	Cover string	`json:"cover"`
	Page int32 		`json:"page"`
}

type Discount struct{
	Id int32 		`json:"id"`
	Percent int32 	`json:"percent"`
}

type Id struct{
	Id int32 	`json:"id"`
}

type GetBook struct{
	BookId int32
	Author string
	Name string		
	Price float32
	Genre string	
	Cover string	
	Page int32 			
}

type BookMainInformation struct{
	BookId int32
	Author, Name string
	Price float32
}