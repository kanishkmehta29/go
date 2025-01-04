package router

import(
  "github.com/gofiber/fiber/v2"
  "github.com/kanishkmehta29/hrms/handler"
)

func ManageRoutes(app *fiber.App){
	app.Get("employee",handler.GetEmployee)
	app.Post("employee",handler.CreateEmployee)
	app.Put("employee/:id",handler.EditEmployee)
	app.Delete("employee/:id",handler.DeleteEmployee)
}