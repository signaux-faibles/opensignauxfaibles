package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"net/smtp"
	"opensignauxfaibles/dbmongo/lib/engine"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// login
type login struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	BrowserToken string `form:"browserToken" json:"browserToken"`
	CheckCode    string `form:"checkCode" json:"checkCode"`
}

// AdminUser object utilisateur mongodb
type AdminUser struct {
	ID             engine.AdminID `json:"_id" bson:"_id"`
	HashedPassword []byte         `json:"hashedPassword,omitempty" bson:"hashedPassword,omitempty"`
	HashedRecovery []byte         `json:"hashedRecovery,omitempty" bson:"hashedRecovery,omitempty"`
	TimeRecovery   time.Time      `json:"timeRecovery" bson:"timeRecovery"`
	HashedCode     []byte         `json:"hashedCode,omitempty" bson:"hashedCode,omitempty"`
	TimeCode       time.Time      `json:"timeCode,omitempty" bson:"timeCode,omitempty"`
	Cookies        []string       `json:"cookies" bson:"cookies"`
	Level          string         `json:"level" bson:"level"`
	FirstName      string         `json:"firstName" bson:"firstName"`
	LastName       string         `json:"lastName" bson:"lastName"`
	BrowserTokens  []string       `json:"browserTokens" bson:"browserTokens"`
	Regions        []string       `json:"regions" bson:"regions"`
}

func (user AdminUser) save() error {
	err := engine.Db.DBStatus.C("Admin").Update(bson.M{"_id": user.ID}, user)
	return err
}

// type AdminLevel string
const levelAdmin = "admin"
const levelPowerUser = "powerUser"
const levelUser = "user"

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	email := claims["id"].(string)
	user, err := loadUser(email)
	if err != nil {
		c.JSON(500, "Erreur d'identification")
	}
	return &user
}

func loginUserWithCredentials(username string, password string) (AdminUser, error) {
	var user AdminUser
	if err := engine.Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": username}).One(&user); err != nil {
		return AdminUser{}, err
	}
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == nil {
		return user, nil
	}
	return AdminUser{}, err
}

func loadUser(email string) (AdminUser, error) {
	var user AdminUser
	if err := engine.Db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": email}).One(&user); err != nil {
		return AdminUser{}, err
	}
	return user, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*AdminUser); ok && v.Level == "admin" {
		return true
	}
	fmt.Println(data)
	return true
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func unauthorizedHandler(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(AdminUser); ok {
		return jwt.MapClaims{
			"id": v.ID.Key,
		}
	}
	return jwt.MapClaims{}
}

func getCode() int {
	i, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return int(i.Int64())
}

func getRegions() map[string]string {
	regions := map[string]string{
		"ARA": "Auvergne-Rhône-Alpes",
		"BFC": "Bourgogne-Franche-Comté",
		"BRE": "Bretagne",
		"CVL": "Centre-Val de Loire",
		"COR": "Corse",
		"GES": "Grand Est",
		"HDF": "Hauts-de-France",
		"IDF": "Île-de-France",
		"NOR": "Normandie",
		"NAQ": "Nouvelle-Aquitaine",
		"OCC": "Occitanie",
		"PDL": "Pays de la Loire",
		"PAC": "Provence-Alpes-Côte d'Azur",
	}
	return regions
}

//
// @summary Envoie un code de vérification par EMail
// @description Le code n'est envoyé que si le mot de passe est valide.
// @description Le cas échéant, un email avertissant d'une tentative est envoyé.
// @description Pour éviter les tentatives de forçage de mot de passe, ce service ne renvoie jamais d'échec.
// @Tags Authentification
// @accept  json
// @produce  json
// @Param email query string true "Adresse EMail"
// @Param password query string true "Mot de passe"
// @Success 200 {string} string "ok"
// @Router /login/get [post]
func sendMail(sender string, recipient string, title string, body string) error {
	smtpAddress := viper.GetString("smtpAddress")
	smtpConnection, err := smtp.Dial(smtpAddress)
	if err != nil {
		return err
	}
	defer smtpConnection.Close()
	smtpConnection.Mail(sender)
	smtpConnection.Rcpt(recipient)

	emailObject, err := smtpConnection.Data()
	if err != nil {
		return err
	}
	defer emailObject.Close()

	return nil
}

