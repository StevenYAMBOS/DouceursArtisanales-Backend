package controllers

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stevenyambos/hmmm-backend/models"
	"github.com/stevenyambos/hmmm-backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = os.Getenv("SECRET_KEY")

// Middleware pour rendre la route /user-profile privée.
func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        cookie := c.Cookies("jwt")
        if cookie == "" {
            // Aucun token trouvé dans les cookies, renvoie une erreur d'authentification
            c.Status(fiber.StatusUnauthorized)
            return c.JSON(fiber.Map{"message": "Authentification requise pour accéder à cet espace."})
        }

        // Vérifie et parse le token
        _, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
            // Clé secrète pour la vérification de la signature
            return []byte(SecretKey), nil
        })

        if err != nil {
            // Erreur de vérification du token, renvoie une erreur d'authentification
            c.Status(fiber.StatusUnauthorized)
            return c.JSON(fiber.Map{"message": "Token JWT invalide."})
        }

        // Continue vers le gestionnaire de route si l'authentification réussit
        return c.Next()
    }
}

func Home(c *fiber.Ctx) error {
	return c.SendString("Bienvenue sur Douceurs Artisanales !")
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Le _ (trait de soulignement) est l'identifiant vide en Go.
	//Nous pouvons utiliser l'identificateur vide pour déclarer et utiliser les variables inutilisées. Les variables inutilisées sont les variables qui ne sont définies qu'une seule fois dans le programme et qui ne sont pas utilisées. Ces variables affectent la lisibilité du programme et le rendent illisible. La particularité de Go est qu'il s'agit d'un langage lisible et concis. Il ne permet jamais qu'une variable soit définie et jamais utilisée ; lorsque la variable est définie et non utilisée dans un programme, il lève une erreur.

	_, err = database.DB.Exec(`
		INSERT INTO user_account (username, email, password)
		VALUES ($1, $2, $3)
	`, data["username"], data["email"], passwordHash)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Enregistrement de l'utilisateur réussie !"})
}

// Connexion
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT user_id, username, email, password FROM user_account
		WHERE email = $1
	`, data["email"]).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "L'utilisateur n'existe pas."})
	} else if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Mot de passe incorrect.",
		})
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    strconv.Itoa(int(user.ID)),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Erreur lors de la connexion.",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Connexion au compte réussi.",
	})
}

func User(c *fiber.Ctx) error {
    cookie := c.Cookies("jwt")
    token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    })
    if err != nil {
        fmt.Println(err)
        c.Status(fiber.StatusUnauthorized)
        return c.JSON(fiber.Map{
            "message": "Authentification échouée.",
        })
    }
    claims := token.Claims.(*jwt.RegisteredClaims)

    // Utilisation de QueryRow pour récupérer les informations de l'utilisateur
    row := database.DB.QueryRow(`
        SELECT id, username, email FROM user_account WHERE id = $1
    `, claims.Issuer)

    var user models.User
    err = row.Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        fmt.Println(err)
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "message": "Erreur lors de la récupération des informations de l'utilisateur.",
        })
    }

    return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Vous êtes déconnecté.",
	})
}