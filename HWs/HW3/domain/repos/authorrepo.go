package repos

import (
	"errors"
	"fmt"

	"github.com/Metehan1994/HWs/HW3/domain/entities"
	"github.com/Metehan1994/HWs/HW3/models"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

//NewAuthorRepository create a database for author
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

//List gives output for authors
func (a *AuthorRepository) List() {
	var authors []entities.Author
	a.db.Find(&authors)

	for _, author := range authors {
		fmt.Println(author.ToString())
	}
}

//GetByID provides the author info for a given ID
func (a *AuthorRepository) GetByID(ID int) (*entities.Author, error) {
	var author entities.Author
	result := a.db.First(&author, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &author, nil
}

//FindByWord lists the authors with the given word case-insensitively
func (a *AuthorRepository) FindByWord(name string) {
	var authors []entities.Author
	a.db.Where("name ILIKE ? ", "%"+name+"%").Find(&authors)

	for _, author := range authors {
		fmt.Println(author.ToString())
	}
}

//FindByName provides the author with the input of full name
func (a *AuthorRepository) FindByName(name string) {
	var author entities.Author
	a.db.Where("name = ? ", name).Find(&author)

	fmt.Println("found:", author.Name)
}

//Create creates a new author
func (a *AuthorRepository) Create(author entities.Author) error {
	result := a.db.Where("name = ?", author.Name).FirstOrCreate(&author)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//DeleteByName deletes the author with the given full name
func (a *AuthorRepository) DeleteByName(name string) error {
	var author entities.Author
	result := a.db.Unscoped().Where("name = ?", name).Find(&author)
	if result.Error != nil {
		return result.Error
	} else if author.Name != "" && !author.DeletedAt.Valid {
		fmt.Println("Valid author name, deleted:", name)
	} else if author.Name != "" && author.DeletedAt.Valid {
		fmt.Println("It has been already deleted.")
	} else {
		fmt.Println("Invalid author name, no deletion.")
	}
	result = a.db.Where("name = ?", name).Delete(&entities.Author{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//DeleteByID applies a soft delete to an author with given ID
func (a *AuthorRepository) DeleteById(id int) error {
	var author entities.Author
	result := a.db.First(&author, id)
	if result.Error != nil {
		return result.Error
	} else {
		fmt.Println("Valid ID, deleted:", id)
	}
	result = a.db.Delete(&entities.Author{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//GetAuthorsWithBookInformation gives output of authors with their books info
func (a *AuthorRepository) GetAuthorsWithBookInformation() ([]entities.Author, error) {
	var authors []entities.Author
	result := a.db.Preload("Book").Find(&authors)
	if result.Error != nil {
		return nil, result.Error
	}
	return authors, nil
}

//BooksOfAuthors finds the books of author searched
func (a *AuthorRepository) BooksOfAuthors(name string) error {
	var author entities.Author
	result := a.db.Where("name = ?", name).Preload("Book").Find(&author)
	if result.Error != nil {
		return result.Error
	}
	Books := author.Book
	if len(Books) == 0 {
		fmt.Println("No book info for given name.")
	} else {
		fmt.Println("Writer: ", name)
		for i, book := range Books {
			fmt.Printf("Book %d: %s\n", i+1, book.Name)
		}
	}
	return nil
}

//Migrations form an author table in db
func (a *AuthorRepository) Migrations() {
	a.db.AutoMigrate(&entities.Author{})
}

//InsertSampleData creates a list of authors
func (a *AuthorRepository) InsertSampleData(bookList models.BookList) {

	authors := []entities.Author{}
	for _, book := range bookList {
		newAuthor := entities.Author{
			Name: book.Author.AuthorName,
			ID:   uint(book.Author.AuthorID),
		}
		authors = append(authors, newAuthor)
	}

	for _, author := range authors {
		a.db.Unscoped().Where(entities.Author{ID: author.ID}).Attrs(entities.Author{Name: author.Name}).FirstOrCreate(&author)
	}

}