func loginGetHandler(c *gin.Context) {
	var loginVals login

	if err := c.ShouldBind(&loginVals); err != nil {
		c.JSON(401, "requête malformée")
		return
	}

	loginGet(loginVals)

	c.JSON(200, "ok")
}

func loginGet(login login) error {
	email := login.Email
	password := login.Password
	user, err := loginUserWithCredentials(email, password)

	mailTemplate, _ := template.New("loginMail").Parse(`
	Subject: Signaux-Faibles – votre code de vérification
	Content-Type: text/plain; charset=us-ascii; format=flowed
	Content-Transfer-Encoding: 7bit

	Bonjour,
	suite à votre tentative d'identification sur l'applicatif Signaux Faibles, voici votre code de vérification:
	{{.CheckCode}}

	Cordialement,
	l'équipe Signaux-Faibles.

	ps: si vous n'êtes pas à l'origine de cette tentative, nous vous prions d'en faire part à l'adresse contact@signaux-faibles.beta.gouv.fr`)

	if err == nil {
		code := struct{ checkCode string }{}
		code.checkCode = fmt.Sprintf("%06d", getCode())
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(code.checkCode), bcrypt.DefaultCost)
		if err == nil {
			user.HashedCode = hashedCode
			user.TimeCode = time.Now()
			err = user.save()
			if err == nil {
				fmt.Println(code.checkCode)

				smtpAddress := viper.GetString("smtpAddress")
				//smtpUser := viper.GetString("smtpUser")
				//smtpPassword := viper.GetString("smtpPass")

				c, err := smtp.Dial(smtpAddress)
				if err != nil {
					return err
				}
				defer c.Close()

				// Set the sender and recipient.
				c.Mail("Signaux Faibles <do.not.reply@signaux.faibles.fr>")
				c.Rcpt(email)

				// Send the email body.
				wc, err := c.Data()
				if err != nil {
					spew.Dump(err)
				}

				mailTemplate.Execute(wc, code)
				wc.Close()
				return err
			}
		}
	} else {
		fmt.Println("error: " + err.Error())
	}
	return err
}

func loginCheck(email string, password string, checkCode string) error {
	user, err := loginUserWithCredentials(email, password)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedCode, []byte(checkCode))
	if err != nil {
		return err
	}

	user.HashedCode = nil
	user.TimeCode = time.Time{}
	err = user.save()
	return err
}

//
// @summary Rafraichir le jeton d'identification
// @description Nécessite 3 informations: email, mot de passe et jeton de navigateur
// @description Le jeton de navigateur n'a pas de limite de validité et peut être conservé
// @Tags Authentification
// @accept  application/json
// @produce  application/json
// @Param email query string true "Adresse Email" "gnigni"
// @Param password query string true "Mot de Passe"
// @Param browserToken query string true "Token navigateur"
// @Success 200 {string} string "{<br/>'code': 200,<br/>'expire':'2019-01-21T11:16:15+01:00',<br/>'token':'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…'<br/>}"
// @Failure 401 {string} string "{<br/>'code': 401,<br/>'message': 'incorrect Username or Password'<br/>}"
// @Router /login [post]
func dummyLogin() {}

//
// @summary Obtenir un jeton d'identification
// @description Fournit un jeton avec nouvelle date de validité en échange d'un jeton encore valide
// @Tags Authentification
// @accept  application/json
// @produce  application/json
// @Security ApiKeyAuth
// @Success 200 {string} string "{<br/>'code': 200,<br/>   'expire': '2019-01-21T12:16:50+01:00',<br/>'token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9…'<br/>}"
// @Failure 401 {string} string "{<br/>'code': 401,<br/>'message': 'cookie token is empty'<br/>}"
// @Router /api/refreshToken [get]
func dummyRefreshtoken() {}
