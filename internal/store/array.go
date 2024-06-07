package store

import "elf/internal/core"

type ArrayWishlist struct {
	wishlistIdCounter int
	wishlists         []core.Wishlist
}

func NewArrayWishlist() (s *ArrayWishlist) {
	// TODO: Fine tune the value 5
	s = &ArrayWishlist{wishlistIdCounter: 0, wishlists: make([]core.Wishlist, 5)}
	s.Seed()
	return
}

func (s *ArrayWishlist) Seed() {
	products := []core.Product{
		{Name: "iPad", Url: "www.example.com", Price: 100, Currency: "eur"},
		{Name: "Macbook", Url: "www.example.com", Price: 200, Currency: "eur"},
	}
	wishlists := []core.Wishlist{
		{Id: 1, Products: products[:], OwnerId: 0},
		{Id: 2, Products: products[:1], OwnerId: 1},
		{Id: 3, Products: products[0:], OwnerId: 0},
	}
	s.wishlistIdCounter = 4

	s.wishlists = wishlists
}

func (s *ArrayWishlist) Create(name string, ownerId int, products ...core.Product) core.Wishlist {
	w := core.Wishlist{
		Id:       s.nextUserId(),
		Name:     name,
		OwnerId:  ownerId,
		Products: products,
	}

	s.wishlists = append(s.wishlists, w)
	return w
}

func (s *ArrayWishlist) Read(id int) (core.Wishlist, error) {
	for _, w := range s.wishlists {
		if w.Id == id {
			return w, nil
		}
	}
	return core.Wishlist{}, core.ErrWishlistDoesNotExist
}

func (s *ArrayWishlist) ReadAll() []core.Wishlist {
	return s.wishlists
}

func (s *ArrayWishlist) AddProductsTo(id int, products ...core.Product) error {
	wIdx := -1
	for i, w := range s.wishlists {
		if w.Id == id {
			wIdx = i
			break
		}
	}
	if wIdx == -1 {
		return core.ErrWishlistDoesNotExist
	}

	w := s.wishlists[wIdx]
	w = core.Wishlist{
		Id:       w.Id,
		OwnerId:  w.OwnerId,
		Products: append(w.Products, products...),
	}
	s.wishlists[wIdx] = w
	return nil
}

func (s *ArrayWishlist) nextUserId() int {
	id := s.wishlistIdCounter
	s.wishlistIdCounter = s.wishlistIdCounter + 1
	return id
}
