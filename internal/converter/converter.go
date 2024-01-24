package converter

import (
	"github.com/jinzhu/copier"
	"root-histoty-service/internal/dto"
	"root-histoty-service/internal/dto/request"
	"root-histoty-service/internal/model"
)

func CreatePlayerRequestToPlayer(from *request.CreatePlayerRequest) *model.Player {
	target := &model.Player{}

	err := copier.Copy(target, from)
	if err != nil {
		return nil
	}
	return target
}

func PlayerDTOtoPlayer(from *dto.PlayerDTO) *model.Player {
	target := &model.Player{}

	err := copier.Copy(target, from)
	if err != nil {
		return nil
	}
	return target
}

func PlayerToPlayerDTO(from *model.Player) *dto.PlayerDTO {
	target := &dto.PlayerDTO{}

	err := copier.Copy(target, from)
	if err != nil {
		return nil
	}
	return target
}
