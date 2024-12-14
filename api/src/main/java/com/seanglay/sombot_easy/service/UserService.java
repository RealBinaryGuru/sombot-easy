package com.seanglay.sombot_easy.service;

import com.seanglay.sombot_easy.dto.LoginDTO;
import com.seanglay.sombot_easy.dto.RegisterDTO;
import com.seanglay.sombot_easy.mapper.UserMapper;
import com.seanglay.sombot_easy.model.User;
import com.seanglay.sombot_easy.repository.UserRepository;
import com.seanglay.sombot_easy.util.JWTToken;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserService {
    @Autowired
    private UserRepository userRepository;

    @Autowired
    private UserMapper userMapper;

    @Autowired
    private JWTToken jwtToken;

    public String register(RegisterDTO registerDTO) {
        User user = userMapper.registerDtoToUser(registerDTO);
         userRepository.save(user);
         return jwtToken.generateToken(user);
    }

    public String login(LoginDTO loginDTO) {
        User user = userRepository.findByEmail(loginDTO.getEmail()).orElseThrow(() -> new RuntimeException("Invalid credentials"));
        return jwtToken.generateToken(user);
    }
}
