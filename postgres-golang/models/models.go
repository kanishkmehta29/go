package models

import (
	"database/sql"
	"fmt"
	"os"
	"log"
    _ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type Stock struct {
	StockId int64 `json:"stockid"`
	Name string `json:"name"`
	Price int64 `json:"price"`
	Company string `json:"company"`
}

func createConnection() *sql.DB{
  err := godotenv.Load(".env")
  if err != nil{
	log.Fatalf("%v\n",err)
  }
  db,err := sql.Open("postgres",os.Getenv("DATABASE_URL"))
   if err != nil{
	log.Fatalf("%v\n",err)
   }

   err = db.Ping()

   if err != nil{
	log.Fatalf("%v\n",err)
   }

   fmt.Println("Successfully connected to database")
   
   return db
}

func InsertStock(s Stock) int64{
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT into stocks(name,price,company) VALUES($1,$2,$3) RETURNING StockId`
    
	var id int64
	err := db.QueryRow(sqlStatement,s.Name,s.Price,s.Company).Scan(&id)

    if err != nil{
		log.Fatalf("%v\n",err)
	}

	fmt.Printf("Inserted the stock with id %v",id)
    return id
}

func GetStock(id int64) (Stock,error){
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM stocks where StockId = $1`
    
	var s Stock
	row := db.QueryRow(sqlStatement,id)

	err := row.Scan(&s.StockId,&s.Name,&s.Price,&s.Company)

	switch err{
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return s,nil
		
	case nil:
        return s,nil
		
	default:
		return s,err
	}

}

func GetAllStock() ([]Stock,error){
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT * FROM stocks`
    
	var stocks []Stock
	rows,err := db.Query(sqlStatement)

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	for rows.Next(){
		var stock Stock
		err := rows.Scan(&stock.StockId,&stock.Name,&stock.Price,&stock.Company)

		if err != nil{
			log.Fatalf("%v\n",err)
		}

		stocks = append(stocks,stock)
	}
    return stocks,err
}

func UpdateStock(id int64,s Stock) error{

	db := createConnection()
	defer db.Close()

	old_stock,err := GetStock(id)

	if err != nil{
		return err
	}

	new_stock := old_stock

	if s.Name != ""{
		new_stock.Name = s.Name
	}
	if s.Price != 0{
		new_stock.Price = s.Price
	}
	if s.Company != ""{
		new_stock.Company = s.Company
	}

	DeleteStock(id)
	InsertStock(new_stock)

	fmt.Println("Stock Updated Sucessfully")
	return nil

}

func DeleteStock(id int64) error{
    db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks where StockId = $1`
	_,err := db.Exec(sqlStatement,id)
	return err

}