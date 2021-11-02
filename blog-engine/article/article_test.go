package article

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validContent = `Sollicitudin ac orci phasellus egestas tellus. Ultricies mi quis hendrerit dolor
	magna eget est lorem ipsum. Et netus et malesuada fames ac. Euismod quis viverra nibh cras
	pulvinar mattis nunc sed blandit. Aliquam vestibulum morbi blandit cursus risus at. Amet risus
	nullam eget felis eget nunc lobortis. Amet volutpat consequat mauris nunc congue nisi vitae. Sem
	viverra aliquet eget sit amet tellus cras adipiscing enim. Aliquam ultrices sagittis orci a
	scelerisque purus semper eget duis. Interdum velit laoreet id donec ultrices tincidunt.
	Sollicitudin aliquam ultrices sagittis orci a. Aliquet eget sit amet tellus. Quis enim lobortis
	scelerisque fermentum dui faucibus. Dolor sit amet consectetur adipiscing elit ut. Vulputate enim
	nulla aliquet porttitor lacus.`

	validTitle = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
	incididunt ut labore et dolore magna aliqua.`

	longTitle = `et netus et malesuada fames ac turpis egestas integer eget aliquet nibh praesent
	tristique magna sit amet purus gravida quis blandit turpis cursus in hac habitasse platea dictumst
	quisque sagittis purus sit amet volutpat consequat mauris nunc congue nisi vitae suscipit tellus
	mauris a diam maecenas sed enim ut sem viverra aliquet eget sit amet tellus cras adipiscing enim
	eu turpis egestas pretium aenean pharetra magna ac placerat vestibulum lectus mauris ultrices eros
	in cursus turpis massa tincidunt dui ut ornare lectus sit amet est placerat in egestas erat
	imperdiet sed euismod nisi porta lorem mollis aliquam ut porttitor leo`
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		Name  string
		Title string
		Err   error
	}{
		{"Empty", "", ErrEmptyTitle},
		{"TooShort", "too-short", ErrTitleTooShort},
		{"TooLong", longTitle, ErrTitleTooLong},
		{"Valid", validTitle, nil},
	}

	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			assert.Equal(t, item.Err, validateTitle(item.Title))
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		Name string

		Content string
		Err     error
	}{
		{"Empty", "", ErrEmptyContent},
		{"TooShort", "short", ErrContentTooShort},
		{"Valid", validContent, nil},
	}

	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			assert.Equal(t, item.Err, validateContent(item.Content))
		})
	}
}

func TestCreateArticle(t *testing.T) {
	// article yang dianggap nil (karena bukan pointer)
	var nilArticle Article
	require.True(t, nilArticle.IsNil())

	// id yang pasti dianggap valid
	id := uuid.MustParse("836f6aa2-ed56-437a-aee1-eff92cf4ee4d")

	validArticle, err := createArticleWithID(id, validTitle, validContent)
	require.NoError(t, err)

	// tabel test
	tests := []struct {
		// nama test
		Name string
		// input
		Title   string
		Content string
		// output
		Result Article
		Err    error
	}{
		{"EmptyTitleContent", "", "", nilArticle, ErrEmptyTitle},
		{"EmptyTitleOnly", "", validContent, nilArticle, ErrEmptyTitle},
		{"EmptyContent", validTitle, "", nilArticle, ErrEmptyContent},
		{"ShortTitle", "short", validContent, nilArticle, ErrTitleTooShort},
		{"ShortContent", validTitle, "short", nilArticle, ErrContentTooShort},
		{"TooLongTitle", longTitle, validContent, nilArticle, ErrTitleTooLong},
		{"ValidArticle", validTitle, validContent, validArticle, nil},
	}

	// Test dijalankan satu-satu dari tabel
	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			article, err := createArticleWithID(id, item.Title, item.Content)
			assert.Equal(t, item.Err, err)

			if item.Err == nil {
				assert.Equal(t, item.Result.ID, article.ID)
				assert.Equal(t, item.Result.Title, article.Title)
				assert.Equal(t, item.Result.Content, article.Content)
				assert.LessOrEqual(t, time.Since(article.CreatedAt), 5*time.Second)
				assert.False(t, item.Result.IsNil())
			}
		})
	}
}

func TestChangeTitle(t *testing.T) {
	const newValidTitle = "ini title baru yg valid"
	const newShortTitle = "tidak"

	var nilArticle Article

	validArticle, _ := CreateArticle(validTitle, validContent)
	newValidArticle := validArticle
	newValidArticle.Title = newValidTitle

	tests := []struct {
		Name string
		//input
		Article  Article
		NewTitle string
		//ouput
		Err        error
		NewArticle Article
	}{
		{"NilArticle", nilArticle, newValidTitle, ErrNilArticle, nilArticle},
		{"EmptyTitle", validArticle, "", ErrEmptyTitle, validArticle},
		{"TitleTooShort", validArticle, newShortTitle, ErrTitleTooShort, validArticle},
		{"TitleTooLong", validArticle, longTitle, ErrTitleTooLong, validArticle},
		{"validTitle", validArticle, newValidTitle, nil, newValidArticle},
	}
	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			oldArticle := item.Article
			err := oldArticle.ChangeTitle(item.NewTitle)

			assert.Equal(t, item.Err, err)

			if item.Err == nil {
				assert.Equal(t, item.NewArticle, oldArticle)
			}
		})
	}
}

func TestChangeContent(t *testing.T) {
	const newValidContent = "ini content baru yg valid dibuat oleh calvin"
	const newShortContent = "tidak"

	var nilArticle Article

	validArticle, _ := CreateArticle(validTitle, validContent)
	newValidArticle := validArticle
	newValidArticle.Content = newValidContent

	tests := []struct {
		Name string
		//input
		Article    Article
		NewContent string
		//ouput
		Err        error
		NewArticle Article
	}{
		{"NilArticle", nilArticle, newValidContent, ErrNilArticle, nilArticle},
		{"EmptyContent", validArticle, "", ErrEmptyContent, validArticle},
		{"ContentTooShort", validArticle, newShortContent, ErrContentTooShort, validArticle},
		{"validContent", validArticle, newValidContent, nil, newValidArticle},
	}
	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			oldArticle := item.Article
			err := oldArticle.ChangeContent(item.NewContent)

			assert.Equal(t, item.Err, err)

			if item.Err == nil {
				assert.Equal(t, item.NewArticle, oldArticle)
			}
		})
	}
}

func TestEditArticle(t *testing.T) {
	const newValidTitle = "ini title baru yg valid"
	const newShortTitle = "tidak"

	const newValidContent = "ini content baru yg valid dibuat oleh calvin"
	const newShortContent = "tidak"

	var newUpdateTime = time.Now()

	var nilArticle Article
	require.True(t, nilArticle.IsNil())

	validArticle, _ := CreateArticle(validTitle, validContent)
	newValidArticle := validArticle
	newValidArticle.Title = newValidTitle
	newValidArticle.Content = newValidContent
	newValidArticle.UpdatedAt = newUpdateTime
	newValidArticle.Deleted = false
	require.False(t, newValidArticle.Deleted)

	tests := []struct {
		Name string
		//intput
		OldArticle Article
		NewTitle   string
		NewContent string
		//output
		NewArticle Article
		Err        error
	}{
		{"NilArticle", nilArticle, newValidTitle, newValidContent, nilArticle, ErrNilArticle},
		{"EmptyTitleContent", validArticle, "", "", nilArticle, ErrEmptyTitle},
		{"EmptyTitleOnly", validArticle, "", newValidContent, validArticle, ErrEmptyTitle},
		{"EmptyContentOnly", validArticle, newValidTitle, "", validArticle, ErrEmptyContent},
		{"ShortTitle", validArticle, newShortTitle, newValidContent, validArticle, ErrTitleTooShort},
		{"ShortContent", validArticle, newValidTitle, newShortContent, validArticle, ErrContentTooShort},
		{"TooLongTitle", validArticle, longTitle, newValidContent, validArticle, ErrTitleTooLong},
		{"ValidArticle", validArticle, newValidTitle, newValidContent, newValidArticle, nil},
	}
	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			oldArticle := item.OldArticle
			err := oldArticle.EditArticle(item.NewTitle, item.NewContent)
			assert.Equal(t, item.Err, err)

			if item.Err == nil {
				item.NewArticle.UpdatedAt = time.Now()
				oldArticle.UpdatedAt = time.Now()
				assert.True(t, item.NewArticle.IsEqual(oldArticle))

			}
		})
	}
}

func TestSetDeleted(t *testing.T) {

	var nilArticle Article
	require.True(t, nilArticle.IsNil())

	validArticle, _ := CreateArticle(validTitle, validContent)
	deleteArticle := validArticle
	deleteArticle.Deleted = true
	deleteArticle.DeletedAt.Time = time.Now()
	deleteArticle.DeletedAt.Valid = true

	tests := []struct {
		Name string
		//expected
		Article Article
		//output
		Err error
	}{
		{"nilArticle", nilArticle, ErrNilArticle},
		{"validDeleted", deleteArticle, nil},
	}
	for _, item := range tests {
		t.Run(item.Name, func(t *testing.T) {
			article := item.Article
			result := article.SetDeleted()
			require.Equal(t, item.Err, result)

			if item.Err == nil {
				item.Article.UpdatedAt = time.Now()
				assert.True(t, item.Article.IsEqual(deleteArticle))
			}
		})
	}

}
