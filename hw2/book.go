package main

import (
	"errors"

	"package.local/hktn/helper"
)

var (
	ErrNotEnoughStock     = errors.New("Yeterli stok bulunmamaktadır")
	ErrZeroValue          = errors.New("0 dan büyük bir değer giriniz")
	ErrBookAlreadyDeleted = errors.New("Kitap silinemez")
	ErrBookNotFound       = errors.New("Kitap bulunamadı")
)

var book_id int = 0

const (
	max_page     int     = 1000
	min_page     int     = 25
	max_year     int     = 2022
	min_year     int     = 1970
	max_price    float64 = 10.05
	min_price    float64 = 155.83
	max_quantity int     = 20
	min_quantity int     = 0
)

type BookList struct {
	BookList []Book
}

type Book struct {
	Id, NumberOfPages, Year, Quantity int
	Price                             float64
	Name                              string
	Author                            string
	Sku, Isbn                         string
	IsDeleted                         bool
}

type Deletable interface {
	Delete()
}

// satın alma işlemi
func (b *Book) Buy(numberOfPurchases int) error {

	if b.Quantity < numberOfPurchases {
		return ErrNotEnoughStock
	}
	if 0 == numberOfPurchases {
		return ErrZeroValue
	}
	b.Quantity -= numberOfPurchases
	return nil
}

// silme işlemi
func (bl *BookList) Delete(bookId int) error {

	book, index, err := bl.GetBookById(bookId)

	if err != nil {
		return ErrBookNotFound
	}

	if book.IsDeleted {
		return ErrBookAlreadyDeleted
	}

	RemoveBookSliceForIndex(bl, index)
	return nil
}

// slice üzerinden silme
func RemoveBookSliceForIndex(s *BookList, index int) {

	s.BookList = append(s.BookList[:index], s.BookList[index+1:]...)
}

// id göre kitap getir
func (bl *BookList) GetBookById(bookId int) (*Book, int, error) {

	var book *Book
	var index int

	for i := range bl.BookList {
		if bl.BookList[i].Id == bookId {
			book = &bl.BookList[i]
			index = i
			return book, index, nil
		}
	}

	return book, index, ErrBookNotFound
}

// kitaplar
func GetBookSlice() [][]string {

	return [][]string{
		{"Hamlet", "William Shakespeare"},
		{"The Picture of Dorian Gray", "Oscar Wilde"},
		{"Crime and Punishment", "Fyodor Dostoevsky"},
		{"Pride and Prejudice", "Jane Austen"},
		{"Adventures of Huckleberry Finn", "Mark Twain"},
		{"The Count Of Monte Cristo", "Alexandre Dumas"},
		{"Oliver Twist", "Charles Dickens"},
	}
}

// kitaplar BookList atılıyor
func InitializeBookList() BookList {

	var bookList BookList
	var books []Book

	var bookSlice = GetBookSlice()

	for i := range bookSlice {
		var b Book = fillBookValues(bookSlice[i][0], bookSlice[i][1])
		books = append(books, b)
	}
	bookList.BookList = books

	return bookList
}

// kitap alanları türetiliyor
func fillBookValues(name, author string) Book {

	book_id++
	var bf Book
	bf.Id = book_id
	bf.Name = name
	bf.Author = author
	bf.NumberOfPages = helper.RandomIntegerCreator(min_page, max_page)
	bf.Year = helper.RandomIntegerCreator(min_year, max_year)
	bf.Price = helper.RoundFloat(helper.RandomFloatCreator(min_price, max_price))
	bf.Quantity = helper.RandomIntegerCreator(min_quantity, max_quantity)
	bf.Sku = helper.CreateSku(name)
	bf.Isbn = helper.CreateIsbn()
	bf.IsDeleted = helper.RandomBoolCreator()

	return bf
}
