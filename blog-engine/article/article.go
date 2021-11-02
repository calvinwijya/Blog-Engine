package article

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Article struct {
	ID      uuid.UUID
	Title   string
	Content string

	CreatedAt time.Time
	UpdatedAt time.Time

	Deleted   bool
	DeletedAt null.Time

	//Author //type??
}

var (
	ErrNilArticle      = errors.New("article cannot be nil")
	ErrEmptyTitle      = errors.New("title is empty")
	ErrEmptyContent    = errors.New("content is empty")
	ErrTitleTooShort   = errors.New("title too short")
	ErrTitleTooLong    = errors.New("title too long")
	ErrContentTooShort = errors.New("content too short")
	errArticleDeleted  = errors.New("article deleted")
)

func (a Article) IsNil() bool {
	return a.ID == uuid.Nil && len(a.Title) == 0 && len(a.Content) == 0
}

func validateTitle(title string) error {
	const minTitleLength = 10
	const maxTitleLength = 500

	runeCount := utf8.RuneCountInString(title)

	if runeCount == 0 {
		return ErrEmptyTitle
	}

	if runeCount < minTitleLength {
		return ErrTitleTooShort
	}

	if runeCount > maxTitleLength {
		return ErrTitleTooLong
	}

	return nil
}

func validateContent(content string) error {
	const minContentLength = 10

	runeCount := utf8.RuneCountInString(content)

	if runeCount == 0 {
		return ErrEmptyContent
	}

	if runeCount < minContentLength {
		return ErrContentTooShort
	}

	return nil
}

func createArticleWithID(id uuid.UUID, title, content string) (Article, error) {
	var newArticle Article

	if err := validateTitle(title); err != nil {
		return Article{}, err
	}

	if err := validateContent(content); err != nil {
		return Article{}, err
	}
	currentTime := time.Now()
	newArticle = Article{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	return newArticle, nil
}

func CreateArticle(title, content string) (Article, error) {
	newId, err := uuid.NewRandom()

	if err != nil {
		return Article{}, err
	}
	return createArticleWithID(newId, title, content)
}

func (a *Article) ChangeTitle(newTitle string) error {
	if a.IsNil() {
		return ErrNilArticle
	}

	if err := validateTitle(newTitle); err != nil {
		return err
	}
	a.Title = newTitle
	return nil
}

func (a *Article) ChangeContent(newContent string) error {
	if a.IsNil() {
		return ErrNilArticle
	}
	if err := validateContent(newContent); err != nil {
		return err
	}
	a.Content = newContent
	return nil
}

func AbsDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func (a Article) IsEqual(other Article) bool {
	const e = 10 * time.Millisecond

	return a.ID == other.ID && a.Title == other.Title && a.Content == other.Content &&
		AbsDuration(a.CreatedAt.Sub(other.CreatedAt)) <= e &&
		AbsDuration(a.UpdatedAt.Sub(other.UpdatedAt)) <= e

}

func (a *Article) EditArticle(newTitle, newContent string) error {

	if a.IsNil() {
		return ErrNilArticle
	}

	if a.Deleted {
		return errArticleDeleted
	}

	if err := validateTitle(newTitle); err != nil {
		return err
	}

	if err := validateContent(newContent); err != nil {
		return err
	}

	a.Title = newTitle
	a.Content = newContent
	a.UpdatedAt = time.Now()

	return nil
}

func (a *Article) SetDeleted() error {

	if a.IsNil() {
		return ErrNilArticle
	}

	a.Deleted = true
	a.DeletedAt.Time = time.Now()
	a.DeletedAt.Valid = true

	return nil
}
