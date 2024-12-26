//основные сущности и их валидация

package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

const (
	keywordtMaxLenth = 32
	keywordtMinLenth = 1

	categoryMaxLenth = 32
	categoryMinLenth = 1

	descriptionMaxLenth = 255
	descriptionMinLenth = 1

	pattern = `^[a-zA-Zа-яА-Я]+$`
)

type ProductInfo struct {
	ProductId       int      `json:"userId" db:"productId"`
	Category        string   `json:"category" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Status          string   `json:"status" binding:"required"`
	ProductKeyWords []string `json:"productKeyWords" binding:"required"`
}

func ValidateProductId(prId int) error {

	if prId < 0 {
		return errors.New("invalid user id: can`t be less 0")
	}

	return nil
}

func ValidateCategory(category string) error {
	re := regexp.MustCompile(pattern)

	if !re.MatchString(category) {
		return fmt.Errorf("invalid category: %s does not match regexp", category)
	}
	if len(category) > categoryMaxLenth {
		return errors.New("invalid category: too long, max length is " + strconv.Itoa(categoryMaxLenth))
	}

	if len(category) < categoryMinLenth {
		return errors.New("invalid category: too short, min length is " + strconv.Itoa(categoryMinLenth))
	}

	return nil
}

func ValidateDescription(description string) error {

	if len(description) > descriptionMaxLenth {
		return fmt.Errorf("%s %s", "invalid description: too long, max length is",
			strconv.Itoa(descriptionMaxLenth))
	}

	if len(description) < descriptionMinLenth {
		return fmt.Errorf("%s %s", "invalid description: too short, min length is",
			strconv.Itoa(descriptionMinLenth))
	}

	return nil
}

func ValidateProductKeyWord(keyword string) error {
	re := regexp.MustCompile(pattern)

	if !re.MatchString(keyword) {
		return fmt.Errorf("invalid keyword: %s does not match regexp", keyword)
	}

	if len(keyword) > keywordtMaxLenth {
		return errors.New("invalid keyword: too long, max length is " + strconv.Itoa(keywordtMaxLenth))
	}

	if len(keyword) < keywordtMinLenth {
		return errors.New("invalid keyword: too short, min length is " + strconv.Itoa(keywordtMinLenth))
	}

	return nil

}

func (pr *ProductInfo) ValidateProductInfo() error {

	if err := ValidateCategory(pr.Category); err != nil {
		return err
	}

	if err := ValidateDescription(pr.Description); err != nil {
		return err
	}

	if err := ValidateProductId(pr.ProductId); err != nil {
		return err
	}

	for _, keyword := range pr.ProductKeyWords {
		if err := ValidateProductKeyWord(keyword); err != nil {
			return err
		}
	}

	return nil
}
