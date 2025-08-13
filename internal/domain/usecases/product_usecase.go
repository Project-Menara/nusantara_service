package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/domain/repositories"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/workers/payload"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductService struct {
	repo            repositories.ProductRepository
	typeProductRepo repositories.TypeProductRepository
	rdb             *redis.Client
	cloudinary      cloudinary.CloudinaryService
}

func NewProductUsecase(repo repositories.ProductRepository, typeProductRepo repositories.TypeProductRepository, rdb *redis.Client, cloudinary *cloudinary.CloudinaryService) services.ProductService {
	return &ProductService{
		repo:            repo,
		typeProductRepo: typeProductRepo,
		rdb:             rdb,
		cloudinary:      *cloudinary,
	}
}

func Sanitize(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func BuildPublicID(folder, filename string) string {
	return fmt.Sprintf("%s/%s", folder, filename)
}

func PublicIDFromURL(u string) string {
	if u == "" {
		return ""
	}
	const pivot = "/upload/"
	idx := strings.Index(u, pivot)
	if idx == -1 {
		return u // fallback
	}
	sub := u[idx+len(pivot):]

	// buang versi v123456/ jika ada
	if strings.HasPrefix(sub, "v") {
		if s := strings.Index(sub, "/"); s != -1 {
			sub = sub[s+1:]
		}
	}

	// buang ekstensi
	if dot := strings.LastIndex(sub, "."); dot != -1 {
		sub = sub[:dot]
	}
	return sub
}

func (p *ProductService) UploadMany(ctx context.Context, folder, prefix string, files []*multipart.FileHeader, workers int) ([]string, []string, error) {
	if len(files) == 0 {
		return nil, nil, nil
	}
	if workers < 1 {
		workers = 3
	}

	type job struct {
		idx int
		f   *multipart.FileHeader
	}
	type result struct {
		idx      int
		url      string
		publicID string
		err      error
	}

	jobs := make(chan job)
	results := make(chan result)

	worker := func() {
		for j := range jobs {
			ext := strings.ToLower(filepath.Ext(j.f.Filename))
			fileName := fmt.Sprintf("%s_%d%s", prefix, j.idx, ext)

			data, err := p.cloudinary.UploadImage(ctx, j.f, folder, fileName)
			if err != nil {
				results <- result{idx: j.idx, err: err}
				continue
			}
			results <- result{idx: j.idx, url: data.URL, publicID: BuildPublicID(folder, fileName)}
		}
	}

	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		for i, f := range files {
			jobs <- job{idx: i, f: f}
		}
		close(jobs)
	}()

	urls := make([]string, len(files))
	pids := make([]string, len(files))
	var firstErr error
	done := 0
	timeout := time.After(120 * time.Second)

	for done < len(files) {
		select {
		case r := <-results:
			if r.err != nil && firstErr == nil {
				firstErr = r.err
			}
			urls[r.idx] = r.url
			pids[r.idx] = r.publicID
			done++
		case <-timeout:
			return nil, nil, response.NewCustomError(response.ErrBadRequest, "upload timeout", 400)
		case <-ctx.Done():
			return nil, nil, response.NewCustomError(response.ErrBadRequest, ctx.Err().Error(), 400)
		}
	}

	return urls, pids, firstErr
}

