package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) GetById(db *gorm.DB, dto *T, id any) error {
	return db.Where("id = ?", id).Take(dto).Error
}

func (r *Repository[T]) Create(db *gorm.DB, dto *T) error {
	return db.Create(dto).Error
}

func (r *Repository[T]) Update(db *gorm.DB, dto *T) error {
	return db.Save(dto).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, dto *T) error {
	return db.Delete(dto).Error
}

func (r *Repository[T]) GetCountById(db *gorm.DB, id any) (int64, error) {
	var count int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count, err
}
