package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/asaskevich/govalidator"
	"github.com/kanishkmehta29/url-shortner/database"
	"github.com/kanishkmehta29/url-shortner/helpers"
)

type request struct{
	Url string `json:"url"`
	CustomShort string `json:"short"`
	Expiry time.Duration `json:"expiry"`
}

type response struct{
	Url string `json:"url"`
	CustomShort string `json:"short"`
	Expiry time.Duration `json:"expiry"`
	XRateRemaining int `json:"rate_limit"`
	XRateLimiteReset time.Duration `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx)error{
	body := new(request)
	err :=c.BodyParser(&body)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
																"message":"Error parsing request",
																"error":err,
																})
	}

	// implement rate limit
	r2 := database.CreateClient(1)
	defer r2.Close()
	val,err := r2.Get(database.Ctx,c.IP()).Result()
	if err == redis.Nil{
		_ = r2.Set(database.Ctx,c.IP(),os.Getenv("API_QUOTA"),30*time.Minute).Err()
	}else {
		valInt,_ := strconv.Atoi(val)
		if valInt <= 0{
			limit,_ := r2.TTL(database.Ctx,c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"message":"Rate limit exceeded",
				"rate_limit_rest":limit/time.Nanosecond/time.Minute,
			})
		}
	}


	// check if the input is an actual URL
	if !govalidator.IsURL(body.Url){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
																"message":"Invalid URL",
																})
	}


	// check for domain error
	if !helpers.RemoveDomainError(body.Url){
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message":"Domain error",
			})
	}


	// enforce https,SSL
	body.Url = helpers.EnforceHTTP(body.Url)
	var id string

	if body.CustomShort == ""{
		id = uuid.NewString()[:6]
	}else{
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	found_val ,_ := r.Get(database.Ctx,id).Result()
	if found_val != ""{
		c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message":"URL already taken",
		})
	}

	if body.Expiry == 0{
		body.Expiry = 24
	}

	err = r.Set(database.Ctx,id,body.Url,body.Expiry*3600*time.Second).Err()

	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":"error occured while creating short link",
		})
	}

	resp := response{
		Url: body.Url,
		CustomShort: "",
		Expiry: body.Expiry,
		XRateRemaining : 10,
	    XRateLimiteReset : 30,
	}

	r2.Decr(database.Ctx,c.IP())

	val,_ = r2.Get(database.Ctx,c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl,_ := r2.TTL(database.Ctx,c.IP()).Result()
	resp.XRateLimiteReset = ttl/time.Nanosecond/time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Short url created sucessfully",
		"response":resp,
	})


}