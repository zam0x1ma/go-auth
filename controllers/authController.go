package controllers

import (
    "net/http"
    "strconv"
    "time"

    "go-auth/database"
    "go-auth/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
)

const SecretKey = "vU46%^sC+H\\(wpVfiZ?dADedF.SRxcPpW=I`uGXK]a*))QR=S2"

type UserRequest struct {
    Name        string  `json:"name"`
    Email       string  `json:"email"`
    Password    string  `json:"password"`
}

func Register(c *gin.Context) {
    var userRequest UserRequest
    if err := c.ShouldBindJSON(&userRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 14)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    newUser := models.User{
        Name:       userRequest.Name,
        Email:      userRequest.Email,
        Password:   string(password),
    }

    database.DB.Create(&newUser)

    c.JSON(http.StatusOK, gin.H{
        "status": "OK",
    })
}

func Login(c *gin.Context) {
    var userRequest UserRequest
    if err := c.ShouldBindJSON(&userRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    } 

    var user models.User

    database.DB.Where("email = ?", userRequest.Email).First(&user)

    if user.Id == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "unauthorized",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "unauthorized",
        })
        return
    }

    claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        Issuer: strconv.Itoa(int(user.Id)),
        ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    })

    token, err := claims.SignedString([]byte(SecretKey))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
    }

    c.SetCookie("jwt", token, 3600 * 24, "/", "localhost", false, true)

    c.JSON(http.StatusOK, gin.H{
        "status": "OK",
    })
}

func User(c *gin.Context) {
    cookie, err := c.Cookie("jwt")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "unauthorized",
        })
        return 
    }

    token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "unauthorized",
        })
        return 
    }

    claims := token.Claims.(*jwt.StandardClaims)

    var user models.User

    database.DB.Where("id = ?", claims.Issuer).First(&user)

    c.JSON(http.StatusOK, gin.H{
        "name": user.Name,
        "email": user.Email,
    })
}

func Logout(c *gin.Context) {
    c.SetCookie("jwt", "", -1, "/", "localhost", false, true)

    c.JSON(http.StatusOK, gin.H{
        "status": "OK",
    })
}
