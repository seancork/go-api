package main

import (
 "fmt"
 "database/sql"
 "net/http"
 "os"
 _ "github.com/go-sql-driver/mysql" 
 "github.com/gin-gonic/gin"
)

func main() {
     port := os.Getenv("PORT")
	 fmt.Println("running")

          db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	  if err != nil {
            panic("failed to connect database")
       }
	defer db.Close()

	err = db.Ping()
    if err != nil {
	 fmt.Println(err)
}

 type Product struct {
   Id int
   P_desc string
   Price int 
 }
 r := gin.Default()

//Get product
r.GET("/product/:id", func(c *gin.Context) {
	var (
        product Product
		)
    id := c.Param("id")
    	row := db.QueryRow("select id, p_desc, price from product where id = ?", id)
		err := row.Scan(&product.Id, &product.P_desc, &product.Price)
 fmt.Println(err)
    if (err != nil) {
     c.JSON(http.StatusNotFound, gin.H{"status" : http.StatusNotFound, "message" : "No Products found!"})
  } else {
     c.JSON(http.StatusOK, gin.H{"status" : http.StatusOK, "count" : 1, "data" : product})
  }
})

//get all products
r.GET("/all_products", func(c *gin.Context) {
  var (
  	     product Product
  	     products []Product
  	     )
    rows, err := db.Query("SELECT id, p_desc, price FROM product;")
         if err != nil {
	    fmt.Println(err)
         }

    for rows.Next() {
       err = rows.Scan(&product.Id, &product.P_desc, &product.Price)
       products = append(products, product)
        if err != nil {
	    fmt.Println(err)
         }
    }
    defer rows.Close()
    if len(products) <= 0 {
	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No products found!"})
	}else{
  c.JSON(http.StatusOK, gin.H{"status" : http.StatusOK, "data": products, "count": len(products) })}
})

//Delete a product
r.DELETE("/product/:id", func (c *gin.Context) {
	id := c.Param("id")
    
    stmt, err := db.Prepare("delete from product where id= ?;")
    if err != nil {
	   fmt.Println(err)
     }
     _, err = stmt.Exec(id)
     c.JSON(http.StatusOK, gin.H{"status" : http.StatusOK, "message" : "Product deleted successfully!"})
 })

//POST - add new product to database
r.POST("product", func (c *gin.Context) {
  p_desc := c.PostForm("p_desc")
  price := c.PostForm("price")
  stmt, err := db.Prepare("insert into product(p_desc, price) values(?,?);")
  if err != nil {
	   fmt.Println(err)
     }

    res, err := stmt.Exec(p_desc, price)
    lid, err := res.LastInsertId()
    if err != nil {
	   fmt.Println(err)
     }
     defer stmt.Close()

    c.JSON(http.StatusCreated, gin.H{"status" : http.StatusCreated, "message" : "Product created successfully!", "Id": lid})

 //PUT - update a product
 r.PUT("/product/:id", func (c *gin.Context) {
  id := c.Param("id")
  p_desc := c.PostForm("p_desc")
  price := c.PostForm("price")
    
  stmt, err := db.Prepare("update product set p_desc = ?, price = ? where id = ?;")
  if err != nil {
	   fmt.Println(err)
     }
     _,err = stmt.Exec(p_desc, price, id)
    
    if err != nil {
	   fmt.Println(err)
     }

 c.JSON(200, gin.H{"message": "Successfully updated"})
    })
  })
r.Run(":" + port)
}
