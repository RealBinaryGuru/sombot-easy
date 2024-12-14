package com.seanglay.sombot_easy.service;

import com.seanglay.sombot_easy.dto.LoginDTO;
import com.seanglay.sombot_easy.dto.RegisterDTO;
import com.seanglay.sombot_easy.mapper.UserMapper;
import com.seanglay.sombot_easy.model.User;
import com.seanglay.sombot_easy.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserService {
    @Autowired
    private UserRepository userRepository;

    @Autowired
    private UserMapper userMapper;

    public User register(RegisterDTO registerDTO) {
        User user = userMapper.registerDtoToUser(registerDTO);
        return userRepository.save(user);
    }

    public User login(LoginDTO loginDTO) {
        return userRepository.findByEmail(loginDTO.getEmail()).orElseThrow(() -> new RuntimeException("Invalid credentials"));
    }
}