// CreateProduct implements services.ProductService.
func (p *ProductService) CreateProduct(ctx context.Context, userId uuid.UUID, req dto.CreateProductRequest) (*entities.ProductEntity, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, response.NewCustomError(response.ErrBadRequest, "name is required", 400)
	}
	if req.Cover == nil {
		return nil, response.NewCustomError(response.ErrBadRequest, "cover image is required", 400)
	}
	if req.Status < 0 || req.Status > 1 {
		return nil, response.NewCustomError(response.ErrBadRequest, "status must be 0 or 1", 400)
	}

	nameSlug := Sanitize(req.Name)
	folder := "nusantara_service/product"

	var wg sync.WaitGroup
	var upErr error

	type uploadResult struct {
		coverURL string
		coverPID string
		galURLs  []string
		galPIDs  []string
	}
	resultsChan := make(chan uploadResult, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		ext := strings.ToLower(filepath.Ext(req.Cover.Filename))
		filename := fmt.Sprintf("cover_%s%s", nameSlug, ext)

		u, err := p.cloudinary.UploadImage(ctx, req.Cover, folder, filename)
		if err != nil {
			upErr = err
			return
		}
		resultsChan <- uploadResult{
			coverURL: u.URL,
			coverPID: BuildPublicID(folder, filename),
		}
	}()

	go func() {
		defer wg.Done()
		if len(req.Gallery) == 0 {
			resultsChan <- uploadResult{}
			return
		}
		u, p, err := p.UploadMany(ctx, folder, "gallery_"+nameSlug, req.Gallery, 4)
		if err != nil {
			upErr = err
			return
		}
		resultsChan <- uploadResult{
			galURLs: u,
			galPIDs: p,
		}
	}()

	wg.Wait()
	close(resultsChan)

	if upErr != nil {
		return nil, upErr
	}

	var coverUrl string
	var galPIDs []string
	var galUrls []string
	for r := range resultsChan {
		if r.coverURL != "" {
			coverUrl = r.coverURL
		}
		if len(r.galPIDs) > 0 {
			galPIDs = r.galPIDs
			galUrls = r.galURLs
		}
	}

	var createdID uuid.UUID
	err := configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		coverImg := entities.ImageEntity{
			ID:        uuid.New(),
			ImagePath: coverUrl,
		}
		if err := tx.Create(&coverImg).Error; err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create image", 500)
		}

		user, err := p.repo.FindByUserIDSuperAdmin(ctx, userId.String())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewCustomError(response.ErrNotFound, "user not found", 404)
			}
			return response.NewCustomError(response.ErrInternal, "failed to get user", 500)
		}

		typeProduct, err := p.typeProductRepo.FindById(ctx, req.TypeProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewCustomError(response.ErrNotFound, "type product not found", 404)
			}
			return response.NewCustomError(response.ErrInternal, "failed to get type product", 500)
		}

		product := &entities.ProductEntity{
			ID:            uuid.New(),
			Name:          req.Name,
			Code:          req.Code,
			Price:         req.Price,
			Unit:          req.Unit,
			Description:   req.Description,
			Status:        req.Status,
			TypeProductID: typeProduct.ID,
			TypeProduct:   *typeProduct,
			ImageID:       coverImg.ID,
			Image:         coverImg,
			CreatedBy:     userId,
			User:          *user,
		}

		if err := tx.Create(product).Error; err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to create product", 500)
		}
		createdID = product.ID

		if len(galPIDs) > 0 {
			imgs := make([]entities.ImageEntity, 0, len(galUrls))
			for _, url := range galUrls {
				imgs = append(imgs, entities.ImageEntity{
					ID:        uuid.New(),
					ImagePath: url,
				})
			}
			if err := tx.Create(&imgs).Error; err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to create gallery images", 500)
			}

			pivs := make([]entities.ProductImageEntity, 0, len(imgs))
			for _, im := range imgs {
				pivs = append(pivs, entities.ProductImageEntity{
					ID:        uuid.New(),
					ProductID: product.ID,
					Product:   *product,
					ImageID:   im.ID,
					Image:     im,
					AltText:   req.Name,
				})
			}
			if err := tx.Create(&pivs).Error; err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to create product images", 500)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = rabbitmq.PublishToQueue("", rabbitmq.CacheInvalidateQueueName, payload.CacheInvalidateTask{
		Keys: []string{"products:*"},
	})

	return p.repo.GetByID(ctx, createdID)
}

// GetProductByID implements services.ProductService.
func (p *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.ProductEntity, error) {
	cacheKey := fmt.Sprintf("product:%s", id.String())

	cached, err := configs.GetRedis(ctx, cacheKey)
	if err != nil && cached != "" {
		var product entities.ProductEntity
		if err := json.Unmarshal([]byte(cached), &product); err != nil {
			return &product, nil
		}
	}

	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "Product not found", 404)
		}
		return nil, err
	}

	dataCache, _ := json.Marshal(product)
	_ = configs.SetRedis(ctx, cacheKey, dataCache, time.Minute*30)

	return product, nil
}

