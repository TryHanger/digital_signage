package repository

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) CreateTemplate(template *model.Template) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(template).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *TemplateRepository) GetAll() ([]model.Template, error) {
	var templates []model.Template
	err := r.db.Preload("Blocks").Preload("Blocks.Contents").Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) GetByID(id uint) (*model.Template, error) {
	var template model.Template
	err := r.db.Preload("Blocks").Preload("Blocks.Contents").First(&template, id).Error
	return &template, err
}

func (r *TemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.Template{}, id).Error
}

func (r *TemplateRepository) UpdateTemplate(template *model.Template) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// update top-level template fields
		if err := tx.Model(&model.Template{}).Where("id = ?", template.ID).Updates(map[string]interface{}{
			"name":        template.Name,
			"description": template.Description,
		}).Error; err != nil {
			return err
		}

		// load existing blocks (with contents) for diffing
		var existingBlocks []model.TemplateBlock
		if err := tx.Where("template_id = ?", template.ID).Preload("Contents").Find(&existingBlocks).Error; err != nil {
			return err
		}
		existingBlocksMap := make(map[uint]model.TemplateBlock)
		for _, eb := range existingBlocks {
			existingBlocksMap[eb.ID] = eb
		}

		// track incoming block IDs to find deletions
		incomingBlockIDs := make(map[uint]bool)

		// iterate incoming blocks: update existing ones, create new ones
		for i := range template.Blocks {
			b := template.Blocks[i]
			if b.ID != 0 {
				incomingBlockIDs[b.ID] = true
				// if block exists - update fields
				if _, ok := existingBlocksMap[b.ID]; ok {
					if err := tx.Model(&model.TemplateBlock{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
						"name":       b.Name,
						"start_time": b.StartTime,
						"end_time":   b.EndTime,
					}).Error; err != nil {
						return err
					}

					// handle contents for existing block
					// map existing contents by id
					existingContentsMap := make(map[uint]model.TemplateContent)
					for _, ec := range existingBlocksMap[b.ID].Contents {
						existingContentsMap[ec.ID] = ec
					}
					incomingContentIDs := make(map[uint]bool)

					// update or create incoming contents in order
					order := 1
					for j := range b.Contents {
						c := b.Contents[j]
						c.BlockID = b.ID
						c.Order = order
						order++
						if c.ID != 0 {
							incomingContentIDs[c.ID] = true
							if _, okc := existingContentsMap[c.ID]; okc {
								if err := tx.Model(&model.TemplateContent{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
									"content_id": c.ContentID,
									"duration":   c.Duration,
									"order":      c.Order,
									"type":       c.Type,
								}).Error; err != nil {
									return err
								}
							} else {
								// incoming has an id but DB doesn't know it -> insert as new
								c.ID = 0
								if err := tx.Create(&c).Error; err != nil {
									return err
								}
							}
						} else {
							// new content
							c.ID = 0
							if err := tx.Create(&c).Error; err != nil {
								return err
							}
						}
					}

					// delete contents that exist in DB but not in incoming payload
					for _, ec := range existingBlocksMap[b.ID].Contents {
						if !incomingContentIDs[ec.ID] {
							if err := tx.Delete(&model.TemplateContent{}, ec.ID).Error; err != nil {
								return err
							}
						}
					}
					continue
				}
				// if block.ID provided but not found in DB, fallthrough to create new block
			}

			// create new block (no ID or not found)
			newBlock := b
			newBlock.ID = 0
			newBlock.TemplateID = template.ID
			contents := newBlock.Contents
			newBlock.Contents = nil
			if err := tx.Omit("Contents").Create(&newBlock).Error; err != nil {
				return err
			}
			// insert contents for new block
			for j := range contents {
				c := contents[j]
				c.ID = 0
				c.BlockID = newBlock.ID
				if err := tx.Create(&c).Error; err != nil {
					return err
				}
			}
		}

		// delete blocks that exist in DB but missing in incoming payload
		for _, eb := range existingBlocks {
			if !incomingBlockIDs[eb.ID] {
				// delete contents and block
				if err := tx.Where("block_id = ?", eb.ID).Delete(&model.TemplateContent{}).Error; err != nil {
					return err
				}
				if err := tx.Delete(&model.TemplateBlock{}, eb.ID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
