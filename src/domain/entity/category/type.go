package category

import (
	"DX/src/utils"
	"errors"
	"strings"
)

type Type int

const (
	PickUpAndDelivery Type = iota
	Shopping
	Transportation
	MovingServices
	Gardening
	BabySitting
	Laundry
	Cleaning
	ReadingAndWriting
	FashionAndTailoring
	AutoServices
	Catering
)

var Tasks = []Type{PickUpAndDelivery, Shopping, Transportation, MovingServices}
var Services = []Type{Gardening, BabySitting, Laundry, Cleaning, ReadingAndWriting, FashionAndTailoring, AutoServices, Catering}

func IsValidCategory(errandType string, category string) (*Type, error) {
	nCategory := getCategory(category)
	if errandType == "task" {
		if utils.Contains(Tasks, nCategory) {
			return &nCategory, nil
		}
		return nil, errors.New("invalid task errand category")
	} else if errandType == "service" {
		if utils.Contains(Services, nCategory) {
			return &nCategory, nil
		}
		return nil, errors.New("invalid service errand category")
	} else {
		return nil, errors.New("invalid errand type")
	}
}

func getCategory(category string) Type {
	if category == PickUpAndDelivery.Id() {
		return PickUpAndDelivery
	}
	if category == Shopping.Id() {
		return Shopping
	}
	if category == Transportation.Id() {
		return Transportation
	}
	if category == MovingServices.Id() {
		return MovingServices
	}
	if category == Gardening.Id() {
		return Gardening
	}
	if category == BabySitting.Id() {
		return BabySitting
	}
	if category == Laundry.Id() {
		return Laundry
	}
	if category == Cleaning.Id() {
		return Cleaning
	}
	if category == ReadingAndWriting.Id() {
		return ReadingAndWriting
	}
	if category == FashionAndTailoring.Id() {
		return FashionAndTailoring
	}
	if category == AutoServices.Id() {
		return AutoServices
	}
	if category == Catering.Id() {
		return Catering
	}
	return -1
}

func (c Type) Id() string {
	if c == PickUpAndDelivery {
		return "pickup-and-delivery"
	}
	if c == Shopping {
		return "shopping"
	}
	if c == Transportation {
		return "transportation"
	}
	if c == MovingServices {
		return "moving-services"
	}
	if c == Gardening {
		return "gardening"
	}
	if c == BabySitting {
		return "baby-sitting"
	}
	if c == Laundry {
		return "laundry"
	}
	if c == Cleaning {
		return "cleaning"
	}
	if c == ReadingAndWriting {
		return "reading-and-writing"
	}
	if c == FashionAndTailoring {
		return "fashion-and-tailoring"
	}
	if c == AutoServices {
		return "auto-services"
	}
	if c == Catering {
		return "catering"
	}
	return ""
}

func (c Type) String() string {
	return strings.ToTitle(c.Id())
}
