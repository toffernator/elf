package store

import "errors"

type Product struct {
	Name     string
	Url      string
	Price    float32
	Currency string
}

type Wishlist struct {
	Id       int
	Products []Product
}

type User struct {
	Id        int
	FirstName string
	LastName  string
	Wishlists []Wishlist
}

type Store interface {
	GetUserById(id int) (User, error)
	CreateUser(firstName string, lastName string) User
	CreateWishlist(userId int, products []Product) Wishlist
}

type ArrayStore struct {
	userIdCounter     int
	wishlistIdCounter int
	users             []User
}

func (s *ArrayStore) Seed() {
	products := []Product{
		{Name: "iPad", Url: "www.example.com", Price: 100, Currency: "eur"},
		{Name: "Macbook", Url: "www.example.com", Price: 200, Currency: "eur"},
	}
	wishlists := []Wishlist{
		{Id: 1, Products: products},
		{Id: 2, Products: products[:1]},
		{Id: 3, Products: products[0:]},
	}
	users := []User{
		{
			Id:        0,
			FirstName: "parmesan",
			LastName:  "ehrmanntraut",
			Wishlists: wishlists,
		},
		{
			Id:        1,
			FirstName: "fenya",
			LastName:  "mozzerella",
			Wishlists: []Wishlist{
				wishlists[0],
				wishlists[1],
			},
		},
	}
	s.users = append(s.users, users...)
	s.userIdCounter = 2
	s.wishlistIdCounter = 4
}

func (s *ArrayStore) GetUserById(id int) (User, error) {
	for _, user := range s.users {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}

func (s *ArrayStore) CreateUser(firstName string, lastName string) User {
	user := User{
		Id:        s.nextUserId(),
		FirstName: firstName,
		LastName:  lastName,
		Wishlists: []Wishlist{},
	}
	s.users = append(s.users, user)
	return user
}

func (s *ArrayStore) CreateWishlist(userId int, products []Product) (Wishlist, error) {
	user, err := s.GetUserById(userId)
	if err != nil {
		return Wishlist{}, err
	}

	wishlist := Wishlist{
		Id:       s.nextWishlistId(),
		Products: products,
	}
	user.Wishlists = append(user.Wishlists, wishlist)

	return wishlist, nil
}

func (s *ArrayStore) nextUserId() int {
	id := s.userIdCounter
	s.userIdCounter = s.userIdCounter + 1
	return id
}

func (s *ArrayStore) nextWishlistId() int {
	id := s.wishlistIdCounter
	s.userIdCounter = s.userIdCounter + 1
	return id
}
