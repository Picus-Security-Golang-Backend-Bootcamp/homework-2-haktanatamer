package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"package.local/hktn/helper"
)

var choice int
var reset string
var searchWord string
var counter int
var book_list []Book
var book_id int = 0
var selected_book_id int
var number_of_purchases int
var selected_book *Book
var selected_slice_index int

var ErrNotEnoughStock = errors.New("Yeterli stok bulunmamaktadır")
var ErrZeroValue = errors.New("0 dan büyük bir değer giriniz")
var ErrNotAuthorized = errors.New("Kitap zaten silinmiş")

const max_page int = 1000
const min_page int = 25
const max_year int = 2022
const min_year int = 1970
const max_price float64 = 10.05
const min_price float64 = 155.83
const max_quantity int = 20
const min_quantity int = 0

const searchChoice int = 1
const listAllChoice int = 2
const buyChoice int = 3
const deleteChoice int = 4
const exitChoice int = 5
const upperResetChoice string = "R"
const lowerResetChoice string = "r"
const wrongChoice string = "Yanlış seçim. Tekrar deneyiniz"
const notFoundChoice string = "Kitap bulunamadı. Tekrar deneyiniz"
const booksJsonPath string = "books.json"

func init() {

	initializeJsonBook()
}

func main() {

	fmt.Println("Kitaplık uygulamasında kullanabileceğiniz komutlar :")
	fmt.Printf(" search => arama işlemi için %d\n", searchChoice)
	fmt.Printf(" list => listeleme işlemi için %d\n", listAllChoice)
	fmt.Printf(" buy => kitap satın almak için %d\n", buyChoice)
	fmt.Printf(" delete => kitap silmek için %d\n", deleteChoice)
	fmt.Printf(" exit => uygulamadan çıkmak için %d\n", exitChoice)
	fmt.Println("Tuşlarına basınız")

	fmt.Scan(&choice)

	switch choice {
	case searchChoice:
		search()
	case listAllChoice:
		listAll()
	case buyChoice:
		buyBook()
	case deleteChoice:
		deleteBook()
	case exitChoice:
		terminate()
	default:
		helper.AddSeparator()
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		main()
	}
}

// kitap silme işlemleri başladı gerekli kontroller ve silme işlemine yönlendirme
func deleteBook() {

	fmt.Println("Silmek İstediğiniz Kitap Id Giriniz : ")

	if _, err := fmt.Scan(&selected_book_id); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		deleteBook()
	}

	if !checkIdExist() {
		fmt.Println(notFoundChoice)
		deleteBook()
	}

	err := selected_book.Delete()

	if err != nil {
		fmt.Println(err.Error())
		helper.AddSeparator()
		conclude()
	}

	fmt.Println("Kitap silinmiştir")
	helper.AddSeparator()
	conclude()
}

// kitap satın alma işlemleri başladı gerekli kontroller
func buyBook() {

	fmt.Println("Satın Alınacak Kitap Id Giriniz : ")

	if _, err := fmt.Scan(&selected_book_id); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		buyBook()
	}

	if !checkIdExist() {
		fmt.Println(notFoundChoice)
		buyBook()
	}

	purchase()
}

// kitap satın alma işlemleri devam ediyor ve satın alma işlemine yönlendirme
func purchase() {

	fmt.Println("Satın Alınacak Kitap Sayısını Giriniz : ")

	if _, err := fmt.Scan(&number_of_purchases); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		purchase()
	}

	err := selected_book.Buy()

	if err != nil {
		fmt.Println(err.Error())
		helper.AddSeparator()
		purchase()
	}

	fmt.Println("Satın alma işlemi sonrası kitap durumu")
	fmt.Printf("%+v\n", selected_book)
	helper.AddSeparator()
	conclude()
}

// girilen id de kitap kontrolü ve varsa atanması
func checkIdExist() bool {

	for i := range book_list {
		if book_list[i].Id == selected_book_id {
			selected_book = &book_list[i]
			selected_slice_index = i
			fmt.Printf("Seçilen kitap : %+v\n", selected_book)
			helper.AddSeparator()
			return true
		}
	}
	return false
}

