package repository

import "gorm.io/gorm"

type (
	Post interface {
		Fist() error
		Find() error
		Create() error
		Update() error
		Delete() error
	}

	post struct {
		db *gorm.DB
	}
)

// Create implements Post.
func (p *post) Create() error {
	panic("unimplemented")
}

// Delete implements Post.
func (p *post) Delete() error {
	panic("unimplemented")
}

// Find implements Post.
func (p *post) Find() error {
	panic("unimplemented")
}

// Fist implements Post.
func (p *post) Fist() error {
	panic("unimplemented")
}

// Update implements Post.
func (p *post) Update() error {
	panic("unimplemented")
}

func newPost(db *gorm.DB) Post {
	return &post{db: db}
}
