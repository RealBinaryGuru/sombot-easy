package com.seanglay.sombot_easy.controller;

import com.seanglay.sombot_easy.dto.APIResponse;
import com.seanglay.sombot_easy.dto.LoginDTO;
import com.seanglay.sombot_easy.dto.RegisterDTO;
import com.seanglay.sombot_easy.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/auth")
public class Authentication extends BaseController {

    @Autowired
    private UserService userService;

    @PostMapping("/register")
    public ResponseEntity<APIResponse<String>> register(@RequestBody RegisterDTO registerDTO) {
        return ok(userService.register(registerDTO));
    }

    @PostMapping("/login")
    public ResponseEntity<APIResponse<String>> login(@RequestBody LoginDTO loginDTO) {
        return ok(userService.login(loginDTO));
    }
}
