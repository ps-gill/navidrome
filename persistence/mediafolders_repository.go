package persistence

import (
	"context"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/model"
)

type mediaFolderRepository struct {
	ctx context.Context
}

func NewMediaFolderRepository(ctx context.Context, o orm.QueryExecutor) model.MediaFolderRepository {
	return &mediaFolderRepository{ctx}
}

func (r *mediaFolderRepository) Get(id int32) (*model.MediaFolder, error) {
	mediaFolders := hardCoded()
	if mediaFolder, ok := mediaFolders[id]; ok {
		return &mediaFolder, nil
	}
	return nil, fmt.Errorf("media folder with id '%d' not found", id)
}

func (*mediaFolderRepository) GetAll() (model.MediaFolders, error) {
	mediaFolders := hardCoded()
	result := make(model.MediaFolders, len(mediaFolders))
	for i, mediaFolder := range mediaFolders {
		result[i] = mediaFolder
	}
	return result, nil
}

func hardCoded() map[int32]model.MediaFolder {
	mediaFolders := make(map[int32]model.MediaFolder, len(conf.Server.MusicFolder))
	for index, musicFolder := range conf.Server.MusicFolder {
		// unlikely that the number of music folders with be greater than max of int32
		mediaFolder := model.MediaFolder{ID: int32(index), Path: musicFolder}
		mediaFolder.Name = "Music Library"
		mediaFolders[mediaFolder.ID] = mediaFolder
	}
	return mediaFolders
}

var _ model.MediaFolderRepository = (*mediaFolderRepository)(nil)
