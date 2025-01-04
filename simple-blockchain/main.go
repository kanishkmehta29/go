package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Block struct{
	Pos int
	Data BookCheckout
    TimeStamp string
	Hash string
    PrevHash string
}

type Book struct{
	Id string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	PublishDate string `json:"publish_date"`
	Isbn string `json:"isbn"`
}

type BookCheckout struct{
	BookId string `json:"book_id"`
	User string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis bool `json:"is_genesis"`
} 

type Blockchain struct{
	blocks *[]Block
	
}

var bc *Blockchain

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := (*bc.blocks)[len(*bc.blocks)-1]
	var newBlock Block

	newBlock.Data = data
	newBlock.Pos = prevBlock.Pos + 1
	newBlock.TimeStamp = time.Now().String()
	newBlock.PrevHash = prevBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
     
	if validBlock(newBlock,prevBlock){
		*bc.blocks = append(*bc.blocks, newBlock)
	}
	
}

func validBlock(currBlock Block,prevBlock Block) bool{
	return (currBlock.PrevHash == prevBlock.Hash && currBlock.Pos == prevBlock.Pos+1 && currBlock.validateHash())
}

func (blk *Block) validateHash()bool{
	return calculateHash(*blk) == blk.Hash
}

func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Pos, block.Data, block.TimeStamp, block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func newBook(c *fiber.Ctx) error{
	var book Book
	err := c.BodyParser(&book)
    if err != nil{
		c.Status(http.StatusBadRequest)
		c.JSON(fiber.Map{"message":"Error parsing the request","error":err})
		return err
	}
	book.Id = fmt.Sprintf("%d", rand.Intn(100000))
	c.Status(http.StatusOK)
	c.JSON(book)
	return nil
}

func writeBlock(c *fiber.Ctx)error{
	var checkoutitem BookCheckout
	err := c.BodyParser(&checkoutitem)
	if err != nil{
		c.Status(http.StatusBadRequest)
		c.JSON(fiber.Map{"message":"Error parsing the request","error":err})
		return err
	}
	bc.AddBlock(checkoutitem)
	return nil
}

func getBlockchain(c *fiber.Ctx) error {
	c.JSON(bc.blocks)
	c.Status(http.StatusOK)
	return nil
}

func main(){

	bc = 
	&Blockchain{blocks: &[]Block{
		{
			Pos:       0,
			Data:      BookCheckout{IsGenesis: true},
			TimeStamp: time.Now().String(),
			Hash:      calculateHash(Block{Pos: 0, Data: BookCheckout{IsGenesis: true}, TimeStamp: time.Now().String()}),
			PrevHash:  "",
		},
	}}

 app := fiber.New()

 app.Get("/",getBlockchain)
 app.Post("/",writeBlock)
 app.Post("/new",newBook)


 fmt.Println("Starting server on port 8080")
 app.Listen(":8080")

}