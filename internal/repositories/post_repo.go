package repositories

import (
	"blog-service/internal/models"

	"gorm.io/gorm"
)

type PostRepo struct {
	DB *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{DB: db}
}

func (r *PostRepo) Create(p *models.Post) error {
	return r.DB.Create(p).Error
}

func (r *PostRepo) Update(p *models.Post) error {
	return r.DB.Save(p).Error
}

func (r *PostRepo) DeleteByID(id uint) error {
	return r.DB.Delete(&models.Post{}, id).Error
}

func (r *PostRepo) FindByID(id uint) (*models.Post, error) {
	var p models.Post
	if err := r.DB.Preload("Tags").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepo) FindBySlugPublished(slug string) (*models.Post, error) {
	var p models.Post
	if err := r.DB.Preload("Tags").
		Preload("Author").
		Where("slug = ? AND status = ?", slug, models.PostPublished).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepo) FindBySlugAny(slug string) (*models.Post, error) {
	var p models.Post
	if err := r.DB.Preload("Tags").
		Preload("Author").
		Where("slug = ?", slug).
		First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostRepo) SlugExists(slug string) (bool, error) {
	var cnt int64
	if err := r.DB.Model(&models.Post{}).Where("slug = ?", slug).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *PostRepo) ListPublished(page, size int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size

	var total int64
	if err := r.DB.Model(&models.Post{}).
		Where("status = ?", models.PostPublished).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Post
	err := r.DB.
		Select("id,title,slug,status,published_at,author_id,view_count,like_count,comment_count,created_at,updated_at").
		Preload("Tags").
		Preload("Author").
		Where("status = ?", models.PostPublished).
		Order("published_at DESC").
		Offset(offset).Limit(size).
		Find(&items).Error
	return items, total, err
}

func (r *PostRepo) ListAny(status *models.PostStatus, page, size int) ([]models.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size

	q := r.DB.Model(&models.Post{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Post
	err := q.
		Select("id,title,slug,status,published_at,author_id,view_count,like_count,comment_count,created_at,updated_at").
		Preload("Tags").
		Preload("Author").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&items).Error
	return items, total, err
}

func (r *PostRepo) IncViewCount(id uint) error {
	return r.DB.Model(&models.Post{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
