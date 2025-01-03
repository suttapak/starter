package service

import "labostack/internal/repository"

type Post interface {
	Fist() error
	Find() error
	Create() error
	Update() error
	Delete() error
}

type post struct {
	postRepo repository.Post
}

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

func newPost(postRepo repository.Post) Post {
	return &post{
		postRepo: postRepo,
	}
}
