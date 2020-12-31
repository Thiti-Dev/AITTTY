package helpers

import (
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)


var validatorEntity *validator.Validate = validator.New()

var translator = en.New()
var uni = ut.New(translator, translator)
var trans ut.Translator

type dynamicStruct = map[string]interface{}

// InitializeTranslator -> A setting up for translator
func InitializeTranslator(){
	transS, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	trans = transS

	// Registering Default Translation
	if err := en_translations.RegisterDefaultTranslations(validatorEntity, trans); err != nil {
		log.Fatal(err)
	}
	// ────────────────────────────────────────────────────────────────────────────────


	_ = validatorEntity.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validatorEntity.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
}

// ValidateStructAndGetErrorMsg -> uses for validating any struct and returns if there is any error or not according to the tags rule
func ValidateStructAndGetErrorMsg(data interface{}) (bool,dynamicStruct){
	err := validatorEntity.Struct(data)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	if err != nil {

		validateData := make(dynamicStruct)

		for _, e := range err.(validator.ValidationErrors) {
			validateData[e.Field()] = e.Translate(trans)
		}
		return false , validateData
	}else{
		return true , nil
	}
}