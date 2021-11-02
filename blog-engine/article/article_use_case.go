package article

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ArticleFinder interface {
	FindArticleByID(ctx context.Context, id uuid.UUID) (Article, error)
}
type ArticleSaver interface {
	SaveArticle(ctx context.Context, article Article) error
}

type ArticleFinderSaver interface {
	ArticleFinder
	ArticleSaver
}

type ArticleUseCase struct {
	store ArticleFinderSaver
}

// Read Models

type ArticleLister interface {
	ListArticles(ctx context.Context) ([]ArticleBrief, error)
}

type ArticleReader interface {
	ArticleFinder
	ArticleLister
}

var ErrNilStore = errors.New("store cannot be nil")

func NewArticleUseCase(store ArticleFinderSaver) (*ArticleUseCase, error) {
	if store == nil {
		return nil, ErrNilStore
	}

	return &ArticleUseCase{store: store}, nil
}

func (uc *ArticleUseCase) CreateArticle(ctx context.Context, title, content string) (Article, error) {
	newArticle, err := CreateArticle(title, content)

	if err != nil {
		return Article{}, err
	}

	if err := uc.store.SaveArticle(ctx, newArticle); err != nil {
		return Article{}, err
	}

	return newArticle, nil
}

func (uc *ArticleUseCase) EditArticle(ctx context.Context, id uuid.UUID, newTitle, newContent string) error {
	article, err := uc.store.FindArticleByID(ctx, id)
	if err != nil {
		return err
	}

	if err = article.EditArticle(newTitle, newContent); err != nil {
		return err
	}

	return uc.store.SaveArticle(ctx, article)
}

func (uc *ArticleUseCase) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	article, err := uc.store.FindArticleByID(ctx, id)

	if err != nil {
		return err
	}

	if article.Deleted {
		return ErrArticleNotFound
	}
	if err = article.SetDeleted(); err != nil {
		return err
	}
	return uc.store.SaveArticle(ctx, article)
}

func (uc *ArticleUseCase) FindArticleByID(ctx context.Context, id uuid.UUID) (Article, error) {
	article, err := uc.store.FindArticleByID(ctx, id)

	if err != nil {
		return Article{}, err
	}

	if article.Deleted {
		return Article{}, ErrArticleNotFound
	}

	return article, err
}
