package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"nhwc-server/server"
	"github.com/kabukky/httpscerts"
	"net/http"
	"fmt"
)

func main() {
	e := echo.New()
	e.Pre()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.DELETE, echo.POST},
	}))

	e.GET("/ws", server.OnClientConnect)
	//e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//	return func(c echo.Context) error {
	//		fmt.Println(c.Request().Host+c.Request().RequestURI)
	//		if c.Request().RequestURI == "/nhwc" {
	//			c.Redirect(301, "http://" + c.Request().Host+c.Request().RequestURI + "/")
	//			//c.Redirect(301, "http://" + c.Request().Host+"/nhwc1")
	//		} else {
	//			next(c)
	//		}
	//		return nil
	//	}
	//})
	e.Static("/", "client/web-desktop")
	//e.File("/nhwc", "client/web-desktop/index.html")
	err := httpscerts.Check("cert.pem", "key.pem")
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", ":1323")
	}
	//e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":1323"))

	//http.HandleFunc("/ws", handler)
	//log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	//err = http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	//log.Fatal(err)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("666")
	fmt.Fprintf(w, "Hi there!")
}
