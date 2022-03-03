package main

import (
	"fmt"
	"os"

	"package.local/hktn/helper"
)

var (
	bookList   helper.BookList
	reset      string
	searchWord string
	counter    int
)

const (
	searchChoice   int    = 1
	listAllChoice  int    = 2
	buyChoice      int    = 3
	deleteChoice   int    = 4
	exitChoice     int    = 5
	resetChoice    string = "r"
	wrongChoice    string = "Yanlış seçim. Tekrar deneyiniz"
	notFoundChoice string = "Kitap bulunamadı. Tekrar deneyiniz"
)

func init() {

	bookList = helper.InitializeBookList()
}

func main() {

	printUsage()
}

// uygulama seçimlerinin ekranada gösterildiği ve yapıldığı yer
func printUsage() {

	fmt.Println("Kitaplık uygulamasında kullanabileceğiniz komutlar :")
	fmt.Printf(" search => arama işlemi için %d\n", searchChoice)
	fmt.Printf(" list => listeleme işlemi için %d\n", listAllChoice)
	fmt.Printf(" buy => kitap satın almak için %d\n", buyChoice)
	fmt.Printf(" delete => kitap silmek için %d\n", deleteChoice)
	fmt.Printf(" exit => uygulamadan çıkmak için %d\n", exitChoice)
	fmt.Println("Tuşlarına basınız")

	var choice int

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
		printUsage()
	}
}

// kitap silme işlemleri başladı gerekli kontroller ve silme işlemine yönlendirme
func deleteBook() {

	fmt.Println("Silmek İstediğiniz Kitap Id Giriniz : ")

	var selectedBookId int

	if _, err := fmt.Scan(&selectedBookId); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		deleteBook()
	}

	if !checkIdExist(selectedBookId) {
		fmt.Println(notFoundChoice)
		deleteBook()
	}

	var book, _, err = bookList.GetBookById(selectedBookId)

	errorCheck(err)

	err = bookList.Delete(book.Id)

	errorCheck(err)

	fmt.Println("Kitap silinmiştir")
	helper.AddSeparator()
	conclude()

}

// hata kontrol
func errorCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
		helper.AddSeparator()
		conclude()
	}
}

// kitap satın alma işlemleri başladı gerekli kontroller
func buyBook() {

	fmt.Println("Satın Alınacak Kitap Id Giriniz : ")

	var selectedBookId int

	if _, err := fmt.Scan(&selectedBookId); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		buyBook()
	}

	if !checkIdExist(selectedBookId) {
		fmt.Println(notFoundChoice)
		buyBook()
	}

	purchase(selectedBookId)
}

// kitap satın alma işlemleri devam ediyor ve satın alma işlemine yönlendirme
func purchase(bookId int) {

	fmt.Println("Satın Alınacak Kitap Sayısını Giriniz : ")

	var numberOfPurchases int

	if _, err := fmt.Scan(&numberOfPurchases); err != nil {
		fmt.Println(wrongChoice)
		helper.AddSeparator()
		purchase(bookId)
	}

	var book, _, _ = bookList.GetBookById(bookId)

	err := book.Buy(numberOfPurchases)

	if err != nil {
		fmt.Println(err.Error())
		purchase(bookId)
	}

	fmt.Println("Satın alma işlemi sonrası kitap durumu")
	fmt.Printf("%+v\n", book)
	helper.AddSeparator()
	conclude()
}

// girilen id de kitap kontrolü ve varsa atanması
func checkIdExist(bookId int) bool {

	var selectedBook helper.Book
	for i := range bookList.BookList {
		if bookList.BookList[i].Id == bookId {
			selectedBook = bookList.BookList[i]
			fmt.Printf("Seçilen kitap : %+v\n", selectedBook)
			helper.AddSeparator()
			return true
		}
	}
	return true
}

// tüm kitapların listelenmesi
func listAll() {

	fmt.Println("Kitaplar Getiriliyor..")
	helper.AddSeparator()
	for _, line := range bookList.BookList {
		fmt.Printf("%+v\n", line)
	}
	helper.AddSeparator()
	fmt.Printf("%d Tane Kitaplar Getirildi..\n", len(bookList.BookList))
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

	switch helper.StringLower(reset) {
	case resetChoice:
		main()
	default:
		terminate()
	}
}

// kitap arama işlemi, ekrana basılması
func searchByWord(word string) {

	counter = 0
	for _, line := range bookList.BookList {
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
