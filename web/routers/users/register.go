package users

import (
	"net/url"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/belogik/goes"
	"github.com/firstrow/logvoyage/common"
	"github.com/firstrow/logvoyage/web/context"
	"github.com/nu7hatch/gouuid"
)

type EnableValidation struct {
	Valid validation.Validation
}

func (this *EnableValidation) GetError(key string) string {
	for _, err := range this.Valid.Errors {
		if err.Key == key {
			return err.Message
		}
	}
	return ""
}

type registerForm struct {
	Email    string
	Password string
	*EnableValidation
}

func (this *registerForm) SetupValidation() {
	this.Valid.Required(this.Email, "Email")
	this.Valid.Email(this.Email, "Email")
	this.Valid.Required(this.Password, "Password")
	this.Valid.MinSize(this.Password, 5, "Password")
	this.Valid.MaxSize(this.Password, 25, "Password")
}

func Register(ctx *context.Context) {
	ctx.Request.ParseForm()
	form := &registerForm{
		EnableValidation: &EnableValidation{},
	}

	if ctx.Request.Method == "POST" {
		form.Email = ctx.Request.Form.Get("email")
		form.Password = ctx.Request.Form.Get("password")
		form.SetupValidation()

		if !form.EnableValidation.Valid.HasErrors() {
			apiKey, _ := uuid.NewV5(uuid.NamespaceURL, []byte(form.Email))

			doc := goes.Document{
				Index: "users",
				Type:  "user",
				Fields: map[string]string{
					"email":    form.Email,
					"password": com.Sha256(form.Password),
					"apiKey":   apiKey.String(),
				},
			}
			extraArgs := make(url.Values, 0)
			common.GetConnection().Index(doc, extraArgs)
		}
	}

	ctx.HTML("users/register", context.ViewData{
		"form": form,
	})
}
