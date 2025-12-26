package services

import (
	"errors"
	"strings"
	"time"

	"blog-service/internal/models"
	"blog-service/internal/repositories"
	"blog-service/internal/utils/markdown"
	"blog-service/internal/utils/slug"
)

var (
	ErrPostNotFound    = errors.New("post_not_found")
	ErrInvalidStatus   = errors.New("invalid_status")
	ErrTitleRequired   = errors.New("title_required")
	ErrContentRequired = errors.New("content_required")
)

type PostService struct {
	Posts *repositories.PostRepo
	Tags  *repositories.TagRepo
}

type CreatePostInput struct {
	Title     string
	ContentMD string
	Status    models.PostStatus // draft/published
	Tags      []string
	AuthorID  uint
}

func (s *PostService) Create(in CreatePostInput) (*models.Post, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	if strings.TrimSpace(in.ContentMD) == "" {
		return nil, ErrContentRequired
	}
	if in.Status != models.PostDraft && in.Status != models.PostPublished {
		return nil, ErrInvalidStatus
	}

	html, err := markdown.RenderToSafeHTML(in.ContentMD)
	if err != nil {
		return nil, err
	}

	baseSlug := slug.FromTitle(title)
	finalSlug := baseSlug
	exists, err := s.Posts.SlugExists(finalSlug)
	if err != nil {
		return nil, err
	}
	if exists {
		finalSlug = baseSlug + "-" + slug.RandSuffix(3) // 6 hex chars
	}

	var publishedAt *time.Time
	if in.Status == models.PostPublished {
		now := time.Now()
		publishedAt = &now
	}

	p := &models.Post{
		Title:       title,
		Slug:        finalSlug,
		ContentMD:   in.ContentMD,
		ContentHTML: html,
		Status:      in.Status,
		PublishedAt: publishedAt,
		AuthorID:    in.AuthorID,
	}

	if len(in.Tags) > 0 {
		tags, err := s.Tags.GetOrCreateByNames(in.Tags)
		if err != nil {
			return nil, err
		}
		p.Tags = tags
	}

	if err := s.Posts.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

type UpdatePostInput struct {
	Title     *string
	ContentMD *string
	Status    *models.PostStatus
	Tags      *[]string
}

func (s *PostService) Update(postID uint, in UpdatePostInput) (*models.Post, error) {
	p, err := s.Posts.FindByID(postID)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return nil, ErrTitleRequired
		}
		// 如标题变化，slug 可选是否更新；这里默认不改 slug（更稳定）
		p.Title = t
	}

	if in.ContentMD != nil {
		md := *in.ContentMD
		if strings.TrimSpace(md) == "" {
			return nil, ErrContentRequired
		}
		html, err := markdown.RenderToSafeHTML(md)
		if err != nil {
			return nil, err
		}
		p.ContentMD = md
		p.ContentHTML = html
	}

	if in.Status != nil {
		if *in.Status != models.PostDraft && *in.Status != models.PostPublished {
			return nil, ErrInvalidStatus
		}
		// draft -> published 时补发布时间
		if p.Status != models.PostPublished && *in.Status == models.PostPublished {
			now := time.Now()
			p.PublishedAt = &now
		}
		p.Status = *in.Status
	}

	if in.Tags != nil {
		tags, err := s.Tags.GetOrCreateByNames(*in.Tags)
		if err != nil {
			return nil, err
		}
		p.Tags = tags
	}

	if err := s.Posts.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostService) Preview(md string) (string, error) {
	return markdown.RenderToSafeHTML(md)
}