// tüm kitapların listelenmesi
func listAll() {

	fmt.Println("Kitaplar Getiriliyor..")
	helper.AddSeparator()
	for _, line := range book_list {
		fmt.Printf("%+v\n", line)
	}
	helper.AddSeparator()
	fmt.Printf("%d Tane Kitaplar Getirildi..\n", len(book_list))
	helper.AddSeparator()
	conclude()
}

// kitap arama işlemleri başladı gerekli kontroller ve aramaya yönlendirme
func search() {

	fmt.Println("Aranacak Kelimeyi Giriniz : ")
	fmt.Scan(&searchWord)
	helper.AddSeparator()
	fmt.Printf("%s Kelimesi Aranıyor..\n", searchWord)
	helper.AddSeparator()
	searchByWord(searchWord)
	helper.AddSeparator()
	fmt.Printf("%d Tane Kitap Bulundu..\n", counter)
	fmt.Println("Arama Tamamlandı..")
	helper.AddSeparator()
	conclude()
}

// uygulamayı sonlandırma
func terminate() {

	fmt.Println("Uygulama Sonlandırılıyor..")
	os.Exit(3)
}

// uygulamaya devam edilecek mi kontrolleri
func conclude() {

	fmt.Println("Yeni Bir Arama Yapmak İçin R \nUygulamanyı Sonlandırmak İçin Herhangi Bir Tuşa Basınız")

	fmt.Scan(&reset)

	switch reset {
	case lowerResetChoice:
		main()
	case upperResetChoice:
		main()
	default:
		terminate()
	}
}

// kitap arama işlemi, ekrana basılması
func searchByWord(word string) {

	counter = 0
	for _, line := range book_list {
		if helper.StringContains(line.Name, word) {
			counter++
			fmt.Printf("%+v\n", line)
			continue
		}
		if helper.StringContains(line.Author, word) {
			counter++
			fmt.Printf("%+v\n", line)
			continue
		}
		if helper.StringContains(line.Sku, word) {
			counter++
			fmt.Printf("%+v\n", line)
			continue
		}
	}
}

// kitapların yüklenmesi
func initializeJsonBook() {

	file, _ := ioutil.ReadFile(booksJsonPath)

	data := BookList{}

	_ = json.Unmarshal([]byte(file), &data)

	for i := 0; i < len(data.BookList); i++ {
		var b Book = fillBookValues(data.BookList[i])
		book_list = append(book_list, b)
	}
}

// book slice oluşturuluyor
func fillBookValues(db Book) Book {

	book_id++
	var bf Book
	bf.Id = book_id
	bf.Name = db.Name
	bf.Author = db.Author
	bf.NumberOfPages = helper.RandomIntegerCreator(min_page, max_page)
	bf.Year = helper.RandomIntegerCreator(min_year, max_year)
	bf.Price = helper.RoundFloat(helper.RandomFloatCreator(min_price, max_price))
	bf.Quantity = helper.RandomIntegerCreator(min_quantity, max_quantity)
	bf.Sku = helper.CreateSku(db.Name)
	bf.Isbn = helper.CreateIsbn()
	bf.IsDeleted = helper.RandomBoolCreator()

	return bf
}

type BookList struct {
	BookList []Book `json:"books"`
}

type Book struct {
	Id, NumberOfPages, Year, Quantity int
	Price                             float64
	Name                              string `json:"name"`
	Author                            string `json:"author"`
	Sku, Isbn                         string
	IsDeleted                         bool
}

type Deletable interface {
	Delete()
}

// satın alma işlemi
func (b Book) Buy() error {

	if selected_book.Quantity < number_of_purchases {
		return ErrNotEnoughStock
	}
	if 0 == number_of_purchases {
		return ErrZeroValue
	}
	selected_book.Quantity -= number_of_purchases
	return nil
}

// silme işlemi
func (b Book) Delete() error {

	if selected_book.IsDeleted {
		return ErrNotAuthorized
	}
	RemoveBookSliceForIndex(book_list, selected_slice_index)
	return nil
}

// slice üzerinden silme
func RemoveBookSliceForIndex(s []Book, index int) {

	book_list = append(s[:index], s[index+1:]...)
}
