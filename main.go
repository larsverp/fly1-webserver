package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type SetIP struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}

func main() {
	godotenv.Load()
	dbSting := os.Getenv("DB_CONNECT_STRING")

	db, err := sql.Open("mysql", dbSting)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	router.POST("/set-ip", func(c *gin.Context) {
		var request SetIP
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		var id string
		err = db.QueryRow("SELECT id FROM info WHERE id = ?", request.ID).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = db.Exec("INSERT INTO info (id, ip) VALUES (?, ?)", request.ID, request.IP)
				if err != nil {
					c.JSON(500, gin.H{
						"error": err.Error(),
					})
					return
				}
			}
		} else {
			_, err = db.Exec("UPDATE info SET ip = ? WHERE id = ?", request.IP, request.ID)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		db.Close()

		c.JSON(200, gin.H{
			"message": "success",
		})

	})

	router.GET("/get-ip", func(c *gin.Context) {
		id := c.Query("id")
		fmt.Println(id)
		stmt, err := db.Prepare("SELECT ip FROM info WHERE id = ?")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		row := stmt.QueryRow(id)
		var ip string
		err = row.Scan(&ip)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		db.Close()
		c.JSON(200, gin.H{
			"ip": ip,
		})

	})

	router.Run()
}
