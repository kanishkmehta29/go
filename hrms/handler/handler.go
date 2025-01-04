package handler

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanishkmehta29/hrms/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
func CreateEmployee(c *fiber.Ctx) error{
	client := database.Global_client

	collection := client.Database("hrms_database").Collection("employees")
	var new_employee database.Employee

	err := c.BodyParser(&new_employee)
	if err != nil{
		log.Fatalf("Error in parsing body %v\n",err)
	}
    
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
    defer cancel()

    result,err := collection.InsertOne(ctx,new_employee)
	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		c.JSON(fiber.Map{"error": "Failed to create new employee entry"})
	}
	c.Status(fiber.StatusAccepted)
	c.JSON(fiber.Map{"message":"Sucessfully Created new entry","New Entry":result})
	return err
}

func GetEmployee(c *fiber.Ctx) error {
	client := database.Global_client

	collection := client.Database("hrms_database").Collection("employees")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.JSON(fiber.Map{"error": "Failed to fetch employees"})
		return err
	}
	defer cursor.Close(ctx)

	var employees []database.Employee
	err = cursor.All(ctx, &employees)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.JSON(fiber.Map{"error": "Failed to parse employees"})
		return err
	}

	c.Status(fiber.StatusOK)
	c.JSON(employees)
	return nil
}

func DeleteEmployee(c *fiber.Ctx) error{
	client := database.Global_client
	id := c.Params("id")
	collection := client.Database("hrms_database").Collection("employees")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{"error": "Invalid ID format"})
		return err
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.Status(fiber.StatusNotFound)
		c.JSON(fiber.Map{"error": "Employee not found"})
		return err
	}

	c.Status(fiber.StatusOK)
	c.JSON(fiber.Map{
		"message": "record deleted successfully",
		"deleted record": res,
	})
	return err
}

func EditEmployee(c *fiber.Ctx) error{
   client := database.Global_client
   id := c.Params("id")
   objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{"error": "Invalid ID format"})
		return err
	}
   collection := client.Database("hrms_database").Collection("employees")

   ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
   defer cancel() 

	var duplicate_body database.Employee
	err = collection.FindOne(ctx, bson.M{"_id": objid}).Decode(&duplicate_body)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		c.JSON(fiber.Map{"error": "Employee not found"})
		return err
	}

	var updated_body database.Employee
	err = c.BodyParser(&updated_body)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{"error": "Error parsing the request"})
		return err
	}

	if updated_body.Name == "" {
		updated_body.Name = duplicate_body.Name
	}

	if updated_body.Salary == 0 {
		updated_body.Salary = duplicate_body.Salary
	}

	if updated_body.Age == 0 {
		updated_body.Age = duplicate_body.Age
	}

	err = c.BodyParser(&updated_body)
	if err != nil{
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{"error": "Error parsing the request"})
		return err
	}

	update := bson.M{
		"$set": updated_body,
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": objid}, update)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		c.JSON(fiber.Map{
			"message": "Employee not found",
            "error": err,
	})
		return err
	}
    
	c.Status(fiber.StatusOK)
	c.JSON(fiber.Map{
		"message": "record updated sucessfully",
        "updated record":res,
	})
	return err
}