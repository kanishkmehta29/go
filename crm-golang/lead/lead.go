package lead

import (
	"github.com/gofiber/fiber"
	"github.com/kanishkmehta29/crm-golang/database"
	"gorm.io/gorm"
)

type Lead struct{
	gorm.Model
	Name string `json:"name"`
	Company string `json:"company"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func GetLead(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DB
	var lead Lead
	err := db.Find(&lead,"id = ?" ,id)
	if err.Error != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lead not found",
		})
		return
	}
	c.JSON(lead)
}

func GetLeads(c *fiber.Ctx){
	db := database.DB
	var leads []Lead
	db.Find(&leads)
	c.Status(fiber.StatusOK)
	c.JSON(leads)
}

func NewLead(c *fiber.Ctx) {
	db := database.DB
	var lead Lead
	err := c.BodyParser(&lead)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"message": "Not able to create new lead",
			"error": err,
		})
		return
	}
	db.Create(&lead)
	c.JSON(lead)
}

func DeleteLead(c *fiber.Ctx){
	id := c.Params("id")
	db := database.DB
	var lead Lead
	err := db.Delete(&lead,"id = ?",id)
	if err.Error != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"message": "Not able to delete the given lead",
			"error": err,
		})
		return
	}
	c.JSON(fiber.Map{
		"message": "Sucessfully deleted the given lead",
	})
	
}