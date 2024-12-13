package com.seanglay.sombot_easy.mapper;

import com.seanglay.sombot_easy.dto.RegisterDTO;
import com.seanglay.sombot_easy.model.User;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

@Mapper(componentModel = "spring")
public interface UserMapper {
    UserMapper INSTANCE = Mappers.getMapper(UserMapper.class);

    @Mapping(target = "userId", ignore = true)
    User registerDtoToUser(RegisterDTO registerDTO);
}