// GetProductAll implements services.ProductService.
func (p *ProductService) GetProductAll(ctx context.Context, page int, limit int, search string) ([]*entities.ProductEntity, int, error) {
	cachedKey := fmt.Sprintf("products:search:%s:page:%d:limit:%d", search, page, limit)

	cached, err := configs.GetRedis(ctx, cachedKey)
	if err == nil {
		var result struct {
			Data  []*entities.ProductEntity `json:"data"`
			Total int                       `json:"total"`
		}

		_ = json.Unmarshal([]byte(cached), &result)
		return result.Data, result.Total, nil
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	products, total, err := p.repo.GetAll(ctx, offset, limit, search)
	if err != nil {
		return nil, 0, response.NewCustomError(response.ErrNotFound, "Product not found", 404)
	}

	dataCache, _ := json.Marshal(map[string]interface{}{
		"data":  products,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cachedKey, dataCache, time.Minute*30)
	return products, total, nil
}

// UpdateProduct implements services.ProductService.
// UpdateProduct implements services.ProductService.
func (p *ProductService) UpdateProduct(ctx context.Context, userId uuid.UUID, req dto.UpdateProductRequest) (*entities.ProductEntity, error) {
	product, err := p.repo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "product not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	// Update product fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Code != nil {
		product.Code = *req.Code
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.TypeProductID != nil {
		typeProduct, err := p.typeProductRepo.FindById(ctx, *req.TypeProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, response.NewCustomError(response.ErrNotFound, "type product not found", 404)
			}
			return nil, response.NewCustomError(response.ErrInternal, "failed to get type product", 500)
		}
		product.TypeProductID = typeProduct.ID
	}

	// Check if the user exists
	user, err := p.repo.FindByUserIDSuperAdmin(ctx, userId.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewCustomError(response.ErrNotFound, "user not found", 404)
		}
		return nil, response.NewCustomError(response.ErrInternal, "failed to get user", 500)
	}
	product.CreatedBy = uuid.MustParse(user.ID)

	nameSlug := Sanitize(product.Name)
	folder := "nusantara_service/product"
	var toDeletePIDs []string

	// Start a database transaction to ensure atomicity
	err = configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Handle cover image update
		if req.NewCover != nil {
			// Get public ID of old cover image for deletion
			if product.Image.ImagePath != "" {
				oldCoverPID := PublicIDFromURL(product.Image.ImagePath)
				if oldCoverPID != "" {
					toDeletePIDs = append(toDeletePIDs, oldCoverPID)
				}
			}

			// Upload new cover image
			ext := strings.ToLower(filepath.Ext(req.NewCover.Filename))
			filename := fmt.Sprintf("cover_%s%s", nameSlug, ext)
			u, err := p.cloudinary.UploadImage(ctx, req.NewCover, folder, filename)
			if err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to upload new cover image", 500)
			}

			// Delete old cover image entity and create new one
			if err := tx.Delete(&product.Image).Error; err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete old cover image", 500)
			}

			img := entities.ImageEntity{
				ID:        uuid.New(),
				ImagePath: u.URL,
			}
			if err := tx.Create(&img).Error; err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to create new cover image", 500)
			}
			product.ImageID = img.ID
		}

		// Handle gallery image update
		if req.ReplaceGallery {
			// Get public IDs of old gallery images for deletion
			oldPIDs, err := p.repo.GetProductImagePublicIDs(ctx, product.ID)
			if err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to get old product image IDs", 500)
			}
			toDeletePIDs = append(toDeletePIDs, oldPIDs...)

			// Delete old product_images and image entities
			if err := p.repo.DeleteProductImages(ctx, product.ID); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to delete old product images", 500)
			}
		}

		// Upload and create new gallery images if provided
		if len(req.NewGallery) > 0 {
			urls, pids, err := p.UploadMany(ctx, folder, "gallery_"+nameSlug, req.NewGallery, 4)
			if err != nil {
				return err
			}

			imgs := make([]entities.ImageEntity, 0, len(pids))
			for _, url := range urls {
				imgs = append(imgs, entities.ImageEntity{
					ID:        uuid.New(),
					ImagePath: url,
				})
			}
			if err := p.repo.CreateImages(ctx, imgs); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to create new gallery images", 500)
			}

			pivs := make([]entities.ProductImageEntity, 0, len(imgs))
			for _, im := range imgs {
				pivs = append(pivs, entities.ProductImageEntity{
					ID:        uuid.New(),
					ProductID: product.ID,
					ImageID:   im.ID,
					AltText:   product.Name,
				})
			}
			if err := p.repo.CreateProductImages(ctx, pivs); err != nil {
				return response.NewCustomError(response.ErrInternal, "failed to create product images", 500)
			}
		}

		// Save the updated product entity
		if err := p.repo.Update(ctx, product.ID, product); err != nil {
			return response.NewCustomError(response.ErrInternal, "failed to update product", 500)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Asynchronously delete old images from Cloudinary
	if len(toDeletePIDs) > 0 {
		_ = rabbitmq.PublishToQueue("", "image.delete.q", payload.ImageDeleteTask{PublicIDs: toDeletePIDs})
	}

	// Invalidate caches
	_ = rabbitmq.PublishToQueue("", rabbitmq.CacheInvalidateQueueName, payload.CacheInvalidateTask{
		Keys: []string{
			fmt.Sprintf("product:%s", product.ID),
			"products:*",
		},
	})

	return p.repo.GetByID(ctx, product.ID)
}

// Delete implements services.ProductService.
func (p *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewCustomError(response.ErrNotFound, "product not found", 404)
		}
		return response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	var pids []string
	if product.Image.ImagePath != "" {
		pids = append(pids, product.Image.ImagePath)
	}
	for _, pi := range product.ProductImages {
		if pi.Image.ImagePath != "" {
			pids = append(pids, pi.Image.ImagePath)
		}
	}

	if err := p.repo.Delete(ctx, id); err != nil {
		return response.NewCustomError(response.ErrInternal, err.Error(), 500)
	}

	if len(pids) > 0 {
		_ = rabbitmq.PublishToQueue("", "image.delete.q", payload.ImageDeleteTask{PublicIDs: pids})
	}

	_ = rabbitmq.PublishToQueue("", rabbitmq.CacheInvalidateQueueName, payload.CacheInvalidateTask{
		Keys: []string{
			"product:*",
			"products:*",
		},
	})

	return nil
}
