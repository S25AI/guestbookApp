package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func print(data interface{}) {
	entity, _ := json.MarshalIndent(data, "", "  ")
	log.Println(string(entity))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func validateAge(age string) error {
	userAge, err := strconv.Atoi(age)
	if err != nil {
		return errors.New("Укажите корректный возраст")
	}

	if userAge < 10 || userAge > 90 {
		return errors.New("Укажите возраст в диапазоне от 10 до 90 лет")
	}

	return nil
}

func validateUserName(userName string) error {
	if len(userName) == 0 {
		return errors.New("Поле не должно быть пустым")
	}

	if len(userName) < 4 || len(userName) > 50 {
		return errors.New("Поле возраст должно содержать от 4 до 50 символов")
	}

	return nil
}

func validateCity(userCity string) error {
	if len(userCity) == 0 {
		return errors.New("Поле не должно быть пустым")
	}

	return nil
}

func validateUserCreateFormFields(userName string, age string, city string) []string {
	rules := []ValidationRule{
		{
			FieldName: "user_name",
			Value:     userName,
			Validate:  validateUserName,
		},
		{
			FieldName: "age",
			Value:     age,
			Validate:  validateAge,
		},
		{
			FieldName: "city",
			Value:     city,
			Validate:  validateCity,
		},
	}

	fieldErrors := make([]string, 0)

	for _, rule := range rules {
		if err := rule.Validate(rule.Value); err != nil {
			fieldErrors = append(fieldErrors, fmt.Sprintf("%s: %s", rule.FieldName, err.Error()))
		}
	}

	return fieldErrors
}
